package main

import (
	"github.com/VanillaSixtySix/yuriposting/internal/yuriposting"
	"github.com/VanillaSixtySix/yuriposting/internal/yuriposting/bluesky"
	"github.com/VanillaSixtySix/yuriposting/internal/yuriposting/danbooru"
	"github.com/VanillaSixtySix/yuriposting/internal/yuriposting/mastodon"
	"log"
	"os"
)

func main() {
	config, err := yuriposting.LoadConfig("config.json")
	if err != nil {
		log.Fatalln("Failed to load config:", err.Error())
	}
	danbooruAPI := danbooru.NewDanbooruAPI(config)
	blueskyAPI := bluesky.NewBlueskyAPI(config)
	mastodonAPI := mastodon.NewMastodonAPI(config)

	log.Println("[Danbooru] Searching for tags:", config.DanbooruTags)
	post, err := danbooruAPI.GetRandomPost()
	if err != nil {
		log.Fatalln("[Danbooru] Failed to fetch random post:", err.Error())
	}
	postFileSizeTooLarge := post.FileSize > 1000000
	postResTooHigh := post.ImageWidth > 4096 && post.ImageHeight > 4096
	var (
		lqImg       *os.File
		hqImg       *os.File
		fileName    string
		contentType string
	)
	if postFileSizeTooLarge || postResTooHigh {
		lqImg, fileName, contentType, err = danbooruAPI.GetPostImage(post, false)
		if err != nil {
			log.Fatalln("[Danbooru] Failed to fetch LQ post image:", err.Error())
		}
	}
	log.Println("[Danbooru] Fetching HQ image for post ID:", post.Id)
	hqImg, fileName, contentType, err = danbooruAPI.GetPostImage(post, true)
	if err != nil {
		log.Fatalln("[Danbooru] Failed to fetch HQ post image:", err.Error())
	}

	if config.PostToMastodon {
		toUpload := hqImg
		if postResTooHigh {
			toUpload = lqImg
		}
		log.Println("[Mastodon] Uploading media...")
		media, err := mastodonAPI.UploadMedia(toUpload, fileName, post.TagString)
		if err != nil {
			log.Fatalln("[Mastodon] Failed to upload media:", err.Error())
		}
		log.Println("[Mastodon] Media ID:", media.Id)
		if err = mastodonAPI.CreateStatusFromPost(post, media); err != nil {
			log.Fatalln("[Mastodon] Failed to create status with post and media:", err.Error())
		}

		log.Println("[Mastodon] Success!")
	}

	if config.PostToBluesky {
		log.Println("[Bluesky] Creating session:", config.BlueskyIdentifier)
		session, err := blueskyAPI.CreateSession()
		if err != nil {
			log.Fatalln("[Bluesky] Failed to create session:", err.Error())
		}
		toUpload := hqImg
		if postFileSizeTooLarge {
			toUpload = lqImg
		}
		log.Println("[Bluesky] Creating blob...")
		blob, err := blueskyAPI.UploadBlob(session, toUpload, contentType)
		if err != nil {
			log.Fatalln("[Bluesky] Failed to create blob:", err.Error())
		}
		log.Println("[Bluesky] Creating record...")
		createdRecord, err := blueskyAPI.CreateRecordFromPost(post, blob, session)
		if err != nil {
			log.Fatalln("[Bluesky] Failed to create record:", err.Error())
		}

		log.Println("[Bluesky] Success! URI: " + createdRecord.URI)
	}

	_ = (*hqImg).Close()
	err = os.Remove((*hqImg).Name())
	if err != nil {
		log.Fatalln("Failed to remove temporary HQ image")
	}
	if lqImg != nil {
		_ = (*lqImg).Close()
		err = os.Remove((*lqImg).Name())
		if err != nil {
			log.Fatalln("Failed to remove temporary LQ image")
		}
	}
}

package main

import (
	"github.com/VanillaSixtySix/yuriposting/internal/yuriposting"
	"github.com/VanillaSixtySix/yuriposting/internal/yuriposting/bluesky"
	"github.com/VanillaSixtySix/yuriposting/internal/yuriposting/danbooru"
	"github.com/VanillaSixtySix/yuriposting/internal/yuriposting/mastodon"
	"log"
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
	log.Println("[Danbooru] Fetching image for post ID:", post.Id)
	img, fileName, contentType, err := danbooruAPI.GetPostImage(post)
	if err != nil {
		log.Fatalln("[Danbooru] Failed to fetch post image:", err.Error())
	}

	if config.PostToMastodon {
		log.Println("[Mastodon] Uploading media...")
		media, err := mastodonAPI.UploadMedia(img, fileName, post.TagString)
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
		log.Println("[Bluesky] Creating blob...")
		blob, err := blueskyAPI.UploadBlob(session, img, contentType)
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
}

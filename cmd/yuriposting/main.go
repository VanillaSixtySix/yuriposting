package main

import (
	"github.com/VanillaSixtySix/yuriposting/internal/yuriposting"
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
	mastodonAPI := mastodon.NewMastodonAPI(config)

	log.Println("Searching for tags:", config.Tags)
	post, err := danbooruAPI.GetRandomPost()
	if err != nil {
		log.Fatalln("Failed to fetch random Danbooru post:", err.Error())
	}
	log.Println("Post ID:", post.Id)
	img, fileName, err := danbooruAPI.GetPostImage(post)
	if err != nil {
		log.Fatalln("Failed to fetch Danbooru post image:", err.Error())
	}

	log.Println("Uploading media...")
	media, err := mastodonAPI.UploadMedia(img, fileName, post.TagString)
	if err != nil {
		log.Fatalln("Failed to upload media:", err.Error())
	}
	log.Println("Media ID:", media.Id)
	if err = mastodonAPI.CreateStatusFromPost(post, media); err != nil {
		log.Fatalln("Failed to create status with post and media:", err.Error())
	}

	log.Println("Success!")
}

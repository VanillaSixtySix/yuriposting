package main

import (
	"io"
	"log"
	"os"
	"yuriposting/internal/yuriposting"
	"yuriposting/internal/yuriposting/danbooru"
	"yuriposting/internal/yuriposting/mastodon"
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
	if seeker, ok := (*img).(io.Seeker); ok {
		file, err := os.Create(fileName)
		if err != nil {
			log.Fatalln("Failed to create local file:", err.Error())
		}
		if _, err = io.Copy(file, *img); err != nil {
			log.Fatalln("Failed to write local file:", err.Error())
		}
		if err = file.Close(); err != nil {
			log.Fatalln("Failed to close local file:", err.Error())
		}
		if _, err = seeker.Seek(0, io.SeekStart); err != nil {
			log.Fatalln("Failed to seek start of img buffer:", err.Error())
		}
	}

	log.Println("Uploading media...")
	media, err := mastodonAPI.UploadMedia(img, fileName, post.TagString)
	if err != nil {
		return
	}
	log.Println("Media ID:", media.Id)
	if err = mastodonAPI.CreateStatusFromPost(post, media); err != nil {
		log.Fatalln("Failed to create status with post and media:", err.Error())
	}

	log.Println("Success!")
}

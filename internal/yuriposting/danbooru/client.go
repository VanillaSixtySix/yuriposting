package danbooru

import (
	"encoding/json"
	"errors"
	"fmt"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"path"
	"strings"
	"yuriposting/internal/yuriposting"
)

type API struct {
	config     *yuriposting.Config
	authParams string
}

func NewDanbooruAPI(config *yuriposting.Config) *API {
	authParams := fmt.Sprintf("api_key=%s&login=%s", config.DanbooruAPIKey, config.DanbooruUsername)
	return &API{
		config:     config,
		authParams: authParams,
	}
}

func (api *API) GetRandomPost() (*Post, error) {
	tags := api.config.Tags
	searchTags := api.cleanParam(tags)
	reqUrl := fmt.Sprintf("https://danbooru.donmai.us/posts.json?%s&tags=%s&limit=1", api.authParams, searchTags)
	log.Println("URL to GET:", reqUrl)
	res, err := http.Get(reqUrl)
	if err != nil {
		return nil, err
	}
	log.Println("Status code:", res.StatusCode)
	if res.StatusCode != 200 {
		return nil, errors.New("status code not 200")
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("failed to read body")
	}
	posts := make([]Post, 1)
	if err = json.Unmarshal(bodyBytes, &posts); err != nil {
		return nil, err
	}
	if len(posts) < 1 {
		return nil, errors.New("no results for tags")
	}
	log.Println("Received", len(posts), "post(s)")
	return &posts[0], nil
}

func (api *API) GetPostImage(post *Post) (*io.ReadCloser, string, error) {
	fileName := path.Base(post.FileUrl)
	res, err := http.Get(post.FileUrl)
	if err != nil {
		return nil, fileName, err
	}
	log.Println("Status code:", res.StatusCode)
	if res.StatusCode != 200 && res.StatusCode != 202 {
		return nil, fileName, errors.New("status code not 200 or 202")
	}
	return &res.Body, fileName, nil
}

func (api *API) cleanParam(param string) string {
	param = strings.ReplaceAll(param, " ", "+")
	param = strings.ReplaceAll(param, "<", "%3C")
	param = strings.ReplaceAll(param, ">", "%3E")
	return param
}

func FormatTags(tags string) string {
	split := strings.Split(tags, " ")
	for i, tag := range split {
		split[i] = strings.ReplaceAll(tag, "_", " ")
	}
	return strings.Join(split, ", ")
}

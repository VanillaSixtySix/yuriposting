package mastodon

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/VanillaSixtySix/yuriposting/internal/yuriposting"
	"github.com/VanillaSixtySix/yuriposting/internal/yuriposting/danbooru"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type API struct {
	config     *yuriposting.Config
	httpClient *http.Client
}

func NewMastodonAPI(config *yuriposting.Config) *API {
	httpClient := &http.Client{}
	return &API{
		config:     config,
		httpClient: httpClient,
	}
}

func (api *API) UploadMedia(media *os.File, fileName string, tags string) (*UploadedMediaResponse, error) {
	var writerBuf bytes.Buffer
	writer := multipart.NewWriter(&writerBuf)
	filePart, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(filePart, media); err != nil {
		return nil, err
	}
	_, err = media.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	descriptionPart, err := writer.CreateFormField("description")
	if err != nil {
		return nil, err
	}
	if _, err = descriptionPart.Write([]byte("Danbooru tags: " + tags)); err != nil {
		return nil, err
	}

	if err = writer.Close(); err != nil {
		return nil, err
	}

	reqUrl := api.config.MastodonHost + "/api/v2/media"
	req, err := http.NewRequest("POST", reqUrl, &writerBuf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+api.config.MastodonAccessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	log.Println("URL to POST:", reqUrl)
	res, err := api.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("failed to read body: " + err.Error())
	}
	if err = res.Body.Close(); err != nil {
		return nil, err
	}
	if res.StatusCode != 200 && res.StatusCode != 202 {
		bodyStr := string(bodyBytes)
		return nil, errors.New("status code not 200 or 202: " + bodyStr)
	}
	var uploadedMedia UploadedMediaResponse
	if err = json.Unmarshal(bodyBytes, &uploadedMedia); err != nil {
		return nil, err
	}

	return &uploadedMedia, nil
}

func (api *API) CreateStatusFromPost(post *danbooru.Post, uploadedMedia *UploadedMediaResponse) error {
	source := post.Source
	if post.PixivId != nil {
		source = "https://www.pixiv.net/en/artworks/" + strconv.Itoa(*post.PixivId)
	}

	artists := danbooru.FormatTags(post.TagStringArtist)
	copyrights := danbooru.FormatTags(post.TagStringCopyright)
	postFormat := fmt.Sprintf(
		"Artist%s: %s\nMedia: %s\nSource: %s",
		yuriposting.Pluralize(post.TagCountArtist), artists, copyrights, source)

	isSensitive := post.Rating != "g"

	var writerBuf bytes.Buffer
	writer := multipart.NewWriter(&writerBuf)
	if err := writer.WriteField("status", postFormat); err != nil {
		return err
	}
	if err := writer.WriteField("visibility", api.config.MastodonPostVisibility); err != nil {
		return err
	}
	if err := writer.WriteField("media_ids[]", uploadedMedia.Id); err != nil {
		return err
	}
	if err := writer.WriteField("sensitive", strconv.FormatBool(isSensitive)); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	reqUrl := api.config.MastodonHost + "/api/v1/statuses"
	req, err := http.NewRequest("POST", reqUrl, &writerBuf)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+api.config.MastodonAccessToken)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	log.Println("URL to POST:", reqUrl)
	res, err := api.httpClient.Do(req)
	if err != nil {
		return err
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.New("failed to read body: " + err.Error())
	}
	if err = res.Body.Close(); err != nil {
		return errors.New("failed to close body: " + err.Error())
	}
	if res.StatusCode != 200 {
		bodyStr := string(bodyBytes)
		return errors.New("status code not 200: " + bodyStr)
	}
	return nil
}

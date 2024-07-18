package mastodon

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"yuriposting/internal/yuriposting"
	"yuriposting/internal/yuriposting/danbooru"
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

func (api *API) UploadMedia(media *io.ReadCloser, fileName string, tags string) (*UploadedMediaResponse, error) {
	var writerBuf bytes.Buffer
	writer := multipart.NewWriter(&writerBuf)
	filePart, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return nil, err
	}
	if _, err = io.Copy(filePart, *media); err != nil {
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

	if err = (*media).Close(); err != nil {
		return nil, err
	}

	reqUrl := "https://botsin.space/api/v2/media"
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
	log.Println("Status code:", res.StatusCode)
	if res.StatusCode != 200 && res.StatusCode != 202 {
		return nil, errors.New("status code not 200 or 202")
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("failed to read body")
	}
	if err = res.Body.Close(); err != nil {
		return nil, err
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
		pluralize(post.TagCountArtist), artists, copyrights, source)

	isSensitive := post.Rating != "g"

	var writerBuf bytes.Buffer
	writer := multipart.NewWriter(&writerBuf)
	if err := writer.WriteField("status", postFormat); err != nil {
		return err
	}
	if err := writer.WriteField("visibility", api.config.Visibility); err != nil {
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

	reqUrl := "https://botsin.space/api/v1/statuses"
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
	if err = res.Body.Close(); err != nil {
		return err
	}
	log.Println("Status code:", res.StatusCode)
	if res.StatusCode != 200 {
		return errors.New("status code not 200")
	}
	return nil
}

func pluralize(num int) string {
	if num > 1 {
		return "s"
	}
	return ""
}

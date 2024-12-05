package bluesky

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/VanillaSixtySix/yuriposting/internal/yuriposting"
	"github.com/VanillaSixtySix/yuriposting/internal/yuriposting/danbooru"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type API struct {
	config     *yuriposting.Config
	httpClient *http.Client
}

func NewBlueskyAPI(config *yuriposting.Config) *API {
	httpClient := &http.Client{}
	return &API{
		config:     config,
		httpClient: httpClient,
	}
}

func (api *API) CreateSession() (*Session, error) {
	sessionBody := &CreateSessionBody{
		Identifier: api.config.BlueskyIdentifier,
		Password:   api.config.BlueskyAppPassword,
	}
	body, err := json.Marshal(sessionBody)
	if err != nil {
		return nil, errors.New("failed to marshal JSON: " + err.Error())
	}
	bodyReader := bytes.NewReader(body)
	reqUrl := "https://bsky.social/xrpc/com.atproto.server.createSession"
	res, err := http.Post(reqUrl, "application/json", bodyReader)
	if err != nil {
		return nil, err
	}
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("failed to read body: " + err.Error())
	}
	if err = res.Body.Close(); err != nil {
		return nil, errors.New("failed to close body: " + err.Error())
	}
	if res.StatusCode != 200 {
		bodyStr := string(bodyBytes)
		return nil, errors.New("status code not 200: " + bodyStr)
	}
	var session Session
	if err = json.Unmarshal(bodyBytes, &session); err != nil {
		return nil, err
	}
	return &session, err
}

func (api *API) UploadBlob(session *Session, media *io.ReadCloser, contentType string) (*Blob, error) {
	reqUrl := "https://bsky.social/xrpc/com.atproto.repo.uploadBlob"
	req, err := http.NewRequest("POST", reqUrl, *media)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+session.AccessJwt)
	req.Header.Set("Content-Type", contentType)
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
		return nil, errors.New("failed to close body: " + err.Error())
	}
	if res.StatusCode != 200 {
		bodyStr := string(bodyBytes)
		return nil, errors.New("status code not 200: " + bodyStr)
	}
	var uploadedBlob UploadedBlob
	if err = json.Unmarshal(bodyBytes, &uploadedBlob); err != nil {
		return nil, err
	}
	return &uploadedBlob.Blob, nil
}

func (api *API) CreateRecordFromPost(post *danbooru.Post, blob *Blob, session *Session) (*CreatedRecord, error) {
	source := post.Source
	if post.PixivId != nil {
		source = "https://www.pixiv.net/en/artworks/" + strconv.Itoa(*post.PixivId)
	}

	artists := danbooru.FormatTags(post.TagStringArtist)
	copyrights := danbooru.FormatTags(post.TagStringCopyright)
	postFormat := fmt.Sprintf(
		"Artist%s: %s\nMedia: %s\nSource: %s",
		yuriposting.Pluralize(post.TagCountArtist), artists, copyrights, source)

	// TODO: Implement sensitivity based on post rating

	facets := make([]Facet, 0)

	if strings.HasPrefix(source, "http") {
		formatSourceIndex := strings.Index(postFormat, "Source: ") + 8

		sourceLinkByteSlice := ByteSlice{
			ByteStart: formatSourceIndex,
			ByteEnd:   formatSourceIndex + len(source),
		}

		feature := Feature{
			Type: "app.bsky.richtext.facet#link",
			URI:  &source,
		}

		facet := Facet{
			Index:    sourceLinkByteSlice,
			Features: []Feature{feature},
		}

		facets = append(facets, facet)
	}

	now := time.Now().Format(time.RFC3339)

	recordBody := &CreateRecordBody{
		Repo:       api.config.BlueskyIdentifier,
		Collection: "app.bsky.feed.post",
		Record: Record{
			Type:      "app.bsky.feed.post",
			Text:      postFormat,
			Langs:     []string{"en"},
			CreatedAt: now,
			Embed: Embed{
				Type: "app.bsky.embed.images",
				Images: []EmbedImage{
					{
						Alt:   "Danbooru tags: " + post.TagString,
						Image: *blob,
					},
				},
			},
			Facets: facets,
		},
	}
	body, err := json.Marshal(recordBody)
	if err != nil {
		return nil, errors.New("failed to marshal JSON: " + err.Error())
	}
	bodyReader := bytes.NewReader(body)
	reqUrl := "https://bsky.social/xrpc/com.atproto.repo.createRecord"
	req, err := http.NewRequest("POST", reqUrl, bodyReader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+session.AccessJwt)
	req.Header.Set("Content-Type", "application/json")
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
		return nil, errors.New("failed to close body: " + err.Error())
	}
	if res.StatusCode != 200 {
		bodyStr := string(bodyBytes)
		return nil, errors.New("status code not 200: " + bodyStr)
	}
	var createdRecord CreatedRecord
	if err = json.Unmarshal(bodyBytes, &createdRecord); err != nil {
		return nil, err
	}
	return &createdRecord, nil
}

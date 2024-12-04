package bluesky

import (
	"encoding/json"
	"errors"
	"github.com/VanillaSixtySix/yuriposting/internal/yuriposting"
	"io"
	"log"
	"net/http"
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

func (api *API) UploadBlob(media *io.ReadCloser) (*UploadedBlobResponse, error) {
	reqUrl := api.config.BlueskyHost + "/xrpc/com.atproto.repo.uploadBlob"
	req, err := http.NewRequest("POST", reqUrl, *media)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer ")
	req.Header.Set("Content-Type", "application/octet-stream")
	log.Println("URL to POST:", reqUrl)
	res, err := api.httpClient.Do(req)
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
	if err = res.Body.Close(); err != nil {
		return nil, err
	}
	var uploadedBlob UploadedBlobResponse
	if err = json.Unmarshal(bodyBytes, &uploadedBlob); err != nil {
		return nil, err
	}

	return &uploadedBlob, nil
}

package mastodon

type UploadMediaBody struct {
	file        []byte
	description string
}

type UploadedMediaResponse struct {
	Id               string      `json:"id"`
	Type             string      `json:"type"`
	Url              string      `json:"url"`
	PreviewUrl       string      `json:"preview_url"`
	RemoteUrl        interface{} `json:"remote_url"`
	PreviewRemoteUrl interface{} `json:"preview_remote_url"`
	TextUrl          interface{} `json:"text_url"`
	Meta             struct {
		Original UploadMediaBody `json:"original"`
		Small    UploadMediaBody `json:"small"`
	} `json:"meta"`
	Description string `json:"description"`
	BlurHash    string `json:"blurhash"`
}

type UploadedMediaMeta struct {
	Width  int     `json:"width"`
	Height int     `json:"height"`
	Size   string  `json:"size"`
	Aspect float64 `json:"aspect"`
}

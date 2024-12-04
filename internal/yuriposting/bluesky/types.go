package bluesky

type CreateSessionBody struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type CreateSessionResponse struct {
	AccessJwt string `json:"accessJwt"`
}

type UploadedBlobResponse struct {
	Blob Blob `json:"blob"`
}

type Blob struct {
	Type     string  `json:"$type"`
	Ref      BlobRef `json:"ref"`
	MimeType string  `json:"mimeType"`
	Size     int     `json:"size"`
}

type BlobRef struct {
	Link string `json:"$link"`
}

type CreateRecordBody struct {
	Repo       string `json:"repo"`
	Collection string `json:"collection"`
	Record     Record `json:"record"`
}

type Record struct {
	Type      string   `json:"$type"`
	Text      string   `json:"text"`
	Langs     []string `json:"langs"`
	CreatedAt string   `json:"createdAt"`
	Embed     Embed    `json:"embed"`
}

type Embed struct {
	Type   string       `json:"$type"`
	Images []EmbedImage `json:"images"`
}

type EmbedImage struct {
	Alt   string `json:"alt"`
	Image Blob   `json:"image"`
}

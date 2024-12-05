package bluesky

type CreateSessionBody struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type Session struct {
	DID             string `json:"did"`
	Handle          string `json:"handle"`
	Email           string `json:"email"`
	EmailConfirmed  bool   `json:"emailConfirmed"`
	EmailAuthFactor bool   `json:"emailAuthFactor"`
	AccessJwt       string `json:"accessJwt"`
	RefreshJwt      string `json:"refreshJwt"`
	Active          bool   `json:"active"`
}

type DIDDoc struct {
	Context            []string                `json:"@context"`
	Id                 string                  `json:"id"`
	AlsoKnownAs        []string                `json:"alsoKnownAs"`
	VerificationMethod []DIDVerificationMethod `json:"VerificationMethod"`
	Service            []DIDService            `json:"service"`
}

type DIDVerificationMethod struct {
	Id                 string `json:"id"`
	Type               string `json:"type"`
	Controller         string `json:"controller"`
	PublicKeyMultibase string `json:"publicKeyMultibase"`
}

type DIDService struct {
	Id              string `json:"id"`
	Type            string `json:"type"`
	ServiceEndpoint string `json:"serviceEndpoint"`
}

type UploadedBlob struct {
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
	Facets    []Facet  `json:"facets,omitempty"`
}

type Embed struct {
	Type   string       `json:"$type"`
	Images []EmbedImage `json:"images"`
}

type EmbedImage struct {
	Alt   string `json:"alt"`
	Image Blob   `json:"image"`
}

type Facet struct {
	Index    ByteSlice `json:"index"`
	Features []Feature `json:"features"`
}

type ByteSlice struct {
	ByteStart int `json:"byteStart"`
	ByteEnd   int `json:"byteEnd"`
}

type Feature struct {
	Type string `json:"$type"`
	// Mention
	DID *string `json:"did,omitempty"`
	// Link
	URI *string `json:"uri,omitempty"`
	Tag *string `json:"tag,omitempty"`
}

type CreatedRecord struct {
	URI    string `json:"uri"`
	CID    string `json:"cid"`
	Commit struct {
		CID string `json:"cid"`
		Rev string `json:"rev"`
	} `json:"commit"`
	ValidationStatus string `json:"validationStatus"`
}

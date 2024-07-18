package danbooru

type Post struct {
	Id                  int         `json:"id"`
	CreatedAt           string      `json:"created_at"`
	UploaderId          int         `json:"uploader_id"`
	Score               int         `json:"score"`
	Source              string      `json:"source"`
	MD5                 string      `json:"md5"`
	LastCommentBumpedAt string      `json:"last_comment_bumped_at"`
	Rating              string      `json:"rating"`
	ImageWidth          int         `json:"image_width"`
	ImageHeight         int         `json:"image_height"`
	TagString           string      `json:"tag_string"`
	FavoriteCount       int         `json:"fav_count"`
	FileExtension       string      `json:"file_ext"`
	LastNotedAt         *string     `json:"last_noted_at"`
	ParentId            *int        `json:"parent_id"`
	HasChildren         bool        `json:"has_children"`
	ApproverId          *int        `json:"approver_id"`
	TagCountGeneral     int         `json:"tag_count_general"`
	TagCountArtist      int         `json:"tag_count_artist"`
	TagCountCharacter   int         `json:"tag_count_character"`
	TagCountCopyright   int         `json:"tag_count_copyright"`
	FileSize            int         `json:"file_size"`
	UpScore             int         `json:"up_score"`
	DownScore           int         `json:"down_score"`
	IsPending           bool        `json:"is_pending"`
	IsFlagged           bool        `json:"is_flagged"`
	IsDeleted           bool        `json:"is_deleted"`
	TagCount            int         `json:"tag_count"`
	UpdatedAt           string      `json:"updated_at"`
	IsBanned            bool        `json:"is_banned"`
	PixivId             *int        `json:"pixiv_id"`
	LastCommentedAt     string      `json:"last_commented_at"`
	HasActiveChildren   bool        `json:"has_active_children"`
	BitFlags            int         `json:"bit_flags"`
	TagCountMeta        int         `json:"tag_count_meta"`
	HasLarge            bool        `json:"has_large"`
	HasVisibleChildren  bool        `json:"has_visible_children"`
	MediaAsset          *MediaAsset `json:"media_asset"`
	TagStringGeneral    string      `json:"tag_string_general"`
	TagStringCharacter  string      `json:"tag_string_character"`
	TagStringCopyright  string      `json:"tag_string_copyright"`
	TagStringArtist     string      `json:"tag_string_artist"`
	TagStringMeta       string      `json:"tag_string_meta"`
	FileUrl             string      `json:"file_url"`
	LargeFileUrl        string      `json:"large_file_url"`
	PreviewFileUrl      string      `json:"preview_file_url"`
}

type MediaAsset struct {
	Id            int                  `json:"id"`
	CreatedAt     string               `json:"created_at"`
	UpdatedAt     string               `json:"updated_at"`
	MD5           string               `json:"md5"`
	FileExtension string               `json:"file_ext"`
	FileSize      int                  `json:"file_size"`
	ImageWidth    int                  `json:"image_width"`
	ImageHeight   int                  `json:"image_height"`
	Duration      *int                 `json:"duration"`
	Status        string               `json:"status"`
	FileKey       string               `json:"file_key"`
	IsPublic      bool                 `json:"is_public"`
	PixelHash     string               `json:"pixel_hash"`
	Variants      *[]MediaAssetVariant `json:"variants"`
}

type MediaAssetVariant struct {
	Type    string `json:"type"`
	Url     string `json:"url"`
	Width   int    `json:"width"`
	Height  int    `json:"height"`
	FileExt string `json:"file_ext"`
}

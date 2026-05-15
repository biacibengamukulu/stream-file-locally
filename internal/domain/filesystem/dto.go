package filesystem

type FileDTO struct {
	Name          string `json:"name"`
	ContentType   string `json:"content_type"`
	Base64Content string `json:"base64_content"`
}

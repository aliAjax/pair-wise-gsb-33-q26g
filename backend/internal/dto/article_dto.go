package dto

// ArticleCreateRequest is the payload for creating/updating articles.
type ArticleCreateRequest struct {
	Title    string `json:"title" binding:"required,max=255"`
	Content  string `json:"content" binding:"required"`
	Cover    string `json:"cover" binding:"omitempty,max=255"`
	TopicTag string `json:"topic_tag" binding:"required"`
}

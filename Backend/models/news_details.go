// models/news.go
package models

type News struct {
	ID          int    `json:"id" bson:"_id,omitempty"` // works for MongoDB or SQL
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
	Content     string `json:"content" bson:"content"`
	Links       []Link `json:"links" bson:"links"`
}

type Link struct {
	Name string `json:"name" bson:"name"`
	URL  string `json:"url" bson:"url"`
}

package domain

import "gorm.io/gorm"

type ResourceType uint8

const (
	Book ResourceType = iota
	Article
	Blog
	Video
	Song
	Movie
	Generic
)

type Resource struct {
	gorm.Model
	Url          string       `json:"url" form:"url"`
	ResourceType ResourceType `json:"resource_type" form:"resource_type"`
}

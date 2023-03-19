package domain

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type Tag struct {
	ID    uint   `gorm:"primarykey"`
	Name  string `json:"chapter" gorm:"not null;size:100;unique"`
	Ideas []Idea `json:"ideas" gorm:"many2many:ideas_tags"`
}

type IdeasTags struct {
	IdeaID uint `gorm:"primaryKey"`
	TagID  uint `gorm:"primaryKey"`
}

type TagQuery struct {
	ID  uint   `json:"id"`
	And []uint `json:"and"`
	Or  []uint `json:"or"`
	Not []uint `json:"not"`
}

func NewTagQuery(ctx *gin.Context) TagQuery {
	tq := TagQuery{}

	id := ctx.Query("id")
	if id != "" {
		id, _ := strconv.ParseInt(id, 0, 32)
		tq.ID = uint(id)
	}

	andIDs := ctx.Query("and")
	if andIDs != "" {
		IDs := strings.Split(andIDs, "-")
		uintIDs := make([]uint, len(IDs))

		for i := 0; i < len(IDs); i++ {
			intID, _ := strconv.ParseInt(IDs[i], 0, 32)
			uintIDs[i] = uint(intID)
		}

		tq.And = uintIDs
	}

	orIDs := ctx.Query("or")
	if orIDs != "" {
		IDs := strings.Split(orIDs, "-")
		uintIDs := make([]uint, len(IDs))

		for i := 0; i < len(IDs); i++ {
			intID, _ := strconv.ParseInt(IDs[i], 0, 32)
			uintIDs[i] = uint(intID)
		}

		tq.Or = uintIDs
	}

	notIDs := ctx.Query("not")
	if notIDs != "" {
		IDs := strings.Split(notIDs, "-")
		uintIDs := make([]uint, len(IDs))

		for i := 0; i < len(IDs); i++ {
			intID, _ := strconv.ParseInt(IDs[i], 0, 32)
			uintIDs[i] = uint(intID)
		}

		tq.Not = uintIDs
	}

	return tq
}

type TagRepository interface {
	Create(tag *Tag) error
	FetchAll() []Tag
	FetchById(ID string) Tag
	FetchByName(name string) Tag
	DeleteByID(ID string) error
}

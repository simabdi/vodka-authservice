package helper

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"strconv"
)

type PaginateResource struct {
	Data    interface{} `json:"data"`
	Page    int         `json:"page"`
	PerPage int         `json:"per_page"`
	Total   int64       `json:"total"`
}

var TotalRecord int64
var PagePaginate int
var PerPagePaginate int

func Paginate(ctx *fiber.Ctx) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		page, _ := strconv.Atoi(ctx.Query("page"))
		if page <= 0 {
			page = 1
		}

		perPage, _ := strconv.Atoi(ctx.Query("per_page"))
		switch {
		case perPage == 0:
			perPage = 10
		case perPage < 0:
			perPage = 999
		}

		PagePaginate = page
		PerPagePaginate = perPage

		offset := (page - 1) * perPage
		return db.Offset(offset).Limit(perPage)
	}
}

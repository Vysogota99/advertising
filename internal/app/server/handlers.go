package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Vysogota99/advertising/internal/app/models"
	"github.com/gin-gonic/gin"
)

// CreatAdHandler - метод создания объявления
func (r *Router) CreatAdHandler(c *gin.Context) {
	add := models.Ad{}

	if err := c.ShouldBindJSON(&add); err != nil {
		respond(c, http.StatusBadRequest, "", err)
		return
	}

	nRows, err := r.store.Add().Create(add)
	if err != nil {
		respond(c, http.StatusInternalServerError, "", err)
		return
	}

	respond(c, http.StatusOK, map[string]interface{}{
		"id":     nRows,
		"status": "success",
	}, nil)
}

// GetAdHandler - метод получения конкретного объявления
func (r *Router) GetAdHandler(c *gin.Context) {
	type request struct {
		Description bool `form:"description" binding:"omitempty"`
		Photos      bool `form:"photos" binding:"omitempty"`
	}

	req := request{}

	if err := c.ShouldBindQuery(&req); err != nil {
		respond(c, http.StatusBadRequest, "", err)
		return
	}

	id := c.Param("id")
	intID, err := strconv.Atoi(id)
	if err != nil || intID <= 0 {
		respond(c, http.StatusBadRequest, "", fmt.Errorf("Неправильно задан параметр id, ожидалось положительное число"))
		return
	}

	ad, err := r.store.Add().GetOne(intID, req.Description, req.Photos)
	if err != nil {
		respond(c, http.StatusInternalServerError, "", err)
		return
	}

	respond(c, http.StatusOK, ad, nil)
}

// GetAdsHandler - выводит список объявлений
func (r *Router) GetAdsHandler(c *gin.Context) {
	type request struct {
		Page          int    `form:"p" binding:"omitempty,min=1"`
		SortBy        string `form:"sort_by" binding:"omitempty,oneof=created_at price"`
		SortDirection string `form:"sort_direction" binding:"omitempty,oneof=asc desc"`
	}

	req := request{}
	if err := c.ShouldBindQuery(&req); err != nil {
		respond(c, http.StatusBadRequest, "", err)
		return
	}

	if req.Page == 0 {
		req.Page = 1
	}

	result, nPages, err := r.store.Add().GetList(10, req.Page, req.SortBy, req.SortDirection)
	if err != nil {
		respond(c, http.StatusInternalServerError, "", err)
		return
	}

	respond(c, http.StatusOK, map[string]interface{}{
		"n_pages": nPages,
		"data":             result,
	}, nil)
}

func respond(c *gin.Context, code int, result interface{}, err error) {
	if err != nil {
		var msg string
		if err.Error() == "EOF" && code != 500 {
			msg = "Отсутсвует тело запроса"
		} else if code == 400 {
			msg = "Некорректный запрос"
		} else if code == 500 {
			msg = "Внутрення ошибка сервера"
		}

		c.JSON(
			code,
			gin.H{
				"result": map[string]interface{}{
					"error":   err.Error(),
					"message": msg,
					"status":  "fail",
				},
			},
		)
	} else {
		c.JSON(
			code,
			gin.H{
				"result": result,
			},
		)
	}
}

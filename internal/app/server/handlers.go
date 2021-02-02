package server

import (
	"net/http"

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

	respond(c, http.StatusOK, map[string]int{
		"добавлено объявлений": nRows,
	}, nil)
}

func respond(c *gin.Context, code int, result interface{}, err error) {
	if err != nil {
		var msg string
		if err.Error() == "EOF" && code != 500 {
			msg = "Отсутсвует тело запроса"
		} else if code == 400 {
			msg = "Неправильное тело запроса"
		} else if code == 500 {
			msg = "Внутрення ошибка сервера"
		}

		c.JSON(
			code,
			gin.H{
				"error":   err.Error(),
				"message": msg,
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

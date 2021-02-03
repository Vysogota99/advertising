package server

import (
	"github.com/Vysogota99/advertising/internal/app/store"
	"github.com/gin-gonic/gin"
)

// Router ...
type Router struct {
	router     *gin.Engine
	serverPort string
	store      store.Store
}

// NewRouter - helper for initialization http router
func NewRouter(serverPort string, store store.Store) *Router {
	return &Router{
		router:     gin.Default(),
		serverPort: serverPort,
		store:      store,
	}
}

// Setup - найстройка роутера
func (r *Router) Setup() *gin.Engine {
	r.router.POST("/ad", r.CreatAdHandler)
	r.router.GET("/ad/:id", r.GetAdHandler)
	r.router.GET("/ads", r.GetAdsHandler)
	return r.router
}

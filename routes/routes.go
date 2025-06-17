package routes

import (
	"framework/handlers"

	"github.com/gin-gonic/gin"
)

func LoadRoutes(r *gin.Engine) {
	r.LoadHTMLGlob("templates/*")

	r.GET("/", handlers.LoginPage)
	r.POST("/login", handlers.LoginHandler)

	r.GET("/home", handlers.HomePage)
	r.POST("/foto/tambah", handlers.TambahFoto)
	r.POST("/foto/edit", handlers.EditFoto)
	r.GET("/foto/hapus/:id", handlers.HapusFoto)
}

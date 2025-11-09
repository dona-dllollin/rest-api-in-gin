package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	v1 := g.Group("/api/v1")
	v1.Use(app.PromDurationMiddleware())
	{

		v1.GET("/events", app.getAllEvents)
		v1.GET("/events/:id", app.getEvent)
		v1.GET("/events/:id/attendes", app.getAttendeForEvent)
		v1.GET("/attendes/:id/events", app.getEventsByAttende)
		v1.POST("/auth/register", app.registerUser)
		v1.POST("/auth/login", app.login)
		v1.POST("/auth/refresh", app.Refresh)

	}

	authGroup := v1.Group("/")
	authGroup.Use(app.AuthMiddleware())
	{
		authGroup.POST("/events", app.createEvent)
		authGroup.PUT("/events/:id", app.updateEvent)
		authGroup.DELETE("/events/:id", app.deleteEvent)
		authGroup.DELETE("/events/:id/attendes/:userId", app.deleteAttendeFromEvent)
		authGroup.POST("/events/:id/attendes/:userId", app.addAttendeToEvent)
	}

	g.GET("/swagger/*any", func(c *gin.Context) {
		if c.Request.RequestURI == "/swagger/" {
			c.Redirect(302, "/swagger/index.html")
		}
		ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("http://localhost:8080/swagger/doc.json"))(c)
	})

	// g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return g
}

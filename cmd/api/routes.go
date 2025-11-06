package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (app *application) routes() http.Handler {
	g := gin.Default()

	v1 := g.Group("/api/v1")
	{
		v1.POST("/events", app.createEvent)
		v1.GET("/events", app.getAllEvents)
		v1.GET("/events/:id", app.getEvent)
		v1.PUT("/events/:id", app.updateEvent)
		v1.DELETE("/events/:id", app.deleteEvent)
		v1.POST("/events/:id/attendes/:userId", app.addAttendeToEvent)
		v1.GET("/events/:id/attendes", app.getAttendeForEvent)
		v1.DELETE("/events/:id/attendes/:userId", app.deleteAttendeFromEvent)
		v1.GET("/attendes/:id/events", app.getEventsByAttende)

		v1.POST("/auth/register", app.registerUser)
		v1.GET("/auth/login", app.login)
		v1.GET("/auth/register/:id")
		v1.PUT("/auth/register/:id")
		v1.DELETE("/auth/register/:id")

	}

	return g
}

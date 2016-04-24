package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"html/template"
	api "./v1"
)

func main () {
	// Api Version
	apiVersion := "/api/v1"
	e := echo.New()

	// Static files
	e.Static("assets", "public/assets")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())


	// Routes 
	// Page
	e.POST(apiVersion + "/page", api.PageNew)
	e.GET(apiVersion + "/page/:id", api.PageGet)
	e.PUT(apiVersion + "/page/:id", api.PageUpdate)
	e.DELETE(apiVersion + "/page/:id", api.PageDelete)

	// Items
	e.GET(apiVersion + "/item/:id", api.ItemGet)

	// Settings
	e.PUT(apiVersion + "/setting/:key", api.SettingSet)
	e.GET(apiVersion + "/setting/:key", api.SettingGet)

	// List
	e.GET(apiVersion + "/list/page", api.PageList)
	e.GET(apiVersion + "/list/item", api.ItemList)
	e.GET(apiVersion + "/list/setting", api.SettingList)


	// Configure template engine
	e.SetRenderer(&Template{
		templates: template.Must(template.ParseGlob("public/views/*.html")),
	})

	// Dashboard
	e.GET("/", Dashboard)

	// Start server
	e.Run(standard.New(":8080"))
}
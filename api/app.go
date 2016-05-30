package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/labstack/echo/middleware"
	"html/template"
	"os"
	"path"
	"log"
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
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))


	// Routes 

	// Page
	e.POST(apiVersion + "/page", api.PageNew)
	e.GET(apiVersion + "/page/:id", api.PageGet)
	e.PUT(apiVersion + "/page/:id", api.PageUpdate)
	e.DELETE(apiVersion + "/page/:id", api.PageDelete)

	// Items
	e.GET(apiVersion + "/item/:id", api.ItemGet)
	e.PUT(apiVersion + "/item/:id", api.ItemUpdate)

	// Subscribe
	e.POST(apiVersion + "/subscribe", api.Subscribe)
	e.POST(apiVersion + "/unsubscribe", api.Unsubscribe)
	e.POST(apiVersion + "/subscribe/group", api.GroupNew)
	e.GET(apiVersion + "/subscribe/group/:id", api.GroupGet)
	e.PUT(apiVersion + "/subscribe/group/:id", api.GroupUpdate)
	e.DELETE(apiVersion + "/subscribe/group/:id", api.GroupDelete)

	// List
	e.GET(apiVersion + "/list/page", api.PageList)
	e.GET(apiVersion + "/list/item", api.ItemList)
	e.GET(apiVersion + "/list/subscribe", api.SubscribeList)
	e.GET(apiVersion + "/list/subscribe/group", api.GroupList)

	// Dashboard
	e.GET(apiVersion + "/dashboard/stats/tags", api.Tags)
	e.GET(apiVersion + "/dashboard/stats/market", api.MarketValue)

	// Misc method
	e.GET(apiVersion + "/list/category", api.CategoryList)
	e.POST(apiVersion + "/search", api.Search)


	// Get CWD
	cwd, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}

	// Configure template engine
	e.SetRenderer(&Template{
		templates: template.Must(template.ParseGlob(path.Join(cwd, "public/views/*.html"))),
	})

	// Dashboard
	e.GET("/", Dashboard)

	// Start server
	e.Run(standard.New(":8080"))
}
package main

import (
	"io"
	"html/template"
	"github.com/labstack/echo"
)

type Template struct {
	templates *template.Template
}

// Implements templat
func (this *Template) Render(writer io.Writer, name string, data interface{}, c echo.Context) error {
	return this.templates.ExecuteTemplate(writer, name, data)
}
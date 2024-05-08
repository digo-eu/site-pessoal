package main

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/echo/v4"
)

type Template struct {
	tmpl *template.Template
}

func newTemplate() *Template {
    return &Template{
        tmpl: template.Must(template.ParseGlob("views/*.html")),
    }
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.tmpl.ExecuteTemplate(w, name, data)
}

type Reason struct {
	Title string
	Explanation string
}

func newReason(title string, explanation string) Reason {
	return Reason {
		Title: title,
		Explanation: explanation,
	}
}

type Data struct {
	Reasons []Reason 
}

func (d Data) hasReason(title string) bool {
	for _, reason := range d.Reasons {
		if reason.Title == title {
			return true
		}
	}
	return false
}

func newData() Data {
	return Data {
		Reasons: []Reason {
		},
	}
}

type FormData struct {
	Values map[string]string
	Errors map[string]string
}

func newFormData() FormData {
	return FormData {
		Values: make(map[string]string),
		Errors: make(map[string]string),
	}
}

type Page struct {
	Data Data
	Form FormData
}

func newPage() Page {
	return Page {
		Data: newData(),
		Form: newFormData(),
	}
}

func main() {

	e := echo.New()
	
	page := newPage()

	e.Renderer = newTemplate() 
	e.Use(middleware.Logger())

	e.GET("/", func(c echo.Context) error {
		return c.Render(200, "index", page)
	})

	e.POST("/reasons", func(c echo.Context) error {
		title := c.FormValue("title")
		explanation := c.FormValue("explanation")
		formData := newFormData()
		reason := newReason(title, explanation)

		if page.Data.hasReason(title) {
			formData.Errors["title"] = "Esse motivo já foi mencionado, mas eu sei que é tão bom que você quis repetir."

			return c.Render(422, "create-reason", formData)
		}

		page.Data.Reasons = append(page.Data.Reasons, newReason(title, explanation))

		err := c.Render(200, "create-reason", formData)

		if err != nil {
			return err
		}

		return c.Render(200, "oob-contact", reason)
	})


	e.Logger.Fatal(e.Start(":42069"))

}

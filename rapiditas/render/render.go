package render

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type Render struct {
	Renderer   string
	RootPath   string
	Secure     bool
	Port       string
	ServerName string
}

type TemplateData struct {
	IsAuthenticated bool
	IntMap          map[string]int
	StringMap       map[string]string
	FloatMap        map[float64]float64
	Data            map[string]interface{}
	CSRFToken       string
	Port            string
	ServerName      string
	Secure          bool
}

func (r *Render) Page(writer http.ResponseWriter, request *http.Request, view string, variable, data interface{}) error {
	switch strings.ToLower(r.Renderer) {
	case "go":
		return r.GoPage(writer, request, view, data)
	case "jet":
	}
	return nil
}

func (r *Render) GoPage(writer http.ResponseWriter, request *http.Request, view string, data interface{}) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/views/%s.page.tmpl", r.RootPath, view))
	if err != nil {
		return err
	}

	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}
	err = tmpl.Execute(writer, &td)
	if err != nil {
		return err
	}
	return nil
}

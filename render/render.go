package render

import (
	"errors"
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/justinas/nosurf"
	"github.com/kataras/blocks"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"
)

type Render struct {
	Renderer    string
	RootPath    string
	Secure      bool
	Port        string
	ServerName  string
	JetViews    *jet.Set
	BlocksViews *blocks.Blocks
	Session     *scs.SessionManager
}

type TemplateData struct {
	IsAuthenticated bool
	IntMap          map[string]int
	StringMap       map[string]string
	FloatMap        map[string]float64
	Data            map[string]interface{}
	CSRFToken       string
	Port            string
	ServerName      string
	Secure          bool
	Error           string
	Flash           string
}

func (r *Render) DefaultData(templateData *TemplateData, request *http.Request) *TemplateData {
	templateData.Secure = r.Secure
	templateData.ServerName = r.ServerName
	templateData.Port = r.Port
	templateData.CSRFToken = nosurf.Token(request)
	if r.Session.Exists(request.Context(), "userID") {
		templateData.IsAuthenticated = true
	}
	templateData.Error = r.Session.PopString(request.Context(), "error")
	templateData.Flash = r.Session.PopString(request.Context(), "flash")
	return templateData
}

// Page The page render function. You can use it to render pages using go or jet templates.
func (r *Render) Page(writer http.ResponseWriter, request *http.Request, view, layout string, variables interface{}, data map[string]interface{}) error {
	switch strings.ToLower(r.Renderer) {
	case "go":
		return r.GoPage(writer, request, view, data)
	case "jet":
		return r.JetPage(writer, request, view, variables, data)
	case "blocks":
		return r.BlocksPage(writer, request, view, layout, data)

	default:
	}
	return errors.New("No rendering engine available. Please fill the required value (go or jet) in .env file ")
}

// GoPage The default go template engine renderer function.
func (r *Render) GoPage(writer http.ResponseWriter, request *http.Request, view string, data map[string]interface{}) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/views/%s.tmpl", r.RootPath, view))
	if err != nil {
		return err
	}

	td := &TemplateData{}
	if data != nil {
		td.Data = data
	}
	td = r.DefaultData(td, request)
	err = tmpl.Execute(writer, &td)
	if err != nil {
		return err
	}
	return nil
}

// JetPage The jet engine template renderer function.
func (r *Render) JetPage(writer http.ResponseWriter, request *http.Request, view string, variables interface{}, data map[string]interface{}) error {
	var vars jet.VarMap
	if variables == nil {
		vars = make(jet.VarMap)
	} else {
		vars = variables.(jet.VarMap)
	}
	td := &TemplateData{}

	if data != nil {
		td.Data = data
	}

	td = r.DefaultData(td, request)
	t, err := r.JetViews.GetTemplate(fmt.Sprintf("%s.jet", view))
	if err != nil {
		log.Println(err)
		return err
	}
	if err = t.Execute(writer, vars, td); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// BlocksPage The Blocks' engine template renderer function.
func (r *Render) BlocksPage(writer http.ResponseWriter, request *http.Request, view, layout string, data map[string]interface{}) error {
	writer.Header().Set("Content-Type", "text/html; charset=utf-8")
	r.BlocksViews = blocks.New("./views").
		Reload(true).
		Funcs(map[string]interface{}{
			"year": func() int {
				return time.Now().Year()
			},
		})
	err := r.BlocksViews.Load()
	if err != nil {
		return err
	}

	td := &TemplateData{}

	if data != nil {

		data["data"] = r.DefaultData(td, request)
	}

	err = r.BlocksViews.ExecuteTemplate(writer, view, layout, data)
	if err != nil {
		return err
	}
	return nil
}

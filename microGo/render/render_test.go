package render

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRenderPage(t *testing.T) {

	r, err := http.NewRequest("GET", "/test/render", nil)
	if err != nil {
		t.Error(err)
	}
	w := httptest.NewRecorder()

	testRenderer.Renderer = "go"
	testRenderer.RootPath = "./test"
	err = testRenderer.Page(w, r, "home", nil, nil)
	if err != nil {
		t.Error("Unable to render template. ", err)
	}

	err = testRenderer.Page(w, r, "not-exist", nil, nil)
	if err == nil {
		t.Error("Something went wrong. Unable to render a not existing template.", err)
	}
	testRenderer.Renderer = "jet"
	err = testRenderer.Page(w, r, "home", nil, nil)
	if err != nil {
		t.Error("Unable to render template. ", err)
	}
	err = testRenderer.Page(w, r, "not-exist", nil, nil)
	if err == nil {
		t.Error("Something went wrong. Unable to render a not existing  jet template.", err)
	}
	testRenderer.Renderer = ""
	err = testRenderer.Page(w, r, "home", nil, nil)
	if err == nil {
		t.Error("No error while trying to render template with no valid engine. ", err)
	}

}

func TestRenderGoPage(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/url", nil)
	if err != nil {
		t.Error(err)
	}
	testRenderer.Renderer = "go"
	testRenderer.RootPath = "./test"
	err = testRenderer.Page(w, r, "home", nil, nil)
	if err != nil {
		t.Error("Unable to render Go template. ", err)
	}
}
func TestRenderJetPage(t *testing.T) {
	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/url", nil)
	if err != nil {
		t.Error(err)
	}
	testRenderer.Renderer = "jet"
	testRenderer.RootPath = "./test"
	err = testRenderer.Page(w, r, "home", nil, nil)
	if err != nil {
		t.Error("Unable to render Jet template. ", err)
	}
}

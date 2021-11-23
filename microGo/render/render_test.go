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
		t.Error("Unable to render page. ", err)
	}

	testRenderer.Renderer = "jet"
	err = testRenderer.Page(w, r, "home", nil, nil)
	if err != nil {
		t.Error("Unable to render page. ", err)
	}
}

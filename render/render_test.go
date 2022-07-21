package render

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var templateTests = []struct {
	name          string
	renderer      string
	template      string
	errorExpected bool
	errorMSG      string
}{
	{"goTemplate", "go", "home",
		false, "Unable to render go template.",
	},
	{"goTemplateNoTemplate", "go", "not-exists",
		true, "Something went wrong. Unable to render a not existing go template.",
	},
	{"jetTemplate", "jet", "home",
		false, "Unable to render jet template.",
	},
	{"jetTemplateNoTemplate", "jet", "not-exists",
		true, "Something went wrong. Unable to render a not existing jet template.",
	},
	{"invalidRenderEngine", "foo", "home",
		true, "No error while trying to render template with no valid engine.",
	},
}

func TestRenderPage(t *testing.T) {

	for _, task := range templateTests {
		r, err := http.NewRequest("GET", "/test/render", nil)
		if err != nil {
			t.Error(err)
		}
		w := httptest.NewRecorder()

		testRenderer.Renderer = task.renderer
		testRenderer.RootPath = "./test"
		err = testRenderer.Page(w, r, task.template, nil, nil)

		if task.errorExpected {
			if err == nil {
				t.Errorf("%s: %s:", task.errorMSG, err.Error())

			}
		} else {
			if err != nil {
				t.Errorf("%s: %s: %s:", task.name, task.errorMSG, err.Error())
			}
		}

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

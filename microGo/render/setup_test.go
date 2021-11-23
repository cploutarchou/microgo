package render

import (
	"github.com/CloudyKit/jet/v6"
	"os"
	"testing"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./test/views"),
	jet.InDevelopmentMode(),
)

var testRenderer = Render{
	Renderer: "",
	RootPath: "",
	JetViews: views,
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

package render

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/kataras/blocks"
	"os"
	"testing"
	"time"
)

var views = jet.NewSet(
	jet.NewOSFileSystemLoader("./test/views"),
	jet.InDevelopmentMode(),
)
var blocksViews = blocks.New("./test/views").
	Reload(true).
	Funcs(map[string]interface{}{
		"year": func() int {
			return time.Now().Year()
		},
	})

var testRenderer = Render{
	Renderer:    "",
	RootPath:    "",
	JetViews:    views,
	BlocksViews: blocksViews,
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

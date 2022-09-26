package requests

import (
	"net/http"
	"net/url"
)

type Requests struct {
	*http.Request
	*http.Header
	*http.Client
	*http.Transport
	Params url.Values
}

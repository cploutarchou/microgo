package requests

import (
	"bytes"
	"mime"
	"net/http"
	"strings"
)

type Response struct {
	*http.Response
	Error      error
	StatusCode int
}

func (r *Response) ContentType() (string, map[string]string, error) {
	ct := r.Header.Get("content-type")
	filtered, params, err := mime.ParseMediaType(ct)
	if err != nil {
		return "", nil, err
	}
	return filtered, params, nil
}

func (r *Response) String() (string, error) {
	defer r.Body.Close()
	res := new(bytes.Buffer)
	_, err := res.ReadFrom(r.Body)
	if err != nil {
		return "", err
	}
	bodyStr := res.String()
	return bodyStr, nil
}

func (r *Response) Bytes() ([]byte, error) {
	defer r.Body.Close()
	res := new(bytes.Buffer)
	_, err := res.ReadFrom(r.Body)
	if err != nil {
		return nil, err
	}
	bodyBytes := res.Bytes()
	return bodyBytes, nil
}


func (r *Response) JSON() ([]byte, error) {
	res := []byte{}
	for _, content := range r.Header["Content-Type"] {
		t, _, err := mime.ParseMediaType(content)
		if err != nil {
			return nil, err
		}
		if strings.Contains(t, "application/json") {
			res, err = r.Bytes()
			if err != nil {
				return nil, err
			}
		}
	}
	return res, nil
}

func (r *Response) XML() ([]byte, error) {
	res := []byte{}
	for _, content := range r.Header["Content-Type"] {
		t, _, err := mime.ParseMediaType(content)
		if err != nil {
			return nil, err
		}
		if strings.Contains(t, "application/xml") {
			res, err = r.Bytes()
			if err != nil {
				return nil, err
			}
		}
	}
	return res, nil
}

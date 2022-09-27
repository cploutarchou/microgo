package requests

import (
	"bytes"
	"encoding/json"
	"io"
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

// do make http request to the provided
func (r *Requests) do(method, requestURL string, body io.Reader, options []func(*Requests)) (*Requests, error) {
	newRequest, err := http.NewRequest(method, requestURL, body)
	if err != nil {
		return nil, err
	}
	request := &Requests{
		Request:   newRequest,
		Client:    &http.Client{},
		Transport: &http.Transport{},
		Params:    url.Values{},
	}

	// Apply the options parameters to request.
	for _, option := range options {
		option(request)
	}
	parsedURL, err := url.Parse(requestURL)
	if err != nil {
		return nil, err
	}

	// If the request.Params is set, replace raw query with that.
	if len(request.Params) > 0 {
		parsedURL.RawQuery = request.Params.Encode()
	}

	newRequest.URL = parsedURL

	// Parse query values into request Form
	err = newRequest.ParseForm()
	if err != nil {
		return nil, err
	}
	return request, nil
}

// Head sends an HTTP HEAD request to the specified URL,
// with the ability to add query parameters, headers, and timeout, among other options..
func (r *Requests) Head(requestURL string, options ...func(*Requests)) (*Response, error) {
	request, err := r.do("HEAD", requestURL, nil, options)
	if err != nil {
		return nil, err
	}
	resp, err := request.Client.Do(request.Request)
	if err != nil {
		return nil, err
	}

	// Encapsulate *http.Response in *Response
	response := &Response{Response: resp}
	return response, nil
}

// Get sends a GET request to the provided url
func (r *Requests) Get(requestURL string, options ...func(*Requests)) (*Response, error) {
	request, err := r.do("GET", requestURL, nil, options)
	if err != nil {
		return nil, err
	}
	resp, err := request.Client.Do(request.Request)
	if err != nil {
		return nil, err
	}

	// Encapsulate *http.Response in *Response
	response := &Response{Response: resp}
	return response, nil
}

// AsyncGet sends an HTTP GET request to the provided URL and
// returns ch <-chan *http.Response immediately.
func (r *Requests) AsyncGet(requestURL string, options ...func(*Requests)) (<-chan *Response, error) {
	request, err := r.do("GET", requestURL, nil, options)
	if err != nil {
		return nil, err
	}
	ch := make(chan *Response)
	go func() {
		resp, err := request.Client.Do(request.Request)
		response := &Response{}
		if err != nil {
			response.Error = err
			ch <- response
		}
		response.Response = resp
		ch <- response
		close(ch)
	}()
	return ch, nil
}

// Post sends an HTTP POST request to the provided URL
func (r *Requests) Post(requestURL, bodyType string, body io.Reader, options ...func(*Requests)) (*Response, error) {
	request, err := r.do("POST", requestURL, body, options)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", bodyType)
	resp, err := request.Client.Do(request.Request)
	if err != nil {
		return nil, err
	}

	// Encapsulate *http.Response in *Response
	response := &Response{Response: resp}
	return response, nil
}

// AsyncPost sends an HTTP POST request to the provided URL and returns a <-chan *http.Response immediately.
func (r *Requests) AsyncPost(requestURL, bodyType string, body io.Reader, options ...func(*Requests)) (<-chan *Response, error) {
	request, err := r.do("POST", requestURL, body, options)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", bodyType)
	resChannel := make(chan *Response)
	go func() {
		resp, err := request.Client.Do(request.Request)
		// Encapsulate *http.Response in *Response
		response := &Response{}
		if err != nil {
			response.Error = err
			resChannel <- response
		}
		response.Response = resp
		resChannel <- response
		close(resChannel)
	}()
	return resChannel, nil
}

// JSONPost  Marshals request data as JSON and set the content type to "application/json".
func (r *Requests) JSONPost(requestURL string, body interface{}, options ...func(*Requests)) (*Response, error) {
	buff := new(bytes.Buffer)
	err := json.NewEncoder(buff).Encode(body)
	if err != nil {
		return nil, err
	}
	request, err := r.do("POST", requestURL, buff, options)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := request.Client.Do(request.Request)
	if err != nil {
		return nil, err
	}

	// Encapsulate *http.Response in *Response
	response := &Response{Response: resp}
	return response, nil
}

// XMLPost  Marshals request data as XML and set the content type to "application/xml".
func (r *Requests) XMLPost(requestURL string, body interface{}, options ...func(*Requests)) (*Response, error) {
	buff := new(bytes.Buffer)
	err := json.NewEncoder(buff).Encode(body)
	if err != nil {
		return nil, err
	}
	request, err := r.do("POST", requestURL, buff, options)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/xml")
	resp, err := request.Client.Do(request.Request)
	if err != nil {
		return nil, err
	}

	// Encapsulate *http.Response in *Response
	response := &Response{Response: resp}
	return response, nil
}

// Put sends HTTP PUT request to the provided URL.
func (r *Requests) Put(requestURL, bodyType string, body io.Reader, options ...func(*Requests)) (*Response, error) {
	request, err := r.do("PUT", requestURL, body, options)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", bodyType)
	resp, err := request.Client.Do(request.Request)
	if err != nil {
		return nil, err
	}

	// Encapsulate *http.Response in *Response
	response := &Response{Response: resp}
	return response, nil
}

// Patch sends an HTTP PATCH request to the provided URL with optional body to update the data.
func (r *Requests) Patch(requestURL, bodyType string, body io.Reader, options ...func(*Requests)) (*Response, error) {
	request, err := r.do("PATCH", requestURL, body, options)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", bodyType)
	resp, err := request.Client.Do(request.Request)
	if err != nil {
		return nil, err
	}

	// Encapsulate *http.Response in *Response
	response := &Response{Response: resp}
	return response, nil
}

// Delete sends an HTTP DELETE request to the provided URL.
func (r *Requests) Delete(requestURL string, options ...func(*Requests)) (*Response, error) {
	request, err := r.do("DELETE", requestURL, nil, options)
	if err != nil {
		return nil, err
	}
	resp, err := request.Client.Do(request.Request)
	if err != nil {
		return nil, err
	}

	// Encapsulate *http.Response in *Response
	response := &Response{Response: resp}
	return response, nil
}

// Options sends a rarely-used HTTP OPTIONS request to the provided URL.
// Options only permit a single parameter, which is the destination URL string.
func (r *Requests) Options(requestURL string) (*Response, error) {
	request, err := r.do("OPTIONS", requestURL, nil, []func(r *Requests){})
	if err != nil {
		return nil, err
	}
	resp, err := request.Client.Do(request.Request)
	if err != nil {
		return nil, err
	}

	// Encapsulate *http.Response in *Response
	response := &Response{Response: resp}
	return response, nil
}

package microGo

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"path"
	"path/filepath"
)

//WriteJson : Create a JSON response.
func (m *MicroGo) WriteJson(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for key, val := range headers[0] {
			w.Header()[key] = val
		}
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

//WriteXML : Create XML response.
func (m *MicroGo) WriteXML(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := xml.MarshalIndent(data, "", "")
	if err != nil {
		return err
	}
	if len(headers) > 0 {
		for key, val := range headers[0] {
			w.Header()[key] = val
		}
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(status)
	_, err = w.Write(out)
	if err != nil {
		return err
	}
	return nil
}

//SentFile : Send a file on response.
func (m *MicroGo) SentFile(w http.ResponseWriter, r *http.Request, fileLocation, fileName string) error {
	_path := path.Join(fileLocation, fileName)
	fileToServe := filepath.Clean(_path)
	w.Header().Set("Content-Type", fmt.Sprintf("attachment; file=\"%s\"", fileName))
	http.ServeFile(w, r, fileToServe)
	return nil
}

//Error404 : Return Not Found HTTP response .
//
//Status Code : 404
func (m *MicroGo) Error404(w http.ResponseWriter) {
	m.ErrorStatus(w, http.StatusNotFound)
}

//Error500 : Return StatusInternal Server Error HTTP response .
//
//Status Code : 500
func (m *MicroGo) Error500(w http.ResponseWriter) {
	m.ErrorStatus(w, http.StatusInternalServerError)
}

//ErrorUnauthorized : Return Unauthorized response on request error.
//
//Status Code : 401
func (m *MicroGo) ErrorUnauthorized(w http.ResponseWriter) {
	m.ErrorStatus(w, http.StatusUnauthorized)
}

//ErrorForbidden : Return StatusForbidden HTTP response.
//
//Status Code : 403
func (m *MicroGo) ErrorForbidden(w http.ResponseWriter) {
	m.ErrorStatus(w, http.StatusForbidden)
}

//ErrorUnprocessable : Return Unprocessable entity HTTP response.
//
//Status Code 422.
func (m *MicroGo) ErrorUnprocessable(w http.ResponseWriter) {
	m.ErrorStatus(w, http.StatusUnprocessableEntity)
}

//ErrorStatus : Construct Error HTTP response
func (m *MicroGo) ErrorStatus(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)

}

package httpx

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"

	"apiserver/pkg/util/logger"
	"github.com/gorilla/mux"
)

const (
	HTTP_GET    = "GET"
	HTTP_POST   = "POST"
	HTTP_PUT    = "PUT"
	HTTP_DELETE = "DELETE"
)

const (
	OK               = "OK"
	JSON_EMPTY_ARRAY = "[]"
	JSON_EMPTY_OBJ   = "{}"
)

// https://golang.org/src/net/http/status.go
var (
	// 200
	StatusOK = strconv.Itoa(http.StatusOK)

	// 201
	StatusCreated = strconv.Itoa(http.StatusCreated)

	// 204
	StatusNoContent = strconv.Itoa(http.StatusNoContent)

	// 400
	StatusBadRequest = strconv.Itoa(http.StatusBadRequest)

	// 402
	StatusPaymentRequired = strconv.Itoa(http.StatusPaymentRequired)

	// 403
	StatusForbidden = strconv.Itoa(http.StatusForbidden)

	// 404
	StatusNotFound = strconv.Itoa(http.StatusNotFound)

	// 409
	StatusConflict = strconv.Itoa(http.StatusConflict)

	// 500
	StatusInternalServerError = strconv.Itoa(http.StatusInternalServerError)
)

var log = logger.New("")

// -----------------------
// RegisterHttpHandler

type HttpHandler func(*http.Request) (string, interface{})

func RegisterHttpHandler(router *mux.Router, path, method string, handler HttpHandler) {
	h := func(w http.ResponseWriter, r *http.Request) {
		// parseForm
		if err := r.ParseForm(); err != nil {
			log.Warning(err)
		}

		// dump
		bytes, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Warning(err)
		} else {
			log.Debug(string(bytes))
		}

		dump := dumpHttpRequest(r)
		log.Debug(dump)

		t := time.Now()
		status, body := handler(r)
		writeHttpResp(w, dump, status, body, t)
	}
	router.HandleFunc(path, h).Methods(method)
}

func dumpHttpRequest(r *http.Request) string {
	if r.Method == "GET" {
		return fmt.Sprintf("%s %s", r.Method, r.URL.RequestURI())
	}

	if r.Method == "POST" {
		return fmt.Sprintf("%s %s", r.Method, r.URL.RequestURI())
	}

	return fmt.Sprintf("%s %s %s", r.Method, r.URL.RequestURI(), r.Form)
}

// --------------------------------
// response

const httpJsonRespFmt = `{
  "kind":"%s"
  "apiVersion": "1.0",
  "metadata": "%v",
  "status": "%s",
  "message": "%s",
  "reason": "%s",
  "details": "%v",
  "code": "%s",
  "items": %s
}
`

func writeHttpResp(w http.ResponseWriter, dump string, status string, body interface{}, t time.Time) {
	w.Header().Set("Content-Type", "application/json")

	sub := time.Now().Sub(t)
	// empty array
	if body == JSON_EMPTY_ARRAY {
		log.Info(dump, status, sub)
		fmt.Fprintf(w, httpJsonRespFmt, status, `""`, body)
		return
	}

	if body == JSON_EMPTY_OBJ {
		log.Info(dump, status, sub)
		fmt.Fprintf(w, httpJsonRespFmt, status, `""`, body)
		return
	}

	errStr, data := "", JSON_EMPTY_OBJ
	res, err := json.MarshalIndent(body, " ", "    ")
	if err != nil {
		errStr = `"` + err.Error() + `"`
		log.Debug(dump, status, errStr, data, sub)
		fmt.Fprintf(w, httpJsonRespFmt, status, errStr, data)
		return
	}

	// error
	if status != StatusOK && status != StatusCreated && status != StatusNoContent {
		errStr = `"` + string(res) + `"`
		log.Debug(dump, status, errStr, data, sub)
		fmt.Fprintf(w, httpJsonRespFmt, status, errStr, data)
		return
	}

	errStr = `"` + OK + `"`
	data = string(res)

	log.Debug(dump, status, sub)
	fmt.Fprintf(w, httpJsonRespFmt, status, errStr, data)
}

// -----------------------------------------------------------------------------
// HttpListen

func HttpServe(addr string) error {
	// root
	root := mux.NewRouter()

	// api
	api := root.PathPrefix("/api/v1").Subrouter()
	regHttpApiHandles(api)

	// log, admin,
	// ...

	// root
	http.Handle("/", root)

	// listen
	return http.ListenAndServe(addr, nil)
}

func regHttpApiHandles(router *mux.Router) {
	/*RegisterNodeHandler(router)
	RegisterAppHandler(router)
	RegisterContainerHandler(router)
	RegisterDeploymentHandler(router)*/
}

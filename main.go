package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/pprof"
	"strings"
	"time"
)

//Define a map to implement routing table.
var mux map[string]func(http.ResponseWriter, *http.Request)

func main() {
	server := http.Server{
		Addr:        ":8080",
		Handler:     &myHandler{},
		ReadTimeout: 5 * time.Second,
	}

	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/"] = IndexHandler
	mux["/static"] = StaticHandler
	mux["/version"] = VersionHandler
	mux["/headers"] = HeadersHandler
	mux["/environ"] = EnvironHandler
	mux["/memory"] = MemoryHandler
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func writeJSON(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	return e.Encode(v)
}

type ResponseError struct {
	Code        int    `json:"code"`
	Message     string `json:"message"`
	HttpCode    int    `json:"httpcode"`
	HttpMessage string `json:"httpmessage"`
}

type errorResponse struct {
	Error errObj `json:"error"`
}

type errObj struct {
	Message string `json:"message"`
}

type headersResponse struct {
	Headers map[string]string `json:"headers"`
}

type environResponse struct {
	Environment map[string]string `json:"environment"`
}

func writeErrorJSON(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.WriteHeader(http.StatusInternalServerError)
	_ = writeJSON(w, errorResponse{errObj{err.Error()}}) // ignore error, can't do anything
}

func getHeaders(r *http.Request) map[string]string {
	hdr := make(map[string]string, len(r.Header))
	for k, v := range r.Header {
		hdr[k] = v[0]
	}
	return hdr
}

func getCookies(cs []*http.Cookie) map[string]string {
	m := make(map[string]string, len(cs))
	for _, v := range cs {
		m[v.Name] = v.Value
	}
	return m
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Implement route forwarding
	if h, ok := mux[r.URL.String()]; ok {
		//Implement route forwarding with this handler, the corresponding route calls the corresponding func.
		h(w, r)
		return
	}
	io.WriteString(w, "URL: "+r.URL.String())
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "may the force be with you!")
}

func StaticHandler(w http.ResponseWriter, r *http.Request) {
	str, ok := os.LookupEnv("TIGER_STATIC")
	if !ok {
		str = "may the force be with you!"
	}
	io.WriteString(w, str)
}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "tiger v1.0.0")
}

func HeadersHandler(w http.ResponseWriter, r *http.Request) {
	if err := writeJSON(w, headersResponse{getHeaders(r)}); err != nil {
		writeErrorJSON(w, err)
	}
}

func EnvironHandler(w http.ResponseWriter, r *http.Request) {
	getenvironment := func(data []string, getkeyval func(item string) (key, val string)) map[string]string {
		items := make(map[string]string)
		for _, item := range data {
			key, val := getkeyval(item)
			items[key] = val
		}
		return items
	}
	environment := getenvironment(os.Environ(), func(item string) (key, val string) {
		splits := strings.Split(item, "=")
		key = splits[0]
		val = splits[1]
		return
	})

	if err := writeJSON(w, environResponse{environment}); err != nil {
		writeErrorJSON(w, err)
	}
}

func MemoryHandler(w http.ResponseWriter, r *http.Request) {
	// Gather memory allocations profile.
	profile := pprof.Lookup("allocs")

	// Write profile (human readable, via debug: 1) to HTTP response.
	err := profile.WriteTo(w, 1)
	if err != nil {
		log.Printf("Error: Failed to write allocs profile: %v", err)
	}
}

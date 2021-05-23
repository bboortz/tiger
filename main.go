package main

import (
	"io"
	"log"
	"net/http"
	"time"
	"runtime/pprof"
)

//Define a map to implement routing table.
var mux map[string]func(http.ResponseWriter , *http.Request) 

func main(){
	server := http.Server{
		Addr: ":8080",
		Handler: &myHandler{},
		ReadTimeout: 5*time.Second,
	}
	
	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/"] = Index
	mux["/version"] = Version
	mux["/memory"] = Memory
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}	
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request){
	// Implement route forwarding
	if h, ok := mux[r.URL.String()];ok{
	//Implement route forwarding with this handler, the corresponding route calls the corresponding func.
		h(w, r)
		return
	}
	io.WriteString(w, "URL: "+r.URL.String())
}

func Index(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "may the force be with you!")
}

func Version(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "tiger v1.0.0")
}

func Memory(w http.ResponseWriter, r *http.Request) {
        // Gather memory allocations profile.
        profile := pprof.Lookup("allocs")

        // Write profile (human readable, via debug: 1) to HTTP response.
        err := profile.WriteTo(w, 1)
        if err != nil {
                log.Printf("Error: Failed to write allocs profile: %v", err)
        }
}

package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	storageServiceURL, _ := url.Parse(os.Getenv("STORAGE_SERVICE_URL"))
	metadataServiceURL, _ := url.Parse(os.Getenv("METADATA_SERVICE_URL"))

	router.HandleFunc("/storage/{rest:.*}", proxyHandler(storageServiceURL)).
		Methods("GET", "POST", "PUT", "DELETE")
	router.HandleFunc("/metadata/{rest:.*}", proxyHandler(metadataServiceURL)).
		Methods("GET", "POST", "PUT", "DELETE")

	log.Println("API Gateway running on http://localhost:3001")
	log.Fatal(http.ListenAndServe(":3001", router))
}

func proxyHandler(target *url.URL) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		requestID := uuid.New().String()
		r.Header.Add("X-Request-ID", requestID)

		r.Host = target.Host
		r.URL.Scheme = target.Scheme
		r.URL.Host = target.Host
		r.URL.Path = mux.Vars(r)["rest"]
		proxy := httputil.NewSingleHostReverseProxy(target)
		proxy.ServeHTTP(w, r)
	}
}

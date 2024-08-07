package main

import (
	"fmt"
	"net/http"
	"strconv"
)


type apiConfig struct {
	fileserverHits int
}

func (cfg *apiConfig) middlewareMetrics(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {		
	cfg.fileserverHits++
	next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request){
		header := w.Header()
		header["Content-Type"] = []string{"text/plain; charset=utf-8"}
		w.WriteHeader(200)
		val := strconv.Itoa(cfg.fileserverHits)
		fmt.Println(val)
		w.Write([]byte(val))

}


func main(){
	apiconfig := &apiConfig{}
	serverMux := http.NewServeMux()
	server := http.Server{Handler: serverMux}
	server.Addr = "localhost:8080"
	fileserver := http.FileServer(http.Dir("./app"))
	handler := http.StripPrefix("/app",fileserver)
	serverMux.Handle("/app/*", apiconfig.middlewareMetrics(handler))
	serverMux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request){
		header := w.Header()
		header["Content-Type"] = []string{"text/plain; charset=utf-8"}
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	serverMux.HandleFunc("/metrics",apiconfig.metricsHandler)

	err := server.ListenAndServe()
	if err!=nil{
		fmt.Println(err)
	}
}
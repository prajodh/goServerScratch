package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	// "errors"
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
		header["Content-Type"] = []string{"text/html; charset=utf-8"}
		w.WriteHeader(200)
		htmlReturnString := fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", cfg.fileserverHits)
		w.Write([]byte(htmlReturnString))

}

func (cfg *apiConfig) resetHandler(w http.ResponseWriter, R *http.Request){
	header := w.Header()
	header["Content-Type"] = []string{"text/plain; charset=utf-8"}
	w.WriteHeader(200)
	cfg.fileserverHits = 0
	val := strconv.Itoa(cfg.fileserverHits)
	fmt.Println(val)
	w.Write([]byte(val))
}

func writeResponseValidateChrip(v bool, errs string) []byte {
		type returnBodyValidateChrip struct{
			Valid bool
			Errors string
		}
		returnbody := returnBodyValidateChrip{Valid: v, Errors: errs}
		res, err := json.Marshal(returnbody)
		if err != nil{
			fmt.Println(err)
		}
		return res
		
}
func validateChirpHandler(w http.ResponseWriter, r *http.Request){
	header := w.Header()
	type body struct{
		Body string
	}
	decoder := json.NewDecoder(r.Body)
	b:=body{}
	err := decoder.Decode(&b)
	if err != nil{
		res := writeResponseValidateChrip(false, err.Error())
		header["Content-Type"] = []string{"application/json; charset=utf-8"}
		w.WriteHeader(500)
		w.Write(res)
	}
	if len(b.Body) > 140{
		res := writeResponseValidateChrip(false, "chrip is too long")	
		header["Content-Type"] = []string{"application/json; charset=utf-8"}
		w.WriteHeader(400)
		w.Write(res)
		} else{
		res := writeResponseValidateChrip(true, "")		
		header["Content-Type"] = []string{"application/json; charset=utf-8"}
		w.WriteHeader(200)
		w.Write(res)
	}
}


func main(){
	apiconfig := &apiConfig{}
	serverMux := http.NewServeMux()
	server := http.Server{Handler: serverMux}
	server.Addr = "localhost:8080"
	fileserver := http.FileServer(http.Dir("./app"))
	handler := http.StripPrefix("/app",fileserver)
	serverMux.Handle("/app/*", apiconfig.middlewareMetrics(handler))
	serverMux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request){
		header := w.Header()
		header["Content-Type"] = []string{"text/plain; charset=utf-8"}
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	serverMux.HandleFunc("GET /admin/metrics",apiconfig.metricsHandler)
	serverMux.HandleFunc("GET /admin/reset",apiconfig.resetHandler)
	serverMux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)
	err := server.ListenAndServe()
	if err!=nil{
		fmt.Println(err)
	}
}
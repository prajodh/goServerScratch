package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"golang.org/x/crypto/bcrypt"
	"github.com/prajodh/goServerScratch/database"
	// "errors"
)

const databaseUrl string = "./database.json"
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

func writeResponseValidateChrip(v bool, errs string, msg string) []byte {
		type returnBodyValidateChrip struct{
			Valid bool
			Errors string
			CleanedBody string
		}
		returnbody := returnBodyValidateChrip{Valid: v, Errors: errs, CleanedBody: msg}
		res, err := json.Marshal(returnbody)
		if err != nil{
			fmt.Println(err)
		}
		return res
		
}

func replaceProfanity(str string) string{
	dict := map[string]string{"kerfuffle":"****", "sharbert":"****", "fornax":"****"}
	str = strings.ToLower(str)
	arrayStr := strings.Split(str," ")
	for i, x := range(arrayStr){
		val, ok := dict[x]
		if ok{
			arrayStr[i] = val
		}
	}
	return strings.Join(arrayStr, " ")
}


func createChirpHandler(w http.ResponseWriter, r *http.Request){
	DB, err := database.NewDB(databaseUrl)
	if err != nil{
		log.Fatal(err)
	}
	header := w.Header()
	type body struct{
		Body string
	}
	decoder := json.NewDecoder(r.Body)
	b:=body{}
	err = decoder.Decode(&b)
	b.Body = replaceProfanity(b.Body)
	if err != nil{
		res := writeResponseValidateChrip(false, err.Error(), "")
		header["Content-Type"] = []string{"application/json; charset=utf-8"}
		w.WriteHeader(500)
		w.Write(res)
	}
	var res string = b.Body
	if len(b.Body) > 140{
		res := writeResponseValidateChrip(false, "chrip is too long", "")	
		header["Content-Type"] = []string{"application/json; charset=utf-8"}
		w.WriteHeader(400)
		w.Write(res)
		} else{
		storedDBJson, _ := DB.CreateChirp(res)
		res := writeResponseValidateChrip(true, "", string(storedDBJson))		
		header["Content-Type"] = []string{"application/json; charset=utf-8"}
		w.WriteHeader(201)
		w.Write(res)
	}
}

func getChripHandler(w http.ResponseWriter, r *http.Request){
	 DB, err := database.NewDB(databaseUrl)
	 if err!= nil {
		log.Fatal(err)
	 }
	 chripsList , err := DB.GetChirps()
	 if err != nil{
		log.Fatal(err)
	 }
	header := w.Header()
	header["Content-Type"] = []string{"text/plain; charset=utf-8"}
	w.WriteHeader(200)
	mergedData := ""
	for _,v := range(chripsList){
		mergedData += string(v)+" ,"
	}
	w.Write([]byte(mergedData))
}

func getChripsbyIDHandler(w http.ResponseWriter, r *http.Request){
	str := r.PathValue("id")
	id, err := strconv.Atoi(str)
	if err != nil{
		log.Fatal(err)
	}
	DB, err := database.NewDB(databaseUrl)
	if err != nil{
		log.Fatal(err)
	}
	data, err := DB.GetChirps()
	if err != nil{
		log.Fatal(err)
	}
	header := w.Header()
	header["Content-Type"] = []string{"text/plain; charset=utf-8"}
	w.WriteHeader(200)
	w.Write([]byte(data[id]))
	
}

func createUsersHandler(w http.ResponseWriter, r *http.Request){
	header := w.Header()
	header["Content-Type"] = []string{"application/json; charset=utf-8"}
	db, err := database.NewDB(databaseUrl)
	if err != nil{
		log.Fatal(err)
	}
	type user struct{
		Email  string
		Password string
	}
	decoder := json.NewDecoder(r.Body)
	e := user{}
	err = decoder.Decode(&e)
	if err != nil{
		log.Fatal(err)
	}
	encrptedPassword, err := bcrypt.GenerateFromPassword([]byte(e.Password), 10)
	if err != nil{
		log.Fatal(err)
	}
	e.Password = string(encrptedPassword)
	ema, err := db.CreateUsers(e.Email, e.Password)
	if err != nil{
		log.Fatal(err)
	}
	w.WriteHeader(201)
	w.Write([]byte(ema))
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
	serverMux.HandleFunc("POST /api/chirps", createChirpHandler)
	serverMux.HandleFunc("GET /api/chrips", getChripHandler)
	serverMux.HandleFunc("GET /api/chrips/{id}", getChripsbyIDHandler)
	serverMux.HandleFunc("POST /api/users", createUsersHandler)
	err := server.ListenAndServe()
	if err!=nil{
		fmt.Println(err)
	}
}
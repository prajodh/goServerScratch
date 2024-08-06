package main

import ("net/http"
"fmt")


func main(){
	serverMux := http.NewServeMux()
	server := http.Server{Handler: serverMux}
	server.Addr = "localhost:8080"
	serverMux.Handle("/", http.FileServer(http.Dir(".")))
	serverMux.Handle("/assets", http.FileServer(http.Dir("./assets")))
	err := server.ListenAndServe()
	if err!=nil{
		fmt.Println(err)
	}
}
package main

import (
	"crud/servidor"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main(){

	//Gerenciador de roas com o mux
	router := mux.NewRouter()

	router.HandleFunc("/usuarios", servidor.CreateUser).Methods("POST")
	router.HandleFunc("/usuarios", servidor.GetAllUsers).Methods("GET")
	router.HandleFunc("/usuarios/{id}", servidor.GetUser).Methods("GET")
	router.HandleFunc("/usuario/update/{id}", servidor.UpdateUser).Methods("PUT")
	router.HandleFunc("/usuario/del/{id}", servidor.DelUser).Methods("DELETE")

	log.Println("Servidor rodando em http://localhost:5000")
	log.Fatal(http.ListenAndServe(":5000", router))
}
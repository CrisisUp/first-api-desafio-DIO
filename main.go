package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type person struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *address `json:"address,omitempty"`
}

type address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

var people []person

func GetPeople(w http.ResponseWriter, r *http.Request) {
	// Codifica os dados para JSON
	if err := json.NewEncoder(w).Encode(people); err != nil {
		http.Error(w, "Erro ao codificar JSON", http.StatusInternalServerError)
		return
	}
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	http.Error(w, "Contato não encontrado", http.StatusNotFound)
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if person.Firstname == "" || person.Lastname == "" {
		http.Error(w, "O nome e o sobrenome são obrigatórios.", http.StatusBadRequest)
		return
	}
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			// Remove o elemento da lista
			people = append(people[:index], people[index+1:]...)
			// Retorna um status 204 (No Content) indicando sucesso
			w.WriteHeader(http.StatusNoContent)
			return
		}
		// Se o loop terminar sem encontrar o elemento, o contato não existe
		http.Error(w, "Contato não encontrado", http.StatusNotFound)
	}
}

func main() {
	router := mux.NewRouter()
	people = append(people, person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &address{City: "City X", State: "State X"}})
	people = append(people, person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &address{City: "City Z", State: "State Z"}})
	router.HandleFunc("/contato", GetPeople).Methods("GET")
	router.HandleFunc("/contato/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/contato/{id}", CreatePerson).Methods("POST")
	router.HandleFunc("/contato/{id}", DeletePerson).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}

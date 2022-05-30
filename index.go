package main 
import (
	"fmt"
	"encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
)


type Person struct {
	ID        string  `json:"id,omitempty"`
	Firstname string  `json:"firstname,omitempty"`
	Lastname  string  `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}

type Address struct {
	City   string `json:"city,omitempty"`
	State  string `json:"state,omitempty"`
}

var people []Person

func getAllContacts(w http.ResponseWriter, r *http.Request) {
json.NewEncoder(w).Encode(people)
}

func getContact(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}

func createContact(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

func updateContact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Contact Updated"}`))
}

func deleteContact(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
	}

	json.NewEncoder(w).Encode(people)

func main() {

	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	people = append(people, Person{ID: "3", Firstname: "Francis", Lastname: "Sunday"})


	router := mux.NewRouter()
	router.HandleFunc("/contato", getAllContacts).Methods("GET")
	router.HandleFunc("/contato/{id}", getContact).Methods("GET")
	router.HandleFunc("/contato", createContact).Methods("POST")
	router.HandleFunc("/contato/{id}", updateContact).Methods("PUT")
	router.HandleFunc("/contato/{id}", deleteContact).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))
}
package main 
import (
	"encoding/json"
    "log"
    "net/http"
    "github.com/gorilla/mux"
				_ "github.com/go-sql-driver/mysql"
		"database/sql"
)


type Person struct {
	ID        string  `json:"id,omitempty"`
	Email     string  `json:"email,omitempty"`
	Password  string  `json:"password,omitempty"`
	Firstname string  `json:"firstname,omitempty"`
	Lastname  string  `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
	Contact   string  `json:"contact,omitempty"`
	AddUniqueContact  string  `json:"addUniqueContact"`
	Contacts []Person `json:"contacts,omitempty"`
}

type Address struct {
	City   string `json:"city,omitempty"`
	State  string `json:"state,omitempty"`
}


func Error (w http.ResponseWriter, error string, code int) {
	w.WriteHeader(code)
	w.Write([]byte(error))
}



var people []Person


func innerJoin (w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	db := start_database_connection()
	rows, err := db.Query("SELECT contact.id, contact.firstname, contact.lastname, contact.contact, contact.address FROM contact INNER JOIN contact_group ON contact_group.contact_id = contact.id WHERE contact_group.group_id = '" + params["id"] + "'")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var contacts []Person
	for rows.Next() {
		var id string
		var firstname string
		var lastname string
		var contact string
		var address string
		err = rows.Scan(&id, &firstname, &lastname, &contact, &address)
		if err != nil {
			panic(err.Error())
		}
		contacts = append(contacts, Person{ID: id, Firstname: firstname, Lastname: lastname, Contact: contact, Address: &Address{City: address}})
	}
	json.NewEncoder(w).Encode(contacts)
}


func getAllContacts(w http.ResponseWriter, r *http.Request) {
	db := start_database_connection()
	rows, err := db.Query("SELECT * FROM contact ")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	var contacts []Person
	for rows.Next() {
		var id string
		var firstname string
		var lastname string
		var contact string
		var address string
		err = rows.Scan(&id, &firstname, &lastname, &contact, &address)
		if err != nil {
			panic(err.Error())
		}
				
		contacts = append(contacts, Person{ID: id, Firstname: firstname, Lastname: lastname, Contact: contact, Address: &Address{City: address}})
		
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contacts)

}
func existIdInTable(db *sql.DB, id string) bool {
	rows, err := db.Query("SELECT * FROM contact WHERE id = '" + id + "'")
	if err != nil {
		panic(err.Error())
	}

	var contacts []Person
	for rows.Next() {
		var id string
		var firstname string
		var lastname string
		var contact string
		var address string
		err = rows.Scan(&id, &firstname, &lastname, &contact, &address)
		if err != nil {
			panic(err.Error())
		}

		contacts = append(contacts, Person{ID: id, Firstname: firstname, Lastname: lastname, Contact: contact, Address: &Address{City: address}})
	}
	if len(contacts) > 0 {
		return true
	}
	return false
}

func getContact(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	db := start_database_connection()
	rows, err := db.Query("SELECT * FROM contact WHERE id = '" + params["id"] + "'")
	
	if err != nil {

		Error(w, "Internal server Error", 500)
		return
	}
	defer db.Close()
	var contacts []Person
	for rows.Next() {
		var id string
		var firstname string
		var lastname string
		var contact string
		var address string
		err = rows.Scan(&id, &firstname, &lastname, &contact, &address)
		if err != nil {
			panic(err.Error())
		}


		contacts = append(contacts, Person{ID: id, Firstname: firstname, Lastname: lastname, Contact: contact, Address: &Address{City: address} })
	}
	json.NewEncoder(w).Encode(contacts)
}

func createContact(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	db := start_database_connection()	
	sql := "INSERT INTO contact (firstname, lastname, address) VALUES ('" + person.Firstname + "', '" + person.Lastname + "', '" + person.Address.City + "')"
	exec(db, sql)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}


func addUniqueContact(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	db := start_database_connection()
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	isExistParams := existIdInTable(db, params["id"])
	isExistStruct := existIdInTable(db, person.AddUniqueContact)
				defer db.Close()
				if !isExistParams || !isExistStruct {
					Error(w, "Contact not found", 404)
					return
				}
				createRelation:= "INSERT INTO contact_group (group_id, contact_id) VALUES ('" + params["id"] + "', '" + person.AddUniqueContact + "')"
				exec(db, createRelation)
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
}

func start_database_connection () *sql.DB   {
	db, err := sql.Open("mysql", "root:123456789@/contacts")
	if err != nil {
		
		panic(err.Error())
	}
	return db
}



func exec(db *sql.DB, sql string) sql.Result {
	result, err := db.Exec(sql)
	if err != nil {

		log.Fatal(err)
	}
	return result
}

func main() {
	db := start_database_connection()
	exec(db, "CREATE TABLE IF NOT EXISTS contact (id INT AUTO_INCREMENT PRIMARY KEY, firstname VARCHAR(255) ,lastname VARCHAR(255),contact VARCHAR(255) ,address VARCHAR(255))")
	exec(db, "CREATE TABLE IF NOT EXISTS contact_group (contact_id INT, group_id INT)")
	//seeds 
	    // exec(db, "INSERT INTO contact (firstname, lastname, contact, address) VALUES ('John', 'Doe', '99999999-99', 'City X')")
	    // exec(db, "INSERT INTO contact (firstname, lastname, contact, address) VALUES ('Koko', 'Doe', '99999999-99', 'City Z')")
	    // exec(db, "INSERT INTO contact (firstname, lastname, contact, address) VALUES ('Francis', 'Sunday', '99999999-99', 'City Z')")

	    // exec(db, "INSERT INTO contact_group (contact_id, group_id) VALUES (1, 1)")
	    // exec(db, "INSERT INTO contact_group (contact_id, group_id) VALUES (1, 2)")
	    // exec(db, "INSERT INTO contact_group (contact_id, group_id) VALUES (2, 2)")
	    // exec(db, "INSERT INTO contact_group (contact_id, group_id) VALUES (3, 2)")
	defer db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/contato", getAllContacts).Methods("GET")
	router.HandleFunc("/contato/{id}", getContact).Methods("GET")
	router.HandleFunc("/contato", createContact).Methods("POST")
	router.HandleFunc("/addUniqueContact/{id}", addUniqueContact).Methods("POST")
	router.HandleFunc("/contato/{id}", updateContact).Methods("PUT")
	router.HandleFunc("/contato/{id}", deleteContact).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))

}
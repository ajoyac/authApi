package main

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

// App represents the application
type App struct {
	Router *mux.Router
}

// The person Type (more like an object)
type User struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	jwt.StandardClaims
}
type Person struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname,omitempty"`
	Lastname  string   `json:"lastname,omitempty"`
	Address   *Address `json:"address,omitempty"`
}
type Address struct {
	City  string `json:"city,omitempty"`
	State string `json:"state,omitempty"`
}

var people []Person
var users []User

// login api
func login(w http.ResponseWriter, r *http.Request) {
	var user User
	response := make(map[string]string)
	if r.Body == nil {
		json.NewEncoder(w).Encode(map[string]string{"error": "not body found"})
		return
	}
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	for _, item := range users {
		if item.Username == user.Username && item.Password == user.Password {
			token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), item)
			secret := os.Getenv("SECRET")
			tokenstring, err := token.SignedString([]byte(secret))
			if err != nil {
				log.Fatalln(err)
			}
			response["token"] = tokenstring
			json.NewEncoder(w).Encode(response)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}

// Display all from the people var
func GetPeople(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}

// Display a single data
func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, item := range people {
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
	json.NewEncoder(w).Encode(&Person{})
}

// create a new item
func CreatePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = params["id"]
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

// Delete an item
func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(people)
	}
}

func (a *App) Initialize() {
	a.Router = mux.NewRouter()

	//Start DB: fill some values.
	users = append(users, User{ID: "1", Username: "user1", Password: "password"})
	people = append(people, Person{ID: "1", Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: "2", Firstname: "Koko", Lastname: "Doe", Address: &Address{City: "City Z", State: "State Y"}})
	a.initializeRoutes()
}

// Run starts the app and serves on the specified addr
func (a *App) Run() {
	ip := os.Getenv("IP")
	port := os.Getenv("PORT")
	log.Fatal(http.ListenAndServe(ip+":"+port, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/login", login).Methods("POST")
	a.Router.HandleFunc("/people", GetPeople).Methods("GET")
	a.Router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	a.Router.HandleFunc("/people/{id}", CreatePerson).Methods("POST")
	a.Router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")
}

func main() {
	app := App{}
	app.Initialize()
	app.Run()
}

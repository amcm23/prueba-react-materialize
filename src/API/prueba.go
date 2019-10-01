package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Post struct {
	ID        string `json:"id"`
	Nombre    string `json:"nombre"`
	Apellidos string `json:"apellido"`
}

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "falor_fralg:fralg100@gmail.com@tcp(falorente.salesianas.es)/falorente_react")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	/*router := mux.NewRouter().StrictSlash(true)
	//mux.CORSMethodMiddleware(router)

	router.HandleFunc("/posts", getPosts).Methods("GET")
	router.HandleFunc("/posts", createPost).Methods("POST")
	//router.HandleFunc("/posts/{id}", getPost).Methods("GET")
	//router.HandleFunc("/posts/{id}", updatePost).Methods("PUT")
	router.HandleFunc("/posts/{id}", deletePost).Methods("DELETE", "OPTIONS")
	http.ListenAndServe(":8000", handlers.CORS()(router))

	c := cors.New(cors.Options{
		AllowedMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowedOrigins:     []string{"*"},
		AllowCredentials:   true,
		AllowedHeaders:     []string{"Content-Type", "Bearer", "Bearer ", "content-type", "Origin", "Accept"},
		OptionsPassthrough: true,
	})

	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":8000", handler))*/
	r := chi.NewRouter()

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	r.Get("/posts", getPosts)
	r.Delete("/posts/{id}", deletePost)
	r.Post("/posts", createPost)

	http.ListenAndServe(":8000", r)
}

func deletePost(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS,POST,DELETE")
	//w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")

	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM clientes WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Post with ID = %s was deleted", params["id"])
}

func getPosts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	var posts []Post
	result, err := db.Query("SELECT * from clientes")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var post Post
		err := result.Scan(&post.ID, &post.Nombre, &post.Apellidos)
		if err != nil {
			panic(err.Error())
		}
		posts = append(posts, post)
	}
	json.NewEncoder(w).Encode(posts)
}

/*func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}*/
/*func Insert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	if r.Method == "POST" {
		name := r.FormValue("nombre")
		city := r.FormValue("apellido")
		insForm, err := db.Prepare("INSERT INTO clientes(nombre, apellido) VALUES(?,?)")
		if err != nil {
			panic(err.Error())
		}
		insForm.Exec(name, city)
		log.Println("INSERT: Name: " + name + " | City: " + city)
	}

}*/

func createPost(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	//w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	stmt, err := db.Prepare("INSERT INTO clientes(nombre,apellido) VALUES(?,?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	title := keyVal["nombre"]
	apellido := keyVal["apellido"]
	_, err = stmt.Exec(title, apellido)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New post was created")
	log.Println("INSERT: Name: " + title + " | City: " + apellido)
}

func getPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT id, nombre FROM clientes WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var post Post
	for result.Next() {
		err := result.Scan(&post.ID, &post.Nombre)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(post)
}

/*func updatePost(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	stmt, err := db.Prepare("UPDATE clientes SET nombre = ? WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	newTitle := keyVal["nombre"]
	_, err = stmt.Exec(newTitle, params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Post with ID = %s was updated", params["id"])
}*/

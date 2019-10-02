package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/cors"
)

type Cliente struct {
	ID              string `json:"id_cliente"`
	Nombre          string `json:"nombre"`
	Apellido        string `json:"apellido"`
	Direccion       string `json:"direccion"`
	FechaNacimiento string `json:"fecha_nacimiento"`
	Telefono        string `json:"telefono"`
	Email           string `json:"email"`
}

type Factura struct {
	ID      string `json:"num_factura"`
	Cliente string `json:"id_cliente"`
	Fecha   string `json:"fecha"`
}

type Detalle struct {
	ID       string `json:"num_detalle"`
	Factura  string `json:"id_factura"`
	Producto string `json:"id_producto"`
	Cantidad string `json:"cantidad"`
	Precio   string `json:"precio"`
}

type Producto struct {
	ID     string `json:"num_factura"`
	Nombre string `json:"nombre"`
	Precio string `json:"precio"`
	Stock  string `json:"stock"`
}

var db *sql.DB
var err error

func main() {
	db, err = sql.Open("mysql", "falor_fralg:fralg100@gmail.com@tcp(falorente.salesianas.es)/falorente_clientes")
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
	r.Get("/clientes", getClientes)
	r.Delete("/clientes/{id_cliente}", deleteCliente)
	r.Post("/clientes", createCliente)

	r.Get("/facturas", getFacturas)
	r.Delete("/facturas/{num_factura}", deleteFactura)
	r.Post("/facturas", createFactura)

	r.Get("/detalles", getDetalles)
	r.Delete("/detalles/{num_detalle}", deleteDetalle)
	r.Post("/detalles", createDetalle)

	r.Get("/productos", getProductos)
	r.Delete("/productos/{id_producto}", deleteProducto)
	r.Post("/productos", createProducto)

	http.ListenAndServe(":8000", r)
}

func deleteCliente(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS,POST,DELETE")
	//w.Header().Set("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
	/*log.Println("DELETE: " + id)
	params := mux.Vars(r)
	stmt, err := db.Prepare("DELETE FROM clientes WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "Cliente with ID = %s was deleted", params["id"])*/
	id := chi.URLParam(r, "id_cliente")

	query, err := db.Prepare("delete from clientes where id_cliente=?")
	catch(err)
	_, er := query.Exec(id)
	catch(er)
	query.Close()
}

func catch(err error) {
	if err != nil {
		panic(err)
	}
}

func getClientes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	var clientes []Cliente
	result, err := db.Query("SELECT * from clientes")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var cliente Cliente
		err := result.Scan(&cliente.ID, &cliente.Nombre, &cliente.Apellido, &cliente.Direccion, &cliente.FechaNacimiento, &cliente.Telefono, &cliente.Email)
		if err != nil {
			panic(err.Error())
		}
		clientes = append(clientes, cliente)
	}
	json.NewEncoder(w).Encode(clientes)
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

func createCliente(w http.ResponseWriter, r *http.Request) {
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	//w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	stmt, err := db.Prepare("INSERT INTO clientes(id_cliente,nombre,apellido,direccion,fecha_nacimiento,telefono,email) VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	id_cliente := keyVal["id_cliente"]
	nombre := keyVal["nombre"]
	apellido := keyVal["apellido"]
	direccion := keyVal["direccion"]
	fecha_nacimiento := keyVal["fecha_nacimiento"]
	telefono := keyVal["telefono"]
	email := keyVal["email"]
	_, err = stmt.Exec(id_cliente, nombre, apellido, direccion, fecha_nacimiento, telefono, email)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New post was created")
}

/*func getPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	result, err := db.Query("SELECT id, nombre FROM clientes WHERE id = ?", params["id"])
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	var post Cliente
	for result.Next() {
		err := result.Scan(&post.ID, &post.Nombre)
		if err != nil {
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(post)
}*/

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
	fmt.Fprintf(w, "Cliente with ID = %s was updated", params["id"])
}*/

func deleteFactura(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "num_factura")

	query, err := db.Prepare("delete from facturas where num_factura=?")
	catch(err)
	_, er := query.Exec(id)
	catch(er)
	query.Close()
}

func getFacturas(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	var facturas []Factura
	result, err := db.Query("SELECT * from facturas")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var factura Factura
		err := result.Scan(&factura.ID, &factura.Cliente, &factura.Fecha)
		if err != nil {
			panic(err.Error())
		}
		facturas = append(facturas, factura)
	}
	json.NewEncoder(w).Encode(facturas)
}

func createFactura(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare("INSERT INTO facturas(id_cliente,fecha) VALUES(?,?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	//num_factura := keyVal["num_factura"]
	id_cliente := keyVal["id_cliente"]
	fecha := keyVal["fecha"]
	_, err = stmt.Exec(id_cliente, fecha)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New post was created")
}

func deleteDetalle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "num_detalle")

	query, err := db.Prepare("delete from detalles where num_detalle=?")
	catch(err)
	_, er := query.Exec(id)
	catch(er)
	query.Close()
}

func getDetalles(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	var detalles []Detalle
	result, err := db.Query("SELECT * from detalles")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var detalle Detalle
		err := result.Scan(&detalle.ID, &detalle.Factura, &detalle.Producto, &detalle.Cantidad, &detalle.Precio)
		if err != nil {
			panic(err.Error())
		}
		detalles = append(detalles, detalle)
	}
	json.NewEncoder(w).Encode(detalles)
}

func createDetalle(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare("INSERT INTO detalles(id_factura,id_producto,cantidad,precio) VALUES(?,?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	//num_factura := keyVal["num_factura"]
	id_factura := keyVal["id_factura"]
	id_producto := keyVal["id_producto"]
	cantidad := keyVal["cantidad"]
	precio := keyVal["precio"]
	_, err = stmt.Exec(id_factura, id_producto, cantidad, precio)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New post was created")
}

func deleteProducto(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id_producto")

	query, err := db.Prepare("delete from productos where id_producto=?")
	catch(err)
	_, er := query.Exec(id)
	catch(er)
	query.Close()
}

func getProductos(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")
	var productos []Producto
	result, err := db.Query("SELECT * from productos")
	if err != nil {
		panic(err.Error())
	}
	defer result.Close()
	for result.Next() {
		var producto Producto
		err := result.Scan(&producto.ID, &producto.Nombre, &producto.Precio, &producto.Stock)
		if err != nil {
			panic(err.Error())
		}
		productos = append(productos, producto)
	}
	json.NewEncoder(w).Encode(productos)
}

func createProducto(w http.ResponseWriter, r *http.Request) {
	stmt, err := db.Prepare("INSERT INTO productos(nombre,precio,stock) VALUES(?,?,?)")
	if err != nil {
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err.Error())
	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	//num_factura := keyVal["num_factura"]
	nombre := keyVal["nombre"]
	precio := keyVal["precio"]
	stock := keyVal["stock"]
	_, err = stmt.Exec(nombre, precio, stock)
	if err != nil {
		panic(err.Error())
	}
	fmt.Fprintf(w, "New post was created")
}

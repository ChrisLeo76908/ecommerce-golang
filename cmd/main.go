/*
@nombre:Christian David Muñoz Tonguino
@fecha: 23/05/2026
@descripción: Sistema de gestión E-Commerce desarrollado en GoLang y MySQL, que permite administrar productos,
carrito de compras y panel administrativo mediante una interfaz web dinámica.
*/

package main

import (
	"ecommerce/internal/carrito"
	"ecommerce/internal/database"
	"ecommerce/internal/productos"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var sesionIniciada bool

func inicio(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/index.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func paginaProductos(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/productos.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	listaProductos := productos.ObtenerProductos()

	tmpl.Execute(w, listaProductos)
}

func agregarProducto(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	listaProductos := productos.ObtenerProductos()

	for _, producto := range listaProductos {

		if id == strconv.Itoa(producto.ID) {
			carrito.AgregarProducto(producto)
		}
	}

	http.Redirect(w, r, "/productos", http.StatusSeeOther)
}

func eliminarProducto(w http.ResponseWriter, r *http.Request) {

	idTexto := r.URL.Query().Get("id")

	id, _ := strconv.Atoi(idTexto)

	carrito.EliminarProducto(id)

	http.Redirect(w, r, "/carrito", http.StatusSeeOther)
}

func finalizarCompra(w http.ResponseWriter, r *http.Request) {

	carrito.VaciarCarrito()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func carritoPagina(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/carrito.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	datos := map[string]interface{}{
		"Productos": carrito.ObtenerCarrito(),
		"Total":     carrito.CalcularTotal(),
	}

	tmpl.Execute(w, datos)
}

func adminPagina(w http.ResponseWriter, r *http.Request) {

	if !sesionIniciada {

		http.Redirect(w, r, "/login", http.StatusSeeOther)

		return
	}

	tmpl, err := template.ParseFiles("templates/admin.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	listaProductos := productos.ObtenerProductos()

	tmpl.Execute(w, listaProductos)
}

func loginPagina(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmpl, err := template.ParseFiles("templates/login.html")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, nil)

		return
	}

	usuario := r.FormValue("usuario")

	password := r.FormValue("password")

	if usuario == "admin" && password == "1234" {

		sesionIniciada = true

		http.Redirect(w, r, "/admin", http.StatusSeeOther)

		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func logout(w http.ResponseWriter, r *http.Request) {

	sesionIniciada = false

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func guardarProducto(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	nombre := r.FormValue("nombre")

	precioTexto := r.FormValue("precio")

	imagen := r.FormValue("imagen")

	precio, _ := strconv.ParseFloat(precioTexto, 64)

	nuevoProducto := productos.Producto{
		ID:     len(productos.ObtenerProductos()) + 1,
		Nombre: nombre,
		Precio: precio,
		Imagen: imagen,
	}

	productos.AgregarNuevoProducto(nuevoProducto)

	http.Redirect(w, r, "/productos", http.StatusSeeOther)
}

func eliminarProductoAdmin(w http.ResponseWriter, r *http.Request) {

	idTexto := r.URL.Query().Get("id")

	id, _ := strconv.Atoi(idTexto)

	productos.EliminarProducto(id)

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func main() {

	db := database.Conexion()

	defer db.Close()

	log.Println("Base de datos conectada correctamente")

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", inicio)
	http.HandleFunc("/productos", paginaProductos)
	http.HandleFunc("/carrito", carritoPagina)
	http.HandleFunc("/agregar", agregarProducto)
	http.HandleFunc("/eliminar", eliminarProducto)
	http.HandleFunc("/finalizar", finalizarCompra)

	http.HandleFunc("/admin", adminPagina)
	http.HandleFunc("/guardar-producto", guardarProducto)
	http.HandleFunc("/eliminar-producto", eliminarProductoAdmin)

	http.HandleFunc("/login", loginPagina)
	http.HandleFunc("/logout", logout)

	log.Println("Servidor ejecutándose en http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}

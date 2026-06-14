/*
@nombre: Christian David Muñoz Tonguino
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
		http.Error(w, "Error al cargar la página de inicio", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

// paginaProductos gestiona la visualización del catálogo.
// Esta función integra búsqueda por nombre, ordenamiento por precio
// y paginación para mostrar los productos de forma organizada.
func paginaProductos(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/productos.html")

	if err != nil {
		http.Error(w, "Error al cargar la página de productos", http.StatusInternalServerError)
		return
	}

	// Se obtiene el texto ingresado por el usuario en el buscador.
	busqueda := r.URL.Query().Get("buscar")

	var listaProductos []productos.Producto

	if busqueda != "" {
		listaProductos, err = productos.BuscarProductos(busqueda)
	} else {
		listaProductos, err = productos.ObtenerProductos()
	}

	if err != nil {
		http.Error(w, "Error al obtener productos desde la base de datos", http.StatusInternalServerError)
		return
	}

	// Se obtiene el criterio de ordenamiento seleccionado por el usuario.
	orden := r.URL.Query().Get("orden")

	listaProductos = productos.OrdenarProductosPorPrecio(listaProductos, orden)

	productosPorPagina := 50
	paginaTexto := r.URL.Query().Get("pagina")
	paginaActual := 1

	if paginaTexto != "" {
		paginaConvertida, err := strconv.Atoi(paginaTexto)

		if err == nil && paginaConvertida > 0 {
			paginaActual = paginaConvertida
		}
	}

	// Se calculan los límites de inicio y fin para mostrar solo los productos de la página actual.
	inicio := (paginaActual - 1) * productosPorPagina
	fin := inicio + productosPorPagina

	if inicio > len(listaProductos) {
		inicio = len(listaProductos)
	}

	if fin > len(listaProductos) {
		fin = len(listaProductos)
	}

	productosPagina := listaProductos[inicio:fin]

	var paginaAnterior int
	var paginaSiguiente int

	if paginaActual > 1 {
		paginaAnterior = paginaActual - 1
	}

	if fin < len(listaProductos) {
		paginaSiguiente = paginaActual + 1
	}

	datos := map[string]interface{}{
		"Productos":       productosPagina,
		"PaginaActual":    paginaActual,
		"PaginaAnterior":  paginaAnterior,
		"PaginaSiguiente": paginaSiguiente,
		"Busqueda":        busqueda,
		"Orden":           orden,
	}
	tmpl.Execute(w, datos)
}

// agregarProducto recibe el ID del producto seleccionado,
// lo busca en la base de datos y lo agrega temporalmente al carrito.
func agregarProducto(w http.ResponseWriter, r *http.Request) {

	id := r.URL.Query().Get("id")

	listaProductos, err := productos.ObtenerProductos()

	if err != nil {
		http.Error(w, "Error al obtener productos", http.StatusInternalServerError)
		return
	}

	for _, producto := range listaProductos {

		if id == strconv.Itoa(producto.ID) {
			carrito.AgregarProducto(producto)
			break
		}
	}

	http.Redirect(w, r, "/productos", http.StatusSeeOther)
}

func eliminarProducto(w http.ResponseWriter, r *http.Request) {

	idTexto := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idTexto)

	if err != nil {
		http.Error(w, "ID de producto inválido", http.StatusBadRequest)
		return
	}

	carrito.EliminarProducto(id)

	http.Redirect(w, r, "/carrito", http.StatusSeeOther)
}

func finalizarCompra(w http.ResponseWriter, r *http.Request) {

	resumen, total := carrito.GenerarResumen()

	tmpl, err := template.ParseFiles("templates/resumen.html")

	if err != nil {
		http.Error(w, "Error al cargar el resumen de compra", http.StatusInternalServerError)
		return
	}

	datos := map[string]interface{}{
		"Productos": resumen,
		"Total":     total,
	}

	carrito.VaciarCarrito()

	tmpl.Execute(w, datos)
}

func carritoPagina(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/carrito.html")

	if err != nil {
		http.Error(w, "Error al cargar el carrito", http.StatusInternalServerError)
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
		http.Error(w, "Error al cargar el panel administrador", http.StatusInternalServerError)
		return
	}

	listaProductos, err := productos.ObtenerProductos()

	if err != nil {
		http.Error(w, "Error al obtener productos para administración", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, listaProductos)
}

func loginPagina(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmpl, err := template.ParseFiles("templates/login.html")

		if err != nil {
			http.Error(w, "Error al cargar login", http.StatusInternalServerError)
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

// guardarProducto valida los datos enviados desde el formulario del administrador.
// Si los datos son correctos, registra el producto en MySQL.
func guardarProducto(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	nombre := r.FormValue("nombre")
	precioTexto := r.FormValue("precio")
	imagen := r.FormValue("imagen")

	precio, err := strconv.ParseFloat(precioTexto, 64)

	if err != nil {
		http.Error(w, "El precio ingresado no es válido", http.StatusBadRequest)
		return
	}

	nuevoProducto := productos.Producto{
		Nombre: nombre,
		Precio: precio,
		Imagen: imagen,
	}

	err = productos.AgregarNuevoProducto(nuevoProducto)

	if err != nil {
		http.Error(w, "Error al guardar producto: "+err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/productos", http.StatusSeeOther)
}

func eliminarProductoAdmin(w http.ResponseWriter, r *http.Request) {

	idTexto := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idTexto)

	if err != nil {
		http.Error(w, "ID inválido para eliminar producto", http.StatusBadRequest)
		return
	}

	err = productos.EliminarProducto(id)

	if err != nil {
		http.Error(w, "Error al eliminar producto: "+err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

// editarProductoPagina carga el formulario de edición.
// Primero valida que exista una sesión de administrador,
// luego busca el producto seleccionado por su ID.
func editarProductoPagina(w http.ResponseWriter, r *http.Request) {

	if !sesionIniciada {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	idTexto := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idTexto)

	if err != nil {
		http.Error(w, "ID inválido para editar producto", http.StatusBadRequest)
		return
	}

	producto, err := productos.ObtenerProductoPorID(id)

	if err != nil {
		http.Error(w, "Producto no encontrado", http.StatusNotFound)
		return
	}

	tmpl, err := template.ParseFiles("templates/editar.html")

	if err != nil {
		http.Error(w, "Error al cargar la página de edición", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, producto)
}

// actualizarProducto recibe los datos modificados desde el formulario,
// valida la información ingresada y actualiza el producto en MySQL.
func actualizarProducto(w http.ResponseWriter, r *http.Request) {

	if !sesionIniciada {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	idTexto := r.FormValue("id")
	nombre := r.FormValue("nombre")
	precioTexto := r.FormValue("precio")
	imagen := r.FormValue("imagen")

	id, err := strconv.Atoi(idTexto)

	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	precio, err := strconv.ParseFloat(precioTexto, 64)

	if err != nil {
		http.Error(w, "Precio inválido", http.StatusBadRequest)
		return
	}

	productoActualizado := productos.Producto{
		ID:     id,
		Nombre: nombre,
		Precio: precio,
		Imagen: imagen,
	}

	err = productos.ActualizarProducto(productoActualizado)

	if err != nil {
		http.Error(w, "Error al actualizar producto: "+err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func resumenCompraPagina(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/resumen.html")

	if err != nil {
		http.Error(w, "Error al cargar el resumen de compra", http.StatusInternalServerError)
		return
	}

	datos := map[string]interface{}{
		"Productos": carrito.UltimoResumen,
		"Total":     carrito.UltimoTotal,
	}

	tmpl.Execute(w, datos)
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
	http.HandleFunc("/resumen", resumenCompraPagina)

	http.HandleFunc("/admin", adminPagina)
	http.HandleFunc("/guardar-producto", guardarProducto)
	http.HandleFunc("/eliminar-producto", eliminarProductoAdmin)
	http.HandleFunc("/editar-producto", editarProductoPagina)
	http.HandleFunc("/actualizar-producto", actualizarProducto)

	http.HandleFunc("/login", loginPagina)
	http.HandleFunc("/logout", logout)

	log.Println("Servidor ejecutándose en http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}

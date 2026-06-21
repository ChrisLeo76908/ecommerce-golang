/*
@nombre: Christian David Muñoz Tonguino
@fecha: 14/06/2026

@descripción:
Sistema de gestión E-Commerce desarrollado en GoLang y MySQL que permite administrar
productos mediante un CRUD completo, realizar búsquedas, ordenamiento por precio,
paginación del catálogo, gestión de carrito de compras y generación de resúmenes
de compra. El proyecto aplica arquitectura modular, encapsulación, manejo de errores
y programación orientada a objetos mediante estructuras, métodos e interfaces.
*/

package main

import (
	"ecommerce/internal/carrito"
	"ecommerce/internal/compras"
	"ecommerce/internal/database"
	"ecommerce/internal/productos"
	"ecommerce/internal/usuarios"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var sesionIniciada bool
var clienteActual usuarios.Usuario

func inicio(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/index.html")

	if err != nil {
		http.Error(w, "Error al cargar la página de inicio", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, datosCliente())
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
		"ClienteActivo":   clienteActual.ID != 0,
		"ClienteNombre":   clienteActual.Nombre,
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

	busqueda := r.URL.Query().Get("buscar")
	orden := r.URL.Query().Get("orden")
	pagina := r.URL.Query().Get("pagina")

	url := "/productos?buscar=" + busqueda + "&orden=" + orden + "&pagina=" + pagina

	http.Redirect(w, r, url, http.StatusSeeOther)

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

	if clienteActual.ID == 0 {
		http.Redirect(w, r, "/login-cliente", http.StatusSeeOther)
		return
	}

	resumen, subtotal, iva, total := carrito.GenerarResumen()

	if len(resumen) == 0 {
		http.Redirect(w, r, "/carrito", http.StatusSeeOther)
		return
	}

	err := compras.RegistrarCompra(clienteActual.ID, resumen, subtotal, iva, total)

	if err != nil {
		http.Error(w, "Error al registrar la compra", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("templates/resumen.html")

	if err != nil {
		http.Error(w, "Error al cargar el resumen de compra", http.StatusInternalServerError)
		return
	}

	datos := map[string]interface{}{
		"Productos":     resumen,
		"Subtotal":      subtotal,
		"IVA":           iva,
		"Total":         total,
		"Cliente":       clienteActual.Nombre,
		"ClienteActivo": true,
		"ClienteNombre": clienteActual.Nombre,
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
		"Productos":     carrito.ObtenerCarrito(),
		"Total":         carrito.CalcularTotal(),
		"ClienteActivo": clienteActual.ID != 0,
		"ClienteNombre": clienteActual.Nombre,
	}

	tmpl.Execute(w, datos)
}

func sesionAdminActiva(r *http.Request) bool {

	cookie, err := r.Cookie("admin")

	if err != nil {
		return false
	}

	return cookie.Value == "activo"
}

func adminPagina(w http.ResponseWriter, r *http.Request) {

	if !sesionAdminActiva(r) {
		http.Redirect(w, r, "/login-admin", http.StatusSeeOther)
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

	orden := r.URL.Query().Get("orden")

	listaProductos = productos.OrdenarProductosAdmin(listaProductos, orden)

	tmpl.Execute(w, listaProductos)

}

func loginPagina(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodGet {

		tmpl, err := template.ParseFiles("templates/login.html")

		if err != nil {
			http.Error(w, "Error al cargar login", http.StatusInternalServerError)
			return
		}

		tmpl.Execute(w, datosCliente())
		return
	}

	usuario := r.FormValue("usuario")
	password := r.FormValue("password")

	if usuario == "admin" && password == "1234" {
		http.SetCookie(w, &http.Cookie{
			Name:  "admin",
			Value: "activo",
			Path:  "/",
		})

		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func logout(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:   "admin",
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	})

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

	http.Redirect(w, r, "/admin", http.StatusSeeOther)
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

	if !sesionAdminActiva(r) {
		http.Redirect(w, r, "/login-admin", http.StatusSeeOther)
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

	if !sesionAdminActiva(r) {
		http.Redirect(w, r, "/login-admin", http.StatusSeeOther)
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

func responderJSON(w http.ResponseWriter, datos interface{}, estado int) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(estado)
	json.NewEncoder(w).Encode(datos)
}

// Servicio 1: devuelve todos los productos registrados en MySQL.
func apiObtenerProductos(w http.ResponseWriter, r *http.Request) {

	listaProductos, err := productos.ObtenerProductos()

	if err != nil {
		responderJSON(w, map[string]string{"error": "No se pudieron obtener los productos"}, http.StatusInternalServerError)
		return
	}

	responderJSON(w, listaProductos, http.StatusOK)
}

// Servicio 2: devuelve un producto específico según su ID.
func apiObtenerProductoPorID(w http.ResponseWriter, r *http.Request) {

	idTexto := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idTexto)

	if err != nil {
		responderJSON(w, map[string]string{"error": "ID inválido"}, http.StatusBadRequest)
		return
	}

	producto, err := productos.ObtenerProductoPorID(id)

	if err != nil {
		responderJSON(w, map[string]string{"error": "Producto no encontrado"}, http.StatusNotFound)
		return
	}

	responderJSON(w, producto, http.StatusOK)
}

// Servicio 3: busca productos por nombre.
func apiBuscarProductos(w http.ResponseWriter, r *http.Request) {

	nombre := r.URL.Query().Get("nombre")

	resultado, err := productos.BuscarProductos(nombre)

	if err != nil {
		responderJSON(w, map[string]string{"error": "Error al buscar productos"}, http.StatusInternalServerError)
		return
	}

	responderJSON(w, resultado, http.StatusOK)
}

// Servicio 4: registra un nuevo producto usando JSON.
func apiAgregarProducto(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		responderJSON(w, map[string]string{"error": "Método no permitido"}, http.StatusMethodNotAllowed)
		return
	}

	var nuevoProducto productos.Producto

	err := json.NewDecoder(r.Body).Decode(&nuevoProducto)

	if err != nil {
		responderJSON(w, map[string]string{"error": "JSON inválido"}, http.StatusBadRequest)
		return
	}

	err = productos.AgregarNuevoProducto(nuevoProducto)

	if err != nil {
		responderJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	responderJSON(w, map[string]string{"mensaje": "Producto agregado correctamente"}, http.StatusCreated)
}

// Servicio 5: actualiza un producto existente usando JSON.
func apiEditarProducto(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		responderJSON(w, map[string]string{"error": "Método no permitido"}, http.StatusMethodNotAllowed)
		return
	}

	var productoActualizado productos.Producto

	err := json.NewDecoder(r.Body).Decode(&productoActualizado)

	if err != nil {
		responderJSON(w, map[string]string{"error": "JSON inválido"}, http.StatusBadRequest)
		return
	}

	err = productos.ActualizarProducto(productoActualizado)

	if err != nil {
		responderJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	responderJSON(w, map[string]string{"mensaje": "Producto actualizado correctamente"}, http.StatusOK)
}

// Servicio 6: elimina un producto según su ID.
func apiEliminarProducto(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodDelete {
		responderJSON(w, map[string]string{"error": "Método no permitido"}, http.StatusMethodNotAllowed)
		return
	}

	idTexto := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idTexto)

	if err != nil {
		responderJSON(w, map[string]string{"error": "ID inválido"}, http.StatusBadRequest)
		return
	}

	err = productos.EliminarProducto(id)

	if err != nil {
		responderJSON(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	responderJSON(w, map[string]string{"mensaje": "Producto eliminado correctamente"}, http.StatusOK)
}

// Servicio 7: devuelve los productos actuales del carrito.
func apiObtenerCarrito(w http.ResponseWriter, r *http.Request) {

	responderJSON(w, carrito.ObtenerCarrito(), http.StatusOK)
}

// Servicio 8: devuelve el último resumen de compra generado.
func apiObtenerResumen(w http.ResponseWriter, r *http.Request) {

	datos := map[string]interface{}{
		"productos": carrito.UltimoResumen,
		"subtotal":  carrito.UltimoSubtotal,
		"iva":       carrito.UltimoIVA,
		"total":     carrito.UltimoTotal,
	}

	responderJSON(w, datos, http.StatusOK)
}

func registroClientePagina(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/registro.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func registrarCliente(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/registro-cliente", http.StatusSeeOther)
		return
	}

	usuario := usuarios.Usuario{
		Nombre:   r.FormValue("nombre"),
		Correo:   r.FormValue("correo"),
		Password: r.FormValue("password"),
	}

	err := usuarios.RegistrarUsuario(usuario)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/login-cliente", http.StatusSeeOther)
}

func loginClientePagina(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/login-cliente.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func validarCliente(w http.ResponseWriter, r *http.Request) {

	correo := r.FormValue("correo")
	password := r.FormValue("password")

	usuario, err := usuarios.ValidarUsuario(correo, password)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	clienteActual = usuario

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func accesoPagina(w http.ResponseWriter, r *http.Request) {

	tmpl, err := template.ParseFiles("templates/acceso.html")

	if err != nil {
		http.Error(w, "Error al cargar página de acceso", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func logoutCliente(w http.ResponseWriter, r *http.Request) {

	clienteActual = usuarios.Usuario{}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func datosCliente() map[string]interface{} {
	return map[string]interface{}{
		"ClienteActivo": clienteActual.ID != 0,
		"ClienteNombre": clienteActual.Nombre,
	}
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

	http.HandleFunc("/login", accesoPagina)
	http.HandleFunc("/login-admin", loginPagina)
	http.HandleFunc("/logout", logout)

	// Servicios Web JSON
	http.HandleFunc("/api/productos", apiObtenerProductos)
	http.HandleFunc("/api/producto", apiObtenerProductoPorID)
	http.HandleFunc("/api/buscar", apiBuscarProductos)
	http.HandleFunc("/api/agregar-producto", apiAgregarProducto)
	http.HandleFunc("/api/editar-producto", apiEditarProducto)
	http.HandleFunc("/api/eliminar-producto", apiEliminarProducto)
	http.HandleFunc("/api/carrito", apiObtenerCarrito)
	http.HandleFunc("/api/resumen", apiObtenerResumen)

	http.HandleFunc("/registro-cliente", registroClientePagina)
	http.HandleFunc("/registrar-cliente", registrarCliente)
	http.HandleFunc("/login-cliente", loginClientePagina)
	http.HandleFunc("/validar-cliente", validarCliente)
	http.HandleFunc("/logout-cliente", logoutCliente)

	log.Println("Servidor ejecutándose en http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}

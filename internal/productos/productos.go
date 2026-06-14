package productos

import (
	"ecommerce/internal/database"
	"errors"
	"sort"
	"strings"
)

type Producto struct {
	ID     int
	Nombre string
	Precio float64
	Imagen string
}

// ObtenerID devuelve el identificador del producto.
func (p Producto) ObtenerID() int {
	return p.ID
}

// ObtenerNombre devuelve el nombre del producto.
func (p Producto) ObtenerNombre() string {
	return p.Nombre
}

// ObtenerPrecio devuelve el precio del producto.
func (p Producto) ObtenerPrecio() float64 {
	return p.Precio
}

// ObtenerImagen devuelve la ruta de imagen del producto.
func (p Producto) ObtenerImagen() string {
	return p.Imagen
}

// ProductoRepository define las operaciones que debe cumplir
// cualquier módulo encargado de gestionar productos.
type ProductoRepository interface {
	ObtenerProductos() ([]Producto, error)
	AgregarNuevoProducto(producto Producto) error
	EliminarProducto(id int) error
}

// ValidarProducto aplica encapsulación al controlar que los datos
// del producto sean correctos antes de guardarlos en la base de datos.
func ValidarProducto(producto Producto) error {

	if producto.Nombre == "" {
		return errors.New("el nombre del producto no puede estar vacío")
	}

	if producto.Precio <= 0 {
		return errors.New("el precio del producto debe ser mayor a cero")
	}

	if producto.Imagen == "" {
		producto.Imagen = "/static/img/default.jpg"
	}

	return nil
}

// ObtenerProductos consulta la base de datos MySQL y devuelve
// todos los productos registrados en el sistema.
func ObtenerProductos() ([]Producto, error) {

	db := database.Conexion()
	defer db.Close()

	rows, err := db.Query("SELECT id, nombre, precio, imagen FROM productos")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var listaProductos []Producto

	for rows.Next() {

		var producto Producto

		err := rows.Scan(
			&producto.ID,
			&producto.Nombre,
			&producto.Precio,
			&producto.Imagen,
		)

		if err != nil {
			return nil, err
		}

		listaProductos = append(listaProductos, producto)
	}

	return listaProductos, nil
}

// AgregarNuevoProducto valida y registra un nuevo producto
// dentro de la tabla productos de MySQL.
func AgregarNuevoProducto(producto Producto) error {

	err := ValidarProducto(producto)

	if err != nil {
		return err
	}

	db := database.Conexion()
	defer db.Close()

	query := `
	INSERT INTO productos(nombre, precio, imagen)
	VALUES (?, ?, ?)
	`

	_, err = db.Exec(
		query,
		producto.Nombre,
		producto.Precio,
		producto.Imagen,
	)

	if err != nil {
		return err
	}

	return nil
}

// EliminarProducto elimina un producto de la base de datos
// utilizando su ID como identificador principal.
func EliminarProducto(id int) error {

	if id <= 0 {
		return errors.New("el ID del producto no es válido")
	}

	db := database.Conexion()
	defer db.Close()

	query := "DELETE FROM productos WHERE id = ?"

	_, err := db.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}

// BuscarProductos filtra los productos por nombre.
// La búsqueda no distingue entre mayúsculas y minúsculas,
// permitiendo encontrar coincidencias parciales.
func BuscarProductos(texto string) ([]Producto, error) {

	productos, err := ObtenerProductos()

	if err != nil {
		return nil, err
	}

	var resultado []Producto

	for _, producto := range productos {

		if strings.Contains(
			strings.ToLower(producto.Nombre),
			strings.ToLower(texto),
		) {
			resultado = append(resultado, producto)
		}
	}

	return resultado, nil
}

// OrdenarProductosPorPrecio organiza la lista de productos
// de menor a mayor o de mayor a menor según el criterio recibido.
func OrdenarProductosPorPrecio(lista []Producto, orden string) []Producto {

	if orden == "menor" {

		sort.Slice(lista, func(i, j int) bool {
			return lista[i].Precio < lista[j].Precio
		})

	} else if orden == "mayor" {

		sort.Slice(lista, func(i, j int) bool {
			return lista[i].Precio > lista[j].Precio
		})
	}

	return lista
}

// ObtenerProductoPorID busca un producto específico por su ID.
func ObtenerProductoPorID(id int) (Producto, error) {

	if id <= 0 {
		return Producto{}, errors.New("el ID del producto no es válido")
	}

	db := database.Conexion()
	defer db.Close()

	query := "SELECT id, nombre, precio, imagen FROM productos WHERE id = ?"

	var producto Producto

	err := db.QueryRow(query, id).Scan(
		&producto.ID,
		&producto.Nombre,
		&producto.Precio,
		&producto.Imagen,
	)

	if err != nil {
		return Producto{}, err
	}

	return producto, nil
}

// ActualizarProducto modifica los datos de un producto existente.
// Antes de ejecutar la actualización, valida que el producto tenga
// un ID correcto y datos válidos.
func ActualizarProducto(producto Producto) error {

	if producto.ID <= 0 {
		return errors.New("el ID del producto no es válido")
	}

	err := ValidarProducto(producto)

	if err != nil {
		return err
	}

	db := database.Conexion()
	defer db.Close()

	query := `
	UPDATE productos
	SET nombre = ?, precio = ?, imagen = ?
	WHERE id = ?
	`

	_, err = db.Exec(
		query,
		producto.Nombre,
		producto.Precio,
		producto.Imagen,
		producto.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

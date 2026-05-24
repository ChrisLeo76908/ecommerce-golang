package productos

import (
	"ecommerce/internal/database"
)

type Producto struct {
	ID     int
	Nombre string
	Precio float64
	Imagen string
}

func ObtenerProductos() []Producto {

	db := database.Conexion()

	rows, err := db.Query("SELECT * FROM productos")

	if err != nil {
		return nil
	}

	var productos []Producto

	for rows.Next() {

		var producto Producto

		rows.Scan(
			&producto.ID,
			&producto.Nombre,
			&producto.Precio,
			&producto.Imagen,
		)

		productos = append(productos, producto)
	}

	return productos
}

func AgregarNuevoProducto(producto Producto) {

	db := database.Conexion()

	query := `
	INSERT INTO productos(nombre, precio, imagen)
	VALUES (?, ?, ?)
	`

	db.Exec(
		query,
		producto.Nombre,
		producto.Precio,
		producto.Imagen,
	)
}

func EliminarProducto(id int) {

	db := database.Conexion()

	query := "DELETE FROM productos WHERE id = ?"

	db.Exec(query, id)
}

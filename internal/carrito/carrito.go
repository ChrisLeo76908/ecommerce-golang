package carrito

import "ecommerce/internal/productos"

var Carrito []productos.Producto

func AgregarProducto(producto productos.Producto) {
	Carrito = append(Carrito, producto)
}

func ObtenerCarrito() []productos.Producto {
	return Carrito
}

func CalcularTotal() float64 {

	var total float64

	for _, producto := range Carrito {
		total += producto.Precio
	}

	return total
}

func EliminarProducto(id int) {

	var nuevoCarrito []productos.Producto

	eliminado := false

	for _, producto := range Carrito {

		if producto.ID == id && !eliminado {

			eliminado = true

			continue
		}

		nuevoCarrito = append(nuevoCarrito, producto)
	}

	Carrito = nuevoCarrito
}

func VaciarCarrito() {
	Carrito = []productos.Producto{}
}

package carrito

import "ecommerce/internal/productos"

var Carrito []productos.Producto

type ResumenItem struct {
	ID       int
	Nombre   string
	Precio   float64
	Cantidad int
	Subtotal float64
}

var UltimoResumen []ResumenItem
var UltimoTotal float64

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

// GenerarResumen agrupa productos repetidos y calcula cantidad, subtotal y total.
func GenerarResumen() ([]ResumenItem, float64) {

	resumenMapa := make(map[int]ResumenItem)

	var total float64

	for _, producto := range Carrito {

		item, existe := resumenMapa[producto.ID]

		if existe {
			item.Cantidad++
			item.Subtotal = item.Precio * float64(item.Cantidad)
		} else {
			item = ResumenItem{
				ID:       producto.ID,
				Nombre:   producto.Nombre,
				Precio:   producto.Precio,
				Cantidad: 1,
				Subtotal: producto.Precio,
			}
		}

		resumenMapa[producto.ID] = item
		total += producto.Precio
	}

	var resumen []ResumenItem

	for _, item := range resumenMapa {
		resumen = append(resumen, item)
	}

	UltimoResumen = resumen
	UltimoTotal = total

	return resumen, total
}

func VaciarCarrito() {
	Carrito = []productos.Producto{}
}

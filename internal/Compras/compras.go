package compras

import (
	"ecommerce/internal/carrito"
	"ecommerce/internal/database"
)

func RegistrarCompra(usuarioID int, resumen []carrito.ResumenItem, subtotal float64, iva float64, total float64) error {

	db := database.Conexion()
	defer db.Close()

	resultado, err := db.Exec(
		"INSERT INTO compras(usuario_id, subtotal, iva, total) VALUES (?, ?, ?, ?)",
		usuarioID,
		subtotal,
		iva,
		total,
	)

	if err != nil {
		return err
	}

	compraID, err := resultado.LastInsertId()

	if err != nil {
		return err
	}

	for _, item := range resumen {

		_, err = db.Exec(
			"INSERT INTO detalle_compras(compra_id, producto_id, cantidad, precio_unitario, subtotal) VALUES (?, ?, ?, ?, ?)",
			compraID,
			item.ID,
			item.Cantidad,
			item.Precio,
			item.Subtotal,
		)

		if err != nil {
			return err
		}
	}

	return nil
}

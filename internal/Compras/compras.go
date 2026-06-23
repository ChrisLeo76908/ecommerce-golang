package compras

import (
	"ecommerce/internal/carrito"
	"ecommerce/internal/database"
)

type Compra struct {
	ID       int
	Cliente  string
	Fecha    string
	Subtotal float64
	IVA      float64
	Total    float64
	Detalles []DetalleCompra
}

type DetalleCompra struct {
	Producto       string
	Cantidad       int
	PrecioUnitario float64
	Subtotal       float64
}

type ClienteCompra struct {
	ID     int
	Nombre string
	Correo string
}

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

func ObtenerComprasPorUsuario(usuarioID int) ([]Compra, error) {

	db := database.Conexion()
	defer db.Close()

	rows, err := db.Query(`
		SELECT id, fecha, subtotal, iva, total
		FROM compras
		WHERE usuario_id = ?
		ORDER BY fecha DESC
	`, usuarioID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var compras []Compra

	for rows.Next() {

		var compra Compra

		err := rows.Scan(
			&compra.ID,
			&compra.Fecha,
			&compra.Subtotal,
			&compra.IVA,
			&compra.Total,
		)

		if err != nil {
			return nil, err
		}

		compra.Detalles, err = obtenerDetallesCompra(compra.ID)

		if err != nil {
			return nil, err
		}

		compras = append(compras, compra)
	}

	return compras, nil
}

func ObtenerTodasLasCompras() ([]Compra, error) {

	db := database.Conexion()
	defer db.Close()

	rows, err := db.Query(`
		SELECT c.id, u.nombre, c.fecha, c.subtotal, c.iva, c.total
		FROM compras c
		INNER JOIN usuarios u ON c.usuario_id = u.id
		ORDER BY c.fecha DESC
	`)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var compras []Compra

	for rows.Next() {

		var compra Compra

		err := rows.Scan(
			&compra.ID,
			&compra.Cliente,
			&compra.Fecha,
			&compra.Subtotal,
			&compra.IVA,
			&compra.Total,
		)

		if err != nil {
			return nil, err
		}

		compra.Detalles, err = obtenerDetallesCompra(compra.ID)

		if err != nil {
			return nil, err
		}

		compras = append(compras, compra)
	}

	return compras, nil
}

func obtenerDetallesCompra(compraID int) ([]DetalleCompra, error) {

	db := database.Conexion()
	defer db.Close()

	rows, err := db.Query(`
		SELECT p.nombre, d.cantidad, d.precio_unitario, d.subtotal
		FROM detalle_compras d
		INNER JOIN productos p ON d.producto_id = p.id
		WHERE d.compra_id = ?
	`, compraID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var detalles []DetalleCompra

	for rows.Next() {

		var detalle DetalleCompra

		err := rows.Scan(
			&detalle.Producto,
			&detalle.Cantidad,
			&detalle.PrecioUnitario,
			&detalle.Subtotal,
		)

		if err != nil {
			return nil, err
		}

		detalles = append(detalles, detalle)
	}

	return detalles, nil
}

func ObtenerClientesConCompras(busqueda string) ([]ClienteCompra, error) {

	db := database.Conexion()
	defer db.Close()

	query := `
		SELECT DISTINCT u.id, u.nombre, u.correo
		FROM usuarios u
		INNER JOIN compras c ON u.id = c.usuario_id
		WHERE u.nombre LIKE ?
		ORDER BY u.nombre ASC
	`

	rows, err := db.Query(query, "%"+busqueda+"%")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var clientes []ClienteCompra

	for rows.Next() {

		var cliente ClienteCompra

		err := rows.Scan(
			&cliente.ID,
			&cliente.Nombre,
			&cliente.Correo,
		)

		if err != nil {
			return nil, err
		}

		clientes = append(clientes, cliente)
	}

	return clientes, nil
}

func ObtenerComprasPorUsuarioYFecha(usuarioID int, fechaInicio string, fechaFin string) ([]Compra, error) {

	db := database.Conexion()
	defer db.Close()

	query := `
		SELECT id, fecha, subtotal, iva, total
		FROM compras
		WHERE usuario_id = ?
	`

	args := []interface{}{usuarioID}

	if fechaInicio != "" {
		query += " AND DATE(fecha) >= ?"
		args = append(args, fechaInicio)
	}

	if fechaFin != "" {
		query += " AND DATE(fecha) <= ?"
		args = append(args, fechaFin)
	}

	query += " ORDER BY fecha DESC"

	rows, err := db.Query(query, args...)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var compras []Compra

	for rows.Next() {

		var compra Compra

		err := rows.Scan(
			&compra.ID,
			&compra.Fecha,
			&compra.Subtotal,
			&compra.IVA,
			&compra.Total,
		)

		if err != nil {
			return nil, err
		}

		compra.Detalles, err = obtenerDetallesCompra(compra.ID)

		if err != nil {
			return nil, err
		}

		compras = append(compras, compra)
	}

	return compras, nil
}

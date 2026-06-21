package usuarios

import (
	"database/sql"
	"ecommerce/internal/database"
	"errors"
)

type Usuario struct {
	ID       int
	Nombre   string
	Correo   string
	Password string
}

func RegistrarUsuario(usuario Usuario) error {

	if usuario.Nombre == "" || usuario.Correo == "" || usuario.Password == "" {
		return errors.New("todos los campos son obligatorios")
	}

	db := database.Conexion()
	defer db.Close()

	_, err := db.Exec(
		"INSERT INTO usuarios(nombre, correo, password) VALUES (?, ?, ?)",
		usuario.Nombre,
		usuario.Correo,
		usuario.Password,
	)

	return err
}

func ValidarUsuario(correo string, password string) (Usuario, error) {

	db := database.Conexion()
	defer db.Close()

	var usuario Usuario

	err := db.QueryRow(
		"SELECT id, nombre, correo, password FROM usuarios WHERE correo = ? AND password = ?",
		correo,
		password,
	).Scan(
		&usuario.ID,
		&usuario.Nombre,
		&usuario.Correo,
		&usuario.Password,
	)

	if err == sql.ErrNoRows {
		return Usuario{}, errors.New("correo o contraseña incorrectos")
	}

	if err != nil {
		return Usuario{}, err
	}

	return usuario, nil
}

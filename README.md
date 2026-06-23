# Sistema de Gestión E-Commerce

## Información del Proyecto

**Autor:** Christian David Muñoz Tonguino
**Fecha:** 23/06/2026

## Objetivo del Programa

Desarrollar una aplicación web de comercio electrónico utilizando GoLang y MySQL que permita la gestión de productos, clientes, compras y administración del sistema mediante una arquitectura modular, aplicando los conocimientos adquiridos durante la asignatura de Programación Orientada a Objetos.

## Descripción General

El sistema permite administrar una tienda virtual especializada en productos tecnológicos orientados al gaming y alto rendimiento, incluyendo periféricos, monitores, sillas ergonómicas y accesorios informáticos.

La aplicación integra funcionalidades para clientes y administradores, permitiendo gestionar productos, realizar compras, consultar historiales y consumir servicios web basados en JSON.

---

## Funcionalidades Principales

### Gestión de Productos

* Registro de nuevos productos.
* Edición de productos existentes.
* Eliminación de productos.
* Visualización de catálogo.
* Búsqueda de productos por nombre.
* Ordenamiento por precio.
* Paginación de resultados.

### Gestión de Clientes

* Registro de usuarios.
* Inicio de sesión de clientes.
* Cierre de sesión.
* Validación de credenciales.

### Gestión de Compras

* Carrito de compras dinámico.
* Agregar productos al carrito.
* Eliminar productos del carrito.
* Resumen de compra.
* Cálculo automático de subtotal.
* Cálculo automático de IVA (15%).
* Cálculo de total a pagar.
* Registro de compras en MySQL.

### Historial de Compras

#### Cliente

* Consulta de compras realizadas.
* Visualización del detalle de cada compra.

#### Administrador

* Consulta de clientes que han realizado compras.
* Búsqueda de clientes por nombre.
* Visualización del historial individual de cada cliente.
* Filtro de compras por rango de fechas.

### Panel Administrativo

* Acceso protegido mediante autenticación.
* Gestión completa de productos.
* Acceso a historial general de compras.

### Servicios Web JSON

El sistema incorpora servicios web que permiten intercambiar información mediante serialización JSON.

Ejemplos:

* Listado de productos.
* Consulta de producto por ID.
* Consulta del carrito.
* Consulta del total del carrito.
* Búsqueda de productos.
* Productos ordenados por precio.
* Información de clientes.
* Historial de compras.

---

## Tecnologías Utilizadas

* GoLang
* MySQL
* HTML5
* CSS3
* JSON
* Git
* GitHub
* Visual Studio Code

---

## Librerías Utilizadas

* net/http
* html/template
* database/sql
* encoding/json
* strconv
* github.com/go-sql-driver/mysql

---

## Arquitectura del Proyecto

### Módulos Principales

* productos → gestión de productos
* carrito → gestión del carrito
* usuarios → gestión de clientes
* compras → gestión de compras
* database → conexión MySQL
* templates → vistas HTML
* static → archivos CSS e imágenes

---

## Conceptos Aplicados

### Unidad 1

* Programación funcional.
* Funciones y modularización.

### Unidad 2

* Estructuras y métodos.
* Encapsulación.
* Interfaces.

### Unidad 3

* Manejo de bases de datos MySQL.
* Persistencia de información.

### Unidad 4

* Servicios Web.
* Serialización JSON.
* Aplicaciones Web con GoLang.

---

## Ejecución del Proyecto

1. Clonar el repositorio.
2. Abrir el proyecto en Visual Studio Code.
3. Configurar la base de datos MySQL.
4. Ejecutar:

```bash
go run cmd/main.go
```

5. Abrir el navegador en:

```text
http://localhost:8080
```

---

## Repositorio

El código fuente completo del proyecto se encuentra almacenado en GitHub como parte de la entrega final de la asignatura.


```bash
go run cmd/main.go

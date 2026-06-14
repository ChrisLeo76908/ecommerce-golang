# Sistema de Gestión E-Commerce

## Información del proyecto

**Autor:** Christian David Muñoz Tonguino  
**Fecha:** 23/05/2026  
**Descripción:** Sistema de gestión E-Commerce desarrollado en GoLang y MySQL, que permite administrar productos, carrito de compras y panel administrativo mediante una interfaz web dinámica.

Proyecto desarrollado en GoLang y MySQL como parte de una práctica académica orientada al desarrollo de un sistema de gestión E-Commerce utilizando programación funcional y arquitectura modular.

## Descripción

El sistema permite gestionar una tienda virtual mediante una interfaz web dinámica, ofreciendo funcionalidades básicas de administración de productos y carrito de compras.

## Funcionalidades principales

- Gestión completa de productos: agregar, visualizar, editar y eliminar.
- Conexión con base de datos MySQL para almacenar productos.
- Carrito de compras dinámico.
- Resumen final de compra agrupando productos repetidos.
- Cálculo de cantidades, subtotales y total general.
- Búsqueda de productos por nombre.
- Ordenamiento de productos por precio.
- Paginación del catálogo de productos.
- Visualización de productos en cuadrícula.
- Panel administrador protegido con inicio de sesión.
- Implementación de encapsulación mediante métodos en Producto.
- Implementación de la interfaz ProductoRepository.
- Manejo de errores y validaciones.
- Comentarios en funcionalidades complejas.

## Tecnologías utilizadas

- GoLang
- MySQL
- HTML5
- CSS3
- GitHub
- Visual Studio Code

## Librerías utilizadas

- net/http
- html/template
- database/sql
- github.com/go-sql-driver/mysql

## Arquitectura del proyecto

El proyecto está dividido en módulos independientes:

- `productos` → gestión de productos
- `carrito` → gestión del carrito de compras
- `database` → conexión con MySQL
- `templates` → vistas HTML
- `static` → archivos CSS e imágenes

## Ejecución del proyecto

1. Clonar repositorio
2. Abrir el proyecto en Visual Studio Code
3. Configurar MySQL
4. Ejecutar:

```bash
go run cmd/main.go

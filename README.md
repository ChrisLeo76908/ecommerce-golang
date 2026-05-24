# Sistema de Gestión E-Commerce

## Información del proyecto

**Autor:** Christian David Muñoz Tonguino  
**Fecha:** 23/05/2026  
**Descripción:** Sistema de gestión E-Commerce desarrollado en GoLang y MySQL, que permite administrar productos, carrito de compras y panel administrativo mediante una interfaz web dinámica.

Proyecto desarrollado en GoLang y MySQL como parte de una práctica académica orientada al desarrollo de un sistema de gestión E-Commerce utilizando programación funcional y arquitectura modular.

## Descripción

El sistema permite gestionar una tienda virtual mediante una interfaz web dinámica, ofreciendo funcionalidades básicas de administración de productos y carrito de compras.

## Funcionalidades principales

- Visualización de productos
- Carrito de compras dinámico
- Agregar productos al carrito
- Eliminar productos del carrito
- Cálculo automático del total
- Panel administrador
- Inicio de sesión para administrador
- Agregar nuevos productos
- Eliminar productos del inventario
- Conexión con base de datos MySQL

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

package main

import (
	"fmt"

	"RDSI/web"
)

// ConectarBaseDeDatos avisa cuando haya base de datos real.
// Por ahora los datos viven en memoria y se cargan desde datos.json.
func ConectarBaseDeDatos() {
	fmt.Println("Base de datos: Pendiente (se usa almacen en memoria)")
}

// main arranca el servidor web del sistema Rayitos de Sol.
func main() {
	ConectarBaseDeDatos()

	if err := web.CargarDatos("datos.json"); err != nil {
		fmt.Println("Aviso:", err)
		fmt.Println("Se inicia con el almacen vacio.")
	}

	fmt.Println("Servidor de Rayitos de Sol corriendo en http://localhost:8080")
	if err := web.Servir(); err != nil {
		fmt.Println("Error en el servidor:", err)
	}
}

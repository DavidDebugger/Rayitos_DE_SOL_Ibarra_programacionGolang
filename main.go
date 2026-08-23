// main.go
// =======
// Archivo de arranque del programa.
//
// Este archivo es el mas pequeno a proposito: solo importa el paquete "web"
// (que esta en la carpeta web/) y le dice que encienda el servidor.
// Toda la logica del sistema vive dentro del paquete web, separada por
// archivos segun su funcion (modelo, almacen, servicios, portal, etc.).
//
// Para ejecutar el programa:
//   go run .
// (el punto significa "compila y ejecuta todo el paquete actual").

package main

import "RDSI/web"

func main() {
	// Iniciamos el servidor web del centro de terapias "Rayitos de Sol".
	// El servidor quedara escuchando en el puerto 8080 hasta que lo detengamos
	// (Ctrl+C en la terminal, o bien: fuser -k 8080/tcp).
	web.IniciarServidor(":8080")
}

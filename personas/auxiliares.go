// archivo: personas/auxiliares.go
// funciones auxiliares para leer entrada del usuario

package personas

import "fmt"

// leerTexto lee una linea de texto con espacios
func leerTexto(mensaje string) string {
	fmt.Print(mensaje)
	var texto string
	fmt.Scanf("%[^\n]\n", &texto)
	return texto
}

// leerNumero lee un numero entero
func leerNumero(mensaje string) int {
	fmt.Print(mensaje)
	var num int
	fmt.Scanf("%d\n", &num)
	return num
}

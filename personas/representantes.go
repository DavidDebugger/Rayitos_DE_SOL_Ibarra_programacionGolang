package personas

import (
	"fmt"
)

// Representante guarda los datos de un representante legal
type Representante struct {
	ID            int
	Nombre        string
	Telefono      string
	Correo        string
	Direccion     string
	Relacion      string
	Observaciones string
}

// guarda todos los representantes
var ListaRepresentantes []Representante

// contador para IDs
var contadorIDRep = 1

// RegistrarRepresentante pide los datos y guarda un representante
func RegistrarRepresentante() {
	fmt.Println("____ Registrar Representante ____")

	var rep Representante
	rep.ID = contadorIDRep
	contadorIDRep++

	rep.Nombre = leerTexto("Nombre completo: ")
	rep.Telefono = leerTexto("Telefono: ")
	rep.Correo = leerTexto("Correo: ")
	rep.Direccion = leerTexto("Direccion: ")
	rep.Relacion = leerTexto("Relacion con estudiante: ")
	rep.Observaciones = leerTexto("Observaciones: ")

	ListaRepresentantes = append(ListaRepresentantes, rep)
	fmt.Printf("Representante '%s' registrado con ID %d\n", rep.Nombre, rep.ID)
}

// ListarRepresentantes muestra todos los representantes
func ListarRepresentantes() {
	fmt.Println("____ Lista de Representantes ____")

	if len(ListaRepresentantes) == 0 {
		fmt.Println("No hay representantes registrados")
		return
	}

	for _, rep := range ListaRepresentantes {
		fmt.Printf("[%d] %s - %s - %s\n", rep.ID, rep.Nombre, rep.Relacion, rep.Telefono)
	}
}

// ActualizarRepresentante modifica datos de un representante
func ActualizarRepresentante() {
	fmt.Println("____ Actualizar Representante ____")

	if len(ListaRepresentantes) == 0 {
		fmt.Println("No hay representantes registrados")
		return
	}

	id := leerNumero("ID del representante: ")

	for i := range ListaRepresentantes {
		if ListaRepresentantes[i].ID == id {
			rep := &ListaRepresentantes[i]
			fmt.Printf("Representante: %s\n", rep.Nombre)

			nuevoNombre := leerTexto("Nuevo nombre (Enter para no cambiar): ")
			if nuevoNombre != "" {
				rep.Nombre = nuevoNombre
			}

			nuevoTelefono := leerTexto("Nuevo telefono (Enter para no cambiar): ")
			if nuevoTelefono != "" {
				rep.Telefono = nuevoTelefono
			}

			nuevoCorreo := leerTexto("Nuevo correo (Enter para no cambiar): ")
			if nuevoCorreo != "" {
				rep.Correo = nuevoCorreo
			}

			nuevaDireccion := leerTexto("Nueva direccion (Enter para no cambiar): ")
			if nuevaDireccion != "" {
				rep.Direccion = nuevaDireccion
			}

			fmt.Println("Representante actualizado")
			return
		}
	}
	fmt.Printf("No se encontro representante con ID %d\n", id)
}

// EliminarRepresentante borra un representante
func EliminarRepresentante() {
	fmt.Println("____ Eliminar Representante ____")

	if len(ListaRepresentantes) == 0 {
		fmt.Println("No hay representantes registrados")
		return
	}

	id := leerNumero("ID del representante a eliminar: ")

	for i := range ListaRepresentantes {
		if ListaRepresentantes[i].ID == id {
			fmt.Printf("Representante '%s' eliminado\n", ListaRepresentantes[i].Nombre)
			ListaRepresentantes = append(ListaRepresentantes[:i], ListaRepresentantes[i+1:]...)
			return
		}
	}
	fmt.Printf("No se encontro representante con ID %d\n", id)
}

// BuscarRepresentantePorID busca un representante por ID
func BuscarRepresentantePorID(id int) *Representante {
	for i := range ListaRepresentantes {
		if ListaRepresentantes[i].ID == id {
			return &ListaRepresentantes[i]
		}
	}
	return nil
}

// DetalleRepresentante muestra todos los datos de un representante
func DetalleRepresentante() {
	fmt.Println("____ Detalle de Representante ____")

	if len(ListaRepresentantes) == 0 {
		fmt.Println("No hay representantes registrados")
		return
	}

	id := leerNumero("ID del representante: ")
	rep := BuscarRepresentantePorID(id)
	if rep == nil {
		fmt.Printf("No se encontro representante con ID %d\n", id)
		return
	}

	fmt.Printf("ID: %d\n", rep.ID)
	fmt.Printf("Nombre: %s\n", rep.Nombre)
	fmt.Printf("Telefono: %s\n", rep.Telefono)
	fmt.Printf("Correo: %s\n", rep.Correo)
	fmt.Printf("Direccion: %s\n", rep.Direccion)
	fmt.Printf("Relacion: %s\n", rep.Relacion)
	fmt.Printf("Observaciones: %s\n", rep.Observaciones)
}


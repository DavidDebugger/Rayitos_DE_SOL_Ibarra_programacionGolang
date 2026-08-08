package personas

import (
	"fmt"
)

// Terapeuta guarda los datos de un terapeuta
type Terapeuta struct {
	ID            int
	Nombre        string
	Especialidad  string
	Telefono      string
	Correo        string
	DiasTrabajo   string
	Horario       string
	Observaciones string
}

// guarda todos los terapeutas
var ListaTerapeutas []Terapeuta

// contador para IDs
var contadorIDTer = 1

// RegistrarTerapeuta pide los datos y guarda un terapeuta
func RegistrarTerapeuta() {
	fmt.Println("____ Registrar Terapeuta ____")

	var ter Terapeuta
	ter.ID = contadorIDTer
	contadorIDTer++

	ter.Nombre = leerTexto("Nombre completo: ")
	ter.Especialidad = leerTexto("Especialidad: ")
	ter.Telefono = leerTexto("Telefono: ")
	ter.Correo = leerTexto("Correo: ")
	ter.DiasTrabajo = leerTexto("Dias de trabajo: ")
	ter.Horario = leerTexto("Horario: ")
	ter.Observaciones = leerTexto("Observaciones: ")

	ListaTerapeutas = append(ListaTerapeutas, ter)
	fmt.Printf("Terapeuta '%s' registrado con ID %d\n", ter.Nombre, ter.ID)
}

// ListarTerapeutas muestra todos los terapeutas
func ListarTerapeutas() {
	fmt.Println("____ Lista de Terapeutas ____")

	if len(ListaTerapeutas) == 0 {
		fmt.Println("No hay terapeutas registrados")
		return
	}

	for _, ter := range ListaTerapeutas {
		fmt.Printf("[%d] %s - %s - %s\n", ter.ID, ter.Nombre, ter.Especialidad, ter.Telefono)
	}
}

// ActualizarTerapeuta modifica datos de un terapeuta
func ActualizarTerapeuta() {
	fmt.Println("____ Actualizar Terapeuta ____")

	if len(ListaTerapeutas) == 0 {
		fmt.Println("No hay terapeutas registrados")
		return
	}

	id := leerNumero("ID del terapeuta: ")

	for i := range ListaTerapeutas {
		if ListaTerapeutas[i].ID == id {
			ter := &ListaTerapeutas[i]
			fmt.Printf("Terapeuta: %s\n", ter.Nombre)

			nuevoNombre := leerTexto("Nuevo nombre (Enter para no cambiar): ")
			if nuevoNombre != "" {
				ter.Nombre = nuevoNombre
			}

			nuevaEspecialidad := leerTexto("Nueva especialidad (Enter para no cambiar): ")
			if nuevaEspecialidad != "" {
				ter.Especialidad = nuevaEspecialidad
			}

			nuevoTelefono := leerTexto("Nuevo telefono (Enter para no cambiar): ")
			if nuevoTelefono != "" {
				ter.Telefono = nuevoTelefono
			}

			nuevoCorreo := leerTexto("Nuevo correo (Enter para no cambiar): ")
			if nuevoCorreo != "" {
				ter.Correo = nuevoCorreo
			}

			fmt.Println("Terapeuta actualizado")
			return
		}
	}
	fmt.Printf("No se encontro terapeuta con ID %d\n", id)
}

// EliminarTerapeuta borra un terapeuta
func EliminarTerapeuta() {
	fmt.Println("____ Eliminar Terapeuta ____")

	if len(ListaTerapeutas) == 0 {
		fmt.Println("No hay terapeutas registrados")
		return
	}

	id := leerNumero("ID del terapeuta a eliminar: ")

	for i := range ListaTerapeutas {
		if ListaTerapeutas[i].ID == id {
			fmt.Printf("Terapeuta '%s' eliminado\n", ListaTerapeutas[i].Nombre)
			ListaTerapeutas = append(ListaTerapeutas[:i], ListaTerapeutas[i+1:]...)
			return
		}
	}
	fmt.Printf("No se encontro terapeuta con ID %d\n", id)
}

// BuscarTerapeutaPorID busca un terapeuta por ID
func BuscarTerapeutaPorID(id int) *Terapeuta {
	for i := range ListaTerapeutas {
		if ListaTerapeutas[i].ID == id {
			return &ListaTerapeutas[i]
		}
	}
	return nil
}

// DetalleTerapeuta muestra todos los datos de un terapeuta
func DetalleTerapeuta() {
	fmt.Println("____ Detalle de Terapeuta ____")

	if len(ListaTerapeutas) == 0 {
		fmt.Println("No hay terapeutas registrados")
		return
	}

	id := leerNumero("ID del terapeuta: ")
	ter := BuscarTerapeutaPorID(id)
	if ter == nil {
		fmt.Printf("No se encontro terapeuta con ID %d\n", id)
		return
	}

	fmt.Printf("ID: %d\n", ter.ID)
	fmt.Printf("Nombre: %s\n", ter.Nombre)
	fmt.Printf("Especialidad: %s\n", ter.Especialidad)
	fmt.Printf("Telefono: %s\n", ter.Telefono)
	fmt.Printf("Correo: %s\n", ter.Correo)
	fmt.Printf("Dias de trabajo: %s\n", ter.DiasTrabajo)
	fmt.Printf("Horario: %s\n", ter.Horario)
	fmt.Printf("Observaciones: %s\n", ter.Observaciones)
}

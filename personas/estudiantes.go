// archivo: personas/estudiantes.go
// funciones para registrar, listar, actualizar y eliminar estudiantes

package personas

import (
	"fmt"
)

// Estudiante guarda los datos de un estudiante
type Estudiante struct {
	ID            int
	Nombre        string
	Edad          int
	Grado         string
	Representante string
	Terapeuta     string
	PlanTerapia   string
	FechaIngreso  string
	Observaciones string
}

// guarda todos los estudiantes
var ListaEstudiantes []Estudiante

// contador para IDs
var contadorIDEst = 1

// RegistrarEstudiante pide los datos y guarda un nuevo estudiante
func RegistrarEstudiante() {
	fmt.Println("____ Registrar Estudiante ____")

	var est Estudiante
	est.ID = contadorIDEst
	contadorIDEst++

	est.Nombre = leerTexto("Nombre completo: ")
	est.Edad = leerNumero("Edad: ")
	est.Grado = leerTexto("Grado escolar: ")
	est.Representante = leerTexto("Representante: ")
	est.Terapeuta = leerTexto("Terapeuta: ")
	est.PlanTerapia = leerTexto("Plan de terapia: ")
	est.FechaIngreso = leerTexto("Fecha de ingreso: ")
	est.Observaciones = leerTexto("Observaciones: ")

	ListaEstudiantes = append(ListaEstudiantes, est)
	fmt.Printf("Estudiante '%s' registrado con ID %d\n", est.Nombre, est.ID)
}

// ListarEstudiantes muestra todos los estudiantes
func ListarEstudiantes() {
	fmt.Println("____ Lista de Estudiantes ____")

	if len(ListaEstudiantes) == 0 {
		fmt.Println("No hay estudiantes registrados")
		return
	}

	for _, est := range ListaEstudiantes {
		fmt.Printf("[%d] %s - %d anios - %s\n", est.ID, est.Nombre, est.Edad, est.Grado)
	}
}

// BuscarEstudiantePorID busca un estudiante por su ID
func BuscarEstudiantePorID(id int) *Estudiante {
	for i := range ListaEstudiantes {
		if ListaEstudiantes[i].ID == id {
			return &ListaEstudiantes[i]
		}
	}
	return nil
}

// ActualizarEstudiante modifica los datos de un estudiante
func ActualizarEstudiante() {
	fmt.Println("____ Actualizar Estudiante ____")

	if len(ListaEstudiantes) == 0 {
		fmt.Println("No hay estudiantes registrados")
		return
	}

	id := leerNumero("ID del estudiante: ")
	est := BuscarEstudiantePorID(id)
	if est == nil {
		fmt.Printf("No se encontro estudiante con ID %d\n", id)
		return
	}

	fmt.Printf("Estudiante: %s\n", est.Nombre)
	nuevoNombre := leerTexto("Nuevo nombre (Enter para no cambiar): ")
	if nuevoNombre != "" {
		est.Nombre = nuevoNombre
	}

	nuevaEdad := leerNumero("Nueva edad (0 para no cambiar): ")
	if nuevaEdad > 0 {
		est.Edad = nuevaEdad
	}

	nuevoGrado := leerTexto("Nuevo grado (Enter para no cambiar): ")
	if nuevoGrado != "" {
		est.Grado = nuevoGrado
	}

	nuevoTerapeuta := leerTexto("Nuevo terapeuta (Enter para no cambiar): ")
	if nuevoTerapeuta != "" {
		est.Terapeuta = nuevoTerapeuta
	}

	fmt.Println("Estudiante actualizado")
}

// EliminarEstudiante borra un estudiante
func EliminarEstudiante() {
	fmt.Println("____ Eliminar Estudiante ____")

	if len(ListaEstudiantes) == 0 {
		fmt.Println("No hay estudiantes registrados")
		return
	}

	id := leerNumero("ID del estudiante a eliminar: ")

	for i := range ListaEstudiantes {
		if ListaEstudiantes[i].ID == id {
			fmt.Printf("Estudiante '%s' eliminado\n", ListaEstudiantes[i].Nombre)
			ListaEstudiantes = append(ListaEstudiantes[:i], ListaEstudiantes[i+1:]...)
			return
		}
	}
	fmt.Printf("No se encontro estudiante con ID %d\n", id)
}

// DetalleEstudiante muestra todos los datos de un estudiante
func DetalleEstudiante() {
	fmt.Println("____ Detalle de Estudiante ____")

	if len(ListaEstudiantes) == 0 {
		fmt.Println("No hay estudiantes registrados")
		return
	}

	id := leerNumero("ID del estudiante: ")
	est := BuscarEstudiantePorID(id)
	if est == nil {
		fmt.Printf("No se encontro estudiante con ID %d\n", id)
		return
	}

	fmt.Printf("ID: %d\n", est.ID)
	fmt.Printf("Nombre: %s\n", est.Nombre)
	fmt.Printf("Edad: %d\n", est.Edad)
	fmt.Printf("Grado: %s\n", est.Grado)
	fmt.Printf("Representante: %s\n", est.Representante)
	fmt.Printf("Terapeuta: %s\n", est.Terapeuta)
	fmt.Printf("Plan: %s\n", est.PlanTerapia)
	fmt.Printf("Ingreso: %s\n", est.FechaIngreso)
	fmt.Printf("Observaciones: %s\n", est.Observaciones)
}

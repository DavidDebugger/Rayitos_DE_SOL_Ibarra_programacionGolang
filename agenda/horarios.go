package agenda

import (
	"fmt"
)

// Cita guarda los datos de una cita programada
type Cita struct {
	ID          int
	Estudiante  string
	Terapeuta   string
	Fecha       string
	Hora        string
	Duracion    string
	Descripcion string
	Plan        string
	Estado      string
}

// guarda todas las citas
var ListaCitas []Cita

// contador para IDs
var contadorIDCita = 1

// RegistrarCita pide datos y guarda una cita
func RegistrarCita() {
	fmt.Println("____ Registrar Cita ____")

	var cita Cita
	cita.ID = contadorIDCita
	contadorIDCita++

	cita.Estudiante = leerTexto("Estudiante: ")
	cita.Terapeuta = leerTexto("Terapeuta: ")
	cita.Fecha = leerTexto("Fecha: ")
	cita.Hora = leerTexto("Hora: ")
	cita.Duracion = leerTexto("Duracion: ")
	cita.Descripcion = leerTexto("Descripcion: ")
	cita.Plan = leerTexto("Plan: ")
	cita.Estado = "activa"

	ListaCitas = append(ListaCitas, cita)
	fmt.Printf("Cita #%d registrada\n", cita.ID)
}

// ListarCitas muestra todas las citas
func ListarCitas() {
	fmt.Println("____ Lista de Citas ____")

	if len(ListaCitas) == 0 {
		fmt.Println("No hay citas registradas")
		return
	}

	for _, cita := range ListaCitas {
		fmt.Printf("[%d] %s - %s - %s %s - %s\n",
			cita.ID, cita.Estudiante, cita.Terapeuta, cita.Fecha, cita.Hora, cita.Estado)
	}
}

// ListarCitasPorFecha muestra citas de una fecha
func ListarCitasPorFecha() {
	fmt.Println("____ Citas por Fecha ____")

	if len(ListaCitas) == 0 {
		fmt.Println("No hay citas registradas")
		return
	}

	fecha := leerTexto("Fecha a consultar: ")
	var encontradas int

	for _, cita := range ListaCitas {
		if cita.Fecha == fecha {
			fmt.Printf("[%d] %s - %s - %s - %s\n", cita.ID, cita.Estudiante, cita.Terapeuta, cita.Hora, cita.Descripcion)
			encontradas++
		}
	}

	if encontradas == 0 {
		fmt.Printf("No hay citas para la fecha %s\n", fecha)
	}
}

// ActualizarCita cambia datos de una cita
func ActualizarCita() {
	fmt.Println("____ Actualizar Cita ____")

	if len(ListaCitas) == 0 {
		fmt.Println("No hay citas registradas")
		return
	}

	id := leerNumero("ID de la cita: ")

	for i := range ListaCitas {
		if ListaCitas[i].ID == id {
			cita := &ListaCitas[i]
			fmt.Printf("Cita: %s con %s\n", cita.Estudiante, cita.Terapeuta)

			nuevaFecha := leerTexto("Nueva fecha (Enter para no cambiar): ")
			if nuevaFecha != "" {
				cita.Fecha = nuevaFecha
			}

			nuevaHora := leerTexto("Nueva hora (Enter para no cambiar): ")
			if nuevaHora != "" {
				cita.Hora = nuevaHora
			}

			nuevoEstado := leerTexto("Nuevo estado (Enter para no cambiar): ")
			if nuevoEstado != "" {
				cita.Estado = nuevoEstado
			}

			fmt.Println("Cita actualizada")
			return
		}
	}
	fmt.Printf("No se encontro cita con ID %d\n", id)
}

// EliminarCita borra una cita
func EliminarCita() {
	fmt.Println("____ Eliminar Cita ____")

	if len(ListaCitas) == 0 {
		fmt.Println("No hay citas registradas")
		return
	}

	id := leerNumero("ID de la cita a eliminar: ")

	for i := range ListaCitas {
		if ListaCitas[i].ID == id {
			fmt.Printf("Cita #%d eliminada\n", ListaCitas[i].ID)
			ListaCitas = append(ListaCitas[:i], ListaCitas[i+1:]...)
			return
		}
	}
	fmt.Printf("No se encontro cita con ID %d\n", id)
}

// BuscarCitaPorID busca una cita por ID
func BuscarCitaPorID(id int) *Cita {
	for i := range ListaCitas {
		if ListaCitas[i].ID == id {
			return &ListaCitas[i]
		}
	}
	return nil
}

// DetalleCita muestra todos los datos de una cita
func DetalleCita() {
	fmt.Println("____ Detalle de Cita ____")

	if len(ListaCitas) == 0 {
		fmt.Println("No hay citas registradas")
		return
	}

	id := leerNumero("ID de la cita: ")
	cita := BuscarCitaPorID(id)
	if cita == nil {
		fmt.Printf("No se encontro cita con ID %d\n", id)
		return
	}

	fmt.Printf("ID: %d\n", cita.ID)
	fmt.Printf("Estudiante: %s\n", cita.Estudiante)
	fmt.Printf("Terapeuta: %s\n", cita.Terapeuta)
	fmt.Printf("Fecha: %s\n", cita.Fecha)
	fmt.Printf("Hora: %s\n", cita.Hora)
	fmt.Printf("Duracion: %s\n", cita.Duracion)
	fmt.Printf("Descripcion: %s\n", cita.Descripcion)
	fmt.Printf("Plan: %s\n", cita.Plan)
	fmt.Printf("Estado: %s\n", cita.Estado)
}

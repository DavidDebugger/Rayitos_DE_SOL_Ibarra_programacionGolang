package agenda

import (
	"fmt"
)

// Asistencia guarda el registro de asistencia de un estudiante
type Asistencia struct {
	ID          int
	Estudiante  string
	Fecha       string
	HoraEntrada string
	HoraSalida  string
	Estado      string
	Observaciones string
}

// guarda todas las asistencias
var ListaAsistencias []Asistencia

// contador para IDs
var contadorIDAsis = 1

// RegistrarAsistencia registra la entrada de un estudiante
func RegistrarAsistencia() {
	fmt.Println("____ Registrar Asistencia ____")

	var asis Asistencia
	asis.ID = contadorIDAsis
	contadorIDAsis++

	asis.Estudiante = leerTexto("Estudiante: ")
	asis.Fecha = leerTexto("Fecha: ")
	asis.HoraEntrada = leerTexto("Hora de entrada: ")
	asis.HoraSalida = leerTexto("Hora de salida: ")
	asis.Estado = leerTexto("Estado (presente/tarde/ausente): ")
	asis.Observaciones = leerTexto("Observaciones: ")

	ListaAsistencias = append(ListaAsistencias, asis)
	fmt.Printf("Asistencia registrada para '%s'\n", asis.Estudiante)
}

// ListarAsistencias muestra todas las asistencias
func ListarAsistencias() {
	fmt.Println("____ Lista de Asistencias ____")

	if len(ListaAsistencias) == 0 {
		fmt.Println("No hay asistencias registradas")
		return
	}

	for _, asis := range ListaAsistencias {
		fmt.Printf("[%d] %s - %s - %s - Entrada: %s\n",
			asis.ID, asis.Estudiante, asis.Fecha, asis.Estado, asis.HoraEntrada)
	}
}

// ListarAsistenciasPorEstudiante muestra asistencias de un estudiante
func ListarAsistenciasPorEstudiante() {
	fmt.Println("____ Asistencias por Estudiante ____")

	if len(ListaAsistencias) == 0 {
		fmt.Println("No hay asistencias registradas")
		return
	}

	nombre := leerTexto("Nombre del estudiante: ")
	var encontradas int

	for _, asis := range ListaAsistencias {
		if asis.Estudiante == nombre {
			fmt.Printf("[%d] %s - %s - Entrada: %s - Salida: %s\n",
				asis.ID, asis.Fecha, asis.Estado, asis.HoraEntrada, asis.HoraSalida)
			encontradas++
		}
	}

	if encontradas == 0 {
		fmt.Printf("No hay asistencias para '%s'\n", nombre)
	}
}

// ActualizarAsistencia modifica una asistencia
func ActualizarAsistencia() {
	fmt.Println("____ Actualizar Asistencia ____")

	if len(ListaAsistencias) == 0 {
		fmt.Println("No hay asistencias registradas")
		return
	}

	id := leerNumero("ID de la asistencia: ")

	for i := range ListaAsistencias {
		if ListaAsistencias[i].ID == id {
			asis := &ListaAsistencias[i]
			fmt.Printf("Asistencia: %s - %s\n", asis.Estudiante, asis.Fecha)

			nuevaFecha := leerTexto("Nueva fecha (Enter para no cambiar): ")
			if nuevaFecha != "" {
				asis.Fecha = nuevaFecha
			}

			nuevaEntrada := leerTexto("Nueva hora entrada (Enter para no cambiar): ")
			if nuevaEntrada != "" {
				asis.HoraEntrada = nuevaEntrada
			}

			nuevaSalida := leerTexto("Nueva hora salida (Enter para no cambiar): ")
			if nuevaSalida != "" {
				asis.HoraSalida = nuevaSalida
			}

			nuevoEstado := leerTexto("Nuevo estado (Enter para no cambiar): ")
			if nuevoEstado != "" {
				asis.Estado = nuevoEstado
			}

			fmt.Println("Asistencia actualizada")
			return
		}
	}
	fmt.Printf("No se encontro asistencia con ID %d\n", id)
}

// EliminarAsistencia borra un registro de asistencia
func EliminarAsistencia() {
	fmt.Println("____ Eliminar Asistencia ____")

	if len(ListaAsistencias) == 0 {
		fmt.Println("No hay asistencias registradas")
		return
	}

	id := leerNumero("ID de la asistencia a eliminar: ")

	for i := range ListaAsistencias {
		if ListaAsistencias[i].ID == id {
			fmt.Printf("Asistencia #%d eliminada\n", ListaAsistencias[i].ID)
			ListaAsistencias = append(ListaAsistencias[:i], ListaAsistencias[i+1:]...)
			return
		}
	}
	fmt.Printf("No se encontro asistencia con ID %d\n", id)
}

// DetalleAsistencia muestra todos los datos de una asistencia
func DetalleAsistencia() {
	fmt.Println("____ Detalle de Asistencia ____")

	if len(ListaAsistencias) == 0 {
		fmt.Println("No hay asistencias registradas")
		return
	}

	id := leerNumero("ID de la asistencia: ")

	for _, asis := range ListaAsistencias {
		if asis.ID == id {
			fmt.Printf("ID: %d\n", asis.ID)
			fmt.Printf("Estudiante: %s\n", asis.Estudiante)
			fmt.Printf("Fecha: %s\n", asis.Fecha)
			fmt.Printf("Hora entrada: %s\n", asis.HoraEntrada)
			fmt.Printf("Hora salida: %s\n", asis.HoraSalida)
			fmt.Printf("Estado: %s\n", asis.Estado)
			fmt.Printf("Observaciones: %s\n", asis.Observaciones)
			return
		}
	}
	fmt.Printf("No se encontro asistencia con ID %d\n", id)
}

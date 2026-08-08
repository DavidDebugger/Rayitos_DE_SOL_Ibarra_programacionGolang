package agenda

import (
	"fmt"
	"strings"
)

// Asistencia guarda el registro de asistencia de un estudiante
type Asistencia struct {
	ID            int
	Estudiante    string
	Fecha         string
	HoraEntrada   string
	HoraSalida    string
	Estado        string
	Observaciones string
}

// guarda todas las asistencias
var ListaAsistencias []Asistencia

// contador para IDs
var contadorIDAsis = 1

// RegistrarAsistencia registra la entrada de un estudiante
func RegistrarAsistencia() error {
	fmt.Println("____ Registrar Asistencia ____")

	var asis Asistencia

	asis.Estudiante = leerTexto("Estudiante (nombre): ")
	if strings.TrimSpace(asis.Estudiante) == "" {
		return fmt.Errorf("%w: el nombre del estudiante no puede estar vacio", ErrEntradaInvalida)
	}

	asis.Fecha = leerTexto("Fecha (ej: 2026-08-10): ")
	asis.HoraEntrada = leerTexto("Hora de entrada (ej: 08:00): ")
	asis.HoraSalida = leerTexto("Hora de salida (ej: 09:00): ")
	asis.Estado = leerTexto("Estado (presente/tarde/ausente): ")
	asis.Observaciones = leerTexto("Observaciones (opcional): ")

	asis.ID = contadorIDAsis
	contadorIDAsis++
	ListaAsistencias = append(ListaAsistencias, asis)

	fmt.Printf("Asistencia registrada para '%s'\n", asis.Estudiante)
	fmt.Println("Puede verla en: Consultar > 7.Asistencias")
	return nil
}

// ListarAsistencias muestra todas las asistencias
func ListarAsistencias() {
	fmt.Println("")
	fmt.Println("____ Lista de Asistencias ____")

	if len(ListaAsistencias) == 0 {
		fmt.Println("")
		fmt.Println("Aun no hay asistencias registradas")
		return
	}

	fmt.Println("")
	fmt.Printf("Total de asistencias registradas: %d\n", len(ListaAsistencias))
	fmt.Println("")

	for _, asis := range ListaAsistencias {
		fmt.Printf("[%d] %s - %s\n", asis.ID, asis.Estudiante, asis.Fecha)
		fmt.Printf("    Entrada: %s | Salida: %s\n", asis.HoraEntrada, asis.HoraSalida)
		fmt.Printf("    Estado: %s\n", asis.Estado)
		fmt.Println("")
	}
}

// ListarAsistenciasPorEstudiante muestra las asistencias de un estudiante
func ListarAsistenciasPorEstudiante() {
	fmt.Println("")
	fmt.Println("____ Asistencias por Estudiante ____")

	if len(ListaAsistencias) == 0 {
		fmt.Println("")
		fmt.Println("Aun no hay asistencias registradas")
		return
	}

	nombre := leerTexto("Ingrese el nombre del estudiante a consultar (ej: primer nombre): ")
	if strings.TrimSpace(nombre) == "" {
		fmt.Println("Error:", fmt.Errorf("%w: el nombre no puede estar vacio", ErrEntradaInvalida))
		return
	}

	var encontradas int

	fmt.Println("")
	fmt.Printf("Asistencias de '%s':\n", nombre)

	for _, asis := range ListaAsistencias {
		if asis.Estudiante == nombre {
			fmt.Printf("[%d] %s - Entrada: %s - Salida: %s\n",
				asis.ID, asis.Fecha, asis.HoraEntrada, asis.HoraSalida)
			fmt.Printf("    Estado: %s\n", asis.Estado)
			encontradas++
			fmt.Println("")
		}
	}

	if encontradas == 0 {
		fmt.Println("")
		fmt.Printf("No se encontraron asistencias para '%s'.\n", nombre)
		fmt.Println("Revise el nombre escrito (una sola palabra) e intente nuevamente.")
	}
}

// ActualizarAsistencia modifica una asistencia
func ActualizarAsistencia() error {
	fmt.Println("---- Actualizar Asistencia ----")

	id, err := leerNumero("ID de la asistencia: ")
	if err != nil {
		return err
	}

	for i := range ListaAsistencias {
		if ListaAsistencias[i].ID == id {
			asis := &ListaAsistencias[i]
			fmt.Printf("Asistencia: %s - %s\n", asis.Estudiante, asis.Fecha)

			nuevaFecha := leerTexto("Nueva fecha (Enter para no cambiar): ")
			if strings.TrimSpace(nuevaFecha) != "" {
				asis.Fecha = nuevaFecha
			}

			nuevaEntrada := leerTexto("Nueva hora entrada (Enter para no cambiar): ")
			if strings.TrimSpace(nuevaEntrada) != "" {
				asis.HoraEntrada = nuevaEntrada
			}

			nuevaSalida := leerTexto("Nueva hora salida (Enter para no cambiar): ")
			if strings.TrimSpace(nuevaSalida) != "" {
				asis.HoraSalida = nuevaSalida
			}

			nuevoEstado := leerTexto("Nuevo estado (Enter para no cambiar): ")
			if strings.TrimSpace(nuevoEstado) != "" {
				asis.Estado = nuevoEstado
			}

			fmt.Println("Asistencia actualizada")
			return nil
		}
	}
	return fmt.Errorf("%w: asistencia con ID %d", ErrNoEncontrado, id)
}

// EliminarAsistencia borra un registro de asistencia
func EliminarAsistencia() error {
	fmt.Println("---- Eliminar Asistencia ----")

	id, err := leerNumero("ID de la asistencia a eliminar: ")
	if err != nil {
		return err
	}

	for i := range ListaAsistencias {
		if ListaAsistencias[i].ID == id {
			fmt.Printf("Asistencia #%d eliminada\n", ListaAsistencias[i].ID)
			ListaAsistencias = append(ListaAsistencias[:i], ListaAsistencias[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("%w: asistencia con ID %d", ErrNoEncontrado, id)
}

// DetalleAsistencia muestra todos los datos de una asistencia
func DetalleAsistencia() error {
	fmt.Println("")
	fmt.Println("____ Detalle de Asistencia ____")

	id, err := leerNumero("ID de la asistencia a consultar: ")
	if err != nil {
		return err
	}

	for _, asis := range ListaAsistencias {
		if asis.ID == id {
			fmt.Println("")
			fmt.Printf("Datos de la asistencia ID %d:\n", asis.ID)
			fmt.Println("--------------------------------------------------")
			fmt.Printf("Estudiante: %s\n", asis.Estudiante)
			fmt.Printf("Fecha: %s\n", asis.Fecha)
			fmt.Printf("Hora de entrada: %s\n", asis.HoraEntrada)
			fmt.Printf("Hora de salida: %s\n", asis.HoraSalida)
			fmt.Printf("Estado: %s\n", asis.Estado)
			fmt.Printf("Observaciones: %s\n", asis.Observaciones)
			return nil
		}
	}
	return fmt.Errorf("%w: asistencia con ID %d", ErrNoEncontrado, id)
}

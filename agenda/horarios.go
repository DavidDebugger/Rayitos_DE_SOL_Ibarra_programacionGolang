package agenda

import (
	"fmt"
	"strings"
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
func RegistrarCita() error {
	fmt.Println("____ Registrar Cita ____")

	var cita Cita

	cita.Estudiante = leerTexto("Estudiante (nombre): ")
	if strings.TrimSpace(cita.Estudiante) == "" {
		return fmt.Errorf("%w: el nombre del estudiante no puede estar vacio", ErrEntradaInvalida)
	}

	cita.Terapeuta = leerTexto("Terapeuta (nombre): ")
	cita.Fecha = leerTexto("Fecha (ej: 2026-08-10): ")
	cita.Hora = leerTexto("Hora (ej: 10:00): ")
	cita.Duracion = leerTexto("Duracion en minutos (ej: 45): ")
	cita.Descripcion = leerTexto("Descripcion de la cita: ")
	cita.Plan = leerTexto("Plan de terapia: ")
	cita.Estado = "activa"

	cita.ID = contadorIDCita
	contadorIDCita++
	ListaCitas = append(ListaCitas, cita)

	fmt.Printf("Cita #%d registrada\n", cita.ID)
	fmt.Println("Puede verla en: Consultar > 5.Citas")
	return nil
}

// ListarCitas muestra todas las citas
func ListarCitas() {
	fmt.Println("")
	fmt.Println("____ Lista de Citas ____")

	if len(ListaCitas) == 0 {
		fmt.Println("")
		fmt.Println("Aun no hay citas registradas")
		return
	}

	fmt.Println("")
	fmt.Printf("Total de citas registradas: %d\n", len(ListaCitas))
	fmt.Println("")

	for _, cita := range ListaCitas {
		fmt.Printf("[%d] Cita del %s a las %s\n", cita.ID, cita.Fecha, cita.Hora)
		fmt.Printf("    Estudiante: %s\n", cita.Estudiante)
		fmt.Printf("    Terapeuta: %s\n", cita.Terapeuta)
		fmt.Printf("    Duracion: %s minutos\n", cita.Duracion)
		fmt.Printf("    Descripcion: %s\n", cita.Descripcion)
		fmt.Printf("    Plan de terapia: %s\n", cita.Plan)
		fmt.Printf("    Estado: %s\n", cita.Estado)
		fmt.Println("")
	}
}

// ListarCitasPorFecha muestra citas de una fecha
func ListarCitasPorFecha() {
	fmt.Println("")
	fmt.Println("____ Citas por Fecha ____")

	if len(ListaCitas) == 0 {
		fmt.Println("")
		fmt.Println("Aun no hay citas registradas")
		return
	}

	fecha := leerTexto("Ingrese la fecha a consultar (ej: 2026-08-10): ")
	if strings.TrimSpace(fecha) == "" {
		fmt.Println("Error:", fmt.Errorf("%w: la fecha no puede estar vacia", ErrEntradaInvalida))
		return
	}

	var encontradas int

	fmt.Println("")
	fmt.Printf("Citas encontradas para el dia %s:\n", fecha)

	for _, cita := range ListaCitas {
		if cita.Fecha == fecha {
			fmt.Printf("[%d] %s - Terapeuta: %s - Hora: %s\n",
				cita.ID, cita.Estudiante, cita.Terapeuta, cita.Hora)
			fmt.Printf("    Estado: %s | Duracion: %s min | Plan: %s\n",
				cita.Estado, cita.Duracion, cita.Plan)
			encontradas++
			fmt.Println("")
		}
	}

	if encontradas == 0 {
		fmt.Println("")
		fmt.Printf("No se encontraron citas para la fecha %s.\n", fecha)
		fmt.Println("Revise la fecha escrita (ej: 2026-08-10) e intente nuevamente.")
	}
}

// ActualizarCita cambia datos de una cita
func ActualizarCita() error {
	fmt.Println("---- Actualizar Cita ----")

	id, err := leerNumero("ID de la cita: ")
	if err != nil {
		return err
	}

	cita, err := BuscarCitaPorID(id)
	if err != nil {
		return err
	}

	fmt.Printf("Cita: %s con %s\n", cita.Estudiante, cita.Terapeuta)

	nuevaFecha := leerTexto("Nueva fecha (Enter para no cambiar): ")
	if strings.TrimSpace(nuevaFecha) != "" {
		cita.Fecha = nuevaFecha
	}

	nuevaHora := leerTexto("Nueva hora (Enter para no cambiar): ")
	if strings.TrimSpace(nuevaHora) != "" {
		cita.Hora = nuevaHora
	}

	nuevoEstado := leerTexto("Nuevo estado (Enter para no cambiar): ")
	if strings.TrimSpace(nuevoEstado) != "" {
		cita.Estado = nuevoEstado
	}

	fmt.Println("Cita actualizada")
	return nil
}

// EliminarCita borra una cita
func EliminarCita() error {
	fmt.Println("---- Eliminar Cita ----")

	id, err := leerNumero("ID de la cita a eliminar: ")
	if err != nil {
		return err
	}

	for i := range ListaCitas {
		if ListaCitas[i].ID == id {
			fmt.Printf("Cita #%d eliminada\n", ListaCitas[i].ID)
			ListaCitas = append(ListaCitas[:i], ListaCitas[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("%w: cita con ID %d", ErrNoEncontrado, id)
}

// BuscarCitaPorID busca una cita por ID
// Regresa un error cuando no existe
func BuscarCitaPorID(id int) (*Cita, error) {
	for i := range ListaCitas {
		if ListaCitas[i].ID == id {
			return &ListaCitas[i], nil
		}
	}
	return nil, fmt.Errorf("%w: cita con ID %d", ErrNoEncontrado, id)
}

// DetalleCita muestra todos los datos de una cita
func DetalleCita() error {
	fmt.Println("")
	fmt.Println("____ Detalle de Cita ____")

	id, err := leerNumero("ID de la cita a consultar: ")
	if err != nil {
		return err
	}

	cita, err := BuscarCitaPorID(id)
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Printf("Datos de la cita ID %d:\n", cita.ID)
	fmt.Println("--------------------------------------------------")
	fmt.Printf("Estudiante: %s\n", cita.Estudiante)
	fmt.Printf("Terapeuta: %s\n", cita.Terapeuta)
	fmt.Printf("Fecha: %s\n", cita.Fecha)
	fmt.Printf("Hora: %s\n", cita.Hora)
	fmt.Printf("Duracion: %s minutos\n", cita.Duracion)
	fmt.Printf("Descripcion: %s\n", cita.Descripcion)
	fmt.Printf("Plan: %s\n", cita.Plan)
	fmt.Printf("Estado: %s\n", cita.Estado)
	return nil
}

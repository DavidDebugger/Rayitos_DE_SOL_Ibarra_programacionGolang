package personas

import (
	"fmt"
	"strings"
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
func RegistrarTerapeuta() error {
	fmt.Println("____ Registrar Terapeuta ____")

	var ter Terapeuta

	ter.Nombre = leerNombreCompleto("Nombre completo:")
	if strings.TrimSpace(ter.Nombre) == "" {
		return fmt.Errorf("%w: el nombre no puede estar vacio", ErrEntradaInvalida)
	}

	ter.Especialidad = leerTexto("Especialidad (ej: Fisioterapia): ")
	ter.Telefono = leerTexto("Telefono (ej: 0991234567): ")
	ter.Correo = leerTexto("Correo (ej: nombre@mail.com): ")
	ter.DiasTrabajo = leerTexto("Dias de trabajo (ej: Lunes, Martes): ")
	ter.Horario = leerTexto("Horario (ej: 9:00-17:00): ")
	ter.Observaciones = leerTexto("Observaciones (opcional): ")

	ter.ID = contadorIDTer
	contadorIDTer++
	ListaTerapeutas = append(ListaTerapeutas, ter)

	fmt.Printf("Terapeuta '%s' registrado con ID %d\n", ter.Nombre, ter.ID)
	fmt.Println("Puede verlo en: Consultar > 3.Terapeutas")
	return nil
}

// ListarTerapeutas muestra todos los terapeutas
func ListarTerapeutas() {
	fmt.Println("")
	fmt.Println("____ Lista de Terapeutas ____")

	if len(ListaTerapeutas) == 0 {
		fmt.Println("")
		fmt.Println("Aun no hay terapeutas registrados")
		return
	}

	fmt.Println("")
	fmt.Printf("Total de terapeutas registrados: %d\n", len(ListaTerapeutas))
	fmt.Println("")

	for _, ter := range ListaTerapeutas {
		fmt.Printf("[%d] %s\n", ter.ID, ter.Nombre)
		fmt.Printf("    Especialidad: %s\n", ter.Especialidad)
		fmt.Printf("    Telefono: %s | Correo: %s\n", ter.Telefono, ter.Correo)
		fmt.Printf("    Dias de trabajo: %s | Horario: %s\n", ter.DiasTrabajo, ter.Horario)
		fmt.Println("")
	}
}

// ActualizarTerapeuta modifica datos de un terapeuta
func ActualizarTerapeuta() error {
	fmt.Println("____ Actualizar Terapeuta ____")

	id, err := leerNumero("ID del terapeuta: ")
	if err != nil {
		return err
	}

	ter, err := BuscarTerapeutaPorID(id)
	if err != nil {
		return err
	}

	fmt.Printf("Terapeuta: %s\n", ter.Nombre)

	nuevoNombre := leerTexto("Nuevo nombre (Enter para no cambiar): ")
	if strings.TrimSpace(nuevoNombre) != "" {
		ter.Nombre = nuevoNombre
	}

	nuevaEspecialidad := leerTexto("Nueva especialidad (Enter para no cambiar): ")
	if strings.TrimSpace(nuevaEspecialidad) != "" {
		ter.Especialidad = nuevaEspecialidad
	}

	nuevoTelefono := leerTexto("Nuevo telefono (Enter para no cambiar): ")
	if strings.TrimSpace(nuevoTelefono) != "" {
		ter.Telefono = nuevoTelefono
	}

	nuevoCorreo := leerTexto("Nuevo correo (Enter para no cambiar): ")
	if strings.TrimSpace(nuevoCorreo) != "" {
		ter.Correo = nuevoCorreo
	}

	fmt.Println("Terapeuta actualizado")
	return nil
}

// EliminarTerapeuta borra un terapeuta
func EliminarTerapeuta() error {
	fmt.Println("____ Eliminar Terapeuta ____")

	id, err := leerNumero("ID del terapeuta a eliminar: ")
	if err != nil {
		return err
	}

	for i := range ListaTerapeutas {
		if ListaTerapeutas[i].ID == id {
			fmt.Printf("Terapeuta '%s' eliminado\n", ListaTerapeutas[i].Nombre)
			ListaTerapeutas = append(ListaTerapeutas[:i], ListaTerapeutas[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("%w: terapeuta con ID %d", ErrNoEncontrado, id)
}

// BuscarTerapeutaPorID busca un terapeuta por ID
// Regresa un error cuando no existe
func BuscarTerapeutaPorID(id int) (*Terapeuta, error) {
	for i := range ListaTerapeutas {
		if ListaTerapeutas[i].ID == id {
			return &ListaTerapeutas[i], nil
		}
	}
	return nil, fmt.Errorf("%w: terapeuta con ID %d", ErrNoEncontrado, id)
}

// DetalleTerapeuta muestra todos los datos de un terapeuta
func DetalleTerapeuta() error {
	fmt.Println("")
	fmt.Println("____ Detalle de Terapeuta ____")

	id, err := leerNumero("ID del terapeuta a consultar: ")
	if err != nil {
		return err
	}

	ter, err := BuscarTerapeutaPorID(id)
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Printf("Datos del terapeuta ID %d:\n", ter.ID)
	fmt.Println("--------------------------------------------------")
	fmt.Printf("Nombre: %s\n", ter.Nombre)
	fmt.Printf("Especialidad: %s\n", ter.Especialidad)
	fmt.Printf("Telefono: %s\n", ter.Telefono)
	fmt.Printf("Correo: %s\n", ter.Correo)
	fmt.Printf("Dias de trabajo: %s\n", ter.DiasTrabajo)
	fmt.Printf("Horario: %s\n", ter.Horario)
	fmt.Printf("Observaciones: %s\n", ter.Observaciones)
	return nil
}

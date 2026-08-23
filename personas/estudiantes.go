package personas

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// errores comunes del modulo de personas
var (
	ErrEntradaInvalida = errors.New("entrada no valida")
	ErrNoEncontrado    = errors.New("no se encontro el registro")
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

// leerTexto lee una palabra
func leerTexto(mensaje string) string {
	fmt.Print(mensaje)
	var texto string
	fmt.Scanln(&texto)
	return texto
}

// leerNumero lee un numero y valida que sea numero
// Regresa el numero y un error si lo escrito es invalido
func leerNumero(mensaje string) (int, error) {
	fmt.Print(mensaje)
	var num int
	_, err := fmt.Scanln(&num)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return 0, io.EOF
		}
		if strings.Contains(err.Error(), "newline") {
			return num, nil
		}
		return 0, fmt.Errorf("%w: se esperaba un numero entero", ErrEntradaInvalida)
	}
	return num, nil
}

// leerNombreCompleto pide nombre y apellido por separado y los junta
func leerNombreCompleto(mensaje string) string {
	fmt.Println(mensaje)
	nombre := leerTexto("  Nombre(s): ")
	apellido := leerTexto("  Apellido(s): ")
	return nombre + " " + apellido
}

// RegistrarEstudiante pide los datos y guarda un nuevo estudiante
func RegistrarEstudiante() error {
	fmt.Println("____ Registrar Estudiante ____")

	var est Estudiante

	est.Nombre = leerNombreCompleto("Nombre completo:")
	if strings.TrimSpace(est.Nombre) == "" {
		return fmt.Errorf("%w: el nombre no puede estar vacio", ErrEntradaInvalida)
	}

	edad, err := leerNumero("Edad: ")
	if err != nil {
		return err
	}
	est.Edad = edad

	est.Grado = leerTexto("Grado escolar (ej: 2do Grado): ")
	est.Representante = leerTexto("Representante (nombre del familiar): ")
	if strings.TrimSpace(est.Representante) == "" {
		return fmt.Errorf("%w: el representante no puede estar vacio", ErrEntradaInvalida)
	}
	est.Terapeuta = leerTexto("Terapeuta asignado (nombre): ")
	est.PlanTerapia = leerTexto("Plan de terapia asignado (nombre del plan): ")
	est.FechaIngreso = leerTexto("Fecha de ingreso (ej: 2026-08-01): ")
	est.Observaciones = leerTexto("Observaciones (opcional, ej: alergias, notas medicas): ")

	est.ID = contadorIDEst
	contadorIDEst++
	ListaEstudiantes = append(ListaEstudiantes, est)

	fmt.Printf("Estudiante '%s' registrado con ID %d\n", est.Nombre, est.ID)
	fmt.Println("Puede verlo en: Consultar > 1.Estudiantes")
	return nil
}

// ListarEstudiantes muestra todos los estudiantes
func ListarEstudiantes() {
	fmt.Println("")
	fmt.Println("____ Lista de Estudiantes ____")

	if len(ListaEstudiantes) == 0 {
		fmt.Println("")
		fmt.Println("Aun no hay estudiantes registrados")
		return
	}

	fmt.Println("")
	fmt.Printf("Total de estudiantes registrados: %d\n", len(ListaEstudiantes))
	fmt.Println("")

	for _, est := range ListaEstudiantes {
		fmt.Printf("[%d] %s\n", est.ID, est.Nombre)
		fmt.Printf("    Edad: %d años | Grado: %s\n", est.Edad, est.Grado)
		fmt.Printf("    Representante: %s\n", est.Representante)
		fmt.Printf("    Terapeuta: %s | Plan: %s\n", est.Terapeuta, est.PlanTerapia)
		fmt.Println("")
	}
}

// BuscarEstudiantePorID busca un estudiante por su ID
// Regresa un error si no existe
func BuscarEstudiantePorID(id int) (*Estudiante, error) {
	for i := range ListaEstudiantes {
		if ListaEstudiantes[i].ID == id {
			return &ListaEstudiantes[i], nil
		}
	}
	return nil, fmt.Errorf("%w: estudiante con ID %d", ErrNoEncontrado, id)
}

// ActualizarEstudiante modifica los datos de un estudiante
func ActualizarEstudiante() error {
	fmt.Println("____ Actualizar Estudiante ____")

	id, err := leerNumero("ID del estudiante: ")
	if err != nil {
		return err
	}

	est, err := BuscarEstudiantePorID(id)
	if err != nil {
		return err
	}

	fmt.Printf("Estudiante: %s\n", est.Nombre)
	nuevoNombre := leerTexto("Nuevo nombre (Enter para no cambiar): ")
	if strings.TrimSpace(nuevoNombre) != "" {
		est.Nombre = nuevoNombre
	}

	nuevaEdad, err := leerNumero("Nueva edad (0 para no cambiar): ")
	if err != nil {
		return err
	}
	if nuevaEdad > 0 {
		est.Edad = nuevaEdad
	}

	nuevoGrado := leerTexto("Nuevo grado (Enter para no cambiar): ")
	if strings.TrimSpace(nuevoGrado) != "" {
		est.Grado = nuevoGrado
	}

	nuevoTerapeuta := leerTexto("Nuevo terapeuta (Enter para no cambiar): ")
	if strings.TrimSpace(nuevoTerapeuta) != "" {
		est.Terapeuta = nuevoTerapeuta
	}

	fmt.Println("Estudiante actualizado")
	return nil
}

// EliminarEstudiante borra un estudiante
func EliminarEstudiante() error {
	fmt.Println("____ Eliminar Estudiante ____")

	id, err := leerNumero("ID del estudiante a eliminar: ")
	if err != nil {
		return err
	}

	for i := range ListaEstudiantes {
		if ListaEstudiantes[i].ID == id {
			fmt.Printf("Estudiante '%s' eliminado\n", ListaEstudiantes[i].Nombre)
			ListaEstudiantes = append(ListaEstudiantes[:i], ListaEstudiantes[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("%w: estudiante con ID %d", ErrNoEncontrado, id)
}

// DetalleEstudiante muestra todos los datos de un estudiante
func DetalleEstudiante() error {
	fmt.Println("")
	fmt.Println("____ Detalle de Estudiante ____")

	id, err := leerNumero("ID del estudiante a consultar: ")
	if err != nil {
		return err
	}

	est, err := BuscarEstudiantePorID(id)
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Printf("Datos del estudiante ID %d:\n", est.ID)
	fmt.Println("--------------------------------------------------")
	fmt.Printf("Nombre: %s\n", est.Nombre)
	fmt.Printf("Edad: %d años\n", est.Edad)
	fmt.Printf("Grado: %s\n", est.Grado)
	fmt.Printf("Representante: %s\n", est.Representante)
	fmt.Printf("Terapeuta: %s\n", est.Terapeuta)
	fmt.Printf("Plan: %s\n", est.PlanTerapia)
	fmt.Printf("Ingreso: %s\n", est.FechaIngreso)
	fmt.Printf("Observaciones: %s\n", est.Observaciones)
	return nil
}

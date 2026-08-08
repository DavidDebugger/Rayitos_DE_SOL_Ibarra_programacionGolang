package agenda

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

// errores comunes del modulo de agenda
var (
	ErrEntradaInvalida = errors.New("entrada no valida")
	ErrNoEncontrado    = errors.New("no se encontro el registro")
)

// Plan guarda los datos de un plan de terapia
type Plan struct {
	ID             int
	Nombre         string
	Descripcion    string
	Duracion       int
	Frecuencia     int
	DuracionSesion int
	Observaciones  string
}

// guarda todos los planes
var ListaPlanes []Plan

// contador para IDs
var contadorIDPlan = 1

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

// RegistrarPlan pide datos y guarda un plan
func RegistrarPlan() error {
	fmt.Println("____ Registrar Plan de Terapia ____")

	var plan Plan

	plan.Nombre = leerTexto("Nombre del plan (ej: Fisioterapia): ")
	if strings.TrimSpace(plan.Nombre) == "" {
		return fmt.Errorf("%w: el nombre del plan no puede estar vacio", ErrEntradaInvalida)
	}

	plan.Descripcion = leerTexto("Descripcion: ")

	dur, err := leerNumero("Duracion total del plan en semanas (ej: 8): ")
	if err != nil {
		return err
	}
	plan.Duracion = dur

	frec, err := leerNumero("Veces por semana (ej: 3): ")
	if err != nil {
		return err
	}
	plan.Frecuencia = frec

	sesion, err := leerNumero("Minutos por sesion (ej: 45): ")
	if err != nil {
		return err
	}
	plan.DuracionSesion = sesion

	plan.Observaciones = leerTexto("Notas adicionales (opcional): ")

	plan.ID = contadorIDPlan
	contadorIDPlan++
	ListaPlanes = append(ListaPlanes, plan)

	fmt.Printf("Plan '%s' registrado con ID %d\n", plan.Nombre, plan.ID)
	fmt.Println("Puede verlo en: Consultar > 4.Planes")
	return nil
}

// ListarPlanes muestra todos los planes
func ListarPlanes() {
	fmt.Println("")
	fmt.Println("____ Lista de Planes ____")

	if len(ListaPlanes) == 0 {
		fmt.Println("")
		fmt.Println("Aun no hay planes registrados")
		return
	}

	fmt.Println("")
	fmt.Printf("Total de planes registrados: %d\n", len(ListaPlanes))
	fmt.Println("")

	for _, plan := range ListaPlanes {
		fmt.Printf("[%d] %s\n", plan.ID, plan.Nombre)
		fmt.Printf("    Descripcion: %s\n", plan.Descripcion)
		fmt.Printf("    Frecuencia: %d vez/veces por semana\n", plan.Frecuencia)
		fmt.Printf("    Duracion del plan: %d semana(s)\n", plan.Duracion)
		fmt.Printf("    Duracion de cada sesion: %d minutos\n", plan.DuracionSesion)
		fmt.Println("")
	}
}

// ActualizarPlan modifica datos de un plan
func ActualizarPlan() error {
	fmt.Println("---- Actualizar Plan ----")

	id, err := leerNumero("ID del plan: ")
	if err != nil {
		return err
	}

	plan, err := BuscarPlanPorID(id)
	if err != nil {
		return err
	}

	fmt.Printf("Plan: %s\n", plan.Nombre)

	nuevoNombre := leerTexto("Nuevo nombre (Enter para no cambiar): ")
	if strings.TrimSpace(nuevoNombre) != "" {
		plan.Nombre = nuevoNombre
	}

	nuevaDescripcion := leerTexto("Nueva descripcion (Enter para no cambiar): ")
	if strings.TrimSpace(nuevaDescripcion) != "" {
		plan.Descripcion = nuevaDescripcion
	}

	nuevaDuracion, err := leerNumero("Nueva duracion en semanas (0 para no cambiar): ")
	if err != nil {
		return err
	}
	if nuevaDuracion > 0 {
		plan.Duracion = nuevaDuracion
	}

	nuevaFrecuencia, err := leerNumero("Nuevas veces por semana (0 para no cambiar): ")
	if err != nil {
		return err
	}
	if nuevaFrecuencia > 0 {
		plan.Frecuencia = nuevaFrecuencia
	}

	nuevaSesion, err := leerNumero("Nuevos minutos por sesion (0 para no cambiar): ")
	if err != nil {
		return err
	}
	if nuevaSesion > 0 {
		plan.DuracionSesion = nuevaSesion
	}

	fmt.Println("Plan actualizado")
	return nil
}

// EliminarPlan borra un plan
func EliminarPlan() error {
	fmt.Println("---- Eliminar Plan ----")

	id, err := leerNumero("ID del plan a eliminar: ")
	if err != nil {
		return err
	}

	for i := range ListaPlanes {
		if ListaPlanes[i].ID == id {
			fmt.Printf("Plan '%s' eliminado\n", ListaPlanes[i].Nombre)
			ListaPlanes = append(ListaPlanes[:i], ListaPlanes[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("%w: plan con ID %d", ErrNoEncontrado, id)
}

// BuscarPlanPorID busca un plan por ID
// Regresa un error cuando no existe
func BuscarPlanPorID(id int) (*Plan, error) {
	for i := range ListaPlanes {
		if ListaPlanes[i].ID == id {
			return &ListaPlanes[i], nil
		}
	}
	return nil, fmt.Errorf("%w: plan con ID %d", ErrNoEncontrado, id)
}

// DetallePlan muestra todos los datos de un plan
func DetallePlan() error {
	fmt.Println("")
	fmt.Println("____ Detalle de Plan ____")

	id, err := leerNumero("ID del plan a consultar: ")
	if err != nil {
		return err
	}

	plan, err := BuscarPlanPorID(id)
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Printf("Datos del plan de terapia ID %d:\n", plan.ID)
	fmt.Println("--------------------------------------------------")
	fmt.Printf("Nombre: %s\n", plan.Nombre)
	fmt.Printf("Descripcion: %s\n", plan.Descripcion)
	fmt.Printf("Duracion: %d semana(s)\n", plan.Duracion)
	fmt.Printf("Frecuencia: %d vez/veces por semana\n", plan.Frecuencia)
	fmt.Printf("Duracion de cada sesion: %d minutos\n", plan.DuracionSesion)
	fmt.Printf("Observaciones: %s\n", plan.Observaciones)
	return nil
}

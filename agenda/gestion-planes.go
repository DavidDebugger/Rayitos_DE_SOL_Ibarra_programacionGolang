package agenda

import (
	"fmt"
)

// Plan guarda los datos de un plan de terapia
type Plan struct {
	ID             int
	Nombre         string
	Descripcion    string
	Duracion       string
	Frecuencia     string
	DuracionSesion string
	Observaciones  string
}

// guarda todos los planes
var ListaPlanes []Plan

// contador para IDs
var contadorIDPlan = 1

// RegistrarPlan pide datos y guarda un plan
func RegistrarPlan() {
	fmt.Println("____ Registrar Plan de Terapia ____")

	var plan Plan
	plan.ID = contadorIDPlan
	contadorIDPlan++

	plan.Nombre = leerTexto("Nombre del plan: ")
	plan.Descripcion = leerTexto("Descripcion: ")
	plan.Duracion = leerTexto("Duracion total: ")
	plan.Frecuencia = leerTexto("Frecuencia por semana: ")
	plan.DuracionSesion = leerTexto("Duracion de cada sesion: ")
	plan.Observaciones = leerTexto("Observaciones: ")

	ListaPlanes = append(ListaPlanes, plan)
	fmt.Printf("Plan '%s' registrado con ID %d\n", plan.Nombre, plan.ID)
}

// ListarPlanes muestra todos los planes
func ListarPlanes() {
	fmt.Println("____ Lista de Planes ____")

	if len(ListaPlanes) == 0 {
		fmt.Println("No hay planes registrados")
		return
	}

	for _, plan := range ListaPlanes {
		fmt.Printf("[%d] %s - %s - %s\n", plan.ID, plan.Nombre, plan.Descripcion, plan.Frecuencia)
	}
}

// ActualizarPlan modifica datos de un plan
func ActualizarPlan() {
	fmt.Println("____ Actualizar Plan ____")

	if len(ListaPlanes) == 0 {
		fmt.Println("No hay planes registrados")
		return
	}

	id := leerNumero("ID del plan: ")

	for i := range ListaPlanes {
		if ListaPlanes[i].ID == id {
			plan := &ListaPlanes[i]
			fmt.Printf("Plan: %s\n", plan.Nombre)

			nuevoNombre := leerTexto("Nuevo nombre (Enter para no cambiar): ")
			if nuevoNombre != "" {
				plan.Nombre = nuevoNombre
			}

			nuevaDescripcion := leerTexto("Nueva descripcion (Enter para no cambiar): ")
			if nuevaDescripcion != "" {
				plan.Descripcion = nuevaDescripcion
			}

			nuevaDuracion := leerTexto("Nueva duracion (Enter para no cambiar): ")
			if nuevaDuracion != "" {
				plan.Duracion = nuevaDuracion
			}

			nuevaFrecuencia := leerTexto("Nueva frecuencia (Enter para no cambiar): ")
			if nuevaFrecuencia != "" {
				plan.Frecuencia = nuevaFrecuencia
			}

			fmt.Println("Plan actualizado")
			return
		}
	}
	fmt.Printf("No se encontro plan con ID %d\n", id)
}

// EliminarPlan borra un plan
func EliminarPlan() {
	fmt.Println("____ Eliminar Plan ____")

	if len(ListaPlanes) == 0 {
		fmt.Println("No hay planes registrados")
		return
	}

	id := leerNumero("ID del plan a eliminar: ")

	for i := range ListaPlanes {
		if ListaPlanes[i].ID == id {
			fmt.Printf("Plan '%s' eliminado\n", ListaPlanes[i].Nombre)
			ListaPlanes = append(ListaPlanes[:i], ListaPlanes[i+1:]...)
			return
		}
	}
	fmt.Printf("No se encontro plan con ID %d\n", id)
}

// BuscarPlanPorID busca un plan por ID
func BuscarPlanPorID(id int) *Plan {
	for i := range ListaPlanes {
		if ListaPlanes[i].ID == id {
			return &ListaPlanes[i]
		}
	}
	return nil
}

// DetallePlan muestra todos los datos de un plan
func DetallePlan() {
	fmt.Println("____ Detalle de Plan ____")

	if len(ListaPlanes) == 0 {
		fmt.Println("No hay planes registrados")
		return
	}

	id := leerNumero("ID del plan: ")
	plan := BuscarPlanPorID(id)
	if plan == nil {
		fmt.Printf("No se encontro plan con ID %d\n", id)
		return
	}

	fmt.Printf("ID: %d\n", plan.ID)
	fmt.Printf("Nombre: %s\n", plan.Nombre)
	fmt.Printf("Descripcion: %s\n", plan.Descripcion)
	fmt.Printf("Duracion: %s\n", plan.Duracion)
	fmt.Printf("Frecuencia: %s\n", plan.Frecuencia)
	fmt.Printf("Duracion sesion: %s\n", plan.DuracionSesion)
	fmt.Printf("Observaciones: %s\n", plan.Observaciones)
}

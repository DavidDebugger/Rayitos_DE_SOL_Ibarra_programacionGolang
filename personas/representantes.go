package personas

import (
	"fmt"
	"strings"
)

// Representante guarda los datos de un representante legal
type Representante struct {
	ID            int
	Nombre        string
	Telefono      string
	Correo        string
	Direccion     string
	Relacion      string
	Observaciones string
}

// guarda todos los representantes
var ListaRepresentantes []Representante

// contador para IDs
var contadorIDRep = 1

// RegistrarRepresentante pide los datos y guarda un representante
func RegistrarRepresentante() error {
	fmt.Println("____ Registrar Representante ____")

	var rep Representante
	rep.ID = contadorIDRep
	contadorIDRep++

	rep.Nombre = leerNombreCompleto("Nombre completo:")
	if strings.TrimSpace(rep.Nombre) == "" {
		return fmt.Errorf("%w: el nombre no puede estar vacio", ErrEntradaInvalida)
	}

	rep.Telefono = leerTexto("Telefono (ej: 0991234567): ")
	rep.Correo = leerTexto("Correo (ej: nombre@mail.com): ")
	rep.Direccion = leerTexto("Direccion: ")
	rep.Relacion = leerTexto("Relacion con estudiante (ej: Madre, Padre): ")
	rep.Observaciones = leerTexto("Observaciones (opcional): ")

	ListaRepresentantes = append(ListaRepresentantes, rep)
	fmt.Printf("Representante '%s' registrado con ID %d\n", rep.Nombre, rep.ID)
	fmt.Println("Puede verlo en: Consultar > 2.Representantes")
	return nil
}

// ListarRepresentantes muestra todos los representantes
func ListarRepresentantes() {
	fmt.Println("")
	fmt.Println("____ Lista de Representantes ____")

	if len(ListaRepresentantes) == 0 {
		fmt.Println("")
		fmt.Println("Aun no hay representantes registrados")
		return
	}

	fmt.Println("")
	fmt.Printf("Total de representantes registrados: %d\n", len(ListaRepresentantes))
	fmt.Println("")

	for _, rep := range ListaRepresentantes {
		fmt.Printf("[%d] %s\n", rep.ID, rep.Nombre)
		fmt.Printf("    Relacion con estudiante: %s\n", rep.Relacion)
		fmt.Printf("    Telefono: %s | Correo: %s\n", rep.Telefono, rep.Correo)
		fmt.Printf("    Direccion: %s\n", rep.Direccion)
		fmt.Println("")
	}
}

// ActualizarRepresentante modifica datos de un representante
func ActualizarRepresentante() error {
	fmt.Println("____ Actualizar Representante ____")

	id, err := leerNumero("ID del representante: ")
	if err != nil {
		return err
	}

	rep, err := BuscarRepresentantePorID(id)
	if err != nil {
		return err
	}

	fmt.Printf("Representante: %s\n", rep.Nombre)

	nuevoNombre := leerTexto("Nuevo nombre (Enter para no cambiar): ")
	if strings.TrimSpace(nuevoNombre) != "" {
		rep.Nombre = nuevoNombre
	}

	nuevoTelefono := leerTexto("Nuevo telefono (Enter para no cambiar): ")
	if strings.TrimSpace(nuevoTelefono) != "" {
		rep.Telefono = nuevoTelefono
	}

	nuevoCorreo := leerTexto("Nuevo correo (Enter para no cambiar): ")
	if strings.TrimSpace(nuevoCorreo) != "" {
		rep.Correo = nuevoCorreo
	}

	nuevaDireccion := leerTexto("Nueva direccion (Enter para no cambiar): ")
	if strings.TrimSpace(nuevaDireccion) != "" {
		rep.Direccion = nuevaDireccion
	}

	fmt.Println("Representante actualizado")
	return nil
}

// EliminarRepresentante borra un representante
func EliminarRepresentante() error {
	fmt.Println("____ Eliminar Representante ____")

	id, err := leerNumero("ID del representante a eliminar: ")
	if err != nil {
		return err
	}

	for i := range ListaRepresentantes {
		if ListaRepresentantes[i].ID == id {
			fmt.Printf("Representante '%s' eliminado\n", ListaRepresentantes[i].Nombre)
			ListaRepresentantes = append(ListaRepresentantes[:i], ListaRepresentantes[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("%w: representante con ID %d", ErrNoEncontrado, id)
}

// BuscarRepresentantePorID busca un representante por su ID
// Regresa un error cuando no existe
func BuscarRepresentantePorID(id int) (*Representante, error) {
	for i := range ListaRepresentantes {
		if ListaRepresentantes[i].ID == id {
			return &ListaRepresentantes[i], nil
		}
	}
	return nil, fmt.Errorf("%w: representante con ID %d", ErrNoEncontrado, id)
}

// DetalleRepresentante muestra todos los datos de un representante
func DetalleRepresentante() error {
	fmt.Println("")
	fmt.Println("____ Detalle de Representante ____")

	id, err := leerNumero("ID del representante a consultar: ")
	if err != nil {
		return err
	}

	rep, err := BuscarRepresentantePorID(id)
	if err != nil {
		return err
	}

	fmt.Println("")
	fmt.Printf("Datos del representante ID %d:\n", rep.ID)
	fmt.Println("--------------------------------------------------")
	fmt.Printf("Nombre: %s\n", rep.Nombre)
	fmt.Printf("Telefono: %s\n", rep.Telefono)
	fmt.Printf("Correo: %s\n", rep.Correo)
	fmt.Printf("Direccion: %s\n", rep.Direccion)
	fmt.Printf("Relacion: %s\n", rep.Relacion)
	fmt.Printf("Observaciones: %s\n", rep.Observaciones)
	return nil
}

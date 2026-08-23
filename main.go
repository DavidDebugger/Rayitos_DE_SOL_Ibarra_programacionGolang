package main

import (
	"errors"
	"fmt"
	"io"
	"strings"

	"RDSI/agenda"
	"RDSI/personas"
)

// ConectarBaseDeDatos conecta... cuando haya base de datos real
func ConectarBaseDeDatos() {
	fmt.Println("Base de datos: Pendiente (no se usa GORM por ahora)")
}

// leerNumero lee un numero entero y valida la entrada
// Regresa el numero y un error si lo escrito no es un numero
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
		return 0, fmt.Errorf("entrada invalida: se esperaba un numero entero")
	}
	return num, nil
}

// mostrarMenuPrincipal imprime el menu principal
func mostrarMenuPrincipal() {
	fmt.Println("")
	fmt.Println("")
	fmt.Println("___ Menu Principal ___")
	fmt.Println("1. Registrar")
	fmt.Println("2. Consultar")
	fmt.Println("3. Actualizar")
	fmt.Println("4. Eliminar")
	fmt.Println("5. Salir")
	fmt.Print("Opcion: ")
}

// mostrarMenuRegistrar muestra las opciones para registrar
func mostrarMenuRegistrar() {
	fmt.Println("")
	fmt.Println("")
	fmt.Println("___ Registrar ___")
	fmt.Println("1. Estudiante")
	fmt.Println("2. Representante")
	fmt.Println("3. Terapeuta")
	fmt.Println("4. Plan de terapia")
	fmt.Println("5. Cita")
	fmt.Println("6. Asistencia")
	fmt.Println("7. Regresar")
	fmt.Print("Opcion: ")
}

// mostrarMenuConsultar muestra las opciones para consultar
func mostrarMenuConsultar() {
	fmt.Println("")
	fmt.Println("")
	fmt.Println("___ Consultar ___")
	fmt.Println("1. Estudiantes")
	fmt.Println("2. Representantes")
	fmt.Println("3. Terapeutas")
	fmt.Println("4. Planes")
	fmt.Println("5. Citas")
	fmt.Println("6. Citas por fecha")
	fmt.Println("7. Asistencias")
	fmt.Println("8. Asistencias por estudiante")
	fmt.Println("9. Detalle estudiante")
	fmt.Println("10. Detalle cita")
	fmt.Println("11. Regresar")
	fmt.Print("Opcion: ")
}

// mostrarMenuActualizar muestra las opciones para actualizar
func mostrarMenuActualizar() {
	fmt.Println("")
	fmt.Println("")
	fmt.Println("___ Actualizar ___")
	fmt.Println("1. Estudiante")
	fmt.Println("2. Representante")
	fmt.Println("3. Terapeuta")
	fmt.Println("4. Plan")
	fmt.Println("5. Cita")
	fmt.Println("6. Asistencia")
	fmt.Println("7. Regresar")
	fmt.Print("Opcion: ")
}

// mostrarMenuEliminar muestra las opciones para eliminar
func mostrarMenuEliminar() {
	fmt.Println("")
	fmt.Println("")
	fmt.Println("___ Eliminar ___")
	fmt.Println("1. Estudiante")
	fmt.Println("2. Representante")
	fmt.Println("3. Terapeuta")
	fmt.Println("4. Plan")
	fmt.Println("5. Cita")
	fmt.Println("6. Asistencia")
	fmt.Println("7. Regresar")
	fmt.Print("Opcion: ")
}

// submenuRegistrar maneja el menu de registrar
func submenuRegistrar() {
	for {
		mostrarMenuRegistrar()
		opcion, err := leerNumero("")
		if errors.Is(err, io.EOF) {
			fmt.Println("")
			return
		}
		if err != nil {
			fmt.Println("Error de opcion:", err)
			continue
		}

		switch opcion {
		case 1:
			if err := personas.RegistrarEstudiante(); err != nil {
				fmt.Println("Error:", err)
			}
		case 2:
			if err := personas.RegistrarRepresentante(); err != nil {
				fmt.Println("Error:", err)
			}
		case 3:
			if err := personas.RegistrarTerapeuta(); err != nil {
				fmt.Println("Error:", err)
			}
		case 4:
			if err := agenda.RegistrarPlan(); err != nil {
				fmt.Println("Error:", err)
			}
		case 5:
			if err := agenda.RegistrarCita(); err != nil {
				fmt.Println("Error:", err)
			}
		case 6:
			if err := agenda.RegistrarAsistencia(); err != nil {
				fmt.Println("Error:", err)
			}
		case 7:
			return
		default:
			fmt.Println("Opcion no valida")
		}
	}
}

// submenuConsultar maneja el menu de consultar
func submenuConsultar() {
	for {
		mostrarMenuConsultar()
		opcion, err := leerNumero("")
		if errors.Is(err, io.EOF) {
			fmt.Println("")
			return
		}
		if err != nil {
			fmt.Println("Error de opcion:", err)
			continue
		}

		switch opcion {
		case 1:
			personas.ListarEstudiantes()
		case 2:
			personas.ListarRepresentantes()
		case 3:
			personas.ListarTerapeutas()
		case 4:
			agenda.ListarPlanes()
		case 5:
			agenda.ListarCitas()
		case 6:
			agenda.ListarCitasPorFecha()
		case 7:
			agenda.ListarAsistencias()
		case 8:
			agenda.ListarAsistenciasPorEstudiante()
		case 9:
			if err := personas.DetalleEstudiante(); err != nil {
				fmt.Println("Error:", err)
			}
		case 10:
			if err := agenda.DetalleCita(); err != nil {
				fmt.Println("Error:", err)
			}
		case 11:
			return
		default:
			fmt.Println("Opcion no valida")
		}
	}
}

// submenuActualizar maneja el menu de actualizar
func submenuActualizar() {
	for {
		mostrarMenuActualizar()
		opcion, err := leerNumero("")
		if errors.Is(err, io.EOF) {
			fmt.Println("")
			return
		}
		if err != nil {
			fmt.Println("Error de opcion:", err)
			continue
		}

		switch opcion {
		case 1:
			if err := personas.ActualizarEstudiante(); err != nil {
				fmt.Println("Error:", err)
			}
		case 2:
			if err := personas.ActualizarRepresentante(); err != nil {
				fmt.Println("Error:", err)
			}
		case 3:
			if err := personas.ActualizarTerapeuta(); err != nil {
				fmt.Println("Error:", err)
			}
		case 4:
			if err := agenda.ActualizarPlan(); err != nil {
				fmt.Println("Error:", err)
			}
		case 5:
			if err := agenda.ActualizarCita(); err != nil {
				fmt.Println("Error:", err)
			}
		case 6:
			if err := agenda.ActualizarAsistencia(); err != nil {
				fmt.Println("Error:", err)
			}
		case 7:
			return
		default:
			fmt.Println("Opcion no valida")
		}
	}
}

// submenuEliminar maneja el menu de eliminar
func submenuEliminar() {
	for {
		mostrarMenuEliminar()
		opcion, err := leerNumero("")
		if errors.Is(err, io.EOF) {
			fmt.Println("")
			return
		}
		if err != nil {
			fmt.Println("Error de opcion:", err)
			continue
		}

		switch opcion {
		case 1:
			if err := personas.EliminarEstudiante(); err != nil {
				fmt.Println("Error:", err)
			}
		case 2:
			if err := personas.EliminarRepresentante(); err != nil {
				fmt.Println("Error:", err)
			}
		case 3:
			if err := personas.EliminarTerapeuta(); err != nil {
				fmt.Println("Error:", err)
			}
		case 4:
			if err := agenda.EliminarPlan(); err != nil {
				fmt.Println("Error:", err)
			}
		case 5:
			if err := agenda.EliminarCita(); err != nil {
				fmt.Println("Error:", err)
			}
		case 6:
			if err := agenda.EliminarAsistencia(); err != nil {
				fmt.Println("Error:", err)
			}
		case 7:
			return
		default:
			fmt.Println("Opcion no valida")
		}
	}
}

// main es la funcion principal del programa
func main() {
	ConectarBaseDeDatos()

	fmt.Println("")
	fmt.Println("___ Bienvenido al Sistema Rayitos de Sol ___")

	for {
		mostrarMenuPrincipal()
		opcion, err := leerNumero("")
		if errors.Is(err, io.EOF) {
			fmt.Println("")
			return
		}
		if err != nil {
			fmt.Println("Error de opcion:", err)
			continue
		}

		switch opcion {
		case 1:
			submenuRegistrar()
		case 2:
			submenuConsultar()
		case 3:
			submenuActualizar()
		case 4:
			submenuEliminar()
		case 5:
			fmt.Println("Gracias por usar el sistema. ¡Hasta luego!")
			return
		default:
			fmt.Println("Opcion no valida. Intente de nuevo.")
		}
	}
}

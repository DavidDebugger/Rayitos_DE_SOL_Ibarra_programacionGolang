package main

import (
	"fmt"

	"RDSI/agenda"
	"RDSI/personas"
)

// mostrarMenuPrincipal imprime el menu principal
func mostrarMenuPrincipal() {
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

// leerNumero lee un numero entero
func leerNumero(mensaje string) int {
	fmt.Print(mensaje)
	var num int
	fmt.Scanf("%d\n", &num)
	return num
}

// submenuRegistrar maneja el menu de registrar
func submenuRegistrar() {
	for {
		mostrarMenuRegistrar()
		opcion := leerNumero("")

		switch opcion {
		case 1:
			personas.RegistrarEstudiante()
		case 2:
			personas.RegistrarRepresentante()
		case 3:
			personas.RegistrarTerapeuta()
		case 4:
			agenda.RegistrarPlan()
		case 5:
			agenda.RegistrarCita()
		case 6:
			agenda.RegistrarAsistencia()
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
		opcion := leerNumero("")

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
			personas.DetalleEstudiante()
		case 10:
			agenda.DetalleCita()
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
		opcion := leerNumero("")

		switch opcion {
		case 1:
			personas.ActualizarEstudiante()
		case 2:
			personas.ActualizarRepresentante()
		case 3:
			personas.ActualizarTerapeuta()
		case 4:
			agenda.ActualizarPlan()
		case 5:
			agenda.ActualizarCita()
		case 6:
			agenda.ActualizarAsistencia()
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
		opcion := leerNumero("")

		switch opcion {
		case 1:
			personas.EliminarEstudiante()
		case 2:
			personas.EliminarRepresentante()
		case 3:
			personas.EliminarTerapeuta()
		case 4:
			agenda.EliminarPlan()
		case 5:
			agenda.EliminarCita()
		case 6:
			agenda.EliminarAsistencia()
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

	fmt.Println("___ Bienvenido al Sistema Rayitos de Sol ___")

	for {
		mostrarMenuPrincipal()
		opcion := leerNumero("")

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

package main

import (
	"fmt"
	"log"

	"RDSI/config"

	"gorm.io/gorm"
)

// ConectarBaseDeDatos establece la conexión a PostgreSQL usando config/database.go
// Variables de entorno: DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME
func ConectarBaseDeDatos() *gorm.DB {
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Error al conectar a la base de datos: %v", err)
	}
	fmt.Println("Conexión exitosa a PostgreSQL")
	return db
}

// MigrarBaseDeDatos ejecuta AutoMigrate de los modelos del proyecto.
// Los modelos son los structs de tus paquetes (personas, agenda, contabilidad, vista).
// El orden respeta las dependencias entre tablas (ver docs sección 2.5).
func MigrarBaseDeDatos(db *gorm.DB) error {
	models := []interface{}{
		// --- PERSONAS ---
		// &personas.Usuario{},
		// &personas.Representante{},
		// &personas.Terapeuta{},

		// --- AGENDA ---
		// &agenda.Plan{},
		// &personas.Estudiante{}, // depende de Representante y Plan
		// &agenda.Horario{},      // depende de Terapeuta
		// &agenda.Cita{},         // depende de Estudiante y Terapeuta
		// &agenda.Asistencia{},   // depende de Cita
		// &agenda.Valoracion{},   // depende de Estudiante y Terapeuta
		// &agenda.Recuperacion{}, // depende de Cita

		// --- CONTABILIDAD ---
		// &contabilidad.Factura{},            // depende de Representante
		// &contabilidad.Pago{},               // depende de Factura y Representante
		// &contabilidad.PagoEmpleado{},       // depende de Terapeuta
		// &contabilidad.Inventario{},
		// &contabilidad.MovimientoInventario{}, // depende de Inventario
		// &contabilidad.Contabilidad{},

		// --- VISTA ---
		// &vista.Dashboard{},
	}

	if len(models) == 0 {
		return fmt.Errorf("no hay modelos registrados para migrar (agrega tus structs en MigrarBaseDeDatos)")
	}

	if err := db.AutoMigrate(models...); err != nil {
		return fmt.Errorf("error en la migración: %w", err)
	}

	fmt.Println("Migración completada exitosamente")
	return nil
}

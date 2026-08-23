# Cosas Faltantes - Rayitos de Sol

Este documento lista funcionalidades, mejoras y correcciones que aun no estan implementadas en el sistema `RDSI` (Rayitos de Sol).

## Pendientes

### 1. Contabilidad
- **Estado actual:** El paquete `contabilidad/` contiene archivos vacios (`contabilidad.go`, `facturacion.go`, `inventario.go`, `pago-empleados.go`) que solo tienen la declaracion `package contabilidad` y no implementan nada.
- **Falta:** Implementar funcionalidades de facturacion, inventario y pago de empleados. Estos podrian incluir:
  - Calculo de costos por sesion/terapia.
  - Registro de pagos de pacientes y terapeutas.
  - Gestion de inventario de insumos.

### 2. Base de datos y GORM
- **Estado actual:** `base-datos.go` solo imprime un mensaje de placeholder. No hay persistencia real.
- **Falta:**
  - Conectar a una base de datos real (PostgreSQL) usando GORM.
  - Migrar las estructuras (`Estudiante`, `Representante`, `Terapeuta`, `Plan`, `Cita`, `Asistencia`) a tablas.
  - Implementar funciones de guardado y carga de datos (`db.Save`, `db.Find`).

### 3. Confirmacion de eliminacion
- **Estado actual:** Las funciones `Eliminar*` borran directamente sin preguntar.
- **Falta:** Agregar una confirmacion antes de borrar (ej: "¿Esta seguro? (s/n)").

### 4. Menu de buscar separado
- **Estado actual:** Las busquedas estan integradas en funciones como `BuscarEstudiantePorID` pero no hay un menu visible para buscar.
- **Falta:** Crear un menu "Buscar" en el menu principal o en cada submenu, permitiendo buscar por ID, nombre, etc.

### 5. Validar que el estudiante/terapeuta exista en RegistrarCita y RegistrarAsistencia
- **Estado actual:** `RegistrarCita` y `RegistrarAsistencia` aceptan cualquier nombre como string sin validar si existe realmente en `personas.ListaEstudiantes` o `personas.ListaTerapeutas`.
- **Falta:**
  - En `RegistrarCita`: validar que el estudiante y terapeuta ingresados existan.
  - Si no existen, preguntar si desea crearlos primero o abortar.
  - En `RegistrarAsistencia`: validar que el nombre del estudiante exista en `personas.ListaEstudiantes`.

### 6. Campo "Estado" de cita con opciones predefinidas
- **Estado actual:** `RegistrarCita` asigna `cita.Estado = "activa"` automaticamente.
- **Falta:** Permitir al usuario elegir el estado de la cita de una lista: activa, completada, cancelada. Esto podria usarse en `ActualizarCita` con un menu de seleccion.

---

*Este archivo se ira actualizando conforme se identifiquen nuevas carencias o se completan las existentes.*

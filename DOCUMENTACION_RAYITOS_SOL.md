# DOCUMENTACIÓN DEL PROYECTO
# Sistema de Gestión del Centro de Terapias Infantiles "Rayitos de Sol"

> Documento maestro de planificación y diseño técnico.  
> Estado: **Diseño inicial / planificación de desarrollo**.

---

## 1. Información general

### 1.1 Nombre del sistema

**Sistema de Gestión Integral para el Centro de Terapias Infantiles "Rayitos de Sol"**

### 1.2 Propósito

Desarrollar un sistema informático que centralice y automatice los procesos administrativos y operativos del centro de terapias.

El sistema permitirá administrar personas, agenda de terapias, asistencia, planes, facturación, pagos, inventario, contabilidad y estadísticas mediante un dashboard.

### 1.3 Lenguaje

**Go (Golang)**

### 1.4 Base de datos

**PostgreSQL**

### 1.5 Backend

**Go + Gin**

### 1.6 ORM

**GORM**

### 1.7 Arquitectura de comunicación

```text
Frontend
   |
   v
Go + Gin
   |
   v
GORM
   |
   v
PostgreSQL
```

La API será de tipo **REST**.

---

# 2. Contexto del problema

El centro de terapias necesita centralizar la información generada durante sus actividades administrativas y operativas.

Entre los procesos que se busca mejorar se encuentran:

- Registro y administración de estudiantes.
- Administración de representantes.
- Gestión de terapeutas.
- Organización de horarios.
- Gestión de citas.
- Registro de asistencia.
- Administración de planes de terapia.
- Registro de valoraciones.
- Gestión de recuperaciones.
- Facturación.
- Control de pagos.
- Pago a empleados/terapeutas.
- Control de inventario.
- Registro de ingresos y gastos.
- Generación de indicadores administrativos.

El sistema busca reducir la dependencia de procesos manuales, mejorar la disponibilidad de información y facilitar la toma de decisiones.

---

# 3. Objetivo general

Desarrollar un sistema de gestión para el Centro de Terapias Infantiles "Rayitos de Sol" utilizando Go, Gin, GORM y PostgreSQL, que permita automatizar y centralizar los procesos administrativos y operativos del establecimiento.

---

# 4. Objetivos específicos

- Gestionar información de estudiantes.
- Administrar representantes.
- Registrar y administrar terapeutas.
- Gestionar horarios.
- Crear y administrar citas.
- Registrar asistencia de estudiantes.
- Administrar planes de terapia.
- Registrar valoraciones.
- Gestionar recuperaciones de sesiones.
- Generar facturas.
- Registrar pagos.
- Administrar pagos de empleados.
- Registrar ingresos y gastos.
- Controlar inventario.
- Registrar movimientos de inventario.
- Generar indicadores mediante un dashboard.
- Mantener la integridad de la información mediante PostgreSQL.
- Exponer las funcionalidades mediante una API REST.

---

# 5. Alcance

El sistema se dividirá en cuatro áreas funcionales principales:

```text
PERSONAS
AGENDA
CONTABILIDAD
VISTA
```

Cada área tendrá responsabilidades independientes, pero existirá comunicación entre ellas mediante las relaciones de negocio.

---

# 6. Organización funcional

## 6.1 Personas

Responsable de la información de las personas relacionadas con el centro.

Incluye:

- Estudiantes.
- Representantes.
- Terapeutas.

---

## 6.2 Agenda

Responsable de la planificación y seguimiento de las terapias.

Incluye:

- Citas / agenda.
- Horarios.
- Asistencia.
- Planes.
- Valoraciones.
- Recuperaciones.

---

## 6.3 Contabilidad

Responsable de las operaciones económicas y de inventario.

Incluye:

- Facturación.
- Pagos.
- Pagos de empleados.
- Inventario.
- Movimientos de inventario.
- Ingresos.
- Gastos.
- Contabilidad.

---

## 6.4 Vista

Responsable de presentar información consolidada al usuario.

Incluye:

- Dashboard.
- Indicadores.
- Estadísticas.

El dashboard no debe considerarse una entidad principal de la base de datos. Su función es consultar y presentar información de los demás módulos.

---

# 7. Arquitectura general

```text
                         USUARIO
                            |
                            v
                     FRONTEND WEB
                            |
                            v
                    API REST - GIN
                            |
                            v
                     LÓGICA EN GO
                            |
                            v
                         GORM
                            |
                            v
                      POSTGRESQL
```

## 7.1 Responsabilidades

### Frontend

Presenta la interfaz gráfica y consume la API REST.

### Gin

Gestiona las peticiones HTTP, rutas, parámetros y respuestas de la API.

### Go

Contiene las estructuras, métodos y lógica de negocio.

### GORM

Actúa como ORM entre las estructuras de Go y PostgreSQL.

### PostgreSQL

Almacena la información persistente del sistema.

---

# 8. Organización del proyecto

La organización principal seguirá el criterio definido por el proyecto: **módulos funcionales por dominio**.

```text
rayitos-sol/
│
├── main.go
│
├── config/
│   └── database.go
│
├── personas/
│   ├── estudiante.go
│   ├── representante.go
│   ├── terapeuta.go
│   └── usuario.go
│
├── agenda/
│   ├── agenda.go
│   ├── horario.go
│   ├── asistencia.go
│   ├── plan.go
│   ├── valoracion.go
│   └── recuperacion.go
│
├── contabilidad/
│   ├── facturacion.go
│   ├── pago.go
│   ├── pago_empleado.go
│   ├── inventario.go
│   ├── movimiento.go
│   └── contabilidad.go
│
├── vista/
│   └── dashboard.go
│
├── database/
│   └── schema.sql
│
├── routes/
│   └── routes.go
│
└── go.mod
```

> La estructura puede evolucionar durante la implementación, pero no se debe perder la separación principal por dominios: `personas`, `agenda`, `contabilidad` y `vista`.

---

# 9. Módulo Personas

## 9.1 Estudiante

Representa al niño o niña que recibe servicios de terapia.

### Atributos iniciales

| Atributo | Tipo Go | Descripción |
|---|---|---|
| ID | int | Identificador único |
| Nombres | string | Nombres |
| Apellidos | string | Apellidos |
| FechaNacimiento | time.Time | Fecha de nacimiento |
| Genero | string | Género |
| Diagnostico | string | Diagnóstico |
| Direccion | string | Dirección |
| Telefono | string | Contacto |
| Estado | string | Estado del estudiante |
| RepresentanteID | int | Representante asociado |
| PlanID | int | Plan asociado |

### Métodos previstos

```text
Registrar()
Actualizar()
Eliminar()
Buscar()
CalcularEdad()
AsignarPlan()
CambiarEstado()
ConsultarHistorial()
```

---

## 9.2 Representante

Persona responsable del estudiante.

### Atributos iniciales

| Atributo | Tipo Go | Descripción |
|---|---|---|
| ID | int | Identificador |
| Cedula | string | Identificación |
| Nombres | string | Nombres |
| Apellidos | string | Apellidos |
| Telefono | string | Teléfono |
| Correo | string | Correo |
| Direccion | string | Dirección |

### Métodos previstos

```text
Registrar()
Actualizar()
Eliminar()
Buscar()
ObtenerEstudiantes()
```

---

## 9.3 Terapeuta

Profesional que presta los servicios de terapia.

### Atributos iniciales

| Atributo | Tipo Go | Descripción |
|---|---|---|
| ID | int | Identificador |
| Cedula | string | Identificación |
| Nombres | string | Nombres |
| Apellidos | string | Apellidos |
| Especialidad | string | Especialidad |
| Telefono | string | Teléfono |
| Correo | string | Correo |
| Estado | string | Estado |

### Métodos previstos

```text
Registrar()
Actualizar()
Eliminar()
ConsultarHorario()
ConsultarHorasDisponibles()
```

---

## 9.4 Usuario

Se contempla como entidad del sistema para autenticación y control de acceso.

La definición detallada de sus atributos y roles se realizará en la etapa de autenticación y seguridad.

---

# 10. Módulo Agenda

## 10.1 Agenda / Cita

Representa una sesión programada.

### Atributos iniciales

| Atributo | Tipo Go | Descripción |
|---|---|---|
| ID | int | Identificador |
| EstudianteID | int | Estudiante |
| TerapeutaID | int | Terapeuta |
| Fecha | time.Time | Fecha |
| Hora | string | Hora |
| Tipo | string | Tipo de sesión |
| Estado | string | Estado de la cita |

### Métodos previstos

```text
Agendar()
Cancelar()
Reprogramar()
ConsultarAgenda()
```

---

## 10.2 Horario

Representa la disponibilidad de un terapeuta.

### Atributos iniciales

| Atributo | Tipo Go | Descripción |
|---|---|---|
| ID | int | Identificador |
| TerapeutaID | int | Terapeuta |
| Dia | string | Día |
| HoraInicio | string | Inicio |
| HoraFin | string | Fin |
| Disponible | bool | Disponibilidad |

### Métodos previstos

```text
CrearHorario()
ModificarHorario()
ConsultarDisponibilidad()
```

---

## 10.3 Asistencia

Registra la asistencia asociada a una cita.

### Atributos iniciales

| Atributo | Tipo Go | Descripción |
|---|---|---|
| ID | int | Identificador |
| AgendaID | int | Cita |
| Asistio | bool | Asistencia |
| Observacion | string | Observaciones |

### Métodos previstos

```text
RegistrarAsistencia()
RegistrarInasistencia()
```

---

## 10.4 Plan

Representa un plan de terapia contratado.

### Atributos iniciales

| Atributo | Tipo Go | Descripción |
|---|---|---|
| ID | int | Identificador |
| Nombre | string | Nombre del plan |
| NumeroSesiones | int | Cantidad de sesiones |
| Precio | float64 | Precio |
| Estado | string | Estado |

### Métodos previstos

```text
CrearPlan()
EditarPlan()
Suspender()
Reactivar()
ConsumirSesion()
```

---

## 10.5 Valoración

Registra información relacionada con la valoración del estudiante.

### Atributos iniciales

| Atributo | Tipo Go | Descripción |
|---|---|---|
| ID | int | Identificador |
| EstudianteID | int | Estudiante |
| TerapeutaID | int | Terapeuta |
| Fecha | time.Time | Fecha |
| Observacion | string | Observación |
| Diagnostico | string | Diagnóstico |

### Métodos previstos

```text
Registrar()
Editar()
Consultar()
```

---

## 10.6 Recuperación

Permite gestionar una sesión que necesita ser recuperada o reprogramada.

### Atributos iniciales

| Atributo | Tipo Go | Descripción |
|---|---|---|
| ID | int | Identificador |
| AgendaOriginalID | int | Cita original |
| NuevaFecha | time.Time | Nueva fecha |
| Motivo | string | Motivo |

### Métodos previstos

```text
Registrar()
Reprogramar()
```

---

# 11. Módulo Contabilidad

## 11.1 Facturación

Registra las facturas generadas por los servicios.

### Atributos iniciales

| Atributo | Tipo Go | Descripción |
|---|---|---|
| ID | int | Identificador |
| RepresentanteID | int | Representante |
| Fecha | time.Time | Fecha |
| Total | float64 | Total |
| Estado | string | Estado |
| MetodoPago | string | Método de pago |

### Métodos previstos

```text
GenerarFactura()
RegistrarPago()
ConsultarFacturas()
```

---

## 11.2 Pago

Registra pagos realizados por los representantes.

Este componente se relacionará con la facturación y los servicios/planes correspondientes.

### Operaciones previstas

```text
RegistrarPago()
ConsultarPagos()
ConsultarHistorial()
```

---

## 11.3 Pago de empleado

Registra los pagos realizados a terapeutas u otros empleados.

### Atributos iniciales

| Atributo | Tipo Go | Descripción |
|---|---|---|
| ID | int | Identificador |
| TerapeutaID | int | Terapeuta |
| Fecha | time.Time | Fecha |
| Sueldo | float64 | Sueldo |
| Descuento | float64 | Descuento |
| Total | float64 | Total |

### Métodos previstos

```text
RegistrarPago()
CalcularTotal()
ConsultarHistorial()
```

---

## 11.4 Inventario

Administra los productos, materiales o recursos utilizados por el centro.

### Atributos iniciales

| Atributo | Tipo Go | Descripción |
|---|---|---|
| ID | int | Identificador |
| Nombre | string | Nombre |
| Descripcion | string | Descripción |
| Stock | int | Stock actual |
| StockMinimo | int | Nivel mínimo |

### Métodos previstos

```text
Registrar()
Actualizar()
Entrada()
Salida()
ConsultarStock()
```

---

## 11.5 Movimiento de inventario

Registra entradas y salidas del inventario.

### Atributos iniciales

| Atributo | Tipo Go | Descripción |
|---|---|---|
| ID | int | Identificador |
| InventarioID | int | Producto |
| Tipo | string | Entrada o salida |
| Cantidad | int | Cantidad |
| Fecha | time.Time | Fecha |

### Métodos previstos

```text
RegistrarMovimiento()
```

---

## 11.6 Contabilidad

Representa los movimientos financieros generales del centro.

### Atributos iniciales

| Atributo | Tipo Go | Descripción |
|---|---|---|
| ID | int | Identificador |
| Fecha | time.Time | Fecha |
| Concepto | string | Concepto |
| Tipo | string | Ingreso o gasto |
| Valor | float64 | Valor |

### Métodos previstos

```text
RegistrarIngreso()
RegistrarGasto()
CalcularBalance()
GenerarBalanceMensual()
```

---

# 12. Módulo Vista

## 12.1 Dashboard

El dashboard no almacena información propia. Consulta información de los diferentes módulos y la presenta de manera consolidada.

### Indicadores previstos

- Estudiantes activos.
- Representantes registrados.
- Terapeutas activos.
- Citas del día.
- Valoraciones.
- Recuperaciones.
- Sesiones pendientes.
- Ingresos del mes.
- Gastos del mes.
- Balance/ganancia.
- Horas ocupadas.
- Horas libres.
- Inventario bajo.

### Métodos previstos

```text
ObtenerIndicadores()
MostrarDashboard()
GenerarEstadisticas()
```

---

# 13. Modelo conceptual de relaciones

```text
Representante (1)
       |
       | 1:N
       v
Estudiante
       |
       | N:1
       v
      Plan

Estudiante (1) -------- N Cita N -------- (1) Terapeuta
                           |
                           |
                           +---- 0..1 Asistencia
                           |
                           +---- 0..1 Recuperación

Estudiante (1) -------- N Valoración N -------- (1) Terapeuta

Terapeuta (1) -------- N Horario

Representante (1) -------- N Factura

Terapeuta (1) -------- N PagoEmpleado

Inventario (1) -------- N Movimiento
```

---

# 14. Relaciones principales

| Entidad origen | Cardinalidad | Entidad destino |
|---|---:|---|
| Representante | 1:N | Estudiante |
| Estudiante | N:1 | Plan |
| Estudiante | 1:N | Cita |
| Terapeuta | 1:N | Cita |
| Terapeuta | 1:N | Horario |
| Cita | 1:0..1 | Asistencia |
| Cita | 1:0..1 | Recuperación |
| Estudiante | 1:N | Valoración |
| Terapeuta | 1:N | Valoración |
| Representante | 1:N | Factura |
| Terapeuta | 1:N | PagoEmpleado |
| Inventario | 1:N | Movimiento |

---

# 15. Base de datos PostgreSQL

Las tablas principales previstas son:

## Personas

```text
usuarios
representantes
estudiantes
terapeutas
```

## Agenda

```text
citas
horarios
asistencias
planes
valoraciones
recuperaciones
```

## Contabilidad

```text
facturas
pagos
pagos_empleados
inventario
movimientos_inventario
contabilidad
```

El dashboard no tendrá una tabla propia inicialmente.

---

# 16. Integridad de datos

PostgreSQL será responsable de garantizar la integridad referencial mediante:

- Primary Keys.
- Foreign Keys.
- Restricciones `NOT NULL`.
- Restricciones `UNIQUE`.
- Restricciones `CHECK`.
- Tipos de datos apropiados.
- Transacciones cuando una operación afecte varias tablas.

Ejemplos conceptuales:

```text
Un estudiante no debería poder existir con un representante inexistente.

Una cita no debería poder registrarse con un terapeuta inexistente.

Un movimiento de inventario debe pertenecer a un producto existente.

Una factura debe pertenecer a un representante válido.
```

---

# 17. Go y programación orientada a objetos

Go no implementa clases de la misma manera que Java o C#.

El proyecto utilizará:

- `struct`.
- Métodos con receivers.
- Encapsulación mediante visibilidad.
- Interfaces cuando sean necesarias.
- Composición.
- Separación de responsabilidades.

Ejemplo conceptual:

```go
type Estudiante struct {
    ID              int
    Nombres         string
    Apellidos       string
    FechaNacimiento time.Time
    Genero          string
    Diagnostico     string
    Estado          string
    RepresentanteID int
    PlanID          int
}

func (e *Estudiante) CambiarEstado(estado string) {
    e.Estado = estado
}
```

---

# 18. GORM

GORM será utilizado como ORM para relacionar las estructuras de Go con las tablas PostgreSQL.

Ejemplo conceptual:

```go
type Estudiante struct {
    ID              uint
    Nombres         string
    Apellidos       string
    FechaNacimiento time.Time
    Genero          string
    Diagnostico     string
    Estado          string

    RepresentanteID uint
    Representante   Representante

    PlanID uint
    Plan   Plan
}
```

Las relaciones se definirán mediante los campos correspondientes y las convenciones/etiquetas de GORM cuando sea necesario.

---

# 19. Gin

Gin será utilizado para construir la API REST.

Ejemplo conceptual de rutas:

```text
GET    /api/estudiantes
GET    /api/estudiantes/:id
POST   /api/estudiantes
PUT    /api/estudiantes/:id
DELETE /api/estudiantes/:id
```

Otros grupos de rutas previstos:

```text
/api/representantes
/api/terapeutas
/api/citas
/api/horarios
/api/asistencias
/api/planes
/api/valoraciones
/api/recuperaciones
/api/facturas
/api/pagos
/api/pagos-empleados
/api/inventario
/api/movimientos
/api/contabilidad
/api/dashboard
```

---

# 20. Flujo de una operación

Ejemplo: registrar un estudiante.

```text
Frontend
   |
   | POST /api/estudiantes
   v
Gin / Router
   |
   v
Controlador / Handler
   |
   v
Lógica del módulo Personas
   |
   v
GORM
   |
   v
PostgreSQL
```

El proceso inverso se utiliza para devolver la respuesta al frontend.

---

# 21. Configuración de base de datos

La conexión se centralizará en:

```text
config/database.go
```

La configuración no debe quedar dispersa por los diferentes módulos.

Se recomienda utilizar variables de entorno para datos sensibles como:

```text
DB_HOST
DB_PORT
DB_USER
DB_PASSWORD
DB_NAME
```

No se deben guardar contraseñas reales directamente en el código fuente.

---

# 22. Plan de desarrollo

El proyecto se desarrollará progresivamente.

## Fase 1 — Diseño

- Requisitos.
- Entidades.
- Relaciones.
- Modelo de base de datos.
- Diagrama de clases.
- Arquitectura.

## Fase 2 — Inicialización

- Crear proyecto Go.
- Crear `go.mod`.
- Instalar Gin.
- Instalar GORM.
- Instalar driver PostgreSQL.
- Configurar variables de entorno.

## Fase 3 — Base de datos

- Crear base de datos PostgreSQL.
- Crear tablas o migraciones.
- Definir claves primarias.
- Definir claves foráneas.
- Definir restricciones.
- Probar relaciones.

## Fase 4 — Personas

Implementar:

- Estudiantes.
- Representantes.
- Terapeutas.
- Usuarios, si corresponde.

## Fase 5 — Agenda

Implementar:

- Horarios.
- Citas.
- Asistencia.
- Planes.
- Valoraciones.
- Recuperaciones.

## Fase 6 — Contabilidad

Implementar:

- Facturación.
- Pagos.
- Pagos de empleados.
- Inventario.
- Movimientos.
- Ingresos.
- Gastos.

## Fase 7 — API

Crear endpoints REST para los módulos.

## Fase 8 — Dashboard

Crear consultas para indicadores y estadísticas.

## Fase 9 — Seguridad

Implementar:

- Autenticación.
- Autorización.
- Roles.
- Validación de datos.
- Protección de credenciales.
- Manejo seguro de errores.

## Fase 10 — Pruebas

- Pruebas unitarias.
- Pruebas de integración.
- Pruebas de API.
- Validación de base de datos.
- Pruebas de flujo completo.

## Fase 11 — Despliegue

Evaluar:

```text
Frontend
    -> Cloudflare Pages u otra plataforma

Backend
    -> Servidor/VPS o plataforma compatible con Go

PostgreSQL
    -> Servicio administrado o servidor propio
```

---

# 23. Estructura documental recomendada

Además del código, el proyecto puede mantenerse con:

```text
Rayitos-Sol/
│
├── README.md
│
├── docs/
│   ├── DOCUMENTACION.md
│   ├── DIAGRAMA_CLASES.md
│   ├── ARQUITECTURA.md
│   └── API.md
│
├── backend/
│   └── codigo-go
│
└── database/
    └── schema.sql
```

---

# 24. Principios de desarrollo

Durante la implementación se buscará mantener:

### Separación de responsabilidades

Cada módulo debe encargarse de su propio dominio.

### Reutilización

Evitar duplicar lógica innecesariamente.

### Integridad

Las operaciones que involucren varias tablas deben mantener consistencia.

### Seguridad

Las credenciales y datos sensibles no deben estar directamente en el código.

### Mantenibilidad

El código debe poder modificarse sin afectar innecesariamente otros módulos.

### Escalabilidad

La estructura debe permitir agregar funcionalidades posteriormente.

---

# 25. Decisiones técnicas actuales

| Aspecto | Decisión |
|---|---|
| Lenguaje | Go |
| Framework backend | Gin |
| ORM | GORM |
| Base de datos | PostgreSQL |
| Arquitectura API | REST |
| Organización principal | Por módulos funcionales |
| Módulos | Personas, Agenda, Contabilidad, Vista |
| Interfaz | Aplicación web |
| Dashboard | Sí |
| Integridad relacional | PostgreSQL |
| ORM-relación | GORM |

---

# 26. Estado actual del proyecto

```text
[✓] Idea general
[✓] Definición del sistema
[✓] Objetivo general
[✓] Objetivos específicos
[✓] Módulos principales
[✓] Stack tecnológico
[✓] Arquitectura general
[✓] Organización inicial del proyecto
[✓] Entidades iniciales
[✓] Relaciones principales
[ ] Modelo físico definitivo de PostgreSQL
[ ] Diagrama ER definitivo
[ ] Diagrama de clases definitivo
[ ] Inicialización del proyecto Go
[ ] Configuración GORM
[ ] Configuración PostgreSQL
[ ] Implementación Personas
[ ] Implementación Agenda
[ ] Implementación Contabilidad
[ ] API REST completa
[ ] Dashboard
[ ] Autenticación y autorización
[ ] Pruebas
[ ] Despliegue
```

---

# 27. Próximo paso

El siguiente paso recomendado es **cerrar el diseño de la base de datos antes de comenzar a programar**.

Orden:

```text
Entidades
   ↓
Atributos definitivos
   ↓
Relaciones
   ↓
Modelo ER
   ↓
PostgreSQL
   ↓
Structs de Go
   ↓
GORM
   ↓
API Gin
```

No se recomienda comenzar todavía con todos los CRUD hasta terminar las relaciones principales, porque cambios posteriores en las entidades pueden obligar a modificar simultáneamente la base de datos, los structs y los endpoints.

---

# 28. Resumen ejecutivo

**Rayitos de Sol** será un sistema web de gestión para un centro de terapias infantiles.

La solución utilizará:

```text
Go
 ├── Gin
 └── GORM
       |
       v
   PostgreSQL
```

La aplicación se organizará funcionalmente en:

```text
PERSONAS
 ├── Estudiantes
 ├── Representantes
 └── Terapeutas

AGENDA
 ├── Citas
 ├── Horarios
 ├── Asistencia
 ├── Planes
 ├── Valoraciones
 └── Recuperaciones

CONTABILIDAD
 ├── Facturación
 ├── Pagos
 ├── Pagos de empleados
 ├── Inventario
 ├── Movimientos
 ├── Ingresos
 └── Gastos

VISTA
 └── Dashboard
```

El desarrollo se realizará de manera incremental, comenzando por el diseño definitivo de las entidades y la base de datos, para posteriormente trasladar ese modelo a Go mediante structs y GORM y exponerlo mediante una API REST con Gin.

# Documentación de Base de Datos y Frontend
## Sistema de Gestión "Rayitos de Sol"

> Este documento describe el modelo de base de datos PostgreSQL y la interfaz web (frontend) del sistema.  
>
---

## 1. Base de Datos PostgreSQL

La base de datos se encuentra definida en **`database/schema.sql`** y contiene **16 tablas** organizadas en tres módulos: Personas, Agenda, y Contabilidad, más una tabla para el Dashboard.

### 1.1 Módulo Personas

#### `usuarios`
Tabla de autenticación y control de acceso (definición básica; los roles detallados se implementarán en la fase de seguridad).

| Campo | Tipo | Restricciones |
|---|---|---|
| `id` | SERIAL | PRIMARY KEY |
| `nombre_usuario` | VARCHAR(50) | UNIQUE, NOT NULL |
| `contrasena` | VARCHAR(255) | NOT NULL |
| `email` | VARCHAR(100) | UNIQUE, NOT NULL |
| `rol` | VARCHAR(50) | DEFAULT 'administrador' |
| `estado` | VARCHAR(20) | DEFAULT 'activo' |
| `creado_en` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |
| `actualizado_en` | TIMESTAMP | |

#### `representantes`
Persona responsable del estudiante.

| Campo | Tipo | Restricciones |
|---|---|---|
| `id` | SERIAL | PRIMARY KEY |
| `cedula` | VARCHAR(20) | UNIQUE, NOT NULL |
| `nombres` | VARCHAR(100) | NOT NULL |
| `apellidos` | VARCHAR(100) | NOT NULL |
| `telefono` | VARCHAR(20) | |
| `correo` | VARCHAR(100) | |
| `direccion` | TEXT | |
| `creado_en` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |
| `actualizado_en` | TIMESTAMP | |

#### `terapeutas`
Profesional que presta servicios de terapia.

| Campo | Tipo | Restricciones |
|---|---|---|
| `id` | SERIAL | PRIMARY KEY |
| `cedula` | VARCHAR(20) | UNIQUE, NOT NULL |
| `nombres` | VARCHAR(100) | NOT NULL |
| `apellidos` | VARCHAR(100) | NOT NULL |
| `especialidad` | VARCHAR(100) | |
| `telefono` | VARCHAR(20) | |
| `correo` | VARCHAR(100) | |
| `estado` | VARCHAR(20) | DEFAULT 'activo' |
| `creado_en` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |
| `actualizado_en` | TIMESTAMP | |

#### `planes`
Plan de terapia contratado.

| Campo | Tipo | Restricciones |
|---|---|---|
| `id` | SERIAL | PRIMARY KEY |
| `nombre` | VARCHAR(100) | NOT NULL |
| `numero_sesiones` | INTEGER | CHECK (numero_sesiones >= 0) |
| `precio` | NUMERIC(10,2) | CHECK (precio >= 0) |
| `estado` | VARCHAR(20) | DEFAULT 'activo' |
| `creado_en` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |
| `actualizado_en` | TIMESTAMP | |

#### `estudiantes`
Niño/a que recibe terapia. Referencia a `representantes` y `planes`.

| Campo | Tipo | Restricciones |
|---|---|---|
| `id` | SERIAL | PRIMARY KEY |
| `nombres` | VARCHAR(100) | NOT NULL |
| `apellidos` | VARCHAR(100) | NOT NULL |
| `fecha_nacimiento` | DATE | NOT NULL |
| `genero` | VARCHAR(20) | |
| `diagnostico` | TEXT | |
| `direccion` | TEXT | |
| `telefono` | VARCHAR(20) | |
| `estado` | VARCHAR(20) | DEFAULT 'activo' |
| `representante_id` | INTEGER | NOT NULL, FK → representantes(id) ON DELETE RESTRICT |
| `plan_id` | INTEGER | FK → planes(id) ON DELETE SET NULL |
| `creado_en` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |
| `actualizado_en` | TIMESTAMP | |

### 1.2 Módulo Agenda

#### `horarios`
Disponibilidad de un terapeuta.

| Campo | Tipo | Restricciones |
|---|---|---|
| `id` | SERIAL | PRIMARY KEY |
| `terapeuta_id` | INTEGER | NOT NULL, FK → terapeutas(id) ON DELETE CASCADE |
| `dia` | VARCHAR(20) | NOT NULL |
| `hora_inicio` | TIME | NOT NULL |
| `hora_fin` | TIME | NOT NULL |
| `disponible` | BOOLEAN | DEFAULT TRUE |
| `creado_en` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |
| `actualizado_en` | TIMESTAMP | |

> **CHECK**: `hora_fin > hora_inicio`

#### `citas`
Sesión programada.

| Campo | Tipo | Restricciones |
|---|---|---|
| `id` | SERIAL | PRIMARY KEY |
| `estudiante_id` | INTEGER | NOT NULL, FK → estudiantes(id) ON DELETE RESTRICT |
| `terapeuta_id` | INTEGER | NOT NULL, FK → terapeutas(id) ON DELETE RESTRICT |
| `fecha` | DATE | NOT NULL |
| `hora` | TIME | NOT NULL |
| `tipo` | VARCHAR(50) | NOT NULL |
| `estado` | VARCHAR(20) | DEFAULT 'programada' |
| `creado_en` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |
| `actualizado_en` | TIMESTAMP | |

> **CHECK**: `estado IN ('programada', 'confirmada', 'cancelada', 'completada', 'reprogramada')`

#### `asistencias`
Registro de asistencia asociado a una cita.

| Campo | Tipo | Restricciones |
|---|---|---|
| `id` | SERIAL | PRIMARY KEY |
| `cita_id` | INTEGER | NOT NULL, FK → citas(id) ON DELETE CASCADE |
| `asistio` | BOOLEAN | NOT NULL |
| `observacion` | TEXT | |
| `creado_en` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |
| `actualizado_en` | TIMESTAMP | |

#### `valoraciones`
Informe de valoración del estudiante por el terapeuta.

| Campo | Tipo | Restricciones |
|---|---|---|
| `id` | SERIAL | PRIMARY KEY |
| `estudiante_id` | INTEGER | NOT NULL, FK → estudiantes(id) ON DELETE RESTRICT |
| `terapeuta_id` | INTEGER | NOT NULL, FK → terapeutas(id) ON DELETE RESTRICT |
| `fecha` | DATE | NOT NULL |
| `observacion` | TEXT | |
| `diagnostico` | TEXT | |
| `creado_en` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |
| `actualizado_en` | TIMESTAMP | |

#### `recuperaciones`
Gestión de sesiones recuperadas.

| Campo | Tipo | Restricciones |
|---|---|---|
| `id` | SERIAL | PRIMARY KEY |
| `cita_original_id` | INTEGER | NOT NULL, FK → citas(id) ON DELETE RESTRICT |
| `nueva_fecha` | DATE | NOT NULL |
| `motivo` | TEXT | |
| `estado` | VARCHAR(20) | DEFAULT 'pendiente' |
| `creado_en` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |
| `actualizado_en` | TIMESTAMP | |

> **CHECK**: `estado IN ('pendiente', 'programada', 'completada', 'cancelada')`

### 1.3 Módulo Contabilidad

#### `facturas`
Factura generada por servicios.

| Campo | Tipo | Restricciones |
|---|---|---|
| `id` | SERIAL | PRIMARY KEY |
| `representante_id` | INTEGER | NOT NULL, FK → representantes(id) ON DELETE RESTRICT |
| `fecha` | DATE | NOT NULL |
| `total` | NUMERIC(10,2) | CHECK (total >= 0) |
| `estado` | VARCHAR(20) | DEFAULT 'pendiente' |
| `metodo_pago` | VARCHAR(50) | |
| `creado_en` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |
| `actualizado_en` | TIMESTAMP | |

> **CHECK**: `estado IN ('pendiente', 'pagada', 'cancelada', 'anulada')`

#### `pagos`
Pagos realizados por representantes contra facturas.

| Campo | Tipo | Restricciones |
|---|---|---|
| `id` | SERIAL | PRIMARY KEY |
| `factura_id` | INTEGER | NOT NULL, FK → facturas(id) ON DELETE RESTRICT |
| `representante_id` | INTEGER | NOT NULL, FK → representantes(id) ON DELETE RESTRICT |
| `fecha` | DATE | NOT NULL |
| `monto` | NUMERIC(10,2) | CHECK (monto >= 0) |
| `metodo_pago` | VARCHAR(50) | |
| `observacion` | TEXT | |
| `creado_en` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |
| `actualizado_en` | TIMESTAMP | |

#### `pagos_empleados`
Pagos realizados a terapeutas.

| Campo | Tipo | Restricciones |
|---|---|---|
| `id` | SERIAL | PRIMARY KEY |
| `terapeuta_id` | INTEGER | NOT NULL, FK → terapeutas(id) ON DELETE RESTRICT |
| `fecha` | DATE | NOT NULL |
| `sueldo` | NUMERIC(10,2) | CHECK (sueldo >= 0) |
| `descuento` | NUMERIC(10,2) | DEFAULT 0, CHECK (descuento >= 0) |
| `total` | NUMERIC(10,2) | CHECK (total >= 0) |
| `creado_en` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |
| `actualizado_en` | TIMESTAMP | |

#### `inventario`
Productos, materiales o recursos del centro.

| Campo | Tipo | Restricciones |
|---|---|---|
| `id` | SERIAL | PRIMARY KEY |
| `nombre` | VARCHAR(100) | NOT NULL |
| `descripcion` | TEXT | |
| `stock` | INTEGER | DEFAULT 0, CHECK (stock >= 0) |
| `stock_minimo` | INTEGER | DEFAULT 0, CHECK (stock_minimo >= 0) |
| `creado_en` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |
| `actualizado_en` | TIMESTAMP | |

#### `movimientos_inventario`
Entradas y salidas de inventario.

| Campo | Tipo | Restricciones |
|---|---|---|
| `id` | SERIAL | PRIMARY KEY |
| `inventario_id` | INTEGER | NOT NULL, FK → inventario(id) ON DELETE RESTRICT |
| `tipo` | VARCHAR(20) | CHECK (tipo IN ('entrada', 'salida')) |
| `cantidad` | INTEGER | CHECK (cantidad > 0) |
| `fecha` | DATE | NOT NULL |
| `creado_en` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |
| `actualizado_en` | TIMESTAMP | |

#### `contabilidad`
Movimientos financieros generales (ingresos y gastos).

| Campo | Tipo | Restricciones |
|---|---|---|
| `id` | SERIAL | PRIMARY KEY |
| `fecha` | DATE | NOT NULL |
| `concepto` | VARCHAR(200) | NOT NULL |
| `tipo` | VARCHAR(20) | CHECK (tipo IN ('ingreso', 'gasto')) |
| `valor` | NUMERIC(10,2) | CHECK (valor >= 0) |
| `creado_en` | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |
| `actualizado_en` | TIMESTAMP | |

#### `dashboard`
No es una tabla propia. El dashboard es una capa de consultas que agrega información de los módulos anteriores. Se define como tabla para almacenar snapshots de indicadores si se requiere histórico.

| Campo | Tipo | Descripción |
|---|---|---|
| `id` | SERIAL | PRIMARY KEY |
| `creado_en` | TIMESTAMP | Fecha del snapshot |
| `estudiantes_activos` | INTEGER | |
| `representantes_total` | INTEGER | |
| `terapeutas_activos` | INTEGER | |
| `citas_del_dia` | INTEGER | |
| `valoraciones_total` | INTEGER | |
| `recuperaciones_total` | INTEGER | |
| `sesiones_pendientes` | INTEGER | |
| `ingresos_del_mes` | NUMERIC(10,2) | |
| `gastos_del_mes` | NUMERIC(10,2) | |
| `balance_ganancia` | NUMERIC(10,2) | |
| `horas_ocupadas` | NUMERIC(10,2) | |
| `horas_libres` | NUMERIC(10,2) | |
| `inventario_bajo` | INTEGER | |

### 1.4 Índices

Se crean los siguientes índices para optimizar consultas frecuentes:

- `idx_estudiantes_representante` — estudiantes(representante_id)
- `idx_estudiantes_plan` — estudiantes(plan_id)
- `idx_citas_estudiante` — citas(estudiante_id)
- `idx_citas_terapeuta` — citas(terapeuta_id)
- `idx_citas_fecha` — citas(fecha)
- `idx_horarios_terapeuta` — horarios(terapeuta_id)
- `idx_valoraciones_estudiante` — valoraciones(estudiante_id)
- `idx_valoraciones_terapeuta` — valoraciones(terapeuta_id)
- `idx_facturas_representante` — facturas(representante_id)
- `idx_pagos_factura` — pagos(factura_id)
- `idx_pagos_empleados_terapeuta` — pagos_empleados(terapeuta_id)
- `idx_movimientos_inventario` — movimientos_inventario(inventario_id)
- `idx_contabilidad_fecha` — contabilidad(fecha)
- `idx_contabilidad_tipo` — contabilidad(tipo)

### 1.5 Diagrama de Relaciones

```
Representante (1) ──── N ──── Estudiante (N) ──── 1 ──── Plan
                        │
                        ├── N ── Cita (N) ──── 1 ──── Terapeuta
                        │          │
                        │          ├── 0..1 ── Asistencia
                        │          └── 0..1 ── Recuperación
                        │
                        ├── N ── Valoración (N) ──── 1 ──── Terapeuta
                        │
                        └── N ── Factura
                                  │
                                  └── N ── Pago

Terapeuta (1) ──── N ──── Horario
Terapeuta (1) ──── N ──── PagoEmpleado

Inventario (1) ──── N ──── MovimientoInventario

Contabilidad (ingresos / gastos)
```

### 1.6 Integridad de Datos

PostgreSQL garantiza la integridad mediante:

1. **Primary Keys**: Todas las tablas usan `SERIAL PRIMARY KEY` (autoincremental).
2. **Foreign Keys**: Restricciones `ON DELETE RESTRICT` protegen registros referenciados; `ON DELETE CASCADE` elimina hijos automáticamente (horarios → terapeuta).
3. **NOT NULL**: Campos obligatorios marcados explícitamente.
4. **UNIQUE**: `cedula`, `nombre_usuario`, `email` son únicos.
5. **CHECK**: Estados válidos, valores no negativos, horas de horario coherentes.
6. **Tipos adecuados**: `NUMERIC(10,2)` para moneda, `TIME` para horas, `DATE` para fechas.
7. **Transacciones**: Las operaciones multi-tabla (ej. registrar pago + factura) deben usar transacciones.

---

## 2. Interfaz Web (Frontend)

> **IMPORTANTE**: El frontend ya está construido en la carpeta **`frontend/`** (15 páginas HTML + sidebar compartido, con Tailwind CSS vía CDN, CSS y JS propios en `static/`). Consumen la API REST mediante `fetch()` con un wrapper central (`app.js`). Tu tarea es implementar el backend Go + Gin que exponga los endpoints descritos en la sección 2.2.

### 2.1 Arquitectura

```
Usuario (browser)
    │
    ▼
Frontend Web (frontend/ → HTML + JS + Tailwind)
    │
    │ HTTP REST
    ▼
Go + Gin (API)
    │
    ▼
GORM → PostgreSQL
```

El frontend consume la API REST expuesta por Gin. No accede directamente a la base de datos.

### 2.1.1 Estructura del frontend

```text
frontend/
├── dashboard.html              → indicadores y gráficas /api/dashboard
├── sidebar.html                → menú de navegación compartido (se carga con fetch)
├── static/
│   ├── css/app.css             → sistema de diseño (sidebar, cards, tablas, badges, modales, toasts)
│   └── js/app.js               → api(), toast(), abrirModal(), badge(), llenarSelect(), cargarSidebar()...
├── personas/
│   ├── estudiantes.html        → CRUD → /api/personas/estudiantes
│   ├── representantes.html     → CRUD → /api/personas/representantes
│   └── terapeutas.html         → CRUD → /api/personas/terapeutas
├── agenda/
│   ├── citas.html              → CRUD + cambio de estado → /api/agenda/citas
│   ├── horarios.html           → CRUD → /api/agenda/horarios
│   ├── asistencias.html        → CRUD + filtro → /api/agenda/asistencias
│   ├── planes.html             → CRUD → /api/agenda/planes
│   ├── valoraciones.html       → CRUD → /api/agenda/valoraciones
│   └── recuperaciones.html     → CRUD + estado → /api/agenda/recuperaciones
└── contabilidad/
    ├── facturas.html           → CRUD + registrar pago → /api/contabilidad/facturas
    ├── pagos.html              → CRUD → /api/contabilidad/pagos
    ├── pagos_empleados.html    → CRUD → /api/contabilidad/pagos-empleados
    ├── inventario.html         → CRUD + bajo stock → /api/contabilidad/inventario
    ├── movimientos.html        → CRUD + filtro → /api/contabilidad/movimientos
    └── contabilidad.html       → ingresos/gastos + balance mensual → /api/contabilidad/contabilidad
```

Cada página carga datos al iniciar con `fetch()` (centralizado en `api()`), dibuja la tabla con datos crudos de la API y usa modales para crear registros. El sidebar se inyecta desde `sidebar.html` con `cargarSidebar()`. El servidor Gin debe servir:

- `r.Static("/static", "./frontend/static")` → CSS y JS.
- `r.Static("/personas", "./frontend/personas")`, `r.Static("/agenda", "./frontend/agenda")`, `r.Static("/contabilidad", "./frontend/contabilidad")` o alternativamente `r.Static("/", "./frontend")` y las páginas se abren directamente con `r.GET("/dashboard", ...)` sirviendo `frontend/dashboard.html`.

### 2.1.2 Convenciones de datos consumidas

- Fechas: string `YYYY-MM-DD`. Horas: `HH:MM:SS` o `HH:MM`.
- Montos: número decimal (USD).
- Badges por estado: `activo`→verde, `inactivo`→gris, `pendiente`→ámbar, `pagada`/`completada`/`asistio`→verde, `vencida`/`cancelada`/`falto`/`inactivo`→rojo, `programada`/`confirmada`→azul, `realizada`→verde.
- Los selects de formularios cargan opciones desde los GET listados (ej. estudiantes.html carga representantes y planes para poblar sus selects).

### 2.2 Endpoints de la API REST

> Los endpoints que el frontend consume **exactamente** (verificado con los `fetch()` de los HTML). El backend debe exponer estas rutas:

| Recurso | Endpoint | Método HTTP | Usado por |
|---|---|---|---|
| Estudiantes | `/api/personas/estudiantes` | GET | estudiantes.html |
| | `/api/personas/estudiantes` | POST | estudiantes.html |
| | `/api/personas/estudiantes/:id` | PUT | estudiantes.html |
| | `/api/personas/estudiantes/:id` | DELETE | estudiantes.html |
| Representantes | `/api/personas/representantes` | GET | estudiantes.html, facturas.html, pagos.html |
| | `/api/personas/representantes` | POST | representantes.html |
| | `/api/personas/representantes/:id` | DELETE | representantes.html |
| Terapeutas | `/api/personas/terapeutas` | GET | citas.html, horarios.html, pagos_empleados.html |
| | `/api/personas/terapeutas` | POST | terapeutas.html |
| | `/api/personas/terapeutas/:id` | DELETE | terapeutas.html |
| Citas | `/api/agenda/citas` | GET | citas.html, asistencias.html, recuperaciones.html |
| | `/api/agenda/citas` | POST | citas.html |
| | `/api/agenda/citas/:id` | PUT | citas.html (cambio de estado) |
| Horarios | `/api/agenda/horarios` | GET | horarios.html |
| | `/api/agenda/horarios` | POST | horarios.html |
| | `/api/agenda/horarios/:id` | DELETE | horarios.html |
| Planes | `/api/agenda/planes` | GET | estudiantes.html, planes.html |
| | `/api/agenda/planes` | POST | planes.html |
| | `/api/agenda/planes/:id` | DELETE | planes.html |
| Asistencias | `/api/agenda/asistencias` | GET | asistencias.html |
| | `/api/agenda/asistencias` | POST | asistencias.html |
| | `/api/agenda/asistencias/:id` | DELETE | asistencias.html |
| Valoraciones | `/api/agenda/valoraciones` | GET | valoraciones.html |
| | `/api/agenda/valoraciones` | POST | valoraciones.html |
| | `/api/agenda/valoraciones/:id` | DELETE | valoraciones.html |
| Recuperaciones | `/api/agenda/recuperaciones` | GET | recuperaciones.html |
| | `/api/agenda/recuperaciones` | POST | recuperaciones.html |
| | `/api/agenda/recuperaciones/:id` | PUT | recuperaciones.html (estado) |
| | `/api/agenda/recuperaciones/:id/reprogramar` | PUT | recuperaciones.html |
| Facturas | `/api/contabilidad/facturas` | GET | facturas.html, pagos.html |
| | `/api/contabilidad/facturas` | POST | facturas.html |
| | `/api/contabilidad/facturas/:id/registrar-pago` | POST | facturas.html |
| Pagos | `/api/contabilidad/pagos` | GET | pagos.html |
| | `/api/contabilidad/pagos` | POST | pagos.html |
| | `/api/contabilidad/pagos/:id` | DELETE | pagos.html |
| Pagos Empleados | `/api/contabilidad/pagos-empleados` | GET | pagos_empleados.html |
| | `/api/contabilidad/pagos-empleados` | POST | pagos_empleados.html |
| | `/api/contabilidad/pagos-empleados/:id` | DELETE | pagos_empleados.html |
| Inventario | `/api/contabilidad/inventario` | GET | inventario.html, movimientos.html |
| | `/api/contabilidad/inventario` | POST | inventario.html |
| | `/api/contabilidad/inventario/bajo-stock` | GET | inventario.html |
| | `/api/contabilidad/inventario/:id` | PUT | inventario.html |
| | `/api/contabilidad/inventario/:id` | DELETE | inventario.html |
| Movimientos | `/api/contabilidad/movimientos` | GET | movimientos.html |
| | `/api/contabilidad/movimientos` | POST | movimientos.html |
| | `/api/contabilidad/movimientos/inventario/:id` | GET | movimientos.html |
| | `/api/contabilidad/movimientos/:id` | DELETE | movimientos.html |
| Contabilidad | `/api/contabilidad/contabilidad` | GET | contabilidad.html |
| | `/api/contabilidad/contabilidad/registrar-ingreso` | POST | contabilidad.html |
| | `/api/contabilidad/contabilidad/registrar-gasto` | POST | contabilidad.html |
| | `/api/contabilidad/contabilidad/balance-mensual` | GET | contabilidad.html |
| | `/api/contabilidad/contabilidad/:id` | DELETE | contabilidad.html |
| Dashboard | `/api/dashboard` | GET | (indicadores) |
| | `/api/dashboard/estadisticas` | GET | (estadísticas) |

### 2.3 Vistas/Interfaces del Frontend

#### Dashboard Principal
Muestra indicadores clave en 8 stat-cards (cargadas desde `/api/dashboard`), 3 indicadores rápidos (hoy, semana, citas pendientes) y 2 gráficas Chart.js desde `/api/dashboard/estadisticas`: barras de ingresos vs gastos (6 meses) y dona de citas por estado. Incluye topbar con fecha del día en español.

#### Módulo Personas
1. **Estudiantes**: tabla con filtros por nombre, representante, plan, estado. Formulario de registro con selector de representante y plan.
2. **Representantes**: tabla con buscador por cédula o nombre. Formulario de registro completo.
3. **Terapeutas**: tabla con filtros por especialidad y estado. Formulario de registro.

#### Módulo Agenda
1. **Citas**: tabla con buscador y filtro por estado. Crear cita con selects de estudiante/terapeuta, fecha, hora y tipo. Cambiar estado (completada/cancelada) con botones por fila.
2. **Horarios**: tabla por terapeuta con días y rangos de hora. Marcar disponibilidad.
3. **Planes**: tabla con precios y número de sesiones. Formulario de creación/edición.
4. **Asistencias**: registrado al momento de la cita. Botones de asistió/no asistió.
5. **Valoraciones**: formulario con observaciones y diagnóstico. Historial por estudiante.
6. **Recuperaciones**: listar sesiones pendientes de recuperación. Reprogramar con nueva fecha.

#### Módulo Contabilidad
1. **Facturación**: generar factura seleccionando representante. Ver historial de facturas con estado y método de pago.
2. **Pagos**: registrar pagos contra facturas. Ver historial por representante.
3. **Pagos Empleados**: registrar sueldos con descuentos. Ver historial por terapeuta.
4. **Inventario**: tabla con stock actual y mínimo. Color rojo cuando stock está bajo. Entradas/salidas con motivo.
5. **Contabilidad**: formulario rápido para registrar ingresos y gastos. Reporte mensual de balance.

### 2.4 Variables de Entorno del Backend

El backend utiliza las siguientes variables de entorno para la conexión a PostgreSQL:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=tu_password
DB_NAME=rayitos_sol
```

### 2.5 Configuración de la Base de Datos

#### Script de Inicialización

Para crear la base de datos desde cero:

```bash
# 1. Crear la base de datos en PostgreSQL
createdb rayitos_sol

# 2. Ejecutar el schema
psql -U postgres -d rayitos_sol -f database/schema.sql
```

#### Migraciones con GORM

El backend usará `db.AutoMigrate()` con los modelos definidos en los paquetes `personas`, `agenda`, `contabilidad`, y `vista`. El orden de migración respeta las dependencias:

1. `usuarios`, `representantes`, `terapeutas`
2. `planes`
3. `estudiantes` (depende de representantes y planes)
4. `horarios` (depende de terapeutas)
5. `citas` (depende de estudiantes y terapeutas)
6. `asistencias` (depende de citas)
7. `valoraciones` (depende de estudiantes y terapeutas)
8. `recuperaciones` (depende de citas)
9. `facturas` (depende de representantes)
10. `pagos` (depende de facturas y representantes)
11. `pagos_empleados` (depende de terapeutas)
12. `inventario`
13. `movimientos_inventario` (depende de inventario)
14. `contabilidad`
15. `dashboard`

### 2.6 Consideraciones de UI/UX

- **Responsive**: La interfaz debe adaptarse a tablets y desktops.
- **Formulario en modal**: El registro/edición de entidades se hará mediante modales.
- **Validación frontend**: Verificar campos requeridos, formato de email, cédula, etc.
- **Feedback visual**: Mensajes de éxito/error al realizar operaciones.
- **Calendario de citas**: Vista de calendario con colores por estado de cita.
- **Permisos por rol**: El dashboard y ciertas acciones requerirán autenticación (fase 9).

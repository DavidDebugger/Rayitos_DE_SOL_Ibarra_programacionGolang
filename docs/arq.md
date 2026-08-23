# Arquitectura del Sistema Rayitos de Sol

## Introduccion
Este documento describe la arquitectura del sistema de gestion para el centro de terapia "Rayitos de Sol".

## Estructura de Directorios
```
rayitos/
├── main.go              # Archivo principal con el menu
├── base-datos.go        # Placeholder para GORM (no implementado)
├── util/
│   └── entrada.go       # Funciones auxiliares para leer entrada
├── personas/
│   ├── estudiantes.go   # CRUD de estudiantes
│   ├── representantes.go # CRUD de representantes
│   └── terapeutas.go    # CRUD de terapeutas
├── agenda/
│   ├── gestion-planes.go  # CRUD de planes de terapia
│   ├── horarios.go        # CRUD de citas/horarios
│   └── asistencia-estudiantes.go # CRUD de asistencias
├── config/
│   └── database.go      # Configuracion de base de datos (PostgreSQL/GORM)
├── contabilidad/        # Pendiente de implementar
├── vista/               # Dashboard (pendiente)
└── frontend/            # Interfaz web (pendiente)
```

## Arquitectura por Capas

### 1. Capa de Presentacion (main.go)
- Menu principal con 5 opciones: Registrar, Consultar, Actualizar, Eliminar, Salir
- Submenus para cada accion
- No tiene logica de negocio, solo navegacion

### 2. Capa de Utilidades (util/entrada.go)
- `LeerTexto(mensaje)`: Lee una linea con espacios
- `LeerNumero(mensaje)`: Lee un numero entero
- Ambas usan `bufio.Reader` global para consistencia

### 3. Capa de Negocio (personas/ y agenda/)
- Cada paquete tiene funciones CRUD completas
- Usan variables globales (slices) para almacenar datos
- Funciones simples con comentarios claros

## Funcionalidades Implementadas

### Personas
- **Estudiantes**: Registrar, Listar, Buscar, Actualizar, Eliminar, Detalle
- **Representantes**: Registrar, Listar, Actualizar, Eliminar, Detalle
- **Terapeutas**: Registrar, Listar, Actualizar, Eliminar, Detalle

### Agenda
- **Planes**: Registrar, Listar, Actualizar, Eliminar, Detalle
- **Citas**: Registrar, Listar, ListarPorFecha, Actualizar, Eliminar, Detalle
- **Asistencias**: Registrar, Listar, ListarPorEstudiante, Actualizar, Eliminar, Detalle

## Pendientes
- Base de datos (GORM + PostgreSQL)
- Frontend/web
- Contabilidad (facturacion, inventario, pagos)
- Dashboard de estadisticas

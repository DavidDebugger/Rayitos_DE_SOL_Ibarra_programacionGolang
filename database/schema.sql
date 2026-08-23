-- ============================================================
-- SCHEMA SQL - Sistema Rayitos de Sol
-- Generado desde DOCUMENTACION_RAYITOS_SOL.md
-- ============================================================

-- ============================================================
-- TABLAS DEL MÓDULO PERSONAS
-- ============================================================

-- Tabla: usuarios
-- Entidad contemplada para autenticación y control de acceso.
-- La definición detallada de roles se realizará en la etapa de seguridad.
CREATE TABLE usuarios (
    id              SERIAL PRIMARY KEY,
    nombre_usuario  VARCHAR(50) NOT NULL UNIQUE,
    contrasena      VARCHAR(255) NOT NULL,
    email           VARCHAR(100) NOT NULL UNIQUE,
    rol             VARCHAR(50) NOT NULL DEFAULT 'administrador',
    estado          VARCHAR(20) NOT NULL DEFAULT 'activo',
    creado_en       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    actualizado_en  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla: representantes
CREATE TABLE representantes (
    id          SERIAL PRIMARY KEY,
    cedula      VARCHAR(20) NOT NULL UNIQUE,
    nombres     VARCHAR(100) NOT NULL,
    apellidos   VARCHAR(100) NOT NULL,
    telefono    VARCHAR(20),
    correo      VARCHAR(100),
    direccion   TEXT,
    creado_en   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    actualizado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla: terapeutas
CREATE TABLE terapeutas (
    id          SERIAL PRIMARY KEY,
    cedula      VARCHAR(20) NOT NULL UNIQUE,
    nombres     VARCHAR(100) NOT NULL,
    apellidos   VARCHAR(100) NOT NULL,
    especialidad VARCHAR(100),
    telefono    VARCHAR(20),
    correo      VARCHAR(100),
    estado      VARCHAR(20) NOT NULL DEFAULT 'activo',
    creado_en   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    actualizado_en TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla: planes
CREATE TABLE planes (
    id              SERIAL PRIMARY KEY,
    nombre          VARCHAR(100) NOT NULL,
    numero_sesiones INTEGER NOT NULL CHECK (numero_sesiones >= 0),
    precio          NUMERIC(10,2) NOT NULL CHECK (precio >= 0),
    estado          VARCHAR(20) NOT NULL DEFAULT 'activo',
    creado_en       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    actualizado_en  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla: estudiantes
CREATE TABLE estudiantes (
    id              SERIAL PRIMARY KEY,
    nombres         VARCHAR(100) NOT NULL,
    apellidos       VARCHAR(100) NOT NULL,
    fecha_nacimiento DATE NOT NULL,
    genero          VARCHAR(20),
    diagnostico     TEXT,
    direccion       TEXT,
    telefono        VARCHAR(20),
    estado          VARCHAR(20) NOT NULL DEFAULT 'activo',
    representante_id INTEGER NOT NULL REFERENCES representantes(id) ON DELETE RESTRICT,
    plan_id         INTEGER REFERENCES planes(id) ON DELETE SET NULL,
    creado_en       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    actualizado_en  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ============================================================
-- TABLAS DEL MÓDULO AGENDA
-- ============================================================

-- Tabla: horarios
CREATE TABLE horarios (
    id              SERIAL PRIMARY KEY,
    terapeuta_id    INTEGER NOT NULL REFERENCES terapeutas(id) ON DELETE CASCADE,
    dia             VARCHAR(20) NOT NULL,
    hora_inicio     TIME NOT NULL,
    hora_fin        TIME NOT NULL,
    disponible      BOOLEAN NOT NULL DEFAULT TRUE,
    creado_en       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    actualizado_en  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_horario_fin CHECK (hora_fin > hora_inicio)
);

-- Tabla: citas (agenda)
CREATE TABLE citas (
    id              SERIAL PRIMARY KEY,
    estudiante_id   INTEGER NOT NULL REFERENCES estudiantes(id) ON DELETE RESTRICT,
    terapeuta_id    INTEGER NOT NULL REFERENCES terapeutas(id) ON DELETE RESTRICT,
    fecha           DATE NOT NULL,
    hora            TIME NOT NULL,
    tipo            VARCHAR(50) NOT NULL,
    estado          VARCHAR(20) NOT NULL DEFAULT 'programada',
    creado_en       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    actualizado_en  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_estado_cita CHECK (estado IN ('programada', 'confirmada', 'cancelada', 'completada', 'reprogramada'))
);

-- Tabla: asistencias
CREATE TABLE asistencias (
    id              SERIAL PRIMARY KEY,
    cita_id         INTEGER NOT NULL REFERENCES citas(id) ON DELETE CASCADE,
    asistio         BOOLEAN NOT NULL,
    observacion     TEXT,
    creado_en       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    actualizado_en  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla: valoraciones
CREATE TABLE valoraciones (
    id              SERIAL PRIMARY KEY,
    estudiante_id   INTEGER NOT NULL REFERENCES estudiantes(id) ON DELETE RESTRICT,
    terapeuta_id    INTEGER NOT NULL REFERENCES terapeutas(id) ON DELETE RESTRICT,
    fecha           DATE NOT NULL,
    observacion     TEXT,
    diagnostico     TEXT,
    creado_en       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    actualizado_en  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla: recuperaciones
CREATE TABLE recuperaciones (
    id                  SERIAL PRIMARY KEY,
    cita_original_id    INTEGER NOT NULL REFERENCES citas(id) ON DELETE RESTRICT,
    nueva_fecha         DATE NOT NULL,
    motivo              TEXT,
    estado              VARCHAR(20) NOT NULL DEFAULT 'pendiente',
    creado_en           TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    actualizado_en      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_estado_recuperacion CHECK (estado IN ('pendiente', 'programada', 'completada', 'cancelada'))
);

-- ============================================================
-- TABLAS DEL MÓDULO CONTABILIDAD
-- ============================================================

-- Tabla: facturas
CREATE TABLE facturas (
    id              SERIAL PRIMARY KEY,
    representante_id INTEGER NOT NULL REFERENCES representantes(id) ON DELETE RESTRICT,
    fecha           DATE NOT NULL,
    total           NUMERIC(10,2) NOT NULL CHECK (total >= 0),
    estado          VARCHAR(20) NOT NULL DEFAULT 'pendiente',
    metodo_pago     VARCHAR(50),
    creado_en       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    actualizado_en  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_estado_factura CHECK (estado IN ('pendiente', 'pagada', 'cancelada', 'anulada'))
);

-- Tabla: pagos
CREATE TABLE pagos (
    id              SERIAL PRIMARY KEY,
    factura_id      INTEGER NOT NULL REFERENCES facturas(id) ON DELETE RESTRICT,
    representante_id INTEGER NOT NULL REFERENCES representantes(id) ON DELETE RESTRICT,
    fecha           DATE NOT NULL,
    monto           NUMERIC(10,2) NOT NULL CHECK (monto >= 0),
    metodo_pago     VARCHAR(50),
    observacion     TEXT,
    creado_en       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    actualizado_en  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla: pagos_empleados
CREATE TABLE pagos_empleados (
    id              SERIAL PRIMARY KEY,
    terapeuta_id    INTEGER NOT NULL REFERENCES terapeutas(id) ON DELETE RESTRICT,
    fecha           DATE NOT NULL,
    sueldo          NUMERIC(10,2) NOT NULL CHECK (sueldo >= 0),
    descuento       NUMERIC(10,2) NOT NULL DEFAULT 0 CHECK (descuento >= 0),
    total           NUMERIC(10,2) NOT NULL CHECK (total >= 0),
    creado_en       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    actualizado_en  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla: inventario
CREATE TABLE inventario (
    id              SERIAL PRIMARY KEY,
    nombre          VARCHAR(100) NOT NULL,
    descripcion     TEXT,
    stock           INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    stock_minimo    INTEGER NOT NULL DEFAULT 0 CHECK (stock_minimo >= 0),
    creado_en       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    actualizado_en  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tabla: movimientos_inventario
CREATE TABLE movimientos_inventario (
    id              SERIAL PRIMARY KEY,
    inventario_id   INTEGER NOT NULL REFERENCES inventario(id) ON DELETE RESTRICT,
    tipo            VARCHAR(20) NOT NULL,
    cantidad        INTEGER NOT NULL CHECK (cantidad > 0),
    fecha           DATE NOT NULL,
    creado_en       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    actualizado_en  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_tipo_movimiento CHECK (tipo IN ('entrada', 'salida'))
);

-- Tabla: contabilidad
CREATE TABLE contabilidad (
    id              SERIAL PRIMARY KEY,
    fecha           DATE NOT NULL,
    concepto        VARCHAR(200) NOT NULL,
    tipo            VARCHAR(20) NOT NULL,
    valor           NUMERIC(10,2) NOT NULL CHECK (valor >= 0),
    creado_en       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    actualizado_en  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_tipo_movimiento_fin CHECK (tipo IN ('ingreso', 'gasto'))
);

-- ============================================================
-- ÍNDICES ADICIONALES
-- ============================================================
CREATE INDEX idx_estudiantes_representante ON estudiantes(representante_id);
CREATE INDEX idx_estudiantes_plan ON estudiantes(plan_id);
CREATE INDEX idx_citas_estudiante ON citas(estudiante_id);
CREATE INDEX idx_citas_terapeuta ON citas(terapeuta_id);
CREATE INDEX idx_citas_fecha ON citas(fecha);
CREATE INDEX idx_horarios_terapeuta ON horarios(terapeuta_id);
CREATE INDEX idx_valoraciones_estudiante ON valoraciones(estudiante_id);
CREATE INDEX idx_valoraciones_terapeuta ON valoraciones(terapeuta_id);
CREATE INDEX idx_facturas_representante ON facturas(representante_id);
CREATE INDEX idx_pagos_factura ON pagos(factura_id);
CREATE INDEX idx_pagos_empleados_terapeuta ON pagos_empleados(terapeuta_id);
CREATE INDEX idx_movimientos_inventario ON movimientos_inventario(inventario_id);
CREATE INDEX idx_contabilidad_fecha ON contabilidad(fecha);
CREATE INDEX idx_contabilidad_tipo ON contabilidad(tipo);

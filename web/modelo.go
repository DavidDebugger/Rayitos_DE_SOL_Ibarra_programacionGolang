package web

// modelo.go define las estructuras de datos del sistema Rayitos de Sol.
// Cada entidad tiene su propio ID numerico para identificarla de forma unica.
// Los nombres de campos y comentarios usan español sin acentos.

// Representante es el adulto responsable de un estudiante.
type Representante struct {
	ID           int
	Nombre       string
	Telefono     string
	Correo       string
	Relacion     string
	EstudianteID int
}

// Terapeuta es el profesional que atiende a los estudiantes.
type Terapeuta struct {
	ID           int
	Nombre       string
	Especialidad string
	Telefono     string
	Correo       string
}

// Plan es el paquete de sesiones que puede contratar un estudiante.
type Plan struct {
	ID       int
	Nombre   string
	Precio   float64
	Sesiones int
}

// Estudiante es el nino o joven que recibe la terapia.
type Estudiante struct {
	ID              int
	Nombre          string
	Edad            int
	Diagnostico     string
	RepresentanteID int
}

// Cita es una sesion agendada entre un estudiante y un terapeuta.
type Cita struct {
	ID           int
	EstudianteID int
	TerapeutaID  int
	Fecha        string
	Hora         string
	Estado       string
}

// Factura es el cobro generado por un plan para un estudiante.
type Factura struct {
	ID           int
	EstudianteID int
	PlanID       int
	Monto        float64
	Fecha        string
	Estado       string
}

// Pago registra un abono hecho sobre una factura.
// Se incluye en el modelo pero no tiene CRUD propio en el portal.
type Pago struct {
	PagoID    int
	FacturaID int
	Monto     float64
	Fecha     string
	Metodo    string
}

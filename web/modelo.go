/*
 modelo.go
 Define el modelo de datos: las entidades del centro de terapias "Rayitos de
 Sol" y como se relacionan entre si.
 RELACIONES (los campos que terminan en "_id" son llaves foraneas que guardan
 el numero de ID de la otra entidad):
   Representante 1 --- * Estudiante   (un tutor tiene muchos alumnos)
   Representante 1 --- * Factura      (un tutor recibe muchas facturas)
   Plan          1 --- * Estudiante   (un plan lo usan muchos alumnos)
   Estudiante    1 --- * Cita         (un alumno tiene muchas citas)
   Terapeuta     1 --- * Cita         (un terapeuta atiende muchas citas)
   Factura       1 --- * Pago          (una factura puede tener pagos)
 Nota: la entidad Pago se define aqui para completar el diseno, pero en esta
 version no tiene servicios ni formularios propios (no se expone en el CRUD).
 Las etiquetas `json:"..."` indican el nombre de cada campo al convertir la
 struct a JSON en los servicios web.
 */

package web

/*
 ESTADOS DE UNA CITA
 Constantes con los valores permitidos para el estado de una cita. Usar
 nombres fijos evita escribir texto suelto (y equivocarse) en varias partes.
 */
const (
	EstadoProgramada   = "Programada"
	EstadoConfirmada   = "Confirmada"
	EstadoCompletada   = "Completada"
	EstadoReprogramada = "Reprogramada"
)

/*
 ESTADOS DE UNA FACTURA
 Valores permitidos para el estado de una factura.
 */
const (
	EstadoPendiente = "Pendiente"
	EstadoPagada    = "Pagada"
	EstadoAnulada   = "Anulada"
)

// Representante: la persona (tutor o padre) que paga las terapias del alumno.
type Representante struct {
	ID        int    `json:"id"`        // numero unico que identifica al representante
	Nombres   string `json:"nombres"`   // primer nombre o nombres
	Apellidos string `json:"apellidos"` // apellido o apellidos
	Cedula    string `json:"cedula"`    // documento de identidad
	Telefono  string `json:"telefono"`  // numero de contacto
}

// Terapeuta: el profesional que atiende al estudiante en cada sesion.
type Terapeuta struct {
	ID        int    `json:"id"`
	Nombres   string `json:"nombres"`
	Apellidos string `json:"apellidos"`
	Cedula    string `json:"cedula"`
	Telefono  string `json:"telefono"`
}

// Estudiante: el nino o nina que recibe la terapia.
// RepresentanteID y PlanID lo conectan con su tutor y con su plan contratado.
type Estudiante struct {
	ID               int    `json:"id"`
	Nombres          string `json:"nombres"`
	Apellidos        string `json:"apellidos"`
	Cedula           string `json:"cedula"`
	Telefono         string `json:"telefono"`
	RepresentanteID  int    `json:"representante_id"` // a que tutor pertenece
	PlanID           int    `json:"plan_id"`          // que plan tiene contratado
}

// Plan: el paquete de sesiones contratado (por ejemplo "Basico" de 60 minutos).
type Plan struct {
	ID          int    `json:"id"`
	Nombre      string `json:"nombre"`      // nombre del plan
	Descripcion string `json:"descripcion"` // detalle breve
	Duracion    int    `json:"duracion"`    // duracion de la sesion en minutos
}

// Cita: una sesion de terapia agendada.
// PlanID se copia automaticamente del estudiante cuando se crea la cita.
type Cita struct {
	ID           int    `json:"id"`
	EstudianteID int    `json:"estudiante_id"` // alumno que asiste
	TerapeutaID  int    `json:"terapeuta_id"`  // profesional que atiende
	PlanID       int    `json:"plan_id"`       // plan que se aplica (se copia del alumno)
	Fecha        string `json:"fecha"`         // fecha de la cita (AAAA-MM-DD)
	Hora         string `json:"hora"`          // hora de la cita (HH:MM)
	Estado       string `json:"estado"`        // uno de los estados de cita definidos arriba
}

// Factura: el cobro que se emite al representante por las terapias.
type Factura struct {
	ID              int    `json:"id"`
	RepresentanteID int    `json:"representante_id"` // a quien se le cobra
	Fecha           string `json:"fecha"`            // fecha de emision
	Monto           int    `json:"monto"`            // valor a pagar
	Estado          string `json:"estado"`           // uno de los estados de factura
}

// Pago: un abono (pago parcial o total) hecho sobre una factura.
// En esta version la entidad se modela pero no se usa en el CRUD.
type Pago struct {
	ID        int    `json:"id"`
	FacturaID int    `json:"factura_id"` // factura a la que corresponde este pago
	Monto     int    `json:"monto"`      // cuanto se pago en este abono
	Fecha     string `json:"fecha"`      // fecha del pago
	Metodo    string `json:"metodo"`     // como se pago (efectivo, tarjeta, etc.)
}

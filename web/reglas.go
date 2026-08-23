package web

// reglas.go tiene las validaciones de cada entidad.
// Cada funcion regresa un error si algo obligatorio esta vacio.

import "errors"

func validarRepresentante(r Representante) error {
	if r.Nombre == "" {
		return errors.New("el nombre del representante es obligatorio")
	}
	if r.Telefono == "" {
		return errors.New("el telefono del representante es obligatorio")
	}
	if r.Relacion == "" {
		return errors.New("la relacion con el estudiante es obligatoria")
	}
	return nil
}

func validarTerapeuta(t Terapeuta) error {
	if t.Nombre == "" {
		return errors.New("el nombre del terapeuta es obligatorio")
	}
	if t.Especialidad == "" {
		return errors.New("la especialidad del terapeuta es obligatoria")
	}
	if t.Telefono == "" {
		return errors.New("el telefono del terapeuta es obligatorio")
	}
	return nil
}

func validarPlan(p Plan) error {
	if p.Nombre == "" {
		return errors.New("el nombre del plan es obligatorio")
	}
	if p.Precio <= 0 {
		return errors.New("el precio del plan debe ser mayor que cero")
	}
	if p.Sesiones <= 0 {
		return errors.New("el plan debe tener al menos una sesion")
	}
	return nil
}

func validarEstudiante(e Estudiante) error {
	if e.Nombre == "" {
		return errors.New("el nombre del estudiante es obligatorio")
	}
	if e.Edad <= 0 {
		return errors.New("la edad del estudiante debe ser mayor que cero")
	}
	if e.Diagnostico == "" {
		return errors.New("el diagnostico del estudiante es obligatorio")
	}
	return nil
}

func validarCita(c Cita) error {
	if c.EstudianteID <= 0 {
		return errors.New("la cita debe tener un estudiante valido")
	}
	if c.TerapeutaID <= 0 {
		return errors.New("la cita debe tener un terapeuta valido")
	}
	if c.Fecha == "" {
		return errors.New("la fecha de la cita es obligatoria")
	}
	if c.Hora == "" {
		return errors.New("la hora de la cita es obligatoria")
	}
	return nil
}

func validarFactura(f Factura) error {
	if f.EstudianteID <= 0 {
		return errors.New("la factura debe tener un estudiante valido")
	}
	if f.PlanID <= 0 {
		return errors.New("la factura debe tener un plan valido")
	}
	if f.Monto <= 0 {
		return errors.New("el monto de la factura debe ser mayor que cero")
	}
	if f.Fecha == "" {
		return errors.New("la fecha de la factura es obligatoria")
	}
	return nil
}

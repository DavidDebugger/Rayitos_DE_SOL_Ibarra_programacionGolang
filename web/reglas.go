/*
 reglas.go
 Aqui vive la logica de negocio (las reglas del sistema): que datos son
 obligatorios, que relaciones deben cumplirse, etc.
 POR QUE ESTA SEPARADO:
   Tanto la API JSON como el portal web necesitan crear, actualizar y borrar
   las mismas entidades con las mismas validaciones. Al tener funciones por
   entidad, ambos (API y portal) las reutilizan y la regla existe en un solo
   lugar. Es una buena practica aunque el proyecto sea pequeño.
 Todas devuelven un error. Si el error es nil, todo salio bien y la entidad
 ya esta guardada o borrada. Si no, el error describe que fallo.
 CONVENCION DE PARAMETROS:
   - crearX(a, *T): recibe la entidad por PUNTERO para asignarle el nuevo id
     dentro de la funcion y que quien llama lo vea en su propia variable.
   - actualizarX(a, id, T): recibe el id por separado (porque viene de la
     ruta) y la entidad por valor, pues el id ya se conoce de antemano.
   - puedeBorrar(a, tipo, id): devuelve (false, mensaje) si el registro esta
     en uso y no se debe borrar.
 */

package web

import "errors"

// ----- Representante -----

func crearRepresentante(a *Almacen, e *Representante) error {
	if e.Nombres == "" || e.Cedula == "" {
		return errors.New("nombres y cedula son obligatorios")
	}
	e.ID = a.nuevoID("representante")
	a.Guardar("representante", e.ID, *e)
	return nil
}

func actualizarRepresentante(a *Almacen, id int, e Representante) error {
	if _, ok := a.Obtener("representante", id); !ok {
		return errors.New("el representante no existe")
	}
	if e.Nombres == "" || e.Cedula == "" {
		return errors.New("nombres y cedula son obligatorios")
	}
	e.ID = id
	a.Guardar("representante", id, e)
	return nil
}

// ----- Terapeuta -----

func crearTerapeuta(a *Almacen, e *Terapeuta) error {
	if e.Nombres == "" || e.Cedula == "" {
		return errors.New("nombres y cedula son obligatorios")
	}
	e.ID = a.nuevoID("terapeuta")
	a.Guardar("terapeuta", e.ID, *e)
	return nil
}

func actualizarTerapeuta(a *Almacen, id int, e Terapeuta) error {
	if _, ok := a.Obtener("terapeuta", id); !ok {
		return errors.New("el terapeuta no existe")
	}
	if e.Nombres == "" || e.Cedula == "" {
		return errors.New("nombres y cedula son obligatorios")
	}
	e.ID = id
	a.Guardar("terapeuta", id, e)
	return nil
}

// ----- Plan -----

func crearPlan(a *Almacen, p *Plan) error {
	if p.Nombre == "" {
		return errors.New("el nombre es obligatorio")
	}
	p.ID = a.nuevoID("plan")
	a.Guardar("plan", p.ID, *p)
	return nil
}

func actualizarPlan(a *Almacen, id int, p Plan) error {
	if _, ok := a.Obtener("plan", id); !ok {
		return errors.New("el plan no existe")
	}
	if p.Nombre == "" {
		return errors.New("el nombre es obligatorio")
	}
	p.ID = id
	a.Guardar("plan", id, p)
	return nil
}

// ----- Estudiante -----

func crearEstudiante(a *Almacen, e *Estudiante) error {
	if e.Nombres == "" || e.Cedula == "" {
		return errors.New("nombres y cedula son obligatorios")
	}
	if !a.Existe("representante", e.RepresentanteID) || !a.Existe("plan", e.PlanID) {
		return errors.New("representante o plan no existen")
	}
	e.ID = a.nuevoID("estudiante")
	a.Guardar("estudiante", e.ID, *e)
	return nil
}

func actualizarEstudiante(a *Almacen, id int, e Estudiante) error {
	if _, ok := a.Obtener("estudiante", id); !ok {
		return errors.New("el estudiante no existe")
	}
	if e.Nombres == "" || e.Cedula == "" {
		return errors.New("nombres y cedula son obligatorios")
	}
	if !a.Existe("representante", e.RepresentanteID) || !a.Existe("plan", e.PlanID) {
		return errors.New("representante o plan no existen")
	}
	e.ID = id
	a.Guardar("estudiante", id, e)
	return nil
}

// ----- Cita -----

func crearCita(a *Almacen, c *Cita) error {
	if !a.Existe("estudiante", c.EstudianteID) || !a.Existe("terapeuta", c.TerapeutaID) {
		return errors.New("debe elegir un estudiante y un terapeuta existentes")
	}
	if v, ok := a.Obtener("estudiante", c.EstudianteID); ok {
		c.PlanID = v.(Estudiante).PlanID
	}
	if c.Estado != "" && c.Estado != EstadoProgramada && c.Estado != EstadoConfirmada && c.Estado != EstadoCompletada && c.Estado != EstadoReprogramada {
		return errors.New("estado de cita no valido")
	}
	c.ID = a.nuevoID("cita")
	a.Guardar("cita", c.ID, *c)
	return nil
}

func actualizarCita(a *Almacen, id int, c Cita) error {
	if _, ok := a.Obtener("cita", id); !ok {
		return errors.New("la cita no existe")
	}
	if !a.Existe("estudiante", c.EstudianteID) || !a.Existe("terapeuta", c.TerapeutaID) {
		return errors.New("debe elegir un estudiante y un terapeuta existentes")
	}
	if v, ok := a.Obtener("estudiante", c.EstudianteID); ok {
		c.PlanID = v.(Estudiante).PlanID
	}
	if c.Estado != "" && c.Estado != EstadoProgramada && c.Estado != EstadoConfirmada && c.Estado != EstadoCompletada && c.Estado != EstadoReprogramada {
		return errors.New("estado de cita no valido")
	}
	c.ID = id
	a.Guardar("cita", id, c)
	return nil
}

// ----- Factura -----

func crearFactura(a *Almacen, f *Factura) error {
	if !a.Existe("representante", f.RepresentanteID) {
		return errors.New("representante no existe")
	}
	if f.Monto <= 0 {
		return errors.New("el monto debe ser mayor a 0")
	}
	f.ID = a.nuevoID("factura")
	a.Guardar("factura", f.ID, *f)
	return nil
}

func actualizarFactura(a *Almacen, id int, f Factura) error {
	if _, ok := a.Obtener("factura", id); !ok {
		return errors.New("la factura no existe")
	}
	if !a.Existe("representante", f.RepresentanteID) {
		return errors.New("representante no existe")
	}
	if f.Monto <= 0 {
		return errors.New("el monto debe ser mayor a 0")
	}
	f.ID = id
	a.Guardar("factura", id, f)
	return nil
}

/*
 puedeBorrar indica si un registro se puede eliminar sin dejar referencias
 rotas. Devuelve (false, mensaje) si esta en uso por otra entidad.
 */
func puedeBorrar(a *Almacen, tipo string, id int) (bool, string) {
	switch tipo {
	case "representante":
		for _, v := range a.Listar("estudiante") {
			if v.(Estudiante).RepresentanteID == id {
				return false, "hay estudiantes con este representante"
			}
		}
	case "plan":
		for _, v := range a.Listar("estudiante") {
			if v.(Estudiante).PlanID == id {
				return false, "hay estudiantes con este plan"
			}
		}
	case "estudiante":
		for _, v := range a.Listar("cita") {
			if v.(Cita).EstudianteID == id {
				return false, "hay citas de este estudiante"
			}
		}
	case "terapeuta":
		for _, v := range a.Listar("cita") {
			if v.(Cita).TerapeutaID == id {
				return false, "hay citas de este terapeuta"
			}
		}
	}
	return true, ""
}

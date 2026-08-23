package web

// portal.go tiene la logica de las paginas del portal (sin HTML ni CSS embebidos).
// El HTML esta en templates/portal.html y el CSS en templates/estilos.css.

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
)

// campoForm describe un campo del formulario de una entidad.
type campoForm struct {
	Nombre   string
	Etiqueta string
	Tipo     string // text, email, number, date, time
}

// camposForm lista los campos editables de cada entidad (sin el ID).
var camposForm = map[string][]campoForm{
	"representantes": {
		{"Nombre", "Nombre", "text"},
		{"Cedula", "Cedula", "text"},
		{"Telefono", "Telefono", "text"},
		{"Correo", "Correo", "email"},
		{"Relacion", "Relacion con el estudiante", "text"},
		{"EstudianteID", "ID del estudiante", "number"},
	},
	"terapeutas": {
		{"Nombre", "Nombre", "text"},
		{"Especialidad", "Especialidad", "text"},
		{"Telefono", "Telefono", "text"},
		{"Correo", "Correo", "email"},
	},
	"planes": {
		{"Nombre", "Nombre", "text"},
		{"Precio", "Precio", "number"},
		{"Sesiones", "Sesiones", "number"},
	},
	"estudiantes": {
		{"Nombre", "Nombre", "text"},
		{"Edad", "Edad", "number"},
		{"Diagnostico", "Diagnostico", "text"},
		{"RepresentanteID", "ID del representante", "number"},
	},
	"citas": {
		{"EstudianteID", "ID del estudiante", "number"},
		{"TerapeutaID", "ID del terapeuta", "number"},
		{"Fecha", "Fecha", "date"},
		{"Hora", "Hora", "time"},
		{"Estado", "Estado", "text"},
	},
	"facturas": {
		{"EstudianteID", "ID del estudiante", "number"},
		{"PlanID", "ID del plan", "number"},
		{"Monto", "Monto", "number"},
		{"Fecha", "Fecha", "date"},
		{"Estado", "Estado", "text"},
	},
}

// columnasEntidad define el orden de las columnas en la tabla de listado.
var columnasEntidad = map[string][]string{
	"representantes": {"ID", "Nombre", "Cedula", "Telefono", "Correo", "Relacion", "EstudianteID"},
	"terapeutas":     {"ID", "Nombre", "Especialidad", "Telefono", "Correo"},
	"planes":         {"ID", "Nombre", "Precio", "Sesiones"},
	"estudiantes":    {"ID", "Nombre", "Edad", "Diagnostico", "RepresentanteID"},
	"citas":          {"ID", "EstudianteID", "TerapeutaID", "Fecha", "Hora", "Estado"},
	"facturas":       {"ID", "EstudianteID", "PlanID", "Monto", "Fecha", "Estado"},
}

// etiquetas es el nombre singular de cada entidad para los titulos.
var etiquetas = map[string]string{
	"representantes": "Representante",
	"terapeutas":     "Terapeuta",
	"planes":         "Plan",
	"estudiantes":    "Estudiante",
	"citas":          "Cita",
	"facturas":       "Factura",
	"servicios":      "Servicios Web",
}

// listaComoMapas convierte la lista de una entidad en mapas para la plantilla.
func listaComoMapas(entidad string) []map[string]any {
	var datos any
	switch entidad {
	case "representantes":
		datos = obtenerRepresentantes()
	case "terapeutas":
		datos = obtenerTerapeutas()
	case "planes":
		datos = obtenerPlanes()
	case "estudiantes":
		datos = obtenerEstudiantes()
	case "citas":
		datos = obtenerCitas()
	case "facturas":
		datos = obtenerFacturas()
	}
	if datos == nil {
		return nil
	}
	b, _ := json.Marshal(datos)
	var salida []map[string]any
	json.Unmarshal(b, &salida)
	return salida
}

// itemComoMapa convierte un solo elemento en mapa para rellenar el formulario.
func itemComoMapa(entidad string, id int) (map[string]any, bool) {
	var datos any
	var ok bool
	switch entidad {
	case "representantes":
		var v Representante
		v, ok = obtenerRepresentante(id)
		datos = v
	case "terapeutas":
		var v Terapeuta
		v, ok = obtenerTerapeuta(id)
		datos = v
	case "planes":
		var v Plan
		v, ok = obtenerPlan(id)
		datos = v
	case "estudiantes":
		var v Estudiante
		v, ok = obtenerEstudiante(id)
		datos = v
	case "citas":
		var v Cita
		v, ok = obtenerCita(id)
		datos = v
	case "facturas":
		var v Factura
		v, ok = obtenerFactura(id)
		datos = v
	}
	if !ok {
		return nil, false
	}
	b, _ := json.Marshal(datos)
	var salida map[string]any
	json.Unmarshal(b, &salida)
	return salida, true
}

// valoresForm reconstruye un mapa con lo que escribio el usuario (para rellenar en error).
func valoresForm(r *http.Request, entidad string) map[string]any {
	m := map[string]any{}
	for _, c := range camposForm[entidad] {
		m[c.Nombre] = r.FormValue(c.Nombre)
	}
	return m
}

// portalInicio es la pagina principal del portal.
func portalInicio(w http.ResponseWriter, r *http.Request) {
	servirPlantilla(w, r, paginaDatos{Vista: "inicio", Titulo: "Inicio", Menu: menuPortal()})
}

// portalListar muestra la tabla de una entidad.
func portalListar(w http.ResponseWriter, r *http.Request, entidad string) {
	label := etiquetas[entidad]
	servirPlantilla(w, r, paginaDatos{
		Vista:      "listar",
		Titulo:     label + "s",
		MenuActivo: entidad,
		Menu:       menuPortal(),
		Label:      label,
		Columnas:   columnasEntidad[entidad],
		Filas:      listaComoMapas(entidad),
		Colspan:    len(columnasEntidad[entidad]) + 1,
	})
}

// portalCrear muestra el formulario vacio.
func portalCrear(w http.ResponseWriter, r *http.Request, entidad string) {
	label := etiquetas[entidad]
	servirPlantilla(w, r, paginaDatos{
		Vista:      "formulario",
		Titulo:     "Crear " + label,
		MenuActivo: entidad,
		Menu:       menuPortal(),
		Label:      label,
		Accion:     "/portal/" + entidad,
		Campos:     camposForm[entidad],
		Valores:    map[string]any{},
	})
}

// portalEditar muestra el formulario con los datos actuales.
func portalEditar(w http.ResponseWriter, r *http.Request, entidad string, id int) {
	m, ok := itemComoMapa(entidad, id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	label := etiquetas[entidad]
	servirPlantilla(w, r, paginaDatos{
		Vista:      "formulario",
		Titulo:     "Editar " + label,
		MenuActivo: entidad,
		Menu:       menuPortal(),
		Label:      label,
		Accion:     fmt.Sprintf("/portal/%s/%d", entidad, id),
		Campos:     camposForm[entidad],
		Valores:    m,
	})
}

// portalGuardar recibe el POST del formulario y crea o actualiza.
func portalGuardar(w http.ResponseWriter, r *http.Request, entidad string, id int) {
	r.ParseForm()
	accion := "/portal/" + entidad
	if id > 0 {
		accion = fmt.Sprintf("/portal/%s/%d", entidad, id)
	}
	var err error
	switch entidad {
	case "representantes":
		rcp := formARepresentante(r)
		if id == 0 {
			if err = validarRepresentante(rcp); err == nil {
				insertarRepresentante(rcp)
			}
		} else {
			if _, ok := obtenerRepresentante(id); !ok {
				http.NotFound(w, r)
				return
			}
			if err = validarRepresentante(rcp); err == nil {
				actualizarRepresentante(id, rcp)
			}
		}
	case "terapeutas":
		ter := formATerapeuta(r)
		if id == 0 {
			if err = validarTerapeuta(ter); err == nil {
				insertarTerapeuta(ter)
			}
		} else {
			if _, ok := obtenerTerapeuta(id); !ok {
				http.NotFound(w, r)
				return
			}
			if err = validarTerapeuta(ter); err == nil {
				actualizarTerapeuta(id, ter)
			}
		}
	case "planes":
		pla := formAPlan(r)
		if id == 0 {
			if err = validarPlan(pla); err == nil {
				insertarPlan(pla)
			}
		} else {
			if _, ok := obtenerPlan(id); !ok {
				http.NotFound(w, r)
				return
			}
			if err = validarPlan(pla); err == nil {
				actualizarPlan(id, pla)
			}
		}
	case "estudiantes":
		est := formAEstudiante(r)
		if id == 0 {
			if err = validarEstudiante(est); err == nil {
				insertarEstudiante(est)
			}
		} else {
			if _, ok := obtenerEstudiante(id); !ok {
				http.NotFound(w, r)
				return
			}
			if err = validarEstudiante(est); err == nil {
				actualizarEstudiante(id, est)
			}
		}
	case "citas":
		cit := formACita(r)
		if id == 0 {
			if err = validarCita(cit); err == nil {
				insertarCita(cit)
			}
		} else {
			if _, ok := obtenerCita(id); !ok {
				http.NotFound(w, r)
				return
			}
			if err = validarCita(cit); err == nil {
				actualizarCita(id, cit)
			}
		}
	case "facturas":
		fac := formAFactura(r)
		if id == 0 {
			if err = validarFactura(fac); err == nil {
				insertarFactura(fac)
			}
		} else {
			if _, ok := obtenerFactura(id); !ok {
				http.NotFound(w, r)
				return
			}
			if err = validarFactura(fac); err == nil {
				actualizarFactura(id, fac)
			}
		}
	}

	if err != nil {
		label := etiquetas[entidad]
		servirPlantilla(w, r, paginaDatos{
			Vista:      "formulario",
			Titulo:     "Error al guardar " + label,
			MenuActivo: entidad,
			Menu:       menuPortal(),
			Label:      label,
			Accion:     accion,
			Campos:     camposForm[entidad],
			Valores:    valoresForm(r, entidad),
			Error:      err.Error(),
		})
		return
	}
	http.Redirect(w, r, "/portal/"+entidad, http.StatusSeeOther)
}

// portalVer muestra el JSON de un registro dentro del portal (frontend).
func portalVer(w http.ResponseWriter, r *http.Request, entidad string, id int) {
	var datos any
	var ok bool
	switch entidad {
	case "representantes":
		var v Representante
		v, ok = obtenerRepresentante(id)
		datos = v
	case "terapeutas":
		var v Terapeuta
		v, ok = obtenerTerapeuta(id)
		datos = v
	case "planes":
		var v Plan
		v, ok = obtenerPlan(id)
		datos = v
	case "estudiantes":
		var v Estudiante
		v, ok = obtenerEstudiante(id)
		datos = v
	case "citas":
		var v Cita
		v, ok = obtenerCita(id)
		datos = v
	case "facturas":
		var v Factura
		v, ok = obtenerFactura(id)
		datos = v
	}
	if !ok {
		http.NotFound(w, r)
		return
	}
	b, _ := json.MarshalIndent(datos, "", "  ")
	label := etiquetas[entidad]
	servirPlantilla(w, r, paginaDatos{
		Vista:      "ver",
		Titulo:     "Ver " + label,
		MenuActivo: entidad,
		Menu:       menuPortal(),
		Label:      label,
		JSON:       string(b),
	})
}

// portalBorrar recibe el POST de borrar y redirige al listado.
func portalBorrar(w http.ResponseWriter, r *http.Request, entidad string, id int) {
	switch entidad {
	case "representantes":
		eliminarRepresentante(id)
	case "terapeutas":
		eliminarTerapeuta(id)
	case "planes":
		eliminarPlan(id)
	case "estudiantes":
		eliminarEstudiante(id)
	case "citas":
		eliminarCita(id)
	case "facturas":
		eliminarFactura(id)
	}
	http.Redirect(w, r, "/portal/"+entidad, http.StatusSeeOther)
}

// portalServicios muestra la vista de los servicios web (la API) del sistema.
// portalServicios muestra los 9 servicios reales de la API y, para cada uno,
// el JSON que la API serializa. En los GET se ejecuta la peticion de solo
// lectura contra la API para mostrar la respuesta real; en POST/PUT se muestra
// el cuerpo JSON de ejemplo que espera la API. Asi se nota la serializacion.
func portalServicios(w http.ResponseWriter, r *http.Request) {
	// jsonLectura ejecuta una peticion GET de solo lectura contra la API
	// (sin mutar el almacen) y devuelve el cuerpo serializado en JSON.
	jsonLectura := func(ruta string) string {
		req := httptest.NewRequest(http.MethodGet, ruta, nil)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		enrutar(rec, req)
		return rec.Body.String()
	}
	serv := []servicioWeb{
		{Metodo: "GET", Ruta: "/api", Desc: "Indice: entidades y servicios disponibles", Ejemplo: jsonLectura("/api")},
		{Metodo: "GET", Ruta: "/api/estadisticas", Desc: "Totales de registros por entidad", Ejemplo: jsonLectura("/api/estadisticas")},
		{Metodo: "GET", Ruta: "/api/planes", Desc: "Listar todos los registros de una entidad", Ejemplo: jsonLectura("/api/planes")},
		{Metodo: "POST", Ruta: "/api/planes", Desc: "Crear un registro de una entidad", Cuerpo: `{"Nombre":"Plan de prueba","Precio":50,"Sesiones":4}`},
		{Metodo: "GET", Ruta: "/api/estudiantes/contar", Desc: "Cantidad de registros de una entidad", Ejemplo: jsonLectura("/api/estudiantes/contar")},
		{Metodo: "GET", Ruta: "/api/facturas/buscar?campo=Estado&valor=Pendiente", Desc: "Filtrar por campo y valor (query)", Ejemplo: jsonLectura("/api/facturas/buscar?campo=Estado&valor=Pendiente")},
		{Metodo: "GET", Ruta: "/api/citas/1", Desc: "Ver un registro por su ID", Ejemplo: jsonLectura("/api/citas/1")},
		{Metodo: "PUT", Ruta: "/api/citas/1", Desc: "Actualizar un registro por su ID", Cuerpo: `{"Estado":"Realizada"}`},
		{Metodo: "DELETE", Ruta: "/api/citas/1", Desc: "Borrar un registro por su ID"},
	}
	servirPlantilla(w, r, paginaDatos{
		Vista:      "servicios",
		Titulo:     "Servicios Web",
		MenuActivo: "servicios",
		Menu:       menuPortal(),
		Servicios:  serv,
	})
}

// portalProbar ejecuta una peticion a la API a partir de parametros en la URL
// y muestra el resultado. Asi los servicios POST/PUT/DELETE se prueban con
// un hipervinculo (el navegador solo envia GET) sin usar JavaScript.
func portalProbar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "metodo no permitido", http.StatusMethodNotAllowed)
		return
	}
	q := r.URL.Query()
	metodo := q.Get("metodo")
	ruta := q.Get("ruta")
	cuerpo := q.Get("cuerpo")
	req := httptest.NewRequest(metodo, ruta, strings.NewReader(cuerpo))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	enrutar(rec, req)
	resultado := fmt.Sprintf("HTTP %d\n\n%s", rec.Code, rec.Body.String())
	servirPlantilla(w, r, paginaDatos{
		Vista:      "probar",
		Titulo:     "Resultado de la API",
		MenuActivo: "servicios",
		Menu:       menuPortal(),
		JSON:       resultado,
	})
}

// ---- Lectura de formularios (convierte el texto del formulario a cada struct) ----

func formARepresentante(r *http.Request) Representante {
	return Representante{
		Nombre:       r.FormValue("Nombre"),
		Cedula:       r.FormValue("Cedula"),
		Telefono:     r.FormValue("Telefono"),
		Correo:       r.FormValue("Correo"),
		Relacion:     r.FormValue("Relacion"),
		EstudianteID: aEntero(r.FormValue("EstudianteID")),
	}
}

func formATerapeuta(r *http.Request) Terapeuta {
	return Terapeuta{
		Nombre:       r.FormValue("Nombre"),
		Especialidad: r.FormValue("Especialidad"),
		Telefono:     r.FormValue("Telefono"),
		Correo:       r.FormValue("Correo"),
	}
}

func formAPlan(r *http.Request) Plan {
	return Plan{
		Nombre:   r.FormValue("Nombre"),
		Precio:   aFlotante(r.FormValue("Precio")),
		Sesiones: aEntero(r.FormValue("Sesiones")),
	}
}

func formAEstudiante(r *http.Request) Estudiante {
	return Estudiante{
		Nombre:          r.FormValue("Nombre"),
		Edad:            aEntero(r.FormValue("Edad")),
		Diagnostico:     r.FormValue("Diagnostico"),
		RepresentanteID: aEntero(r.FormValue("RepresentanteID")),
	}
}

func formACita(r *http.Request) Cita {
	return Cita{
		EstudianteID: aEntero(r.FormValue("EstudianteID")),
		TerapeutaID:  aEntero(r.FormValue("TerapeutaID")),
		Fecha:        r.FormValue("Fecha"),
		Hora:         r.FormValue("Hora"),
		Estado:       r.FormValue("Estado"),
	}
}

func formAFactura(r *http.Request) Factura {
	return Factura{
		EstudianteID: aEntero(r.FormValue("EstudianteID")),
		PlanID:       aEntero(r.FormValue("PlanID")),
		Monto:        aFlotante(r.FormValue("Monto")),
		Fecha:        r.FormValue("Fecha"),
		Estado:       r.FormValue("Estado"),
	}
}

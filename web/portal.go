package web

// portal.go genera las paginas web del portal (HTML del servidor, sin JavaScript).
// El usuario usa formularios normales para crear, editar y borrar.

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

// campoForm describe un campo del formulario de una entidad.
type campoForm struct {
	nombre   string
	etiqueta string
	tipo     string // text, email, number, date, time
}

// camposForm lista los campos editables de cada entidad (sin el ID).
var camposForm = map[string][]campoForm{
	"representantes": {
		{"Nombre", "Nombre", "text"},
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
	"representantes": {"ID", "Nombre", "Telefono", "Correo", "Relacion", "EstudianteID"},
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
}

// escapar evita inyeccion de HTML en los valores que se muestran.
func escapar(s string) string {
	r := strings.NewReplacer("<", "&lt;", ">", "&gt;", "&", "&amp;", "\"", "&quot;")
	return r.Replace(s)
}

// portalEstilos devuelve el CSS comun de todas las paginas.
func portalEstilos() string {
	return `
  * { box-sizing: border-box; }
  body { margin: 0; font-family: Arial, Helvetica, sans-serif; background: #f4f6f8; color: #222; }
  .barra { background: #ffb703; color: #222; padding: 12px 18px; display: flex; justify-content: space-between; align-items: center; }
  .barra a { color: #222; text-decoration: none; font-weight: bold; }
  .cuerpo { display: flex; min-height: calc(100vh - 48px); }
  .menu { width: 200px; background: #023047; padding: 14px; }
  .menu a { display: block; color: #fff; text-decoration: none; padding: 9px 10px; border-radius: 6px; margin-bottom: 4px; }
  .menu a:hover { background: #03506f; }
  .menu a.activo { background: #ffb703; color: #222; }
  .principal { flex: 1; padding: 22px; }
  h1 { margin-top: 0; }
  table { border-collapse: collapse; width: 100%; background: #fff; margin-bottom: 16px; }
  th, td { border: 1px solid #ddd; padding: 8px 10px; text-align: left; }
  th { background: #e9f1f5; }
  a.boton, button { background: #ffb703; color: #222; border: none; padding: 8px 12px; border-radius: 6px; text-decoration: none; cursor: pointer; font-weight: bold; }
  .acciones a, .acciones button { margin-right: 6px; font-size: 13px; background: #e9f1f5; padding: 4px 8px; border-radius: 5px; text-decoration: none; color: #023047; }
  label { display: block; margin-top: 10px; font-weight: bold; }
  input { width: 100%; padding: 8px; margin-top: 4px; border: 1px solid #bbb; border-radius: 6px; }
  .error { background: #ffe3e3; color: #a00; padding: 10px; border-radius: 6px; }
  form { max-width: 460px; background: #fff; padding: 16px; border-radius: 8px; }
`
}

// pagina envuelve el contenido en la plantilla general del portal.
func pagina(titulo, menuActivo, cuerpo string) string {
	menu := ""
	items := []struct {
		plural string
		label  string
	}{
		{"representantes", "Representantes"},
		{"terapeutas", "Terapeutas"},
		{"planes", "Planes"},
		{"estudiantes", "Estudiantes"},
		{"citas", "Citas"},
		{"facturas", "Facturas"},
	}
	for _, it := range items {
		cls := ""
		if it.plural == menuActivo {
			cls = " class=\"activo\""
		}
		menu += fmt.Sprintf("      <a href=\"/portal/%s\"%s>%s</a>\n", it.plural, cls, it.label)
	}
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="es">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Rayitos de Sol - %s</title>
<style>%s</style>
</head>
<body>
<div class="barra"><span>Rayitos de Sol</span><a href="/portal">Inicio</a></div>
<div class="cuerpo">
<nav class="menu">
%s</nav>
<main class="principal">
%s
</main>
</div>
</body>
</html>`, escapar(titulo), portalEstilos(), menu, cuerpo)
}

// portalInicio es la pagina principal del portal.
func portalInicio(w http.ResponseWriter, r *http.Request) {
	cuerpo := "<h1>Bienvenido al portal</h1>" +
		"<p>Elige una entidad en el menu para ver, crear, editar o borrar registros.</p>"
	io.WriteString(w, pagina("Inicio", "", cuerpo))
}

// listaComoMapas convierte la lista de una entidad en mapas para dibujar la tabla.
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
		m[c.nombre] = r.FormValue(c.nombre)
	}
	return m
}

// renderTabla arma la tabla HTML de listado con sus acciones por fila.
func renderTabla(entidad string, filas []map[string]any) string {
	cols := columnasEntidad[entidad]
	html := "<table><thead><tr>"
	for _, c := range cols {
		html += "<th>" + escapar(c) + "</th>"
	}
	html += "<th>Acciones</th></tr></thead><tbody>"
	if len(filas) == 0 {
		html += "<tr><td colspan=\"" + strconv.Itoa(len(cols)+1) + "\">No hay registros.</td></tr>"
	}
	for _, f := range filas {
		html += "<tr>"
		for _, c := range cols {
			html += "<td>" + escapar(fmt.Sprintf("%v", f[c])) + "</td>"
		}
		id := fmt.Sprintf("%v", f["ID"])
		acciones := "<a href=\"/portal/" + entidad + "/" + id + "/editar\">Editar</a>" +
			" <a href=\"/api/" + entidad + "/" + id + "\" target=\"_blank\">Ver</a>" +
			" <form method=\"post\" action=\"/portal/" + entidad + "/" + id + "/borrar\" " +
			"onsubmit=\"return confirm('¿Seguro que quieres borrar?')\" style=\"display:inline\">" +
			"<button type=\"submit\">Borrar</button></form>"
		html += "<td class=\"acciones\">" + acciones + "</td></tr>"
	}
	html += "</tbody></table>"
	return html
}

// portalListar muestra la tabla de una entidad con el boton de crear.
func portalListar(w http.ResponseWriter, r *http.Request, entidad string) {
	filas := listaComoMapas(entidad)
	label := etiquetas[entidad]
	boton := "<a class=\"boton\" href=\"/portal/" + entidad + "/nuevo\">Crear " + escapar(label) + "</a>"
	cuerpo := "<h1>" + escapar(label) + "s</h1>" + boton + renderTabla(entidad, filas)
	io.WriteString(w, pagina(label+"s", entidad, cuerpo))
}

// formulario arma el HTML del formulario para crear o editar.
func formulario(entidad string, datos map[string]any, id int) string {
	accion := "/portal/" + entidad
	if id > 0 {
		accion = fmt.Sprintf("/portal/%s/%d", entidad, id)
	}
	campos := ""
	for _, c := range camposForm[entidad] {
		valor := ""
		if datos != nil {
			if v, ok := datos[c.nombre]; ok {
				valor = fmt.Sprintf("%v", v)
			}
		}
		campos += "<label>" + escapar(c.etiqueta) + "</label>\n" +
			"<input type=\"" + c.tipo + "\" name=\"" + c.nombre + "\" value=\"" + escapar(valor) + "\">\n"
	}
	return "<form method=\"post\" action=\"" + accion + "\">\n" +
		campos +
		"<p><button type=\"submit\">Guardar</button> " +
		"<a href=\"/portal/" + entidad + "\">Cancelar</a></p>\n</form>"
}

// portalCrear muestra el formulario vacio.
func portalCrear(w http.ResponseWriter, r *http.Request, entidad string) {
	label := etiquetas[entidad]
	cuerpo := "<h1>Crear " + escapar(label) + "</h1>" + formulario(entidad, nil, 0)
	io.WriteString(w, pagina("Crear "+label, entidad, cuerpo))
}

// portalEditar muestra el formulario con los datos actuales.
func portalEditar(w http.ResponseWriter, r *http.Request, entidad string, id int) {
	m, ok := itemComoMapa(entidad, id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	label := etiquetas[entidad]
	cuerpo := "<h1>Editar " + escapar(label) + "</h1>" + formulario(entidad, m, id)
	io.WriteString(w, pagina("Editar "+label, entidad, cuerpo))
}

// portalGuardar recibe el POST del formulario y crea o actualiza.
func portalGuardar(w http.ResponseWriter, r *http.Request, entidad string, id int) {
	r.ParseForm()
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
		cuerpo := "<p class=\"error\">" + escapar(err.Error()) + "</p>" +
			formulario(entidad, valoresForm(r, entidad), id)
		io.WriteString(w, pagina("Error al guardar "+label, entidad, cuerpo))
		return
	}
	http.Redirect(w, r, "/portal/"+entidad, http.StatusSeeOther)
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

// ---- Lectura de formularios (convierte el texto del formulario a cada struct) ----

func formARepresentante(r *http.Request) Representante {
	return Representante{
		Nombre:       r.FormValue("Nombre"),
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

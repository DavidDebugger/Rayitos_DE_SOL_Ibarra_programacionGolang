/*
 portal.go
 Implementa el "portal web": una interfaz visual que el usuario maneja con
 clicks desde el navegador, sin tener que escribir JSON.
 QUE ES SSR (Server-Side Rendering):
   El servidor genera el HTML completo y lo envia al navegador, que solo lo
   muestra. No hay JavaScript en el cliente: los formularios se envian con
   POST y el servidor responde con una pagina nueva (o redirige).
 QUE HACE ESTE ARCHIVO:
   - "embed": incluye templates/portal.html dentro del ejecutable, asi no
     dependemos de llevar el archivo suelto al ejecutar.
   - "html/template": rellena la plantilla con los datos (tablas, opciones).
   - "net/http": lee los formularios del portal y responde.
   - "encoding/json": genera el JSON de ejemplo que se muestra en cada tarjeta.
 Operaciones disponibles en el portal:
   portalInicio      -> pagina principal (tablas + formularios de creacion).
   portalCrear       -> crea un registro desde el formulario (POST).
   portalEditar      -> muestra el formulario de edicion precargado (GET).
   portalActualizar  -> aplica los cambios de la edicion (POST).
   portalBorrar      -> elimina un registro (POST), si no deja referencias.
 Todas reutilizan las mismas reglas de negocio que la API (reglas.go).
 */

package web

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

/*
 plantillaPortal contiene el archivo templates/portal.html, incluido dentro
 del binario gracias a //go:embed. Asi la pagina viaja con el ejecutable.
 */
//
//go:embed templates/portal.html
var plantillaPortal embed.FS

/*
 opcion representa una entrada de un <select> (desplegable) del portal.
 Valor es el id que se envia al servidor; Texto es lo que ve el usuario.
 */
type opcion struct {
	Valor string
	Texto string
}

/*
 PortalData es la "bolsa" de datos que pasamos a la plantilla HTML.
 Contiene las listas a mostrar, las opciones de los desplegables y el JSON
 de ejemplo de cada servicio.
 */
type PortalData struct {
	Mensaje            string            // mensaje de exito o error tras crear
	Representantes     []Representante   // lista para la tabla de representantes
	Terapeutas         []Terapeuta       // lista para la tabla de terapeutas
	Planes             []Plan            // lista para la tabla de planes
	Estudiantes        []Estudiante      // lista para la tabla de estudiantes
	Citas              []Cita            // lista para la tabla de citas
	Facturas           []Factura         // lista para la tabla de facturas
	OptsRepresentante  []opcion          // opciones del <select> de representante
	OptsPlan           []opcion          // opciones del <select> de plan
	OptsEstudiante     []opcion          // opciones del <select> de estudiante
	OptsTerapeuta      []opcion          // opciones del <select> de terapeuta
	JsonRepresentantes string            // JSON de ejemplo del servicio de representantes
	JsonTerapeutas     string            // JSON de ejemplo del servicio de terapeutas
	JsonPlanes         string            // JSON de ejemplo del servicio de planes
	JsonEstudiantes    string            // JSON de ejemplo del servicio de estudiantes
	JsonCitas          string            // JSON de ejemplo del servicio de citas
	JsonFacturas       string            // JSON de ejemplo del servicio de facturas
	MontoTotal         int               // suma de los montos de las facturas (para el resumen)
}

// portalInicio renderiza la pagina principal del portal.
func (ui *UI) portalInicio(w http.ResponseWriter, r *http.Request) {
	// Reunimos toda la informacion que la plantilla necesita para pintarse.
	d := PortalData{}

	// Mensaje opcional que llega por la URL (ej. ?msg=Guardado+correctamente).
	d.Mensaje = r.URL.Query().Get("msg")

	// Traemos las listas de cada entidad y las convertimos a su tipo concreto
	// (la plantilla necesita tipos concretos, no la lista generica del almacen).
	d.Representantes = aRepresentantes(ui.a.Listar("representante"))
	d.Terapeutas = aTerapeutas(ui.a.Listar("terapeuta"))
	d.Planes = aPlanes(ui.a.Listar("plan"))
	d.Estudiantes = aEstudiantes(ui.a.Listar("estudiante"))
	d.Citas = aCitas(ui.a.Listar("cita"))
	d.Facturas = aFacturas(ui.a.Listar("factura"))

	// Armamos las opciones de los desplegables segun los datos existentes.
	d.OptsRepresentante = optsRepresentante(ui.a)
	d.OptsPlan = optsPlan(ui.a)
	d.OptsEstudiante = optsEstudiante(ui.a)
	d.OptsTerapeuta = optsTerapeuta(ui.a)

	// Convertimos cada lista a JSON para mostrarlo en la pagina como ejemplo
	// de lo que devuelve la API.
	if b, err := json.Marshal(ui.a.Listar("representante")); err == nil {
		d.JsonRepresentantes = string(b)
	}
	if b, err := json.Marshal(ui.a.Listar("estudiante")); err == nil {
		d.JsonEstudiantes = string(b)
	}
	if b, err := json.Marshal(ui.a.Listar("cita")); err == nil {
		d.JsonCitas = string(b)
	}
	if b, err := json.Marshal(ui.a.Listar("factura")); err == nil {
		d.JsonFacturas = string(b)
	}
	if b, err := json.Marshal(ui.a.Listar("terapeuta")); err == nil {
		d.JsonTerapeutas = string(b)
	}
	if b, err := json.Marshal(ui.a.Listar("plan")); err == nil {
		d.JsonPlanes = string(b)
	}
	// Calculamos el monto total facturado para mostrarlo en el resumen.
	// d.Facturas ya es []Factura (lo armo arriba con aFacturas).
	mt := 0
	for _, f := range d.Facturas {
		mt += f.Monto
	}
	d.MontoTotal = mt

	// Leemos la plantilla empaquetada y la renderizamos con los datos.
	t, err := template.ParseFS(plantillaPortal, "templates/portal.html")
	if err != nil {
		respuestaError(w, http.StatusInternalServerError, "plantilla no disponible")
		return
	}
	t.ExecuteTemplate(w, "portal", d)
}

// portalCrear recibe el formulario del portal y crea la entidad elegida.
func (ui *UI) portalCrear(w http.ResponseWriter, r *http.Request) {
	// Solo aceptamos envios de formulario (POST). Si llega otro metodo, error.
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}
	// ParseForm lee los campos del formulario para poder consultarlos.
	r.ParseForm()

	// El nombre de la entidad viene en la ruta: /portal/crear/cita -> "cita".
	entidad := strings.TrimPrefix(r.URL.Path, "/portal/crear/")

	var err error
	// Segun la entidad, armamos su struct con los campos del formulario y
	// llamamos a la MISMA regla de negocio que usa la API (reglas.go).
	switch entidad {
	case "representante":
		e := Representante{
			Nombres:   formStr(r, "nombres"),
			Apellidos: formStr(r, "apellidos"),
			Cedula:    formStr(r, "cedula"),
			Telefono:  formStr(r, "telefono"),
		}
		err = crearRepresentante(ui.a, &e)
	case "estudiante":
		e := Estudiante{
			Nombres:         formStr(r, "nombres"),
			Apellidos:       formStr(r, "apellidos"),
			Cedula:          formStr(r, "cedula"),
			Telefono:        formStr(r, "telefono"),
			RepresentanteID: formInt(r, "representante_id"),
			PlanID:          formInt(r, "plan_id"),
		}
		err = crearEstudiante(ui.a, &e)
	case "cita":
		c := Cita{
			EstudianteID: formInt(r, "estudiante_id"),
			TerapeutaID:  formInt(r, "terapeuta_id"),
			Fecha:        formStr(r, "fecha"),
			Hora:         formStr(r, "hora"),
			Estado:       formStr(r, "estado"),
		}
		err = crearCita(ui.a, &c)
	case "factura":
		f := Factura{
			RepresentanteID: formInt(r, "representante_id"),
			Fecha:           formStr(r, "fecha"),
			Monto:           formInt(r, "monto"),
			Estado:          formStr(r, "estado"),
		}
		err = crearFactura(ui.a, &f)
	case "plan":
		p := Plan{
			Nombre:      formStr(r, "nombre"),
			Descripcion: formStr(r, "descripcion"),
			Duracion:    formInt(r, "duracion"),
		}
		err = crearPlan(ui.a, &p)
	case "terapeuta":
		t := Terapeuta{
			Nombres:   formStr(r, "nombres"),
			Apellidos: formStr(r, "apellidos"),
			Cedula:    formStr(r, "cedula"),
			Telefono:  formStr(r, "telefono"),
		}
		err = crearTerapeuta(ui.a, &t)
	default:
		// Si la entidad no es ninguna de las conocidas, marcamos error.
		err = errors.New("entidad desconocida")
	}

	// Preparamos un mensaje para mostrar en el portal: exito o la razon del error.
	msg := "Guardado correctamente"
	if err != nil {
		msg = err.Error()
	}
	// Redirigimos al portal, pasando el mensaje en la URL para que se muestre.
	http.Redirect(w, r, "/portal?msg="+url.QueryEscape(msg), http.StatusSeeOther)
}

// formStr lee un campo de texto del formulario y quita espacios sobrantes.
func formStr(r *http.Request, nombre string) string {
	return strings.TrimSpace(r.FormValue(nombre))
}

// formInt lee un campo numerico del formulario y lo convierte a entero.
func formInt(r *http.Request, nombre string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(r.FormValue(nombre)))
	return n
}

/*
 Las siguientes funciones convierten la lista generica del almacen
 ([]interface{}) al tipo concreto que necesita la plantilla HTML.
 */
func aRepresentantes(l []interface{}) []Representante {
	res := make([]Representante, 0, len(l))
	for _, v := range l {
		res = append(res, v.(Representante))
	}
	return res
}
func aEstudiantes(l []interface{}) []Estudiante {
	res := make([]Estudiante, 0, len(l))
	for _, v := range l {
		res = append(res, v.(Estudiante))
	}
	return res
}
func aCitas(l []interface{}) []Cita {
	res := make([]Cita, 0, len(l))
	for _, v := range l {
		res = append(res, v.(Cita))
	}
	return res
}
func aFacturas(l []interface{}) []Factura {
	res := make([]Factura, 0, len(l))
	for _, v := range l {
		res = append(res, v.(Factura))
	}
	return res
}
func aTerapeutas(l []interface{}) []Terapeuta {
	res := make([]Terapeuta, 0, len(l))
	for _, v := range l {
		res = append(res, v.(Terapeuta))
	}
	return res
}
func aPlanes(l []interface{}) []Plan {
	res := make([]Plan, 0, len(l))
	for _, v := range l {
		res = append(res, v.(Plan))
	}
	return res
}

/*
 Las siguientes funciones arman las opciones de los <select> a partir de los
 datos existentes, para que el usuario elija representante, plan, etc.
 */
func optsRepresentante(a *Almacen) []opcion {
	var res []opcion
	for _, v := range a.Listar("representante") {
		r := v.(Representante)
		res = append(res, opcion{Valor: strconv.Itoa(r.ID), Texto: r.Nombres + " " + r.Apellidos})
	}
	return res
}
func optsPlan(a *Almacen) []opcion {
	var res []opcion
	for _, v := range a.Listar("plan") {
		p := v.(Plan)
		res = append(res, opcion{Valor: strconv.Itoa(p.ID), Texto: p.Nombre})
	}
	return res
}
func optsEstudiante(a *Almacen) []opcion {
	var res []opcion
	for _, v := range a.Listar("estudiante") {
		e := v.(Estudiante)
		res = append(res, opcion{Valor: strconv.Itoa(e.ID), Texto: e.Nombres + " " + e.Apellidos})
	}
	return res
}
func optsTerapeuta(a *Almacen) []opcion {
	var res []opcion
	for _, v := range a.Listar("terapeuta") {
		t := v.(Terapeuta)
		res = append(res, opcion{Valor: strconv.Itoa(t.ID), Texto: t.Nombres + " " + t.Apellidos})
	}
	return res
}

// EditData es la "bolsa" de datos para la pagina de edicion de un registro.
type EditData struct {
	Tipo   string         // nombre singular de la entidad (ej. "representante")
	ID     int            // id del registro a editar
	Nombre string         // nombre legible de la entidad (ej. "Representante")
	Form   template.HTML  // formulario de edicion ya armado con los datos actuales
}

/*
 nombreEntidad devuelve el nombre legible (en mayuscula) de cada entidad,
 para mostrarlo en el titulo de la pagina de edicion.
 */
func nombreEntidad(tipo string) string {
	switch tipo {
	case "representante":
		return "Representante"
	case "terapeuta":
		return "Terapeuta"
	case "plan":
		return "Plan"
	case "estudiante":
		return "Estudiante"
	case "cita":
		return "Cita"
	case "factura":
		return "Factura"
	}
	return tipo
}

/*
 portalEditar muestra el formulario para editar un registro existente.
 La ruta es /portal/editar/{entidad}/{id}. El formulario se arma precargado
 con los datos actuales del registro (sin JavaScript).
 */
func (ui *UI) portalEditar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}
	// El tipo viene de la ruta en singular (ej. "cita"); lo usamos directo.
	tipo := r.PathValue("entidad")
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Redirect(w, r, "/portal?msg=id invalido", http.StatusSeeOther)
		return
	}
	// Buscamos el registro; si no existe, volvemos al portal con un aviso.
	v, ok := ui.a.Obtener(tipo, id)
	if !ok {
		http.Redirect(w, r, "/portal?msg=registro no encontrado", http.StatusSeeOther)
		return
	}
	// Construimos el formulario de edicion (HTML ya listo).
	d := EditData{
		Tipo:   tipo,
		ID:     id,
		Nombre: nombreEntidad(tipo),
		Form:   formularioEdicion(ui.a, tipo, id, v),
	}
	t, err := template.ParseFS(plantillaPortal, "templates/portal.html")
	if err != nil {
		respuestaError(w, http.StatusInternalServerError, "plantilla no disponible")
		return
	}
	t.ExecuteTemplate(w, "editar", d)
}

/*
 portalActualizar aplica los cambios del formulario de edicion y vuelve al
 portal. La ruta es /portal/actualizar/{entidad}/{id} (metodo POST).
 */
func (ui *UI) portalActualizar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}
	r.ParseForm()
	tipo := r.PathValue("entidad")
	id, _ := strconv.Atoi(r.PathValue("id"))

	var err error
	// Segun la entidad, armamos su struct con los campos del formulario y
	// llamamos a la MISMA regla de negocio que usa la API (actualizar*).
	switch tipo {
	case "representante":
		e := Representante{
			Nombres:   formStr(r, "nombres"),
			Apellidos: formStr(r, "apellidos"),
			Cedula:    formStr(r, "cedula"),
			Telefono:  formStr(r, "telefono"),
		}
		err = actualizarRepresentante(ui.a, id, e)
	case "terapeuta":
		e := Terapeuta{
			Nombres:   formStr(r, "nombres"),
			Apellidos: formStr(r, "apellidos"),
			Cedula:    formStr(r, "cedula"),
			Telefono:  formStr(r, "telefono"),
		}
		err = actualizarTerapeuta(ui.a, id, e)
	case "plan":
		p := Plan{
			Nombre:      formStr(r, "nombre"),
			Descripcion: formStr(r, "descripcion"),
			Duracion:    formInt(r, "duracion"),
		}
		err = actualizarPlan(ui.a, id, p)
	case "estudiante":
		e := Estudiante{
			Nombres:         formStr(r, "nombres"),
			Apellidos:       formStr(r, "apellidos"),
			Cedula:          formStr(r, "cedula"),
			Telefono:        formStr(r, "telefono"),
			RepresentanteID: formInt(r, "representante_id"),
			PlanID:          formInt(r, "plan_id"),
		}
		err = actualizarEstudiante(ui.a, id, e)
	case "cita":
		c := Cita{
			EstudianteID: formInt(r, "estudiante_id"),
			TerapeutaID:  formInt(r, "terapeuta_id"),
			Fecha:        formStr(r, "fecha"),
			Hora:         formStr(r, "hora"),
			Estado:       formStr(r, "estado"),
		}
		err = actualizarCita(ui.a, id, c)
	case "factura":
		f := Factura{
			RepresentanteID: formInt(r, "representante_id"),
			Fecha:           formStr(r, "fecha"),
			Monto:           formInt(r, "monto"),
			Estado:          formStr(r, "estado"),
		}
		err = actualizarFactura(ui.a, id, f)
	default:
		err = errors.New("entidad desconocida")
	}

	msg := "Actualizado correctamente"
	if err != nil {
		msg = err.Error()
	}
	http.Redirect(w, r, "/portal?msg="+url.QueryEscape(msg), http.StatusSeeOther)
}

/*
 portalBorrar elimina un registro tras comprobar que no deje referencias
 rotas. La ruta es /portal/borrar/{entidad}/{id} (metodo POST, desde un
 formulario).
 */
func (ui *UI) portalBorrar(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo no permitido", http.StatusMethodNotAllowed)
		return
	}
	tipo := r.PathValue("entidad")
	id, _ := strconv.Atoi(r.PathValue("id"))

	// Revisamos que borrarlo no rompa otras entidades (usa reglas.go).
	ok, mensaje := puedeBorrar(ui.a, tipo, id)
	if !ok {
		http.Redirect(w, r, "/portal?msg="+url.QueryEscape(mensaje), http.StatusSeeOther)
		return
	}
	ui.a.Borrar(tipo, id)
	http.Redirect(w, r, "/portal?msg=Borrado correctamente", http.StatusSeeOther)
}

// ----- Ayudantes para construir el formulario de edicion (sin JavaScript) -----

/*
 attr escapa un valor para usarlo dentro de un atributo HTML (ej. value="...").
 Asi evitamos que comillas o signos < > rompan el formulario.
 */
func attr(s string) string {
	return html.EscapeString(s)
}

// inputText arma una fila de texto: etiqueta + caja de texto con su valor.
func inputText(nombre, etiqueta, valor string, oblig bool) string {
	ob := ""
	if oblig {
		ob = ` <span class="oblig">*</span>`
	}
	return fmt.Sprintf(`<label>%s%s</label> <input name="%s" value="%s"><br>`,
		etiqueta, ob, nombre, attr(valor))
}

// inputNum arma una fila numerica: etiqueta + caja de numero con su valor.
func inputNum(nombre, etiqueta string, valor int) string {
	return fmt.Sprintf(`<label>%s</label> <input name="%s" type="number" value="%d"><br>`,
		etiqueta, nombre, valor)
}

/*
 inputSelect arma una fila de lista desplegable, dejando marcada la opcion
 que coincide con el valor actual del registro.
 */
func inputSelect(nombre, etiqueta string, valor int, ops []opcion) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf(`<label>%s</label> <select name="%s">`, etiqueta, nombre))
	for _, o := range ops {
		sel := ""
		if o.Valor == strconv.Itoa(valor) {
			sel = " selected"
		}
		b.WriteString(fmt.Sprintf(`<option value="%s"%s>%s</option>`, o.Valor, sel, attr(o.Texto)))
	}
	b.WriteString("</select><br>")
	return b.String()
}

/*
 formularioEdicion construye el HTML del formulario de edicion, precargado
 con los datos actuales del registro. Devuelve template.HTML (ya es HTML
 seguro, no se vuelve a escapar al renderizar).
 */
func formularioEdicion(a *Almacen, tipo string, id int, v interface{}) template.HTML {
	var b strings.Builder
	// El formulario se envia a actualizar con el tipo y el id en la ruta.
	b.WriteString(fmt.Sprintf(`<form method="POST" action="/portal/actualizar/%s/%d">`, tipo, id))
	switch tipo {
	case "representante", "terapeuta":
		var nom, ape, ced, tel string
		if tipo == "representante" {
			e := v.(Representante)
			nom, ape, ced, tel = e.Nombres, e.Apellidos, e.Cedula, e.Telefono
		} else {
			e := v.(Terapeuta)
			nom, ape, ced, tel = e.Nombres, e.Apellidos, e.Cedula, e.Telefono
		}
		b.WriteString(inputText("nombres", "Nombres", nom, true))
		b.WriteString(inputText("apellidos", "Apellidos", ape, false))
		b.WriteString(inputText("cedula", "Cedula", ced, true))
		b.WriteString(inputText("telefono", "Telefono", tel, false))
	case "plan":
		p := v.(Plan)
		b.WriteString(inputText("nombre", "Nombre", p.Nombre, true))
		b.WriteString(inputText("descripcion", "Descripcion", p.Descripcion, false))
		b.WriteString(inputNum("duracion", "Duracion (min)", p.Duracion))
	case "estudiante":
		e := v.(Estudiante)
		b.WriteString(inputText("nombres", "Nombres", e.Nombres, true))
		b.WriteString(inputText("apellidos", "Apellidos", e.Apellidos, false))
		b.WriteString(inputText("cedula", "Cedula", e.Cedula, true))
		b.WriteString(inputText("telefono", "Telefono", e.Telefono, false))
		b.WriteString(inputSelect("representante_id", "Representante", e.RepresentanteID, optsRepresentante(a)))
		b.WriteString(inputSelect("plan_id", "Plan", e.PlanID, optsPlan(a)))
	case "cita":
		c := v.(Cita)
		b.WriteString(inputSelect("estudiante_id", "Estudiante", c.EstudianteID, optsEstudiante(a)))
		b.WriteString(inputSelect("terapeuta_id", "Terapeuta", c.TerapeutaID, optsTerapeuta(a)))
		b.WriteString(inputText("fecha", "Fecha", c.Fecha, false))
		b.WriteString(inputText("hora", "Hora", c.Hora, false))
		b.WriteString(inputText("estado", "Estado", c.Estado, false))
	case "factura":
		f := v.(Factura)
		b.WriteString(inputSelect("representante_id", "Representante", f.RepresentanteID, optsRepresentante(a)))
		b.WriteString(inputText("fecha", "Fecha", f.Fecha, false))
		b.WriteString(inputNum("monto", "Monto", f.Monto))
		b.WriteString(inputText("estado", "Estado", f.Estado, false))
	}
	b.WriteString(`<button>Guardar cambios</button> <a class="boton" href="/portal">Cancelar</a>`)
	b.WriteString("</form>")
	return template.HTML(b.String())
}

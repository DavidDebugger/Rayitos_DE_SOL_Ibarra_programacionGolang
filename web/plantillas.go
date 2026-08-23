package web

// plantillas.go carga los archivos de vista (HTML y CSS) y los sirve.
// El HTML esta en templates/portal.html y el CSS en templates/estilos.css,
// asi el codigo Go solo se ocupa de la logica.

import (
	"embed"
	"io"
	"net/http"
	"text/template"
)

// recursos contiene los archivos de vista embebidos en el binario.
//
//go:embed templates/portal.html templates/estilos.css
var recursos embed.FS

// plantilla tiene todas las definiciones de portal.html ya parseadas.
var plantilla = func() *template.Template {
	t, err := template.ParseFS(recursos, "templates/portal.html")
	if err != nil {
		panic(err)
	}
	return t
}()

// paginaDatos es la informacion que recibe la plantilla para dibujar una vista.
type paginaDatos struct {
	Vista      string
	Titulo     string
	MenuActivo string
	Menu       []menuItem
	Label      string
	Columnas   []string
	Filas      []map[string]any
	Colspan    int
	Accion     string
	Campos     []campoForm
	Valores    map[string]any
	Error      string
	Servicios  []servicioWeb
	PorEntidad []string
	PortalPag  []paginaInfo
	JSON       string
}

// servicioWeb describe un servicio de la API para la vista de servicios.
// Ejemplo es el JSON que la API serializa: en GET es la respuesta real y en
// POST/PUT es el cuerpo JSON de ejemplo que acepta la API.
type servicioWeb struct {
	Metodo  string
	Ruta    string
	Desc    string
	Cuerpo  string
	Ejemplo string
}

// menuItem es una entrada del menu lateral.
type menuItem struct {
	Plural string
	Label  string
}

// endpoint describe una ruta de la API para la vista de servicios.
type endpoint struct {
	Metodo, Ruta, Desc string
}

// paginaInfo describe una pagina del portal para la vista de servicios.
type paginaInfo struct {
	Ruta, Desc string
}

// menuPortal arma la lista fija del menu lateral.
func menuPortal() []menuItem {
	return []menuItem{
		{"servicios", "Servicios Web"},
		{"representantes", "Representantes"},
		{"terapeutas", "Terapeutas"},
		{"planes", "Planes"},
		{"estudiantes", "Estudiantes"},
		{"citas", "Citas"},
		{"facturas", "Facturas"},
	}
}

// servirPlantilla ejecuta la plantilla que indique datos.Vista (inicio, listar, formulario, servicios).
func servirPlantilla(w http.ResponseWriter, r *http.Request, datos paginaDatos) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := plantilla.ExecuteTemplate(w, datos.Vista, datos); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// servirCSS responde con la hoja de estilos embebida.
func servirCSS(w http.ResponseWriter, r *http.Request) {
	b, err := recursos.ReadFile("templates/estilos.css")
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	io.WriteString(w, string(b))
}

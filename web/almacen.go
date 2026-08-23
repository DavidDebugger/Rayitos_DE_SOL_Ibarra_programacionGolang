package web

// almacen.go guarda los datos en memoria mientras el programa corre.
// Usa un mutex para que varias peticiones no se pisen entre si.
// Los datos se pierden al apagar el servidor (no se escribe datos.json).

import (
	"sort"
	"strconv"
	"strings"
	"sync"
)

// mapas que guardan cada entidad por su ID.
var (
	mu            sync.Mutex
	representantes = map[int]Representante{}
	terapeutas    = map[int]Terapeuta{}
	planes        = map[int]Plan{}
	estudiantes   = map[int]Estudiante{}
	citas         = map[int]Cita{}
	facturas      = map[int]Factura{}
	siguienteID   = 1
)

// pluralDeSingular convierte el nombre singular (de la URL) al plural del endpoint.
var pluralDeSingular = map[string]string{
	"representante": "representantes",
	"terapeuta":     "terapeutas",
	"plan":          "planes",
	"estudiante":    "estudiantes",
	"cita":          "citas",
	"factura":       "facturas",
}

// singularDePlural es la vuelta: del plural al singular, para validaciones.
var singularDePlural = map[string]string{
	"representantes": "representante",
	"terapeutas":     "terapeuta",
	"planes":         "plan",
	"estudiantes":    "estudiante",
	"citas":          "cita",
	"facturas":       "factura",
}

// entidadValida dice si el plural recibido en la URL es uno de los 6 aceptados.
func entidadValida(plural string) bool {
	_, ok := singularDePlural[plural]
	return ok
}

// normalizarEntidad acepta tanto singular como plural y siempre devuelve el plural.
func normalizarEntidad(nombre string) (string, bool) {
	if entidadValida(nombre) {
		return nombre, true
	}
	if plural, ok := pluralDeSingular[nombre]; ok {
		return plural, true
	}
	return "", false
}

// nuevoID entrega un ID distinto cada vez y lo suma al contador.
func nuevoID() int {
	id := siguienteID
	siguienteID++
	return id
}

// ajustarContador sube el contador por encima del ID mas alto conocido.
// Se usa al cargar datos.json para no repetir IDs.
func ajustarContador(id int) {
	if id >= siguienteID {
		siguienteID = id + 1
	}
}

// ---- Representante ----

func obtenerRepresentantes() []Representante {
	mu.Lock()
	defer mu.Unlock()
	lista := make([]Representante, 0, len(representantes))
	for _, v := range representantes {
		lista = append(lista, v)
	}
	sort.Slice(lista, func(i, j int) bool { return lista[i].ID < lista[j].ID })
	return lista
}

func obtenerRepresentante(id int) (Representante, bool) {
	mu.Lock()
	defer mu.Unlock()
	v, ok := representantes[id]
	return v, ok
}

func insertarRepresentante(r Representante) Representante {
	mu.Lock()
	defer mu.Unlock()
	r.ID = nuevoID()
	representantes[r.ID] = r
	return r
}

func actualizarRepresentante(id int, r Representante) (Representante, bool) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := representantes[id]; !ok {
		return Representante{}, false
	}
	r.ID = id
	representantes[id] = r
	return r, true
}

func eliminarRepresentante(id int) bool {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := representantes[id]; !ok {
		return false
	}
	delete(representantes, id)
	return true
}

// ---- Terapeuta ----

func obtenerTerapeutas() []Terapeuta {
	mu.Lock()
	defer mu.Unlock()
	lista := make([]Terapeuta, 0, len(terapeutas))
	for _, v := range terapeutas {
		lista = append(lista, v)
	}
	sort.Slice(lista, func(i, j int) bool { return lista[i].ID < lista[j].ID })
	return lista
}

func obtenerTerapeuta(id int) (Terapeuta, bool) {
	mu.Lock()
	defer mu.Unlock()
	v, ok := terapeutas[id]
	return v, ok
}

func insertarTerapeuta(t Terapeuta) Terapeuta {
	mu.Lock()
	defer mu.Unlock()
	t.ID = nuevoID()
	terapeutas[t.ID] = t
	return t
}

func actualizarTerapeuta(id int, t Terapeuta) (Terapeuta, bool) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := terapeutas[id]; !ok {
		return Terapeuta{}, false
	}
	t.ID = id
	terapeutas[id] = t
	return t, true
}

func eliminarTerapeuta(id int) bool {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := terapeutas[id]; !ok {
		return false
	}
	delete(terapeutas, id)
	return true
}

// ---- Plan ----

func obtenerPlanes() []Plan {
	mu.Lock()
	defer mu.Unlock()
	lista := make([]Plan, 0, len(planes))
	for _, v := range planes {
		lista = append(lista, v)
	}
	sort.Slice(lista, func(i, j int) bool { return lista[i].ID < lista[j].ID })
	return lista
}

func obtenerPlan(id int) (Plan, bool) {
	mu.Lock()
	defer mu.Unlock()
	v, ok := planes[id]
	return v, ok
}

func insertarPlan(p Plan) Plan {
	mu.Lock()
	defer mu.Unlock()
	p.ID = nuevoID()
	planes[p.ID] = p
	return p
}

func actualizarPlan(id int, p Plan) (Plan, bool) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := planes[id]; !ok {
		return Plan{}, false
	}
	p.ID = id
	planes[id] = p
	return p, true
}

func eliminarPlan(id int) bool {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := planes[id]; !ok {
		return false
	}
	delete(planes, id)
	return true
}

// ---- Estudiante ----

func obtenerEstudiantes() []Estudiante {
	mu.Lock()
	defer mu.Unlock()
	lista := make([]Estudiante, 0, len(estudiantes))
	for _, v := range estudiantes {
		lista = append(lista, v)
	}
	sort.Slice(lista, func(i, j int) bool { return lista[i].ID < lista[j].ID })
	return lista
}

func obtenerEstudiante(id int) (Estudiante, bool) {
	mu.Lock()
	defer mu.Unlock()
	v, ok := estudiantes[id]
	return v, ok
}

func insertarEstudiante(e Estudiante) Estudiante {
	mu.Lock()
	defer mu.Unlock()
	e.ID = nuevoID()
	estudiantes[e.ID] = e
	return e
}

func actualizarEstudiante(id int, e Estudiante) (Estudiante, bool) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := estudiantes[id]; !ok {
		return Estudiante{}, false
	}
	e.ID = id
	estudiantes[id] = e
	return e, true
}

func eliminarEstudiante(id int) bool {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := estudiantes[id]; !ok {
		return false
	}
	delete(estudiantes, id)
	return true
}

// ---- Cita ----

func obtenerCitas() []Cita {
	mu.Lock()
	defer mu.Unlock()
	lista := make([]Cita, 0, len(citas))
	for _, v := range citas {
		lista = append(lista, v)
	}
	sort.Slice(lista, func(i, j int) bool { return lista[i].ID < lista[j].ID })
	return lista
}

func obtenerCita(id int) (Cita, bool) {
	mu.Lock()
	defer mu.Unlock()
	v, ok := citas[id]
	return v, ok
}

func insertarCita(c Cita) Cita {
	mu.Lock()
	defer mu.Unlock()
	c.ID = nuevoID()
	citas[c.ID] = c
	return c
}

func actualizarCita(id int, c Cita) (Cita, bool) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := citas[id]; !ok {
		return Cita{}, false
	}
	c.ID = id
	citas[id] = c
	return c, true
}

func eliminarCita(id int) bool {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := citas[id]; !ok {
		return false
	}
	delete(citas, id)
	return true
}

// ---- Factura ----

func obtenerFacturas() []Factura {
	mu.Lock()
	defer mu.Unlock()
	lista := make([]Factura, 0, len(facturas))
	for _, v := range facturas {
		lista = append(lista, v)
	}
	sort.Slice(lista, func(i, j int) bool { return lista[i].ID < lista[j].ID })
	return lista
}

func obtenerFactura(id int) (Factura, bool) {
	mu.Lock()
	defer mu.Unlock()
	v, ok := facturas[id]
	return v, ok
}

func insertarFactura(f Factura) Factura {
	mu.Lock()
	defer mu.Unlock()
	f.ID = nuevoID()
	facturas[f.ID] = f
	return f
}

func actualizarFactura(id int, f Factura) (Factura, bool) {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := facturas[id]; !ok {
		return Factura{}, false
	}
	f.ID = id
	facturas[id] = f
	return f, true
}

func eliminarFactura(id int) bool {
	mu.Lock()
	defer mu.Unlock()
	if _, ok := facturas[id]; !ok {
		return false
	}
	delete(facturas, id)
	return true
}

// aEntero convierte texto a entero; si falla devuelve 0.
func aEntero(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}

// aFlotante convierte texto a decimal; si falla devuelve 0.
func aFlotante(s string) float64 {
	f, _ := strconv.ParseFloat(strings.TrimSpace(s), 64)
	return f
}

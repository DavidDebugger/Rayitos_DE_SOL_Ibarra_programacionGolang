/*
 mux.go
 Este archivo es el "cerebro" que une todo: crea el almacen, carga los datos
 de ejemplo, siembra lo basico y registra las rutas (las paginas del portal y
 los servicios de la API). Luego enciende el servidor.
 "mux" viene de "multiplexor": es el encargado de mirar la URL de cada
 peticion y dirigirla a la funcion correcta.
 */

package web

import (
	"log"
	"net/http"
)

/*
 sembrar crea un Plan y un Terapeuta minimos por defecto, pero solo si el
 almacen esta vacio para esa entidad. Hoy Plan y Terapeuta tambien tienen
 servicios y formularios propios, asi que esto es solo un respaldo por si
 datos.json no existe o no trae ninguno.
 */
func sembrar(a *Almacen) {
	if len(a.Listar("plan")) == 0 {
		p := Plan{ID: a.nuevoID("plan"), Nombre: "Basico", Duracion: 60}
		a.Guardar("plan", p.ID, p)
	}
	if len(a.Listar("terapeuta")) == 0 {
		t := Terapeuta{ID: a.nuevoID("terapeuta"), Nombres: "Terapeuta", Apellidos: "General", Cedula: "000"}
		a.Guardar("terapeuta", t.ID, t)
	}
}

/*
 IniciarServidor arma todo y arranca el servidor en el puerto indicado.
   puerto: texto como ":8080" (el dos puntos indica "todas las direcciones").
 */
func IniciarServidor(puerto string) {
	// 1) Creamos el almacen vacio.
	a := NuevoAlmacen()

	// 2) Intentamos cargar datos.json con ejemplos. Si no existe o falla,
	//    lo avisamos pero seguimos: sembrar() cubre lo esencial.
	if err := cargarDatos(a, "datos.json"); err != nil {
		log.Printf("aviso: no se cargo datos.json: %v", err)
	}

	// 3) Sembramos Plan y Terapeuta por defecto si hiciera falta.
	sembrar(a)

	// 4) Creamos la interfaz (UI) que usaran los handlers.
	ui := &UI{a: a}

	// 5) Creamos el enrutador y registramos las rutas.
	m := http.NewServeMux()

	// --- Rutas del portal web (paginas con clicks, sin JavaScript) ---
	m.HandleFunc("GET /", ui.portalInicio)                          // la pagina principal es el portal
	m.HandleFunc("GET /portal", ui.portalInicio)                    // alias de la pagina principal
	m.HandleFunc("POST /portal/crear/{entidad}", ui.portalCrear)    // crea segun el formulario enviado
	m.HandleFunc("GET /portal/editar/{entidad}/{id}", ui.portalEditar)      // muestra formulario de edicion
	m.HandleFunc("POST /portal/actualizar/{entidad}/{id}", ui.portalActualizar) // aplica los cambios
	m.HandleFunc("POST /portal/borrar/{entidad}/{id}", ui.portalBorrar)        // elimina un registro

	// --- Rutas de la API JSON (servicios GET/POST por entidad) ---
	m.HandleFunc("GET /api/representantes", ui.apiRepresentantes)
	m.HandleFunc("POST /api/representantes", ui.apiRepresentantes)
	m.HandleFunc("GET /api/estudiantes", ui.apiEstudiantes)
	m.HandleFunc("POST /api/estudiantes", ui.apiEstudiantes)
	m.HandleFunc("GET /api/citas", ui.apiCitas)
	m.HandleFunc("POST /api/citas", ui.apiCitas)
	m.HandleFunc("GET /api/facturas", ui.apiFacturas)
	m.HandleFunc("POST /api/facturas", ui.apiFacturas)
	m.HandleFunc("GET /api/planes", ui.apiPlanes)
	m.HandleFunc("POST /api/planes", ui.apiPlanes)
	m.HandleFunc("GET /api/terapeutas", ui.apiTerapeutas)
	m.HandleFunc("POST /api/terapeutas", ui.apiTerapeutas)

	// --- Rutas de la API JSON por id (GET/PUT/DELETE de una entidad) ---
	/*
	 Una sola ruta con {entidad} y {id} atiende las 6 entidades; el nombre
	 plural se traduce al singular dentro de apiPorID (ver api.go).
	 */
	m.HandleFunc("GET /api/{entidad}/{id}", ui.apiPorID)
	m.HandleFunc("PUT /api/{entidad}/{id}", ui.apiPorID)
	m.HandleFunc("DELETE /api/{entidad}/{id}", ui.apiPorID)

	// 6) Arrancamos el servidor. log.Fatal detiene el programa si hay un error
	//    grave (por ejemplo, si el puerto ya esta ocupado).
	log.Fatal(http.ListenAndServe(puerto, m))
}

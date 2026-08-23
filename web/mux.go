package web

// mux.go enlaza las rutas con sus manejadores y arranca el servidor.

import (
	"net/http"
	"strconv"
	"strings"
)

// enrutar decide si la peticion va al portal o a la API.
func enrutar(w http.ResponseWriter, r *http.Request) {
	ruta := r.URL.Path
	switch {
	case ruta == "/" || ruta == "/portal" || ruta == "/portal/":
		portalInicio(w, r)
	case strings.HasPrefix(ruta, "/api/"):
		manejarAPI(w, r, ruta)
	case strings.HasPrefix(ruta, "/portal/"):
		manejarPortal(w, r, ruta)
	default:
		http.NotFound(w, r)
	}
}

// manejarAPI parsea /api/{entidad} y /api/{entidad}/{id}.
func manejarAPI(w http.ResponseWriter, r *http.Request, ruta string) {
	partes := strings.Split(strings.Trim(ruta, "/"), "/") // ["api", entidad, id?]
	if len(partes) < 2 {
		escribirJSON(w, http.StatusNotFound, map[string]string{"error": "ruta no encontrada"})
		return
	}
	entidad, ok := normalizarEntidad(partes[1])
	if !ok {
		escribirJSON(w, http.StatusNotFound, map[string]string{"error": "entidad no valida"})
		return
	}
	if len(partes) == 2 {
		apiColeccion(w, r, entidad)
		return
	}
	id, err := strconv.Atoi(partes[2])
	if err != nil {
		escribirJSON(w, http.StatusBadRequest, map[string]string{"error": "id invalido"})
		return
	}
	apiPorID(w, r, entidad, id)
}

// manejarPortal parsea las rutas del portal web.
func manejarPortal(w http.ResponseWriter, r *http.Request, ruta string) {
	partes := strings.Split(strings.Trim(ruta, "/"), "/") // ["portal", entidad, ...]
	if len(partes) < 2 || partes[1] == "" {
		portalInicio(w, r)
		return
	}
	entidad, ok := normalizarEntidad(partes[1])
	if !ok {
		http.NotFound(w, r)
		return
	}
	if len(partes) == 2 {
		switch r.Method {
		case http.MethodGet:
			portalListar(w, r, entidad)
		case http.MethodPost:
			portalGuardar(w, r, entidad, 0)
		default:
			http.Error(w, "metodo no permitido", http.StatusMethodNotAllowed)
		}
		return
	}
	algo := partes[2]
	if algo == "nuevo" {
		if r.Method == http.MethodGet {
			portalCrear(w, r, entidad)
		} else {
			http.Error(w, "metodo no permitido", http.StatusMethodNotAllowed)
		}
		return
	}
	id, err := strconv.Atoi(algo)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if len(partes) == 3 {
		switch r.Method {
		case http.MethodGet:
			portalEditar(w, r, entidad, id)
		case http.MethodPost:
			portalGuardar(w, r, entidad, id)
		default:
			http.Error(w, "metodo no permitido", http.StatusMethodNotAllowed)
		}
		return
	}
	if len(partes) == 4 && partes[3] == "borrar" && r.Method == http.MethodPost {
		portalBorrar(w, r, entidad, id)
		return
	}
	http.NotFound(w, r)
}

// Servir carga los datos y levanta el servidor en el puerto 8080.
func Servir() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", enrutar)
	return http.ListenAndServe(":8080", mux)
}

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
	case ruta == "/estilos.css":
		servirCSS(w, r)
	case strings.HasPrefix(ruta, "/api"):
		manejarAPI(w, r, ruta)
	case strings.HasPrefix(ruta, "/portal/"):
		manejarPortal(w, r, ruta)
	default:
		http.NotFound(w, r)
	}
}

// manejarAPI parsea /api/{entidad} y /api/{entidad}/{id}.
func manejarAPI(w http.ResponseWriter, r *http.Request, ruta string) {
	partes := strings.Split(strings.Trim(ruta, "/"), "/") // ["api", entidad, ...]

	// /api sin nada mas: indice de entidades
	if len(partes) == 1 {
		apiIndice(w, r)
		return
	}

	// /api/estadisticas: totales por entidad
	if partes[1] == "estadisticas" {
		apiEstadisticas(w, r)
		return
	}

	entidad, ok := normalizarEntidad(partes[1])
	if !ok {
		escribirJSON(w, http.StatusNotFound, map[string]string{"error": "entidad no valida"})
		return
	}

	// /api/{entidad}: coleccion (listar o crear)
	if len(partes) == 2 {
		apiColeccion(w, r, entidad)
		return
	}

	// /api/{entidad}/contar  y  /api/{entidad}/buscar (solo GET)
	if len(partes) == 3 {
		switch partes[2] {
		case "contar":
			apiContar(w, r, entidad)
			return
		case "buscar":
			apiBuscar(w, r, entidad)
			return
		}
	}

	// /api/{entidad}/{id}: uno (ver, actualizar o borrar)
	id, err := strconv.Atoi(partes[2])
	if err != nil {
		escribirJSON(w, http.StatusBadRequest, map[string]string{"error": "ruta no valida"})
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
	if partes[1] == "servicios" {
		if r.Method == http.MethodGet {
			portalServicios(w, r)
		} else {
			http.Error(w, "metodo no permitido", http.StatusMethodNotAllowed)
		}
		return
	}
	if partes[1] == "probar" {
		if r.Method == http.MethodGet {
			portalProbar(w, r)
		} else {
			http.Error(w, "metodo no permitido", http.StatusMethodNotAllowed)
		}
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
	if len(partes) == 4 && partes[3] == "ver" && r.Method == http.MethodGet {
		portalVer(w, r, entidad, id)
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

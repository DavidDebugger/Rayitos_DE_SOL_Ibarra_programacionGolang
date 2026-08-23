package web

/* api.go expone las entidades como API JSON.
   Rutas:
     GET    /api                     - indice
     GET    /api/estadisticas        - totales por entidad
     GET    /api/{entidad}           - lista
     POST   /api/{entidad}           - crea
     GET    /api/{entidad}/contar    - cuenta
     GET    /api/{entidad}/buscar    - filtra por campo y valor
     GET    /api/{entidad}/{id}      - uno
     PUT    /api/{entidad}/{id}      - actualiza
     DELETE /api/{entidad}/{id}      - borra
*/

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// escribirJSON envia un codigo HTTP y un dato como JSON indentado.
func escribirJSON(w http.ResponseWriter, codigo int, datos any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codigo)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	_ = enc.Encode(datos)
}

// leerCuerpoJSON lee el cuerpo de la peticion y lo convierte a JSON.
func leerCuerpoJSON(r *http.Request, destino any) error {
	cuerpo, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(cuerpo, destino)
}

// metodoNoPermitido responde cuando usan un verbo HTTP no soportado.
func metodoNoPermitido(w http.ResponseWriter, r *http.Request) {
	escribirJSON(w, http.StatusMethodNotAllowed, map[string]string{
		"error": "metodo " + r.Method + " no permitido en esta ruta",
	})
}

// apiIndice describe la API: entidades y servicios disponibles.
func apiIndice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		metodoNoPermitido(w, r)
		return
	}
	servicios := []map[string]string{
		{"metodo": "GET", "ruta": "/api", "descripcion": "Indice de la API"},
		{"metodo": "GET", "ruta": "/api/estadisticas", "descripcion": "Totales por entidad"},
		{"metodo": "GET", "ruta": "/api/{entidad}", "descripcion": "Listar registros"},
		{"metodo": "POST", "ruta": "/api/{entidad}", "descripcion": "Crear registro"},
		{"metodo": "GET", "ruta": "/api/{entidad}/contar", "descripcion": "Contar registros"},
		{"metodo": "GET", "ruta": "/api/{entidad}/buscar", "descripcion": "Buscar por campo y valor"},
		{"metodo": "GET", "ruta": "/api/{entidad}/{id}", "descripcion": "Ver un registro"},
		{"metodo": "PUT", "ruta": "/api/{entidad}/{id}", "descripcion": "Actualizar registro"},
		{"metodo": "DELETE", "ruta": "/api/{entidad}/{id}", "descripcion": "Borrar registro"},
	}
	escribirJSON(w, http.StatusOK, map[string]any{
		"entidades": []string{
			"representantes", "terapeutas", "planes", "estudiantes", "citas", "facturas",
		},
		"servicios": servicios,
	})
}

// apiEstadisticas devuelve la cantidad de registros de cada entidad.
func apiEstadisticas(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		metodoNoPermitido(w, r)
		return
	}
	escribirJSON(w, http.StatusOK, map[string]int{
		"representantes": len(obtenerRepresentantes()),
		"terapeutas":     len(obtenerTerapeutas()),
		"planes":         len(obtenerPlanes()),
		"estudiantes":    len(obtenerEstudiantes()),
		"citas":          len(obtenerCitas()),
		"facturas":       len(obtenerFacturas()),
	})
}

// apiContar devuelve el numero de registros de una entidad.
func apiContar(w http.ResponseWriter, r *http.Request, entidad string) {
	if r.Method != http.MethodGet {
		metodoNoPermitido(w, r)
		return
	}
	escribirJSON(w, http.StatusOK, map[string]int{"total": len(listaComoMapas(entidad))})
}

// apiBuscar filtra los registros de una entidad por campo y valor (subcadena).
func apiBuscar(w http.ResponseWriter, r *http.Request, entidad string) {
	if r.Method != http.MethodGet {
		metodoNoPermitido(w, r)
		return
	}
	campo := r.URL.Query().Get("campo")
	valor := strings.ToLower(r.URL.Query().Get("valor"))
	if campo == "" {
		escribirJSON(w, http.StatusBadRequest, map[string]string{"error": "falta el parametro campo"})
		return
	}
	todos := listaComoMapas(entidad)
	filtrados := []map[string]any{}
	for _, f := range todos {
		if strings.Contains(strings.ToLower(fmt.Sprintf("%v", f[campo])), valor) {
			filtrados = append(filtrados, f)
		}
	}
	escribirJSON(w, http.StatusOK, filtrados)
}

// apiColeccion maneja GET (listar) y POST (crear) de una entidad.
func apiColeccion(w http.ResponseWriter, r *http.Request, entidad string) {
	switch entidad {
	case "representantes":
		switch r.Method {
		case http.MethodGet:
			escribirJSON(w, http.StatusOK, obtenerRepresentantes())
		case http.MethodPost:
			var rcp Representante
			if err := leerCuerpoJSON(r, &rcp); err != nil {
				escribirJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON invalido"})
				return
			}
			if err := validarRepresentante(rcp); err != nil {
				escribirJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			escribirJSON(w, http.StatusCreated, insertarRepresentante(rcp))
		default:
			metodoNoPermitido(w, r)
		}
	case "terapeutas":
		switch r.Method {
		case http.MethodGet:
			escribirJSON(w, http.StatusOK, obtenerTerapeutas())
		case http.MethodPost:
			var ter Terapeuta
			if err := leerCuerpoJSON(r, &ter); err != nil {
				escribirJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON invalido"})
				return
			}
			if err := validarTerapeuta(ter); err != nil {
				escribirJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			escribirJSON(w, http.StatusCreated, insertarTerapeuta(ter))
		default:
			metodoNoPermitido(w, r)
		}
	case "planes":
		switch r.Method {
		case http.MethodGet:
			escribirJSON(w, http.StatusOK, obtenerPlanes())
		case http.MethodPost:
			var pla Plan
			if err := leerCuerpoJSON(r, &pla); err != nil {
				escribirJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON invalido"})
				return
			}
			if err := validarPlan(pla); err != nil {
				escribirJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			escribirJSON(w, http.StatusCreated, insertarPlan(pla))
		default:
			metodoNoPermitido(w, r)
		}
	case "estudiantes":
		switch r.Method {
		case http.MethodGet:
			escribirJSON(w, http.StatusOK, obtenerEstudiantes())
		case http.MethodPost:
			var est Estudiante
			if err := leerCuerpoJSON(r, &est); err != nil {
				escribirJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON invalido"})
				return
			}
			if err := validarEstudiante(est); err != nil {
				escribirJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			escribirJSON(w, http.StatusCreated, insertarEstudiante(est))
		default:
			metodoNoPermitido(w, r)
		}
	case "citas":
		switch r.Method {
		case http.MethodGet:
			escribirJSON(w, http.StatusOK, obtenerCitas())
		case http.MethodPost:
			var cit Cita
			if err := leerCuerpoJSON(r, &cit); err != nil {
				escribirJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON invalido"})
				return
			}
			if err := validarCita(cit); err != nil {
				escribirJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			escribirJSON(w, http.StatusCreated, insertarCita(cit))
		default:
			metodoNoPermitido(w, r)
		}
	case "facturas":
		switch r.Method {
		case http.MethodGet:
			escribirJSON(w, http.StatusOK, obtenerFacturas())
		case http.MethodPost:
			var fac Factura
			if err := leerCuerpoJSON(r, &fac); err != nil {
				escribirJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON invalido"})
				return
			}
			if err := validarFactura(fac); err != nil {
				escribirJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			escribirJSON(w, http.StatusCreated, insertarFactura(fac))
		default:
			metodoNoPermitido(w, r)
		}
	}
}

// apiPorID maneja GET, PUT y DELETE de un elemento por su ID.
func apiPorID(w http.ResponseWriter, r *http.Request, entidad string, id int) {
	switch entidad {
	case "representantes":
		switch r.Method {
		case http.MethodGet:
			if v, ok := obtenerRepresentante(id); ok {
				escribirJSON(w, http.StatusOK, v)
			} else {
				escribirJSON(w, http.StatusNotFound, map[string]string{"error": "representante no encontrado"})
			}
		case http.MethodPut:
			var rcp Representante
			if err := leerCuerpoJSON(r, &rcp); err != nil {
				escribirJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON invalido"})
				return
			}
			if err := validarRepresentante(rcp); err != nil {
				escribirJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			if v, ok := actualizarRepresentante(id, rcp); ok {
				escribirJSON(w, http.StatusOK, v)
			} else {
				escribirJSON(w, http.StatusNotFound, map[string]string{"error": "representante no encontrado"})
			}
		case http.MethodDelete:
			if eliminarRepresentante(id) {
				escribirJSON(w, http.StatusOK, map[string]string{"mensaje": "representante eliminado"})
			} else {
				escribirJSON(w, http.StatusNotFound, map[string]string{"error": "representante no encontrado"})
			}
		default:
			metodoNoPermitido(w, r)
		}
	case "terapeutas":
		switch r.Method {
		case http.MethodGet:
			if v, ok := obtenerTerapeuta(id); ok {
				escribirJSON(w, http.StatusOK, v)
			} else {
				escribirJSON(w, http.StatusNotFound, map[string]string{"error": "terapeuta no encontrado"})
			}
		case http.MethodPut:
			var ter Terapeuta
			if err := leerCuerpoJSON(r, &ter); err != nil {
				escribirJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON invalido"})
				return
			}
			if err := validarTerapeuta(ter); err != nil {
				escribirJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			if v, ok := actualizarTerapeuta(id, ter); ok {
				escribirJSON(w, http.StatusOK, v)
			} else {
				escribirJSON(w, http.StatusNotFound, map[string]string{"error": "terapeuta no encontrado"})
			}
		case http.MethodDelete:
			if eliminarTerapeuta(id) {
				escribirJSON(w, http.StatusOK, map[string]string{"mensaje": "terapeuta eliminado"})
			} else {
				escribirJSON(w, http.StatusNotFound, map[string]string{"error": "terapeuta no encontrado"})
			}
		default:
			metodoNoPermitido(w, r)
		}
	case "planes":
		switch r.Method {
		case http.MethodGet:
			if v, ok := obtenerPlan(id); ok {
				escribirJSON(w, http.StatusOK, v)
			} else {
				escribirJSON(w, http.StatusNotFound, map[string]string{"error": "plan no encontrado"})
			}
		case http.MethodPut:
			var pla Plan
			if err := leerCuerpoJSON(r, &pla); err != nil {
				escribirJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON invalido"})
				return
			}
			if err := validarPlan(pla); err != nil {
				escribirJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			if v, ok := actualizarPlan(id, pla); ok {
				escribirJSON(w, http.StatusOK, v)
			} else {
				escribirJSON(w, http.StatusNotFound, map[string]string{"error": "plan no encontrado"})
			}
		case http.MethodDelete:
			if eliminarPlan(id) {
				escribirJSON(w, http.StatusOK, map[string]string{"mensaje": "plan eliminado"})
			} else {
				escribirJSON(w, http.StatusNotFound, map[string]string{"error": "plan no encontrado"})
			}
		default:
			metodoNoPermitido(w, r)
		}
	case "estudiantes":
		switch r.Method {
		case http.MethodGet:
			if v, ok := obtenerEstudiante(id); ok {
				escribirJSON(w, http.StatusOK, v)
			} else {
				escribirJSON(w, http.StatusNotFound, map[string]string{"error": "estudiante no encontrado"})
			}
		case http.MethodPut:
			var est Estudiante
			if err := leerCuerpoJSON(r, &est); err != nil {
				escribirJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON invalido"})
				return
			}
			if err := validarEstudiante(est); err != nil {
				escribirJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			if v, ok := actualizarEstudiante(id, est); ok {
				escribirJSON(w, http.StatusOK, v)
			} else {
				escribirJSON(w, http.StatusNotFound, map[string]string{"error": "estudiante no encontrado"})
			}
		case http.MethodDelete:
			if eliminarEstudiante(id) {
				escribirJSON(w, http.StatusOK, map[string]string{"mensaje": "estudiante eliminado"})
			} else {
				escribirJSON(w, http.StatusNotFound, map[string]string{"error": "estudiante no encontrado"})
			}
		default:
			metodoNoPermitido(w, r)
		}
	case "citas":
		switch r.Method {
		case http.MethodGet:
			if v, ok := obtenerCita(id); ok {
				escribirJSON(w, http.StatusOK, v)
			} else {
				escribirJSON(w, http.StatusNotFound, map[string]string{"error": "cita no encontrada"})
			}
		case http.MethodPut:
			var cit Cita
			if err := leerCuerpoJSON(r, &cit); err != nil {
				escribirJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON invalido"})
				return
			}
			if err := validarCita(cit); err != nil {
				escribirJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			if v, ok := actualizarCita(id, cit); ok {
				escribirJSON(w, http.StatusOK, v)
			} else {
				escribirJSON(w, http.StatusNotFound, map[string]string{"error": "cita no encontrada"})
			}
		case http.MethodDelete:
			if eliminarCita(id) {
				escribirJSON(w, http.StatusOK, map[string]string{"mensaje": "cita eliminada"})
			} else {
				escribirJSON(w, http.StatusNotFound, map[string]string{"error": "cita no encontrada"})
			}
		default:
			metodoNoPermitido(w, r)
		}
	case "facturas":
		switch r.Method {
		case http.MethodGet:
			if v, ok := obtenerFactura(id); ok {
				escribirJSON(w, http.StatusOK, v)
			} else {
				escribirJSON(w, http.StatusNotFound, map[string]string{"error": "factura no encontrada"})
			}
		case http.MethodPut:
			var fac Factura
			if err := leerCuerpoJSON(r, &fac); err != nil {
				escribirJSON(w, http.StatusBadRequest, map[string]string{"error": "JSON invalido"})
				return
			}
			if err := validarFactura(fac); err != nil {
				escribirJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
				return
			}
			if v, ok := actualizarFactura(id, fac); ok {
				escribirJSON(w, http.StatusOK, v)
			} else {
				escribirJSON(w, http.StatusNotFound, map[string]string{"error": "factura no encontrada"})
			}
		case http.MethodDelete:
			if eliminarFactura(id) {
				escribirJSON(w, http.StatusOK, map[string]string{"mensaje": "factura eliminada"})
			} else {
				escribirJSON(w, http.StatusNotFound, map[string]string{"error": "factura no encontrada"})
			}
		default:
			metodoNoPermitido(w, r)
		}
	}
}

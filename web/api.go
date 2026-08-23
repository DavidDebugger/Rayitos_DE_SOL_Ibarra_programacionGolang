package web

// api.go expone las entidades como API JSON.
// Las rutas son:
//   GET  /api/{entidad}        -> lista
//   POST /api/{entidad}        -> crea
//   GET  /api/{entidad}/{id}   -> uno
//   PUT  /api/{entidad}/{id}   -> reemplaza
//   DELETE /api/{entidad}/{id} -> borra

import (
	"encoding/json"
	"io"
	"net/http"
)

// escribirJSON envia un codigo HTTP y un dato cualquiera como JSON.
func escribirJSON(w http.ResponseWriter, codigo int, datos any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codigo)
	_ = json.NewEncoder(w).Encode(datos)
}

// leerCuerpoJSON lee el cuerpo de la peticion y lo convierte a JSON.
func leerCuerpoJSON(r *http.Request, destino any) error {
	cuerpo, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(cuerpo, destino)
}

// metodoNoPermitido responde cuando usan un verbo HTTP que no soportamos.
func metodoNoPermitido(w http.ResponseWriter, r *http.Request) {
	escribirJSON(w, http.StatusMethodNotAllowed, map[string]string{
		"error": "metodo " + r.Method + " no permitido en esta ruta",
	})
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

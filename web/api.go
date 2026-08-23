/*
 api.go
 Implementa los servicios web (API) que intercambian datos en formato JSON
 en lugar de paginas HTML. La API expone las 6 entidades del sistema:
 representantes, terapeutas, planes, estudiantes, citas y facturas.
 Rutas por entidad (el nombre va en plural en la URL):
   GET    /api/{entidad}        -> lista todos los registros
   POST   /api/{entidad}        -> crea un registro (JSON en el cuerpo)
   GET    /api/{entidad}/{id}   -> un registro especifico
   PUT    /api/{entidad}/{id}   -> actualiza un registro
   DELETE /api/{entidad}/{id}   -> borra un registro (si no tiene referencias)
 Metodos HTTP:
   GET  = lee datos sin modificar nada.
   POST = crea un registro a partir del JSON del cuerpo.
   PUT  = reemplaza o actualiza un registro existente.
   DELETE = elimina un registro.
 Toda la validacion vive en reglas.go (crearX / actualizarX / puedeBorrar),
 asi la API y el portal comparten exactamente las mismas reglas.
 */

package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
)

/*
 mapaEntidadSingular traduce el nombre plural usado en la ruta al nombre
 singular que usa el almacen (ej. "representantes" -> "representante").
 Asi una sola funcion puede atender las 6 entidades por id.
 */
var mapaEntidadSingular = map[string]string{
	"representantes": "representante",
	"estudiantes":    "estudiante",
	"citas":          "cita",
	"facturas":       "factura",
	"planes":         "plan",
	"terapeutas":     "terapeuta",
}

/*
 UI (Interfaz de Usuario, aunque aqui es la interfaz del API) guarda el
 almacen para que los handlers puedan leer y guardar datos. Se usa como
 "receptor" para ligar las funciones al enrutador en mux.go.
 */
type UI struct {
	a *Almacen
}

/*
 escribirJSON envia cualquier dato como JSON al cliente.
   1. Avisa que el contenido es de tipo JSON.
   2. Pone el codigo de estado (200 OK, 201 Creado, 400 Error, etc.).
   3. Convierte el dato a JSON y lo escribe en la respuesta.
 */
func escribirJSON(w http.ResponseWriter, estado int, dato interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(estado)
	json.NewEncoder(w).Encode(dato)
}

/*
 leerJSON lee el cuerpo de la peticion (el JSON que envio el cliente) y lo
 convierte a la struct de Golang que le pasemos. Por ejemplo, convierte el
 JSON de un representante en una variable tipo Representante.
 */
func leerJSON(r *http.Request, dato interface{}) error {
	if r.Body == nil {
		return errors.New("cuerpo vacio")
	}
	return json.NewDecoder(r.Body).Decode(dato)
}

/*
 respuestaError devuelve un JSON con forma {"error": "mensaje"} y un estado
 de error, para informar al cliente que algo salio mal.
 */
func respuestaError(w http.ResponseWriter, estado int, mensaje string) {
	escribirJSON(w, estado, map[string]string{"error": mensaje})
}

/*
 apiRepresentantes maneja las peticiones a /api/representantes.
   GET  -> devuelve la lista de representantes en JSON.
   POST -> lee un representante del cuerpo, lo crea y devuelve el nuevo registro.
 */
func (ui *UI) apiRepresentantes(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		escribirJSON(w, http.StatusOK, ui.a.Listar("representante"))
		return
	}
	if r.Method == http.MethodPost {
		var e Representante
		if err := leerJSON(r, &e); err != nil {
			respuestaError(w, http.StatusBadRequest, "JSON invalido")
			return
		}
		// La validacion y el guardado estan en reglas.go (crearRepresentante).
		// Si falla (ej. falta la cedula) devolvemos el error como JSON.
		if err := crearRepresentante(ui.a, &e); err != nil {
			respuestaError(w, http.StatusBadRequest, err.Error())
			return
		}
		escribirJSON(w, http.StatusCreated, e)
		return
	}
	respuestaError(w, http.StatusMethodNotAllowed, "metodo no permitido")
}

// apiEstudiantes maneja /api/estudiantes: lista (GET) o crea (POST) estudiantes.
func (ui *UI) apiEstudiantes(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		escribirJSON(w, http.StatusOK, ui.a.Listar("estudiante"))
		return
	}
	if r.Method == http.MethodPost {
		var e Estudiante
		if err := leerJSON(r, &e); err != nil {
			respuestaError(w, http.StatusBadRequest, "JSON invalido")
			return
		}
		// Validacion y guardado en reglas.go (crearEstudiante).
		if err := crearEstudiante(ui.a, &e); err != nil {
			respuestaError(w, http.StatusBadRequest, err.Error())
			return
		}
		escribirJSON(w, http.StatusCreated, e)
		return
	}
	respuestaError(w, http.StatusMethodNotAllowed, "metodo no permitido")
}

// apiCitas maneja /api/citas: lista (GET) o crea (POST) citas.
func (ui *UI) apiCitas(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		escribirJSON(w, http.StatusOK, ui.a.Listar("cita"))
		return
	}
	if r.Method == http.MethodPost {
		var c Cita
		if err := leerJSON(r, &c); err != nil {
			respuestaError(w, http.StatusBadRequest, "JSON invalido")
			return
		}
		// Validacion y guardado en reglas.go (crearCita).
		if err := crearCita(ui.a, &c); err != nil {
			respuestaError(w, http.StatusBadRequest, err.Error())
			return
		}
		escribirJSON(w, http.StatusCreated, c)
		return
	}
	respuestaError(w, http.StatusMethodNotAllowed, "metodo no permitido")
}

// apiFacturas maneja /api/facturas: lista (GET) o crea (POST) facturas.
func (ui *UI) apiFacturas(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		escribirJSON(w, http.StatusOK, ui.a.Listar("factura"))
		return
	}
	if r.Method == http.MethodPost {
		var f Factura
		if err := leerJSON(r, &f); err != nil {
			respuestaError(w, http.StatusBadRequest, "JSON invalido")
			return
		}
		// Validacion y guardado en reglas.go (crearFactura).
		if err := crearFactura(ui.a, &f); err != nil {
			respuestaError(w, http.StatusBadRequest, err.Error())
			return
		}
		escribirJSON(w, http.StatusCreated, f)
		return
	}
	respuestaError(w, http.StatusMethodNotAllowed, "metodo no permitido")
}

/*
 apiPlanes maneja /api/planes: lista (GET) o crea (POST) planes.
 Los planes tambien forman parte del CRUD (no solo sembrados).
 */
func (ui *UI) apiPlanes(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		escribirJSON(w, http.StatusOK, ui.a.Listar("plan"))
		return
	}
	if r.Method == http.MethodPost {
		var p Plan
		if err := leerJSON(r, &p); err != nil {
			respuestaError(w, http.StatusBadRequest, "JSON invalido")
			return
		}
		// Validacion y guardado en reglas.go (crearPlan).
		if err := crearPlan(ui.a, &p); err != nil {
			respuestaError(w, http.StatusBadRequest, err.Error())
			return
		}
		escribirJSON(w, http.StatusCreated, p)
		return
	}
	respuestaError(w, http.StatusMethodNotAllowed, "metodo no permitido")
}

/*
 apiTerapeutas maneja /api/terapeutas: lista (GET) o crea (POST) terapeutas.
 Los terapeutas tambien forman parte del CRUD (no solo sembrados).
 */
func (ui *UI) apiTerapeutas(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		escribirJSON(w, http.StatusOK, ui.a.Listar("terapeuta"))
		return
	}
	if r.Method == http.MethodPost {
		var t Terapeuta
		if err := leerJSON(r, &t); err != nil {
			respuestaError(w, http.StatusBadRequest, "JSON invalido")
			return
		}
		// Validacion y guardado en reglas.go (crearTerapeuta).
		if err := crearTerapeuta(ui.a, &t); err != nil {
			respuestaError(w, http.StatusBadRequest, err.Error())
			return
		}
		escribirJSON(w, http.StatusCreated, t)
		return
	}
	respuestaError(w, http.StatusMethodNotAllowed, "metodo no permitido")
}

/*
 apiPorID maneja las operaciones por id de una entidad:
   GET    /api/{entidad}/{id}  -> devuelve un registro especifico
   PUT    /api/{entidad}/{id}  -> actualiza un registro (JSON en el cuerpo)
   DELETE /api/{entidad}/{id}  -> borra un registro (si no tiene referencias)
 Todas las entidades usan esta misma funcion; se distinguen por la ruta.
 */
func (ui *UI) apiPorID(w http.ResponseWriter, r *http.Request) {
	// El nombre plural viene de la ruta; lo traducimos al singular del almacen.
	plural := r.PathValue("entidad")
	tipo, ok := mapaEntidadSingular[plural]
	if !ok {
		respuestaError(w, http.StatusNotFound, "entidad desconocida")
		return
	}
	// Convertimos el id de texto a numero; si no es numero, error.
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		respuestaError(w, http.StatusBadRequest, "id invalido")
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Devolvemos el registro; si no existe, 404.
		v, ok := ui.a.Obtener(tipo, id)
		if !ok {
			respuestaError(w, http.StatusNotFound, "no encontrado")
			return
		}
		escribirJSON(w, http.StatusOK, v)

	case http.MethodPut:
		// Actualizamos usando las reglas de negocio (actualizar* de reglas.go).
		// Cada entidad decodifica su propia struct desde el JSON del cuerpo.
		var err2 error
		switch tipo {
		case "representante":
			var e Representante
			if err2 = leerJSON(r, &e); err2 == nil {
				err2 = actualizarRepresentante(ui.a, id, e)
			}
		case "terapeuta":
			var e Terapeuta
			if err2 = leerJSON(r, &e); err2 == nil {
				err2 = actualizarTerapeuta(ui.a, id, e)
			}
		case "plan":
			var p Plan
			if err2 = leerJSON(r, &p); err2 == nil {
				err2 = actualizarPlan(ui.a, id, p)
			}
		case "estudiante":
			var e Estudiante
			if err2 = leerJSON(r, &e); err2 == nil {
				err2 = actualizarEstudiante(ui.a, id, e)
			}
		case "cita":
			var c Cita
			if err2 = leerJSON(r, &c); err2 == nil {
				err2 = actualizarCita(ui.a, id, c)
			}
		case "factura":
			var f Factura
			if err2 = leerJSON(r, &f); err2 == nil {
				err2 = actualizarFactura(ui.a, id, f)
			}
		}
		if err2 != nil {
			respuestaError(w, http.StatusBadRequest, err2.Error())
			return
		}
		// Devolvemos el registro ya actualizado (lo volvemos a leer).
		if v, ok := ui.a.Obtener(tipo, id); ok {
			escribirJSON(w, http.StatusOK, v)
			return
		}
		escribirJSON(w, http.StatusOK, map[string]int{"id": id})

	case http.MethodDelete:
		// Antes de borrar, revisamos que no deje referencias rotas.
		ok2, mensaje := puedeBorrar(ui.a, tipo, id)
		if !ok2 {
			respuestaError(w, http.StatusConflict, mensaje)
			return
		}
		ui.a.Borrar(tipo, id)
		escribirJSON(w, http.StatusOK, map[string]interface{}{"ok": true, "id": id})

	default:
		respuestaError(w, http.StatusMethodNotAllowed, "metodo no permitido")
	}
}

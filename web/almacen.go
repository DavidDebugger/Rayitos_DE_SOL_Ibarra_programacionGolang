/*
 almacen.go
 Implementa el almacen de datos en memoria, que hace las veces de base de datos
 del sistema. Los datos solo existen mientras el servidor esta encendido.
 EN MEMORIA:
   Los datos viven en variables del programa. Al apagar el servidor se pierden.
   Para un proyecto estudiantil esto es suficiente y evita instalar y configurar
   una base de datos real (MySQL, PostgreSQL, etc.).
 CONCURRENCIA:
   Un servidor web atiende muchas peticiones a la vez. Si dos escrituras ocurren
   al mismo tiempo podrian corromperse. Por eso cada metodo abre y cierra un
   candado (sync.Mutex) que serializa el acceso: solo un proceso lee o escribe
   a la vez. Es el mecanismo estandar y seguro en Go.
 */

package web

import "sync"

/*
 Almacen guarda todos los datos del sistema.
   - datos:       mapa de tipo de entidad -> (mapa de id -> dato).
                  Por ejemplo datos["estudiante"][3] es el estudiante con id 3.
   - siguienteID: mapa de tipo de entidad -> proximo id a entregar.
                  Cada entidad lleva su propia cuenta para que los ids no se
                  mezclen entre entidades distintas.
   - mux:         candado que protege el acceso concurrente.
 */
type Almacen struct {
	mux         sync.Mutex                     // candado que serializa el acceso concurrente
	datos       map[string]map[int]interface{} // tipo -> (id -> dato)
	siguienteID map[string]int                 // tipo -> proximo id a entregar
}

/*
 NuevoAlmacen crea un almacen vacio, listo para guardar informacion.
 make() inicializa los mapas, porque un mapa declarado pero no inicializado
 no permite guardar datos.
 */
func NuevoAlmacen() *Almacen {
	return &Almacen{
		datos:       make(map[string]map[int]interface{}),
		siguienteID: make(map[string]int),
	}
}

/*
 nuevoID devuelve el siguiente id disponible para un tipo de entidad
 (por ejemplo "estudiante" o "cita"). Cada tipo tiene su propia cuenta, asi
 los ids de estudiante y de cita pueden empezar en 1 sin solaparse.
 Abre y cierra el candado dentro de la funcion para ser seguro en concurrencia.
 */
func (a *Almacen) nuevoID(tipo string) int {
	a.mux.Lock()
	defer a.mux.Unlock()
	a.siguienteID[tipo]++
	return a.siguienteID[tipo]
}

/*
 Guardar anota un dato dentro de su seccion, usando su id como clave.
 Si la seccion todavia no existe, la crea en el momento con make().
 */
func (a *Almacen) Guardar(tipo string, id int, item interface{}) {
	a.mux.Lock()
	defer a.mux.Unlock()
	if a.datos[tipo] == nil {
		a.datos[tipo] = make(map[int]interface{})
	}
	a.datos[tipo][id] = item
}

/*
 Obtener devuelve un dato por su tipo e id.
 El segundo valor indica si se encontro (true) o no (false).
 */
func (a *Almacen) Obtener(tipo string, id int) (interface{}, bool) {
	a.mux.Lock()
	defer a.mux.Unlock()
	v, ok := a.datos[tipo][id]
	return v, ok
}

/*
 Listar devuelve todos los datos de una seccion como una lista (slice).
 Se usa para pintar las tablas del portal y para resolver nombres al momento
 de mostrar llaves foraneas (por ejemplo, el nombre de un tutor).
 */
func (a *Almacen) Listar(tipo string) []interface{} {
	a.mux.Lock()
	defer a.mux.Unlock()
	var res []interface{}
	for _, v := range a.datos[tipo] {
		res = append(res, v)
	}
	return res
}

// Borrar elimina un registro de la seccion usando su id.
func (a *Almacen) Borrar(tipo string, id int) {
	a.mux.Lock()
	defer a.mux.Unlock()
	delete(a.datos[tipo], id)
}

/*
 Existe indica si ya hay un registro con ese id dentro de la seccion.
 Se usa al validar que una llave foranea (ej. representante_id) apunte a
 algo que realmente existe en el almacen.
 */
func (a *Almacen) Existe(tipo string, id int) bool {
	a.mux.Lock()
	defer a.mux.Unlock()
	_, ok := a.datos[tipo][id]
	return ok
}

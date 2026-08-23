/*
 cargar.go
 Este archivo se encarga de leer el archivo datos.json al momento de iniciar
 el servidor y volcar su contenido dentro del almacen en memoria.
 POR QUE UN ARCHIVO JSON:
   Nos permite tener datos de ejemplo listos para probar el sistema sin tener
   que crear todo a mano desde el portal. Tambien es una forma sencilla de
   "sembrar" informacion inicial en un proyecto estudiantil.
 QUE PASA SI NO EXISTE EL ARCHIVO:
   La funcion devuelve un error. Quien la llama decide que hacer (en nuestro
   caso, avisar y seguir; sembrar() cubre lo mas basico).
 */

package web

import (
	"encoding/json"
	"os"
)

/*
 datosJSON describe exactamente la forma que debe tener el archivo datos.json.
 Cada campo usa una etiqueta json:"..." que indica el nombre de la clave en
 el archivo. Asi Go sabe mapear el texto JSON a estas estructuras.
 */
type datosJSON struct {
	Representantes []Representante `json:"representantes"`
	Terapeutas     []Terapeuta     `json:"terapeutas"`
	Planes         []Plan          `json:"planes"`
	Estudiantes    []Estudiante    `json:"estudiantes"`
	Citas          []Cita          `json:"citas"`
	Facturas       []Factura       `json:"facturas"`
}

/*
 cargarDatos lee el archivo ubicado en "ruta" y guarda cada registro dentro
 del almacen, respetando los ids que trae el archivo. Al final ajusta el
 contador de ids para que las nuevas creaciones usen numeros mayores y no
 choquen con los que ya cargamos.
 */
func cargarDatos(a *Almacen, ruta string) error {
	// 1) Leemos todo el contenido del archivo como bytes (secuencia de bytes).
	b, err := os.ReadFile(ruta)
	if err != nil {
		return err // el archivo no existe o no se puede leer
	}

	// 2) Convertimos el texto JSON a nuestras estructuras de Go (datosJSON).
	var d datosJSON
	if err := json.Unmarshal(b, &d); err != nil {
		return err // el JSON esta mal formado
	}

	// maxID recordara el id mas alto que encontramos en el archivo.
	maxID := 0

	/*
	 guardar es una funcion interna (closure) que anota cada reue anota cada registro en su
	 seccion y actualiza maxID si el registro trae un id mas grande.
	 */
	guardar := func(tipo string, id int, item interface{}) {
		a.Guardar(tipo, id, item)
		if id > maxID {
			maxID = id
		}
	}

	// 3) Guardamos entidad por entidad en el almacen.
	for _, v := range d.Representantes {
		guardar("representante", v.ID, v)
	}
	for _, v := range d.Terapeutas {
		guardar("terapeuta", v.ID, v)
	}
	for _, v := range d.Planes {
		guardar("plan", v.ID, v)
	}
	for _, v := range d.Estudiantes {
		guardar("estudiante", v.ID, v)
	}
	for _, v := range d.Citas {
		guardar("cita", v.ID, v)
	}
	for _, v := range d.Facturas {
		guardar("factura", v.ID, v)
	}

	// 4) Ajustamos el contador de ids para que las proximas creaciones usen
	// numeros mayores al id mas alto que ya cargo el archivo, evitando asi
	// colisiones. Se aplica a todos los tipos conocidos del modelo.
	for _, t := range []string{"representante", "terapeuta", "plan", "estudiante", "cita", "factura", "pago"} {
		if maxID > a.siguienteID[t] {
			a.siguienteID[t] = maxID
		}
	}
	return nil
}

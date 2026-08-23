package web

// cargar.go lee datos.json al iniciar y los mete en el almacen en memoria.
// Solo se lee una vez; durante la ejecucion los cambios viven en memoria.

import (
	"encoding/json"
	"fmt"
	"os"
)

// datosArchivo es el molde del archivo datos.json.
type datosArchivo struct {
	Representantes []Representante `json:"representantes"`
	Terapeutas    []Terapeuta     `json:"terapeutas"`
	Planes        []Plan          `json:"planes"`
	Estudiantes   []Estudiante    `json:"estudiantes"`
	Citas         []Cita          `json:"citas"`
	Facturas      []Factura       `json:"facturas"`
}

// CargarDatos lee el archivo y llena el almacen. Si el archivo no existe,
// simplemente deja el almacen vacio y avisa (no es un error fatal).
func CargarDatos(ruta string) error {
	b, err := os.ReadFile(ruta)
	if err != nil {
		return fmt.Errorf("no se pudo leer %s: %w", ruta, err)
	}

	var d datosArchivo
	if err := json.Unmarshal(b, &d); err != nil {
		return fmt.Errorf("el archivo %s no es JSON valido: %w", ruta, err)
	}

	mu.Lock()
	defer mu.Unlock()

	for _, r := range d.Representantes {
		representantes[r.ID] = r
		ajustarContador(r.ID)
	}
	for _, t := range d.Terapeutas {
		terapeutas[t.ID] = t
		ajustarContador(t.ID)
	}
	for _, p := range d.Planes {
		planes[p.ID] = p
		ajustarContador(p.ID)
	}
	for _, e := range d.Estudiantes {
		estudiantes[e.ID] = e
		ajustarContador(e.ID)
	}
	for _, c := range d.Citas {
		citas[c.ID] = c
		ajustarContador(c.ID)
	}
	for _, f := range d.Facturas {
		facturas[f.ID] = f
		ajustarContador(f.ID)
	}
	return nil
}

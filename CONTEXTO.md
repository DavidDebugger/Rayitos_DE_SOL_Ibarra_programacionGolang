# CONTEXTO.md - Instruccion de contexto para sesiones futuras

Este archivo es una GUIA COMPLETA para que un agente (o tu en otra sesion)
comprenda el proyecto **Rayitos de Sol Ibarra (RDSI)** y lo mantenga igual:
misma estructura, mismas reglas, mismo estilo. Leelo entero antes de tocar codigo.

NO es documentacion de usuario. Es una instruccion de preservacion de contexto.

---

## 1. Que es el proyecto

Sistema de gestion web para un **centro de terapias infantiles** en **Ibarra, Ecuador**.
Todo se usa desde el navegador; NO hay CLI.

Principio rector: **KISS (Keep It Simple, Stupid)** y **Golang como lenguaje principal**.
- Backend y FRONTEND generados en **Golang con SOLO la libreria estandar**
  (`net/http`, `html/template`, `embed`, `encoding/json`, `strconv`, etc.).
- La interfaz se genera en el servidor (**SSR**) con `html/template` y `go:embed`.
  NO hay JavaScript en el cliente: los formularios se envian por `POST` y Go
  valida, guarda y vuelve a pintar la pagina.
- Almacenamiento **en memoria** (se pierde al cerrar el servidor). Sin base de datos.
- Estilos **monocromos** (blanco/negro/gris) dentro de la propia plantilla HTML.
  Sin emojis, sin acentos en codigo ni en textos visibles.
- Ecuador: se usa **cedula** (no DNI) y la **factura se emite al REPRESENTANTE**.

El sistema demuestra tres pilares pedidos por el usuario:
- **Structs**: las entidades estan modeladas como `struct` en `web/modelo.go`.
- **Concurrencia**: el almacen protege el acceso con `sync.Mutex` (`web/almacen.go`).
- **Servicios web**: `web/api.go` expone 8 servicios JSON y `web/portal.go` un portal.

---

## 2. Reglas duras (NUNCA las violes)

1. **Solo libreria estandar de Go.** PROHIBIDO `go get`, frameworks, gin, GORM, etc.
2. **Sin tildes/acentos** en codigo, comentarios, plantillas ni mensajes. Escribe
   `n` para enie, `Cedula`, `Programada`, `Recuperacion`, etc. sin acento.
3. **Sin JavaScript en el cliente.** Toda la interfaz se genera en Go
   (`html/template` + `go:embed` en `web/templates/`). Los formularios usan
   `method="POST"` a rutas de Go; no hay `fetch`, `XMLHttpRequest` ni `onclick`.
4. **Sin emojis** en ningun archivo (codigo, plantillas, markdown).
5. Todo en **espanol**.
6. **Comentarios en espanol, sin acentos ni emojis.** En los comentarios largos
   usa estilo bloque `/* ... */`; en los cortos basta `//`. El codigo de `web/`
   lleva comentarios explicativos (pedidos por el usuario). Mantenlos al modificar.
7. **No borrar codigo Golang en uso** sin preguntar.
8. **No hacer commits** salvo que el usuario lo pida.

---

## 3. Estructura real del proyecto

```
main.go                       -- entrada: import "RDSI/web", llama web.IniciarServidor(":8080")
web/
  modelo.go                   -- structs de las entidades + constantes de estados
  almacen.go                  -- almacen en memoria con sync.Mutex (concurrencia)
  cargar.go                   -- lee datos.json al iniciar y los guarda en el almacen
  reglas.go                   -- crearRepresentante/Estudiante/Cita/Factura (validacion)
  api.go                      -- 8 servicios web JSON (GET y POST de cada entidad)
  portal.go                   -- portal web SSR sin JavaScript (usa la plantilla)
  mux.go                      -- IniciarServidor: crea almacen, carga datos, registra rutas
  templates/
    portal.html               -- plantilla HTML del portal (embebida con go:embed)
datos.json                    -- datos de ejemplo cargados al iniciar el servidor
docs/
  AGENTS.md                   -- reglas resumidas para agentes
CONTEXTO.md                  -- ESTE archivo
```

El backend es un **solo paquete** `web` importado por `main.go` (`import "RDSI/web"`).
Las plantillas se embeben en el binario con `//go:embed templates/portal.html`, por
lo que no dependen del directorio de ejecucion.

---

## 4. Modelo de datos (web/modelo.go)

Todas las structs tienen `ID int` y tags JSON en **snake_case**. Campos de relacion
terminan en `_id` (ej. `representante_id`). Constantes de estados en la cabecera.

- **Representante** { id, nombres, apellidos, cedula, telefono }
- **Terapeuta** { id, nombres, apellidos, cedula, telefono }
- **Estudiante** { id, nombres, apellidos, cedula, telefono,
  `representante_id` (FK), `plan_id` (FK) }
- **Plan** { id, nombre, descripcion, duracion }
- **Cita** { id, `estudiante_id`, `terapeuta_id`, `plan_id`, fecha, hora, estado }
  (el `plan_id` se copia del estudiante al agendar; el formulario NO lo pide)
- **Factura** { id, `representante_id` (FK Representante, NO estudiante), fecha, monto(int), estado }
- **Pago** { id, `factura_id`, monto, fecha, metodo }

Decisiones de negocio (no cambiar sin preguntar):
- La factura va al **representante** (`representante_id`), no al estudiante.
- KISS: personas solo tienen cedula + relaciones, sin correo/direccion.
- Estados de Cita: Programada, Confirmada, Completada, Reprogramada.
- Estados de Factura: Pendiente, Pagada, Anulada.

---

## 5. Backend: patrones a respetar

### Almacen (web/almacen.go) - concurrencia
- `Almacen` guarda `datos map[string]map[int]interface{}` (tipo -> id -> dato) y
  `siguienteID map[string]int` (contador por tipo).
- Cada metodo (`nuevoID`, `Guardar`, `Obtener`, `Listar`, `Borrar`, `Existe`)
  abre y cierra el candado `sync.Mutex` (`mux.Lock()` / `defer mux.Unlock()`).
  Esto hace el almacen seguro ante peticiones simultaneas (concurrencia).
- `NuevoAlmacen()` crea el almacen vacio.

### Carga inicial (web/cargar.go)
- `cargarDatos(a, "datos.json")` lee el JSON y guarda cada registro con su id.
  Luego ajusta `siguienteID` para que las creaciones nuevas no choquen con los ids
  cargados. Devuelve error si el archivo falta o esta mal formado.

### Reglas de negocio (web/reglas.go)
- Para las 6 entidades (`representante`, `terapeuta`, `plan`, `estudiante`, `cita`,
  `factura`) existen `crearX` y `actualizarX` (reciben la entidad por **puntero** para
  asignar el id y reflejarlo en quien llama) que validan campos obligatorios y llaves
  foraneas, asignan id y guardan.
- `puedeBorrar(a, tipo, id)` devuelve `(false, mensaje)` si el registro esta en uso
  (ej. un plan con estudiantes) para evitar referencias rotas; lo usan la API (DELETE)
  y el portal (borrar).
- Estas funciones se reutilizan TANTO en la API (`api.go`) como en el portal
  (`portal.go`), asi la validacion existe en un solo lugar.

### API web (web/api.go) - servicios web
- Tipo `UI` con los handlers. Helpers: `escribirJSON`, `leerJSON`, `respuestaError`.
- Cada entidad (incluidas `planes` y `terapeutas`) tiene `apiX(w, r)` que responde:
  - `GET`  -> lista la entidad en JSON (`ui.a.Listar(...)`).
  - `POST` -> lee el JSON del cuerpo, llama a `crearX`, devuelve 201 o error 400.
- Por id, un solo handler `apiPorID` atiende las 6 entidades via la ruta
  `/api/{entidad}/{id}` (mapea el plural al singular con `mapaEntidadSingular`):
  - `GET`    -> un registro (404 si no existe).
  - `PUT`    -> actualiza con `actualizarX` (refleja el JSON actualizado).
  - `DELETE` -> borra si `puedeBorrar` lo permite (409 si esta en uso).

### Portal SSR (web/portal.go)
- `portalInicio` (GET) construye `PortalData` (listas de las 6 entidades, opciones de
  selects, JSON de ejemplo) y renderiza `templates/portal.html` con `template.ParseFS`.
- `portalCrear` (POST `/portal/crear/{entidad}`) lee el formulario con `formStr`/
  `formInt`, llama a `crearX` y redirige a `/portal?msg=...`.
- `portalEditar` (GET `/portal/editar/{entidad}/{id}`) muestra un formulario precargado
  (armado en `formularioEdicion`, sin JS) para editar un registro.
- `portalActualizar` (POST `/portal/actualizar/{entidad}/{id}`) aplica los cambios con
  `actualizarX` y redirige al portal.
- `portalBorrar` (POST `/portal/borrar/{entidad}/{id}`) borra si `puedeBorrar` lo
  permite; si no, redirige con el aviso del conflicto.
- Helpers: `aRepresentantes` (y similares) convierten la lista generica a tipos
  concretos; `optsRepresentante` (y similares) arman las opciones de los `<select>`.

### Mux (web/mux.go)
- `IniciarServidor(puerto)`:
  1. `NuevoAlmacen()`
  2. `cargarDatos(a, "datos.json")` (avisa si falla, pero sigue)
  3. `sembrar(a)` crea Plan y Terapeuta por defecto si no existen (respaldo)
  4. registra rutas en `http.NewServeMux()`
  5. `http.ListenAndServe(puerto, m)`
- El servidor corre en **:8080**. Arranca con `go run .` (NO `go run main.go`,
  porque hay un paquete `web` aparte y ese comando no lo compilaria).

---

## 6. Rutas (estables)

Paginas (GET, generadas en Go, sin JS):
- `/` y `/portal` -> portal con tablas (6 entidades), formularios de creacion/edicion y
  JSON de ejemplo.
- `/portal/editar/{entidad}/{id}` -> formulario de edicion precargado.

Acciones (POST):
- `/portal/crear/{entidad}`      -> crea la entidad, validando en el servidor.
- `/portal/actualizar/{entidad}/{id}` -> aplica edicion, validando en el servidor.
- `/portal/borrar/{entidad}/{id}`     -> borra si no deja referencias rotas.
- Cada POST redirige a `/portal?msg=...` o `/portal?msg=<error>` (texto en espanol).

API (JSON):
- `GET/POST /api/representantes`, `/api/terapeutas`, `/api/planes`, `/api/estudiantes`,
  `/api/citas`, `/api/facturas`.
- `GET/PUT/DELETE /api/{entidad}/{id}` para operar sobre un registro especifico.
- `{entidad}` es el plural de: representantes, terapeutas, planes, estudiantes, citas,
  facturas.

---

## 7. Plantilla (web/templates/portal.html)

Todo el HTML se genera en Go. No hay JS. La plantilla usa `html/template`
(escapa automaticamente) y se embebe con `//go:embed`. Por cada entidad muestra:
una tabla con los registros, un formulario para crear, y el JSON que devuelve su
servicio `/api`. Los formularios son `method="POST"`.

---

## 8. Flujo de trabajo (siempre al terminar)

```bash
cd /ruta/al/proyecto
go build ./...     # debe compilar sin errores
go vet ./...       # sin advertencias
```

Para probar en navegador: `go run .` y abrir `http://localhost:8080`.
Datos en memoria: se reinician al reiniciar el server (salvo los de datos.json).

### Prueba rapida (curl, sin dejar servidor corriendo)
```bash
go build -o /tmp/rdsi-bin .
( /tmp/rdsi-bin > /tmp/srv.log 2>&1 & )
sleep 2
curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/portal
curl -s -o /dev/null -w "%{http_code}\n" -X POST -d '{"nombres":"Ana","cedula":"1"}' http://localhost:8080/api/representantes
curl -s -o /dev/null -w "%{http_code}\n" -d "nombres=Nino&cedula=3&representante_id=1&plan_id=1" http://localhost:8080/portal/crear/estudiante
fuser -k 8080/tcp   # matar el servidor al terminar
```
Crear representante -> estudiante (con `representante_id` y `plan_id`) -> cita -> factura.
Verificar que referencias inexistentes dan error (redirect 303 con `?error=`).

### CUIDADO: servidor viejo en :8080
Si al arrancar sale `bind: address already in use`, hay un servidor anterior. Matarlo
por puerto (NO usar `pkill -f go` porque mata la propia shell): `fuser -k 8080/tcp`.
No dejar servidores corriendo al terminar la sesion.

---

## 9. Checklist antes de dar por terminado un cambio

- [ ] `go build ./...` y `go vet ./...` pasan.
- [ ] Sin acentos ni emojis en ningun archivo nuevo/modificado.
- [ ] Toda la UI se genera en Go (html/template + go:embed); sin JavaScript en el cliente.
- [ ] Toda relacion `_id` se valida en `crearX` (error si no existe).
- [ ] El almacen sigue protegido por `sync.Mutex` (concurrencia).
- [ ] No se borro codigo Golang en uso.

---

## 10. Resumen de decisiones tomadas (no revertir sin preguntar)

1. Factura emitida al **representante** (`representante_id`), no al estudiante.
2. Personas solo con **cedula** (campo `cedula`), no DNI.
3. Estudiante ligado a **un** representante y **un** plan.
4. Interfaz **generada en Go (SSR)** con `html/template` + `go:embed`; sin JS cliente.
5. Backend 100% libreria estandar; almacenamiento en memoria con `sync.Mutex`.
6. Estructura por carpetas: `main.go` + paquete `web/` (modelo, almacen, carga,
   reglas, api, portal, mux, templates). Datos de ejemplo en `datos.json`.
7. Comentarios largos en estilo bloque `/* ... */`; cortos en `//`.

# AGENTS.md - Guia para agentes que trabajen en este proyecto

Sistema de gestion **Rayitos de Sol Ibarra (RDSI)**: centro de terapia infantil.
Todo funciona desde el navegador. **Golang es el lenguaje principal**: genera tanto
la logica como la interfaz (SSR, sin JavaScript en el cliente).

El proyecto cumple tres pilares pedidos:
- **Structs**: entidades modeladas como `struct` en `web/modelo.go`.
- **Concurrencia**: el almacen usa `sync.Mutex` en `web/almacen.go`.
- **Servicios web**: 8 servicios JSON en `web/api.go` y un portal en `web/portal.go`.

## Comandos utiles (SIEMPRE verificar al terminar)

```bash
go build ./...        # debe compilar sin errores
go vet ./...          # sin advertencias
```

Para ejecutar: `go run .` (NO `go run main.go`, porque existe el paquete `web`).

## Estructura del proyecto

```
main.go                    -- punto de entrada, solo levanta el servidor
web/                       -- unico paquete Go, importado por main.go (import "RDSI/web")
  modelo.go               -- structs de entidades + constantes de estados
  almacen.go              -- almacen en memoria con sync.Mutex (concurrencia)
  cargar.go               -- lee datos.json al iniciar
  reglas.go               -- crearRepresentante/Estudiante/Cita/Factura (validacion)
  api.go                  -- 8 servicios web JSON (GET y POST)
  portal.go               -- portal SSR sin JS (embed de plantilla)
  mux.go                  -- rutas + sembrar + arranque del servidor
  templates/
    portal.html           -- plantilla HTML embebida con go:embed
datos.json                 -- datos de ejemplo cargados al iniciar
```

## Reglas de codigo

1. Solo libreria estandar de Go. PROHIBIDO dependencias externas.
2. Sin tildes (acentos) en codigo, comentarios, plantillas ni mensajes. Usar n en vez de enie.
3. **Sin JavaScript en el cliente.** Toda la UI se genera en Go con `html/template`
   + `go:embed`. Los formularios usan `POST` a rutas de Go; no hay `fetch` ni `onclick`.
4. Sin emojis en ningun archivo.
5. Todo en espanol.
6. Comentarios en espanol, sin acentos ni emojis. En los comentarios largos usar
   bloque `/* ... */`; en los cortos basta `//`. El codigo de `web/` lleva
   comentarios explicativos; mantenerlos al modificar.

## Como funciona la UI (SSR) y la API

- `mux.go` registra las rutas en `http.NewServeMux()`:
  - `GET /` y `GET /portal` -> portal (tablas de las 6 entidades + formularios + JSON).
  - `GET /portal/editar/{entidad}/{id}` -> formulario de edicion precargado.
  - `POST /portal/crear/{entidad}`      -> crea la entidad y redirige con `?msg=`.
  - `POST /portal/actualizar/{entidad}/{id}` -> aplica edicion y redirige.
  - `POST /portal/borrar/{entidad}/{id}`     -> borra si no deja referencias rotas.
  - `GET/POST /api/representantes`, `/api/terapeutas`, `/api/planes`, `/api/estudiantes`,
    `/api/citas`, `/api/facturas` -> servicios JSON de lista/creacion.
  - `GET/PUT/DELETE /api/{entidad}/{id}` -> operaciones por id (una sola ruta, el
    handler `apiPorID` traduce el plural via `mapaEntidadSingular`).
- `portal.go` construye `PortalData` y renderiza `templates/portal.html`.
- `reglas.go` tiene `crearX`, `actualizarX` y `puedeBorrar` para las 6 entidades;
  las usan tanto la API como el portal, asi la regla existe en un solo lugar.
- `cargar.go` lee `datos.json` al iniciar; `sembrar()` crea Plan y Terapeuta por defecto.

## Modelo de datos (web/modelo.go)

- **Representante** { id, nombres, apellidos, cedula, telefono }
- **Terapeuta** { id, nombres, apellidos, cedula, telefono }
- **Estudiante** { id, nombres, apellidos, cedula, telefono, representante_id, plan_id }
- **Plan** { id, nombre, descripcion, duracion }
- **Cita** { id, estudiante_id, terapeuta_id, plan_id, fecha, hora, estado }
- **Factura** { id, representante_id, fecha, monto, estado }
- **Pago** { id, factura_id, monto, fecha, metodo }

La factura se emite al **representante**, no al estudiante.

## Contexto

- El usuario es principiante: explicar con calma, en espanol.
- No hacer commits salvo que el usuario lo pida.
- Datos en memoria: todo se pierde al cerrar el servidor. Sin base de datos real.
- Al terminar, matar el servidor si quedo corriendo: `fuser -k 8080/tcp`.

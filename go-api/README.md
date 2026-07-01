# API Go — Factorización QR (Fiber)

API REST en Go que recibe una matriz rectangular, calcula su **factorización QR**,
solicita estadísticas a la API de Node, persiste el resultado en PostgreSQL
(auditoría) y devuelve todo compuesto.

> **Nota sobre el enunciado:** las secciones introductorias hablan de "rotación
> de matriz", pero la *Funcionalidad requerida* pide **factorización QR** (y la
> operación adicional opera sobre "las matrices" en plural = Q y R). Se
> implementa **QR**, tratando "rotación" como texto heredado del enunciado.

## Arquitectura (capas)

```
cmd/server         → entrypoint + wiring (inyección de dependencias manual)
internal/handler   → HTTP (parse, delega, responde)
internal/service   → caso de uso: orquesta QR + stats + persistencia
internal/matrix    → dominio puro: validación + QR (gonum)
internal/client    → adapter HTTP hacia Node (interface → mockeable)
internal/repository→ persistencia PostgreSQL (pgx) + migraciones
internal/middleware→ JWT, manejo central de errores, logging
internal/dto       → contratos request/response
internal/apperr    → error tipado (status + code + message)
db/migrations      → migraciones SQL (embebidas en el binario)
```

El flujo: `cliente → Go (QR) → Node (stats) → Go compone y persiste → cliente`.
Go actúa como orquestador; el cliente sólo necesita conocer esta API.

## Endpoints

| Método | Ruta | Auth | Descripción |
|---|---|---|---|
| `POST` | `/auth/login` | ❌ | Devuelve un JWT (usuario demo) |
| `POST` | `/api/v1/qr` | ✅ | Calcula QR + estadísticas y persiste |
| `GET`  | `/api/v1/computations` | ✅ | Historial paginado (`?limit=&offset=`) |
| `GET`  | `/api/v1/computations/:id` | ✅ | Recupera un cómputo por id |
| `GET`  | `/health` | ❌ | Readiness (incluye chequeo de DB) |

Autenticación: header `Authorization: Bearer <jwt>`. Errores con envelope único:
`{ "error": { "code", "message", "details" } }`.

### Ejemplo

```bash
# 1) Login
curl -s -X POST localhost:8080/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"demo","password":"demo123"}'

# 2) QR (usar el token del paso anterior)
curl -s -X POST localhost:8080/api/v1/qr \
  -H 'Content-Type: application/json' \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"matrix":[[1,2],[3,4],[5,6]]}'
```

## Decisiones técnicas clave

- **QR reducida (thin):** para `A` de m×n (m ≥ n) → `Q` es m×n y `R` es n×n.
  Mantener `R` cuadrada hace que el chequeo "diagonal" (en Node) sea válido.
- **Se requiere m ≥ n.** `gonum` factoriza con filas ≥ columnas; una matriz
  ancha devuelve `422 VALIDATION_ERROR`. (Extensión posible: transponer + LQ.)
- **QR no es única** → los tests verifican propiedades (`A ≈ Q·R`, `QᵀQ ≈ I`,
  `R` triangular superior), no valores fijos de Q/R.
- **Persistencia síncrona y obligatoria:** si falla el guardado se responde
  `500`. En un dominio de seguros, un log de auditoría que descarta registros
  en silencio es peor que una petición fallida.
- **JWT HS256** con secreto compartido; el token se **propaga** a Node.
- **DI manual** (sin framework) y **interfaces** en `client` y `repository`
  para poder testear el servicio con mocks.

## Configuración (variables de entorno)

Ver [.env.example](.env.example). Requeridas: `JWT_SECRET`, `DATABASE_URL`.
Las migraciones se aplican al arrancar si `AUTO_MIGRATE=true` (por defecto).

## Ejecución local

```bash
cp .env.example .env          # ajustar DATABASE_URL / JWT_SECRET
go run ./cmd/server
```

Requiere un PostgreSQL accesible. El binario embebe las migraciones, así que
no hay pasos manuales de esquema.

## Docker

```bash
docker build -t qr-go-api .
```

La orquestación completa (Go + Node + PostgreSQL) vivirá en el `docker-compose`
de la raíz del repositorio.

## Tests

```bash
go test ./...
```

Cobertura:

- **Dominio** (`internal/matrix`): validación y propiedades de la QR
  (`A ≈ Q·R`, `QᵀQ ≈ I`, `R` triangular).
- **Servicio** (`internal/service`): orquestación con **mocks** del cliente Node
  y del repositorio — happy path, `422` (matriz inválida), `502` (Node caído o
  respuesta incompleta) y `500` (fallo al persistir). Usa clock/id inyectados
  para resultados deterministas.
- **Integración HTTP** (`internal/server`): la app completa vía `app.Test` de
  Fiber (sin puerto ni DB real) — login + JWT, `401` sin token, `200`/`422`/`502`
  en `/qr`, `404` en historial y `/health`.

## Stack

Go 1.26 · Fiber · gonum (QR) · pgx (PostgreSQL) · golang-migrate · golang-jwt · google/uuid · godotenv.

# API Node — Estadísticas (Express + TypeScript)

Servicio de estadísticas que consume la API de Go. Recibe un conjunto de
matrices con nombre y devuelve estadísticas por matriz y un agregado combinado.

Es **genérico**: no sabe nada de QR. La API de Go le envía `q` y `r`, pero el
servicio funcionaría con cualquier conjunto de matrices.

## Arquitectura (capas)

```
src/index.ts          → bootstrap (levanta el servidor)
src/app.ts            → construcción de Express + middlewares
src/config/env.ts     → configuración desde variables de entorno
src/routes            → definición de rutas
src/controllers       → capa HTTP (delega en el servicio)
src/services          → lógica pura de estadísticas (testeable sin HTTP)
src/middlewares       → JWT, validación (zod), manejo central de errores
src/schemas           → esquemas zod (validación + tipos inferidos)
src/errors.ts         → error tipado (status + code + message)
tests                 → unitarios (servicio) + integración (supertest)
```

## Endpoints

| Método | Ruta | Auth | Descripción |
|---|---|---|---|
| `POST` | `/internal/stats` | ✅ | Estadísticas sobre las matrices recibidas |
| `GET`  | `/health` | ❌ | Readiness |

Autenticación: `Authorization: Bearer <jwt>` (HS256, secreto compartido con Go).
Errores con el mismo envelope que la API de Go:
`{ "error": { "code", "message", "details" } }`.

### Contrato de `POST /internal/stats`

Request:
```json
{
  "matrices": [
    { "name": "q", "data": [[0.1, 0.9], [0.5, 0.3]] },
    { "name": "r", "data": [[5.9, 7.4], [0, 0.8]] }
  ]
}
```
Response:
```json
{
  "perMatrix": [
    { "name": "q", "max": 0.9, "min": 0.1, "average": 0.45, "sum": 1.8, "isDiagonal": false },
    { "name": "r", "max": 7.4, "min": 0, "average": 3.525, "sum": 14.1, "isDiagonal": false }
  ],
  "combined": { "max": 7.4, "min": 0, "average": 1.9875, "sum": 15.9, "anyDiagonal": false }
}
```

## Estadísticas calculadas

Por cada matriz y en el combinado: **máximo, mínimo, promedio, suma**. Además:

- `isDiagonal` (por matriz): sólo puede ser `true` en matrices **cuadradas**, y
  los elementos fuera de la diagonal se comparan con **tolerancia** (`1e-9`),
  porque una QR en coma flotante no produce ceros exactos.
- `anyDiagonal` (combinado): `true` si alguna matriz es diagonal.

## Configuración

Ver [.env.example](.env.example). `JWT_SECRET` es obligatorio y **debe coincidir**
con el de la API de Go.

## Ejecución

```bash
npm install
cp .env.example .env      # ajustar JWT_SECRET
npm run dev               # desarrollo (recarga en caliente)
# o
npm run build && npm start
```

## Tests

```bash
npm test
```

- **Unitarios** (`stats.service.test.ts`): lógica de estadísticas y `isDiagonal`
  (incluye tolerancia a ruido de coma flotante y matrices no cuadradas).
- **Integración** (`stats.routes.test.ts`, con supertest): 401 sin token,
  422 con matriz irregular, 400 con JSON malformado, 200 con estadísticas
  correctas y `/health`.

## Docker

```bash
docker build -t stats-node-api .
```

Build multi-stage; imagen final sólo con dependencias de producción y usuario
no-root. La orquestación (Go + Node + PostgreSQL) vive en el `docker-compose`
de la raíz.

## Stack

Node · Express · TypeScript · zod · jsonwebtoken · helmet · cors · pino · Jest · supertest.

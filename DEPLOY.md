# Despliegue en Railway

El proyecto se despliega en Railway como **un proyecto con cuatro servicios**,
todos construidos desde este monorepo con el `Dockerfile` de cada carpeta:

```
PostgreSQL (gestionado)  ──DATABASE_URL──►  go-api
frontend (público) ──HTTP──► go-api (público) ──red privada──► node-api (privado)
```

- **postgres**: base de datos gestionada de Railway.
- **node-api**: servicio de estadísticas, solo accesible en la red privada.
- **go-api**: API principal, con dominio público (la consume el frontend).
- **frontend**: SPA estática (nginx), con dominio público.

> Los contenedores ya están preparados para Railway: el frontend escucha en
> `$PORT` (plantilla nginx) y go-api liga en dual-stack IPv4+IPv6 (la red
> privada de Railway es IPv6).

> ⚠️ **El error #1 en monorepos: la `Root Directory`.** Cada servicio se crea
> desde el MISMO repo, así que DEBES fijar su **Root Directory** a la carpeta
> correcta (`go-api`, `node-api` o `frontend`). Si te equivocas, un servicio
> construye el código de otro (p. ej. el "go-api" corriendo node-api) y las
> rutas dan 404. Verifícalo con `/health` (ver sección **Verificación**).

## Requisitos

- Cuenta en [railway.app](https://railway.app).
- Este repo en GitHub (`xangel1111/qr-stats-fullstack`).

## Pasos

### 1. Crear el proyecto y la base de datos

1. **New Project → Deploy from GitHub repo** → selecciona `qr-stats-fullstack`.
2. **New → Database → Add PostgreSQL**. Railway crea la variable `DATABASE_URL`.

### 2. Servicio `node-api`

- **Add service → GitHub Repo** (el mismo repo).
- **Settings → Source → Root Directory:** `node-api` (detecta el Dockerfile).
- **Variables:**
  - `JWT_SECRET` = *(un secreto fuerte; anótalo)*
  - `PORT` = `3000`
- Sin dominio público (solo lo llama go-api por la red privada).

### 3. Servicio `go-api`

- **Add service → GitHub Repo** (mismo repo). **Root Directory:** `go-api`.
- **Variables:**
  - `JWT_SECRET` = *(EL MISMO valor que en node-api)*
  - `PORT` = `8080`
  - `DATABASE_URL` = `${{Postgres.DATABASE_URL}}`
    - Si las migraciones fallan por SSL, usa `${{Postgres.DATABASE_URL}}?sslmode=disable`.
  - `NODE_STATS_URL` = `http://${{node-api.RAILWAY_PRIVATE_DOMAIN}}:3000`
  - *(opcionales)* `AUTH_USERNAME`, `AUTH_PASSWORD`, `JWT_EXPIRY`
- **Settings → Networking → Generate Domain** → anota la URL pública
  (ej. `https://go-api-production-xxxx.up.railway.app`).
- Al arrancar aplica las migraciones automáticamente (`AUTO_MIGRATE=true`).

### 4. Servicio `frontend`

- **Add service → GitHub Repo** (mismo repo). **Root Directory:** `frontend`.
- **Variables:**
  - `VITE_API_URL` = `https://<dominio-público-de-go-api>` (paso 3)
    - Es *build-time*: Railway lo pasa como `ARG` al Dockerfile. Si cambia la
      URL, hay que **redeploy** (rebuild) del frontend.
- **Settings → Networking → Generate Domain** → dominio público del frontend.
- Abre el dominio del frontend → login `demo` / `demo123` → calcula una matriz.

## Variables por servicio (resumen)

| Servicio | Variable | Valor |
|---|---|---|
| node-api | `JWT_SECRET` | secreto compartido |
| node-api | `PORT` | `3000` |
| go-api | `JWT_SECRET` | **mismo** secreto compartido |
| go-api | `PORT` | `8080` |
| go-api | `DATABASE_URL` | `${{Postgres.DATABASE_URL}}` |
| go-api | `NODE_STATS_URL` | `http://${{node-api.RAILWAY_PRIVATE_DOMAIN}}:3000` |
| frontend | `VITE_API_URL` | URL pública de go-api (build-time) |

> Las referencias `${{Servicio.VARIABLE}}` usan el **nombre del servicio** tal
> como lo nombres en Railway. Ajústalas si usas otros nombres.

## Verificación

Comprueba que cada servicio corre el código correcto antes de usar el frontend:

```bash
# Debe responder con checks.database (es el go-api):
curl https://<dominio-go-api>/health
# → {"status":"ok","checks":{"database":"ok"}}

# El node-api responde SIN checks:
curl https://<dominio-node-api>/health
# → {"status":"ok"}
```

Si el dominio que usa el frontend (`VITE_API_URL`) devuelve `{"status":"ok"}`
**sin** `checks`, ese servicio está corriendo node-api: corrige su **Root
Directory** a `go-api` y añade sus variables (`DATABASE_URL`, `NODE_STATS_URL`).

## Notas

- **Orden:** despliega go-api antes que el frontend (necesitas su URL pública
  para `VITE_API_URL`).
- **Secreto JWT:** debe ser idéntico en go-api y node-api (go-api emite el token
  y node-api lo valida).
- **CORS:** go-api permite todos los orígenes, así que el frontend en su dominio
  de Railway puede llamarlo sin ajustes.
- **Seguridad en producción:** cambia `JWT_SECRET` y las credenciales demo
  (`AUTH_USERNAME` / `AUTH_PASSWORD`) por valores propios.

## Alternativa local

Todo el stack corre en local con un solo comando (ver README):

```bash
docker compose up --build
```

# Frontend — React + TypeScript (Vite)

SPA que consume la API de Go (que a su vez orquesta la API de Node) y muestra la
factorización QR, las estadísticas y el historial de auditoría.

## Funcionalidad

- **Conexión / JWT:** login con las credenciales demo → obtiene el token que se
  usa en las llamadas protegidas.
- **Matriz de entrada:** grilla editable con controles de filas/columnas y un
  botón "Cargar ejemplo".
- **Resultado:** matrices Q y R + tabla de estadísticas (por matriz y combinado).
- **Historial:** últimos cómputos registrados (endpoint de auditoría de Go).

## Estructura

```
src/
├── main.tsx            # bootstrap React
├── App.tsx             # estado y composición de la vista
├── api.ts              # cliente HTTP (login, computeQR, listComputations)
├── types.ts            # tipos que reflejan el contrato de la API
├── format.ts           # formateo numérico (redondeo, -0 → 0)
├── styles.css
└── components/         # MatrixInput, MatrixTable, StatsPanel, History
```

## Ejecución

```bash
npm install
cp .env.example .env      # VITE_API_URL (por defecto http://localhost:8080)
npm run dev               # http://localhost:5173
# build de producción
npm run build && npm run preview
```

## Notas de diseño

- **Sin librería de UI** (Material, etc.): CSS propio y ligero para evitar
  sobreingeniería en un frontend de una sola vista.
- **`VITE_API_URL`** se inyecta en tiempo de build; apunta a donde el navegador
  (en el host) alcanza la API de Go.
- La grilla mantiene las celdas como texto durante la edición y se parsean a
  números al enviar, lo que permite estados intermedios ("-", "1.").

## Docker

```bash
docker build -t qr-frontend .
```

Build multi-stage: compila con Node y sirve los estáticos con nginx (con
fallback SPA). Incluido como servicio en el `docker-compose` de la raíz.

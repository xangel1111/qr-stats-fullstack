# ---- build stage ----
FROM node:22-alpine AS build

WORKDIR /app

COPY package*.json ./
RUN npm ci

COPY . .

# VITE_API_URL is baked into the bundle at build time. It must point to where
# the browser (running on the host) reaches the Go API.
ARG VITE_API_URL=http://localhost:8080
ENV VITE_API_URL=$VITE_API_URL
RUN npm run build

# ---- run stage ----
FROM nginx:1.27-alpine

# Render the listen port from $PORT at startup so the same image works locally
# (defaults to 80) and on platforms that inject a dynamic port (e.g. Railway).
# NGINX_ENVSUBST_FILTER restricts substitution to PORT, leaving nginx's own
# $uri variable untouched.
ENV PORT=80
ENV NGINX_ENVSUBST_FILTER=^PORT$
COPY nginx.conf.template /etc/nginx/templates/default.conf.template
COPY --from=build /app/dist /usr/share/nginx/html

EXPOSE 80

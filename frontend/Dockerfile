# Build context must be repository root: docker compose build frontend
# (requires ../mobile-web for mobile styles)
ARG NODE_IMAGE=node:20-alpine
ARG NGINX_IMAGE=nginx:alpine
FROM ${NODE_IMAGE} AS build
WORKDIR /app
COPY frontend/package*.json ./
RUN npm install
COPY frontend/ ./
COPY mobile-web/ /mobile-web/
RUN npm run build

FROM ${NGINX_IMAGE}
COPY --from=build /app/dist /usr/share/nginx/html
COPY frontend/nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]

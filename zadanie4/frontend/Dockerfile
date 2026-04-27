# ETAP 1: Budowanie (Build Stage)
FROM node:22-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install
COPY . .

RUN npm run build

# ETAP 2: Serwowanie (Production Stage)
FROM nginx:stable-alpine

# Czyścimy domyślne pliki Nginxa
RUN rm -rf /usr/share/nginx/html/*

# Kopiujemy built z poprzedniego etapu (folder /dist) do Nginxa
COPY --from=builder /app/dist /usr/share/nginx/html

# Opcjonalnie: Kopiujemy własną konfigurację Nginxa (ważne dla React Routera!)
# COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
# Stage1: building frontend
FROM node:22-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install
COPY . .

RUN npm run build

# Stage2: nginx serve
FROM nginx:stable-alpine

RUN rm -rf /usr/share/nginx/html/*

COPY --from=builder /app/dist /usr/share/nginx/html

# opcjonalnie: kopiujemy conf nginxa -> P.S Rzeczywiscie wazne xddddd
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
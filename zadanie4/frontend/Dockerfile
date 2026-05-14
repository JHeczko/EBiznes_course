# Stage1: building frontend
FROM node:22-alpine AS builder

WORKDIR /app

COPY package*.json ./
RUN npm ci --ignore-scripts
COPY public ./public
COPY src ./src
COPY eslint.config.js ./
COPY index.html ./
COPY package.json ./
COPY package-lock.json ./
COPY tsconfig.json ./
COPY tsconfig.app.json ./
COPY tsconfig.node.json ./
COPY vite.config.ts ./
# COPY . .

RUN npm run build

# Stage2: nginx serve
FROM nginx:stable-alpine

RUN rm -rf /usr/share/nginx/html/*

COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

RUN addgroup -S appnginx && \
    adduser -S appnginx -G appnginx && \
    touch /var/run/nginx.pid && \
    chown -R appnginx:appnginx /usr/share/nginx/html && \
    chown -R appnginx:appnginx /var/cache/nginx && \
    chown -R appnginx:appnginx /var/log/nginx && \
    chown -R appnginx:appnginx /etc/nginx/conf.d && \
    chown -R appnginx:appnginx /var/run/nginx.pid

USER appnginx

EXPOSE 8080

CMD ["nginx", "-g", "daemon off;"]
# ── Build stage ───────────────────────────────────────────────────────────────
FROM node:22-alpine AS builder

ARG VITE_BFF_URL=http://localhost:9001
ARG VITE_GATEWAY_URL=http://localhost:9002

ENV VITE_BFF_URL=${VITE_BFF_URL}
ENV VITE_GATEWAY_URL=${VITE_GATEWAY_URL}

WORKDIR /src/frontend

RUN corepack enable && corepack prepare pnpm@latest --activate

COPY frontend/package.json frontend/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

COPY frontend ./
COPY contracts/generated/es ../contracts/generated/es

RUN pnpm build

# ── Runtime stage (nginx) ─────────────────────────────────────────────────────
FROM nginx:1.27-alpine AS runtime

COPY --from=builder /src/frontend/dist /usr/share/nginx/html

# SPA fallback: all unknown paths serve index.html
COPY <<'EOF' /etc/nginx/conf.d/default.conf
server {
    listen 80;
    root /usr/share/nginx/html;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /health {
        return 200 "ok\n";
        add_header Content-Type text/plain;
    }
}
EOF

EXPOSE 80

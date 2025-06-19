# syntax=docker/dockerfile:1

# --- Build stage ---
FROM node:20-alpine AS build
WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN npm install -g pnpm && pnpm install --frozen-lockfile
COPY . .
RUN npm run build

# --- Runtime stage ---
FROM node:20-alpine
WORKDIR /app
COPY --from=build /app/dist ./dist
COPY package.json .
RUN npm install -g pnpm && pnpm install --prod --frozen-lockfile
ENV PORT=8080
EXPOSE 8080
CMD ["node", "dist/index.js"]

# Build Stage
FROM oven/bun:1.1 AS builder
WORKDIR /app

# Copy root configurations
COPY package.json bun.lock turbo.json ./

# Copy all packages and apps
COPY apps ./apps
COPY packages ./packages

# Install all dependencies
RUN bun install

# Build the web app via turborepo
RUN bun run build --filter=web

# Runner Stage
FROM node:20-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

# Add non-root user for security
RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

# Copy static files
COPY --from=builder /app/apps/web/public ./apps/web/public

# Copy the standalone output (this contains its own node_modules and required files)
COPY --from=builder --chown=nextjs:nodejs /app/apps/web/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/apps/web/.next/static ./apps/web/.next/static

USER nextjs

EXPOSE 3000

# The standalone build outputs the server file inside the workspace folder
CMD ["node", "apps/web/server.js"]

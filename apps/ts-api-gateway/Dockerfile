# Build stage
FROM oven/bun:1 AS build
WORKDIR /usr/src/app

# Copy package files
COPY package.json bun.lock ./
COPY apps/ts-mobile-app/package.json ./apps/ts-mobile-app/
COPY apps/ts-api-gateway/package.json ./apps/ts-api-gateway/
COPY packages/ts-protos/package.json ./packages/ts-protos/

# Install only production dependencies
# RUN rm bun.lock && bun install
RUN bun install --frozen-lockfile

# Release stage
FROM oven/bun:1 AS release
WORKDIR /usr/src/app

# Copy production dependencies and workspace structure from the build stage
COPY --from=build /usr/src/app ./

# Copy actual source code for the app and its internal workspace dependency
COPY apps/ts-api-gateway ./apps/ts-api-gateway/
COPY packages/ts-protos ./packages/ts-protos/

# Set the final working directory
WORKDIR /usr/src/app/apps/ts-api-gateway

# Run the app
USER bun
EXPOSE 3000/tcp
ENTRYPOINT [ "bun", "run", "src/server.ts" ]

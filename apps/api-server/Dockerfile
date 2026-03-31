# Build stage
FROM oven/bun:1 AS build
WORKDIR /usr/src/app

# Copy package files
COPY package.json bun.lock ./
COPY apps/api-server/package.json ./apps/api-server/
COPY pkgs/protos/package.json ./pkgs/protos/

# Install only production dependencies
RUN rm bun.lock && bun install

# Release stage
FROM oven/bun:1 AS release
WORKDIR /usr/src/app

# Copy production dependencies and workspace structure from the build stage
COPY --from=build /usr/src/app ./

# Copy actual source code for the app and its internal workspace dependency
COPY apps/api-server ./apps/api-server/
COPY pkgs/protos ./pkgs/protos/

# Set the final working directory
WORKDIR /usr/src/app/apps/api-server

# Run the app
USER bun
EXPOSE 3000/tcp
ENTRYPOINT [ "bun", "run", "src/server.ts" ]
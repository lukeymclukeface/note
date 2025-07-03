# Stage 1: Build stage
FROM golang:1.24-alpine AS build
# Install build dependencies for CGO
RUN apk add --no-cache gcc musl-dev sqlite-dev
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o /note-server ./cmd/server

# Stage 2: Runtime stage
FROM alpine:3.20
RUN apk add --no-cache ffmpeg
# unprivileged user
RUN adduser -D -g '' app
USER app
COPY --from=build /note-server /note-server
EXPOSE 8080
ENV PORT=8080
# OpenAI key is now optional and configured through the web interface
HEALTHCHECK --interval=30s --timeout=2s CMD wget -qO- http://localhost:8080/healthz || exit 1
ENTRYPOINT ["/note-server"]

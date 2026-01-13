# Build stage
FROM golang:1.21-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o kanban-server main.go

# Final image
FROM alpine:latest
WORKDIR /app
COPY --from=build /app/kanban-server .
COPY templates ./templates
# Optionally copy a default tasks.json if you want to pre-seed data
# COPY tasks.json ./tasks.json
EXPOSE 8080
CMD ["./kanban-server"]

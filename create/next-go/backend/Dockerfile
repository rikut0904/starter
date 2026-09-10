FROM golang:1.26-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o server ./cmd/server
FROM alpine:3.22
COPY --from=build /app/server /server
EXPOSE 8080
CMD ["/server"]

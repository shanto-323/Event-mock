FROM golang:alpine AS build
WORKDIR /database-mock
COPY database-mock/go.mod database-mock/go.sum ./
COPY database-mock ./
RUN go build -o app .

FROM alpine:3.21
WORKDIR /usr/bin
COPY --from=build /database-mock/app .
EXPOSE 8080
CMD ["./app"]
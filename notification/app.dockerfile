FROM golang:alpine AS build
WORKDIR /notification
COPY notification/go.mod notification/go.sum ./
COPY notification ./
RUN go build -o app .

FROM alpine:3.21
WORKDIR /usr/bin
COPY --from=build /notification/app .
EXPOSE 8080
CMD ["./app"]
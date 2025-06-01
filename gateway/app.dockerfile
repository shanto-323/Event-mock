FROM golang:alpine AS build
WORKDIR /gateway
COPY gateway/go.mod gateway/go.sum ./
COPY gateway ./
RUN go build -o app .

FROM alpine:3.21
WORKDIR /usr/bin
COPY --from=build /gateway/app .
EXPOSE 8080
CMD ["./app"]
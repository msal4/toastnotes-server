FROM golang:1.15.4-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -installsuffix 'static' -o app .

FROM scratch AS final

WORKDIR /app

COPY --from=builder /src/app /src/.env /src/wait-for-it.sh /src/start.sh ./
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 80
EXPOSE 443

VOLUME ["/cert-cache"]

ENTRYPOINT [ "./start.sh" ]

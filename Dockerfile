FROM golang:1.15.4

WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -installsuffix 'static' -o app .

EXPOSE 8080

CMD [ "make", "prod" ]

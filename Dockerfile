FROM golang:1.15.4

RUN apt install ca-certificates git -y

WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -installsuffix 'static' -o app .

EXPOSE $PORT

VOLUME ["/cert-cache"]

ENTRYPOINT ["./start.sh"]


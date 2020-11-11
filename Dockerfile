FROM golang:1.15.4

WORKDIR /app
COPY ./go.mod .
COPY ./go.sum .
RUN go mod download
COPY . .
RUN go build -o toastnotes
ENTRYPOINT [ "./start.sh" ]
EXPOSE 8080
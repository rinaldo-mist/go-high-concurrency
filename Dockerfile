FROM golang:1.25
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go get
RUN go build -o bin .
EXPOSE 8080
ENTRYPOINT ["/app/bin"]

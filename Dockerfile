FROM golang:latest as Builder
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download
COPY . .
RUN go build -o main .
FROM alpine:latest
WORKDIR /root/
COPY --from=Builder /app/main .
EXPOSE 8080
CMD [ "./main" ]
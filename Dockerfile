FROM golang:1.22.2-alpine

RUN mkdir /service-sales
WORKDIR /service-sales

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o app ./cmd

FROM alpine
WORKDIR /service-sales

COPY --from=builder /service-sales/app /service-sales/app

RUN chmod +x app
EXPOSE 8005
CMD ["/service-sales/app"]
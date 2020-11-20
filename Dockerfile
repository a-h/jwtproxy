FROM golang:1.13 as builder

WORKDIR /app

COPY go.* ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-w -s' -a -installsuffix cgo -o /jwtproxy jwtproxy

FROM scratch
COPY --from=builder /jwtproxy /

CMD ["/jwtproxy"]
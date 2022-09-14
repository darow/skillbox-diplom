# syntax = docker/dockerfile:1

FROM golang:alpine AS builder

RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o statuspage ./cmd/statuspage
RUN go build -o simulator ./third_party/simulator

FROM alpine
WORKDIR /app
COPY --from=builder /app/web /app/web
COPY --from=builder /app/statuspage /app/statuspage
COPY --from=builder /app/simulator /app/simulator
COPY --from=builder /app/start.sh /app/start.sh
CMD ["sh", "./start.sh"]
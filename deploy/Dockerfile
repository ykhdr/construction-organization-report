FROM golang:1.22.1-alpine as builder
LABEL authors="ykhdr"

WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go download

COPY ./cmd ./cmd
COPY ./internal ./internal
RUN go build -o construction-ogranization-report ./cmd/main.go

FROM alpine:3.19.1
COPY --from=builder construction-ogranization-report /bin/construction-ogranization-report
ENTRYPOINT ["/bin/construction-ogranization-report"]
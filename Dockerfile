FROM golang:alpine AS builder
RUN apk update \
    && apk add --no-cache \
    g++ \
    git \
    musl-dev \
    go \
    tesseract-ocr-dev \
    tesseract-ocr-data-eng \
    leptonica-dev

WORKDIR /app
COPY . .
RUN go mod download \
    && go build -o /ocr

FROM alpine

RUN apk update \
    && apk add --no-cache \
    tesseract-ocr-dev \
    tesseract-ocr-data-eng \
    leptonica-dev 

COPY --from=builder /ocr /ocr

CMD ["/ocr"]
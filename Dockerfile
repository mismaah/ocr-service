FROM golang:alpine

RUN apk update
RUN apk add \
    g++ \
    git \
    musl-dev \
    go \
    tesseract-ocr-dev

RUN apk add tesseract-ocr-data-eng leptonica-dev

WORKDIR /app

COPY . .
RUN go mod download

RUN go build -o /ocr

CMD ["/ocr"]
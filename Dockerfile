FROM alpine

RUN apk update
RUN apk add \
    g++ \
    git \
    musl-dev \
    go \
    tesseract-ocr-dev

WORKDIR /app

COPY . .
RUN go mod download
RUN go build -o /ocr

CMD ["/ocr"]
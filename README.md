# OCR Service

Uses [gosseract](https://github.com/otiai10/gosseract) (wrapper for [tesseract-ocr](https://github.com/tesseract-ocr/tesseract)) to provide a Service which returns text from an image. Input options are image URL and image upload using a multipart-form.

## Usage

### Pull and start the service

```bash
docker run --name ocr-service --restart unless-stopped -p 8080:8080 -e PORT=8080 -e AUTH_KEY=secret123 ghcr.io/mismaah/ocr-service:main
```

`AUTH_KEY` is an optional environment variable. If this is set, all API calls should have the Authorization header.

### Endpoints

#### Text from URL

`curl --location --request POST 'http://localhost:8123/text/url?url=https://image.png' --header 'Authorization: ••••••'`

#### Text from file upload

`curl --location 'http://localhost:8123/text/upload' --header 'Authorization: ••••••' --form 'file=@"file.jpg"'`

### Response

If successful, returns the text from the image:

```json
{
  "text": "Example text from the image."
}
```

If error, returns the error message:

```json
{
  "error": "Error reading from image."
}
```

# OCR Service

Uses [gosseract](https://github.com/otiai10/gosseract) (wrapper for [tesseract-ocr](https://github.com/tesseract-ocr/tesseract)) to provide a Service which returns text from an image when the url for the image is provided.

## Quick Start

### Clone repository and navigate to directory

` git clone https://github.com/mismaah/ocr-service.git`

` cd ocr-service`

### Build and tag docker image

` docker build . -t ocr-service`

### Run docker container

Environment variables are set during container runtime by passing the .env file to the container. The .env.example file can be renamed to .env and filled prior to this.

AUTH_KEY is an optional variable. Setting this variable prevents access to the service if the key is not provided in the GET request. If this variable is not set, the service is available without authorization.

Port 8080 from the container is mapped to the system's port 8080. This can be any port but make sure that it is the same as the PORT environment variable.

`docker run -p 8080:8080 --env-file .env ocr-service`

An alternative is to set each environment variable individually, for example:

`docker run -p 8080:8080 -e PORT='8080' -e AUTH_KEY='secretkey123' ocr-service`

### Send GET Request

#### Using CURL

`curl -H "Authorization:secretkey123" http://localhost:8080/imagetotext?url=https://image.png`

#### Using Go

```go
url := "http://localhost:8080/imagetotext?url=https://image.png"
request, _ := http.NewRequest("GET", url, nil)
request.Header.Set("Authorization", "secretkey123")
client := &http.Client{}
response, _ := client.Do(request)
```

#### Using JavaScript (axios)

```javascript
axios({
  url: "http://localhost:8080/imagetotext?url=https://image.png",
  method: "GET",
  headers: { Authorization: "secretkey123" },
});
```

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

docker stop ocr-service
docker rm ocr-service
docker build . -t ocr-service
source .env
docker run -p ${PORT}:${PORT} -e PORT=${PORT} -e AUTH_KEY=${AUTH_KEY} --name ocr-service ocr-service
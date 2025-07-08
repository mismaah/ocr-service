docker build . -t ocr-service
docker stop ocr-service
docker rm ocr-service
source .env
docker run -d -p ${PORT}:${PORT} -e PORT=${PORT} -e AUTH_KEY=${AUTH_KEY} --name ocr-service ocr-service
docker logs -f --tail 100 ocr-service
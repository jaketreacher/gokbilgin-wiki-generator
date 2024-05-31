FROM golang:latest
USER root
RUN apt-get update && apt-get install -y tesseract-ocr libtesseract-dev
RUN apt-get install -y tesseract-ocr-tur
WORKDIR /app
COPY app/go.mod app/go.sum ./
RUN go mod download
COPY app/ .

CMD ["go", "run", "."]

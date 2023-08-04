# Start from the official Go image
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the Go web application source code into the container
COPY . .

# Install tesseract-ocr dependencies
RUN apt-get update && \
    apt-get install -y tesseract-ocr libleptonica-dev libtesseract-dev

# Install the tesseract-mar language package
RUN apt-get install -y tesseract-ocr-mar tesseract-ocr-eng tesseract-ocr-san tesseract-ocr-hin

# Build the Go web application
RUN go build -o webapp

# Expose the port on which the Go web application runs
EXPOSE 8080

# Set the command to run the Go web application
CMD ["./webapp"]

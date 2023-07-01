# OCR My Image

OCR My Image is a web frontend for tesseract OCR

![](/screenshot.png)

It has dependency on Tesseract ocr and it must be installed first before installing OCR My Image server.

On fedora, install below dependencies.

```bash
sudo dnf install tesseract tesseract-devl
```

On Ubuntu or debian
```bash
sudo apt install tesseract-ocr libtesseract-dev -y
```

You must install tesseract language pack for the language you want to perform ocr. For example, if you want to perform OCR on Marathi text, use

```bash
sudo apt install tesseract-ocr-mar
```

You can view installed languages using `tesseract --list-langs`

## Install OCR My Image

Use `go install github.com/quaintdev/ocrmyimage@latest` to install the web server on your platform.

If you are on Linux (amd64) you can use pre-built binary from release pages and execute it.

#!/usr/bin/env bash

set -e

OWNER="Nathan-rs"
REPO="pdfgen"

ARCH=$(uname -m)

case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Arquitetura não suportada."
        exit 1
        ;;
esac

VERSION=$(curl -s \
https://api.github.com/repos/$OWNER/$REPO/releases/latest \
| grep tag_name \
| cut -d '"' -f4)

URL="https://github.com/$OWNER/$REPO/releases/download/$VERSION/pdfgen-linux-$ARCH"

echo "Baixando $URL..."

curl -L "$URL" -o pdfgen

chmod +x pdfgen

sudo mv pdfgen /usr/local/bin/pdfgen

echo
echo "Instalação concluída!"
echo

pdfgen --help
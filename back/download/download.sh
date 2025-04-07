#!/bin/bash

# Проверяем наличие аргументов
if [[ $# -eq 0 ]]; then
  echo "Не указаны аргументы. Используйте: $0 \"аргументы\""
  exit 1
fi

DOWNLOAD_DIR="${DIR_DOWNLOAD:-$(pwd)}"

# Выполняем команду docker run с переданными аргументами
docker run --rm \
  -v "{DOWNLOAD_DIR}:/app" \
  -w /app \
  thr3a/yt-dlp "$@"
#!/bin/bash

# Проверяем наличие аргументов
if [[ $# -eq 0 ]]; then
  echo "Не указаны аргументы. Используйте: $0 \"аргументы\""
  exit 1
fi

# Выполняем команду docker run с переданными аргументами
docker run --rm \
  -v "/wv/back/download:/app" \
  -w /app \
  thr3a/yt-dlp "$@"
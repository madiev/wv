#!/bin/bash

if [[ $# -lt 2 ]]; then
  echo "Не указаны необходимые аргументы. Используйте: $0 \"путь_до_папки_download\" \"аргументы_docker_run\""
  exit 1
fi

mount_path=$1
shift

docker run --rm \
  -v "$mount_path:/app" \
  -w /app \
  thr3a/yt-dlp "$@"
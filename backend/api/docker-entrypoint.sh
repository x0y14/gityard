#!/bin/sh

set -e

NUM_OF_WORKERS=${WORKERS:-1}
PORT=${PORT:-8000}

exec uv run gunicorn app.main:app --workers "$NUM_OF_WORKERS" --worker-class uvicorn.workers.UvicornWorker --bind 0.0.0.0:"$PORT"
FROM python:3.13-slim-bookworm

COPY --from=ghcr.io/astral-sh/uv:latest /uv /uvx /bin/

COPY . /project

WORKDIR /project

COPY docker-entrypoint.sh /
RUN chmod +x /docker-entrypoint.sh

RUN uv sync --locked

EXPOSE 8000

ENTRYPOINT [ "/docker-entrypoint.sh" ]
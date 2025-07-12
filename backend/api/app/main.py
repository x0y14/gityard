from fastapi import FastAPI

from app.routers import healthcheck

app = FastAPI()
app.include_router(healthcheck.router)

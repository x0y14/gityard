from contextlib import asynccontextmanager

from fastapi import FastAPI
from sqlmodel import SQLModel

from app.core.db import engine
from app.routes import healthcheck, login, users


def create_db_and_tables():
    SQLModel.metadata.create_all(engine)


@asynccontextmanager
async def lifespan(app: FastAPI):
    # before start
    create_db_and_tables()

    yield
    # after stop

    # do nothing


app = FastAPI(lifespan=lifespan)
app.include_router(healthcheck.router)
app.include_router(login.router)
app.include_router(users.router)

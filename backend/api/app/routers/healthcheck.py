from typing import Annotated

from fastapi import APIRouter, Query
from pydantic import BaseModel, ConfigDict

router = APIRouter()


class HealthCheckRequest(BaseModel):
    """/healthcheckに対するリクエスト"""

    model_config = ConfigDict(frozen=True, extra="forbid")


class HealthCheckResponse(BaseModel):
    """/healthcheckからのレスポンス"""

    model_config = ConfigDict(frozen=True, extra="forbid")


@router.get("/healthcheck", response_model=HealthCheckResponse, status_code=200)
async def healthcheck(query: Annotated[HealthCheckRequest, Query()]):
    _ = query
    return HealthCheckRequest()

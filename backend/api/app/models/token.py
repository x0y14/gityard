# JSON payload containing access token
import uuid
from datetime import timedelta

from sqlmodel import Field, SQLModel


class ShortTermToken(SQLModel):
    access_token: str
    token_type: str = "bearer"


# Contents of JWT token
class ShortTermTokenPayload(SQLModel):
    sub: uuid.UUID | None = None
    term: str = "short"


class LongTermToken(SQLModel, table=True):
    refresh_token: str
    token_type: str = "bearer"
    user_id: uuid.UUID = Field(primary_key=True, foreign_key="user.id")


class LongTermTokenPayload(SQLModel):
    sub: uuid.UUID | None = None
    term: str = "long"


class LongTermTokenCreate(SQLModel):
    user_id: uuid.UUID


class LongTermTokenCreated(LongTermToken):
    expires: timedelta

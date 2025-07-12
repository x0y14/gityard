from datetime import datetime, timedelta, timezone
from typing import Any

import bcrypt
import jwt

from app.core.config import settings

ALGORITHM = "HS256"


def create_access_token(subject: str | Any, expires_delta: timedelta) -> str:
    expire = datetime.now(timezone.utc) + expires_delta
    encoded_jwt = jwt.encode(
        {"exp": expire, "sub": str(subject)}, settings.SECRET_KEY, algorithm=ALGORITHM
    )
    return encoded_jwt


def verify_password(plain_password: str, hashed_password: str) -> bool:
    return bcrypt.checkpw(plain_password.encode(), hashed_password.encode())


def get_password_hash(password: str) -> str:
    return bcrypt.hashpw(password.encode(), bcrypt.gensalt(rounds=10)).decode()

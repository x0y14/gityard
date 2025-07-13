from collections.abc import Generator
from typing import Annotated

import jwt
from fastapi import Cookie, Depends, HTTPException, status
from fastapi.security import OAuth2PasswordBearer
from jwt.exceptions import InvalidTokenError
from pydantic import ValidationError
from sqlmodel import Session

from app import crud
from app.core import security
from app.core.config import settings
from app.core.db import engine
from app.models.token import LongTermTokenPayload, ShortTermTokenPayload
from app.models.user import User

reusable_oauth2 = OAuth2PasswordBearer(
    tokenUrl=f"{settings.API_V1_STR}/login/access-token"
)


def get_db() -> Generator[Session, None, None]:
    with Session(engine) as session:
        yield session


SessionDep = Annotated[Session, Depends(get_db)]
AccessTokenDep = Annotated[str, Depends(reusable_oauth2)]


def get_current_user_from_auth_header(
    session: SessionDep,
    access_token: AccessTokenDep,
) -> User:
    try:
        # access_tokenを用いて検証
        payload = jwt.decode(
            access_token, settings.SECRET_KEY, algorithms=[security.ALGORITHM]
        )
        token_data = ShortTermTokenPayload(**payload)
        if token_data.term != "short":
            raise InvalidTokenError("Non short term token provided")

    except (InvalidTokenError, ValidationError):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Could not validate credentials",
        )
    user = session.get(User, token_data.sub)
    if not user:
        raise HTTPException(status_code=404, detail="User not found")
    return user


CurrentUserFromAuthHeader = Annotated[User, Depends(get_current_user_from_auth_header)]


RefreshTokenDep = Annotated[str, Cookie(alias="refresh_token")]


def get_current_user_from_cookie(
    session: SessionDep, refresh_token: RefreshTokenDep
) -> User:
    try:
        # refresh_tokenを用いて検証
        payload = jwt.decode(
            refresh_token, settings.SECRET_KEY, algorithms=[security.ALGORITHM]
        )
        token_data = LongTermTokenPayload(**payload)
        if token_data.term != "long":
            raise InvalidTokenError("Non long term token provided")
        if token_data.sub is None:
            raise InvalidTokenError("Invalid refresh token provided")

        # 使い回しでないことを確認、最新のものと一致するか?
        # (DBに登録されたものしか使えず、登録されていないものは全て失効扱い)
        active_long_term_token = crud.get_long_term_token(
            session=session, user_id=token_data.sub
        )
        if (active_long_term_token is None) or (
            active_long_term_token.refresh_token != refresh_token
        ):
            raise HTTPException(
                status_code=status.HTTP_401_UNAUTHORIZED,
                detail="Invalid refresh token provided",
            )

    except (InvalidTokenError, ValidationError):
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Could not validate credentials",
        )
    user = session.get(User, token_data.sub)
    if not user:
        raise HTTPException(status_code=404, detail="User not found")
    return user


CurrentUserFromCookie = Annotated[User, Depends(get_current_user_from_cookie)]

from datetime import timedelta
from typing import Annotated

from fastapi import APIRouter, Depends, HTTPException, Response
from fastapi.security import OAuth2PasswordRequestForm

from app import crud
from app.core import security
from app.core.config import settings
from app.deps import CurrentUserFromCookie, SessionDep
from app.models.token import LongTermTokenCreate, ShortTermToken

router = APIRouter(tags=["login"])


@router.post("/login/access-token")
def login_access_token(
    session: SessionDep,
    form_data: Annotated[OAuth2PasswordRequestForm, Depends()],
    response: Response,
) -> ShortTermToken:
    """
    OAuth2 compatible token login, get an access token for future requests
    """
    user = crud.authenticate(
        session=session, email=form_data.username, password=form_data.password
    )
    if not user:
        raise HTTPException(status_code=400, detail="Incorrect email or password")

    # リフレッシュトークンはcookieにつける
    long_term_token_created = crud.update_long_term_token(
        session=session, long_term_token_create=LongTermTokenCreate(user_id=user.id)
    )
    response.set_cookie(
        key="refresh_token",
        value=long_term_token_created.refresh_token,
        max_age=int(long_term_token_created.expires.total_seconds()),
        secure=False,  # http
        samesite="strict",
        httponly=True,
    )

    # アクセストークンはjson bodyで返す。
    access_token_expires = timedelta(minutes=settings.ACCESS_TOKEN_EXPIRE_MINUTES)
    return ShortTermToken(
        access_token=security.create_access_token(
            user.id, expires_delta=access_token_expires
        )
    )


@router.post("/login/refresh")
def refresh_access_token(
    session: SessionDep,
    current_user: CurrentUserFromCookie,
    response: Response,
) -> ShortTermToken:
    # リフレッシュトークンはcookieにつける
    long_term_token_created = crud.update_long_term_token(
        session=session,
        long_term_token_create=LongTermTokenCreate(user_id=current_user.id),
    )
    response.set_cookie(
        key="refresh_token",
        value=long_term_token_created.refresh_token,
        max_age=int(long_term_token_created.expires.total_seconds()),
        secure=False,  # http
        samesite="strict",
        httponly=True,
    )

    # アクセストークンはjson bodyで返す。
    access_token_expires = timedelta(minutes=settings.ACCESS_TOKEN_EXPIRE_MINUTES)
    return ShortTermToken(
        access_token=security.create_access_token(
            current_user.id, expires_delta=access_token_expires
        )
    )

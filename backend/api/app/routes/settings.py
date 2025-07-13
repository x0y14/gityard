from typing import Any

from fastapi import APIRouter, HTTPException, status
from sqlmodel import func, select

from app import crud
from app.deps import CurrentUserFromAuthHeader, SessionDep
from app.models.pubkey import (
    PubKey,
    PubkeyPublic,
    PubkeyRegister,
    PubkeyRegistered,
    PubkeysPublic,
)

router = APIRouter(prefix="/settings", tags=["settings"])


@router.post("/keys", response_model=PubkeyRegistered)
def register_pubkey(
    session: SessionDep,
    current_user: CurrentUserFromAuthHeader,
    pubkey_in: PubkeyRegister,
) -> Any:
    try:
        registered = crud.register_pubkey(
            session=session,
            user_id=current_user.id,
            pubkey_register=PubkeyRegister(
                name=pubkey_in.name, full_text=pubkey_in.full_text
            ),
        )
    except ValueError:
        raise HTTPException(
            status_code=status.HTTP_422_UNPROCESSABLE_ENTITY,
            detail="Invalid key provided",
        )

    return registered


@router.get("/keys", response_model=PubkeysPublic)
def registered_pubkeys(
    session: SessionDep,
    current_user: CurrentUserFromAuthHeader,
    skip: int = 0,
    limit: int = 100,
) -> Any:
    count_statement = (
        select(func.count())
        .select_from(PubKey)
        .where(PubKey.user_id == current_user.id)
    )
    count = session.exec(count_statement).one()

    statement = (
        select(PubKey)
        .where(PubKey.user_id == current_user.id)
        .offset(skip)
        .limit(limit)
    )
    pubkeys = session.exec(statement).all()

    data: list[PubkeyPublic] = [PubkeyPublic.model_validate(pk) for pk in pubkeys]
    return PubkeysPublic(data=data, count=count)

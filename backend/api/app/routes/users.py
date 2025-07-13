from typing import Any

from fastapi import APIRouter, HTTPException

from app import crud
from app.deps import CurrentUserFromAuthHeader, SessionDep
from app.models.user import UserCreate, UserPublic, UserRegister

router = APIRouter(prefix="/users", tags=["users"])


@router.post("/signup", response_model=UserPublic)
def register_user(session: SessionDep, user_in: UserRegister) -> Any:
    """
    Create new user without the need to be logged in.
    """
    user = crud.get_user_by_email(session=session, email=user_in.email)
    if user:
        raise HTTPException(
            status_code=400,
            detail="The user with this email already exists in the system",
        )
    user_create = UserCreate.model_validate(user_in)
    user = crud.create_user(session=session, user_create=user_create)
    return user


@router.get("/me", response_model=UserPublic)
def read_user_me(current_user: CurrentUserFromAuthHeader) -> Any:
    """
    Get current user.
    """
    return current_user

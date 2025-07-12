import uuid

from pydantic import EmailStr
from sqlmodel import Field, SQLModel


class UserBase(SQLModel):
    name: str = Field(default="unknown", max_length=255)
    email: EmailStr = Field(unique=True, index=True, max_length=255)


class User(UserBase, table=True):
    id: uuid.UUID = Field(default_factory=uuid.uuid4, primary_key=True)
    hashed_password: str


class UserCreate(UserBase):
    password: str = Field(min_length=8, max_length=40)


class UserRegister(SQLModel):
    name: str = Field(default="unknown", max_length=255)
    email: EmailStr
    password: str = Field(min_length=8, max_length=40)


class UserPublic(UserBase):
    name: str = Field(default="unknown", max_length=255)
    id: uuid.UUID

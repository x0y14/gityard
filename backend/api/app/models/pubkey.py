import uuid

from sqlmodel import Field, SQLModel


class PubKey(SQLModel, table=True):
    """OpenSSH Format"""

    id: int | None = Field(default=None, primary_key=True)
    name: str = Field(min_length=5, max_length=255)

    full_text: str = Field(max_length=2000)
    fingerprint_sha256: str = Field(max_length=64, unique=True, index=True)

    algorithm: str
    keybody: str = Field(max_length=1500)
    comment: str

    user_id: uuid.UUID = Field(foreign_key="user.id")


class PubkeyRegister(SQLModel):
    name: str
    full_text: str


class PubkeyRegistered(SQLModel):
    fingerprint_sha256: str


class PubkeyPublic(SQLModel):
    name: str
    fingerprint_sha256: str


class PubkeysPublic(SQLModel):
    data: list[PubkeyPublic]
    count: int

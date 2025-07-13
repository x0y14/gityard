import base64
import hashlib
import uuid
from datetime import timedelta

from cryptography.hazmat.primitives import serialization
from sqlmodel import Session, select

from app.core import security
from app.core.config import settings
from app.core.security import get_password_hash, verify_password
from app.models.pubkey import PubKey, PubkeyRegister, PubkeyRegistered
from app.models.token import (
    LongTermToken,
    LongTermTokenCreate,
    LongTermTokenCreated,
    LongTermTokenDelete,
)
from app.models.user import User, UserCreate


def create_user(*, session: Session, user_create: UserCreate) -> User:
    db_obj = User.model_validate(
        user_create, update={"hashed_password": get_password_hash(user_create.password)}
    )
    session.add(db_obj)
    session.commit()
    session.refresh(db_obj)
    return db_obj


def get_user_by_id(*, session: Session, id: uuid.UUID) -> User | None:
    statement = select(User).where(User.id == id)
    session_user = session.exec(statement).first()
    return session_user


def get_user_by_email(*, session: Session, email: str) -> User | None:
    statement = select(User).where(User.email == email)
    session_user = session.exec(statement).first()
    return session_user


def authenticate(*, session: Session, email: str, password: str) -> User | None:
    db_user = get_user_by_email(session=session, email=email)
    if not db_user:
        return None
    if not verify_password(password, db_user.hashed_password):
        return None
    return db_user


def get_long_term_token(
    *, session: Session, user_id: uuid.UUID
) -> LongTermToken | None:
    statement = select(LongTermToken).where(LongTermToken.user_id == user_id)
    session_long_term_token = session.exec(statement).first()
    return session_long_term_token


def update_long_term_token(
    *, session: Session, long_term_token_create: LongTermTokenCreate
) -> LongTermTokenCreated:
    # トークンそのものを作る
    refresh_token_expires = timedelta(minutes=settings.REFRESH_TOKEN_EXPIRE_MINUTES)
    refresh_token = security.create_refresh_token(
        subject=long_term_token_create.user_id,
        expires_delta=refresh_token_expires,
    )

    # 存在確認
    session_long_term_token = get_long_term_token(
        session=session, user_id=long_term_token_create.user_id
    )

    # あったら中身を更新
    if session_long_term_token:
        session_long_term_token.refresh_token = refresh_token
    # なかったら作成
    else:
        session_long_term_token = LongTermToken.model_validate(
            long_term_token_create,
            update={"refresh_token": refresh_token},
        )

    session.add(session_long_term_token)
    session.commit()
    session.refresh(session_long_term_token)
    return LongTermTokenCreated.model_validate(
        session_long_term_token, update={"expires": refresh_token_expires}
    )


def delete_long_term_token(
    *, session: Session, long_term_token_delete: LongTermTokenDelete
) -> bool:
    session_long_term_token = get_long_term_token(
        session=session, user_id=long_term_token_delete.user_id
    )
    if session_long_term_token:
        session.delete(session_long_term_token)
        session.commit()
        return True
    return False


def get_pubkey(
    *, session: Session, user_id: uuid.UUID, fingerprint: str
) -> PubKey | None:
    statement = (
        select(PubKey)
        .where(PubKey.fingerprint == fingerprint)  # and
        .where(PubKey.user_id == user_id)
    )
    session_pubkey = session.exec(statement).first()
    return session_pubkey


def register_pubkey(
    *, session: Session, user_id: uuid.UUID, pubkey_register: PubkeyRegister
) -> PubkeyRegistered:
    try:
        if (  # is pem format?
            "BEGIN" in pubkey_register.full_text
            or "END" in pubkey_register.full_text
            or "---" in pubkey_register.full_text
        ):
            pem_pubkey = serialization.load_pem_public_key(
                pubkey_register.full_text.encode()
            )
            openssh_pubkey_bytes = pem_pubkey.public_bytes(
                encoding=serialization.Encoding.OpenSSH,
                format=serialization.PublicFormat.OpenSSH,
            )
            openssh_public_key_str = openssh_pubkey_bytes.decode("utf-8")
        else:  # openssh format?
            # is valid?
            _ = serialization.load_ssh_public_key(pubkey_register.full_text.encode())
            openssh_public_key_str = pubkey_register.full_text
    except ValueError:
        raise

    item = openssh_public_key_str.split()
    algorithm = item[0]
    keybody = item[1]
    comment = item[2]

    # calc fingerprint
    key_bytes = base64.b64decode(keybody)
    sha256_hash = hashlib.sha256(key_bytes).digest()
    fingerprint_b64 = base64.b64encode(sha256_hash).decode("utf-8").rstrip("=")
    fingerprint = f"SHA256:{fingerprint_b64}"

    # check duplicates
    session_pubkey = get_pubkey(
        session=session, user_id=user_id, fingerprint=fingerprint
    )
    if session_pubkey:
        raise ValueError("registered key")

    db_obj = PubKey(
        name=pubkey_register.name,
        full_text=pubkey_register.full_text,
        fingerprint=fingerprint,
        algorithm=algorithm,
        keybody=keybody,
        comment=comment,
        user_id=user_id,
    )

    session.add(db_obj)
    session.commit()
    session.refresh(db_obj)

    _ = openssh_public_key_str
    return PubkeyRegistered.model_validate(db_obj)


def delete_pubkey(*, session: Session, user_id: uuid.UUID, fingerprint: str) -> bool:
    session_pubkey = get_pubkey(
        session=session, user_id=user_id, fingerprint=fingerprint
    )
    if session_pubkey is None:
        return False

    session.delete(session_pubkey)
    session.commit()
    return True

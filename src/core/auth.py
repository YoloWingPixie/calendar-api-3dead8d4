"""API Key authentication dependency."""

import logging
from typing import Annotated

from fastapi import Depends, Header, HTTPException, status
from sqlalchemy.orm import Session

from src.core.config import settings
from src.core.database import get_db
from src.models.user import User

logger = logging.getLogger(__name__)


async def get_current_user(
    db: Annotated[Session, Depends(get_db)],
    api_key: Annotated[str | None, Header(alias=settings.api_key_header)] = None,
) -> User:
    """
    Get the current user from the provided API key.
    """
    if not api_key:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Not authenticated",
        )

    # A single query handles both the bootstrap key and normal user keys.
    user = db.query(User).filter(User.access_key == api_key).first()

    if not user:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid API Key",
        )

    return user


# Dependency for protected endpoints
CurrentUser = Annotated[User, Depends(get_current_user)]

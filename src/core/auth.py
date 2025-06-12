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
    Get the current user from the provided API key by looking them up in the
    database.
    """
    logger.info(f"Attempting to authenticate with API key: {api_key}")

    if not api_key:
        logger.warning("Authentication failed: API key is missing.")
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Not authenticated",
        )
    user = db.query(User).filter(User.access_key == api_key).first()

    if not user:
        logger.warning("Authentication failed: Invalid API Key provided.")
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid API Key",
        )

    logger.info(f"Successfully authenticated user: {user.username}")
    return user


# Dependency for protected endpoints
CurrentUser = Annotated[User, Depends(get_current_user)]

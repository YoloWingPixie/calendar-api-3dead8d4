"""API Key authentication dependency."""

import logging
from typing import Annotated
from uuid import uuid4

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
        logger.warning("API key missing")
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Not authenticated",
        )

    # The bootstrap key is a valid API key for the root user
    if api_key == settings.bootstrap_admin_key:
        user = db.query(User).filter(User.username == "root").first()
        if user:
            logger.info(f"Authenticated user '{user.username}'")
            return user
        # If root user doesn't exist, create it in memory for this request
        return User(user_id=uuid4(), username="root", access_key=api_key)

    logger.debug("Looking up user by API key")
    user = db.query(User).filter(User.access_key == api_key).first()
    if not user:
        logger.warning("Invalid API key provided")
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Invalid API Key",
        )

    logger.info(f"Authenticated user '{user.username}'")
    return user


# Dependency for protected endpoints
CurrentUser = Annotated[User, Depends(get_current_user)]

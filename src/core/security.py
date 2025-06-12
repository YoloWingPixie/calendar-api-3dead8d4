"""
API Security and Authentication.

This module provides the dependency for API key-based authentication.
"""

import logging

from fastapi import Depends, HTTPException, Security
from fastapi.security import APIKeyHeader
from sqlalchemy.orm import Session
from starlette import status

from src.core.config import settings
from src.core.database import get_db
from src.models import User

logger = logging.getLogger(__name__)

api_key_header = APIKeyHeader(name=settings.api_key_header, auto_error=False)


class APIKeyInvalid(HTTPException):
    """Custom exception for invalid API keys."""

    def __init__(self, detail: str = "Invalid API Key"):
        super().__init__(status_code=status.HTTP_401_UNAUTHORIZED, detail=detail)


class APIKeyMissing(HTTPException):
    """Custom exception for missing API keys."""

    def __init__(self, detail: str = "API Key is missing"):
        super().__init__(status_code=status.HTTP_401_UNAUTHORIZED, detail=detail)


def get_current_user(
    api_key: str | None = Security(api_key_header),
    db: Session = Depends(get_db),  # noqa: B008
) -> User:
    """
    Dependency to get the current user from the provided API key.

    Args:
        api_key: The API key from the request header.
        db: The database session dependency.

    Returns:
        The authenticated User object.

    Raises:
        APIKeyMissing: If the API key is not provided in the request header.
        APIKeyInvalid: If the provided API key does not match any user.
    """
    if not api_key:
        logger.warning("API key missing from request.")
        raise APIKeyMissing()

    user = db.query(User).filter(User.access_key == api_key).first()
    if not user:
        logger.warning("Invalid API key provided: %s", api_key)
        raise APIKeyInvalid()

    logger.info("Successfully authenticated user: %s", user.username)
    return user

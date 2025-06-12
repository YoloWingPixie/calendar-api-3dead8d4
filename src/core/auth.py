"""API Key authentication dependency."""

from typing import Annotated

from fastapi import Depends, Header, HTTPException, status
from sqlalchemy.orm import Session

from src.core.config import settings
from src.core.database import get_db
from src.models.user import User


async def get_current_user(
    db: Annotated[Session, Depends(get_db)],
    api_key: Annotated[str | None, Header(alias=settings.api_key_header)] = None,
) -> User:
    """
    Validate API key and return the authenticated user.

    Args:
        api_key: API key from X-API-Key header
        db: Database session

    Returns:
        User: The authenticated user

    Raises:
        HTTPException: If API key is missing or invalid
    """
    if not api_key:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail={
                "error": "Unauthorized",
                "detail": "Missing API key",
                "code": "UNAUTHORIZED",
            },
        )

    # Check database for user with this API key
    # This includes the root user seeded with BOOTSTRAP_ADMIN_KEY
    user = db.query(User).filter(User.access_key == api_key).first()
    if not user:
        raise HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail={
                "error": "Unauthorized",
                "detail": "Invalid API key",
                "code": "UNAUTHORIZED",
            },
        )

    return user


# Dependency for protected endpoints
CurrentUser = Annotated[User, Depends(get_current_user)]

"""User management endpoints."""

import logging
import secrets
from typing import Annotated

from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.exc import IntegrityError
from sqlalchemy.orm import Session

from src.core.auth import CurrentUser
from src.core.database import get_db
from src.models.user import User
from src.schemas.user import UserCreate, UserWithAccessKey

router = APIRouter(prefix="/users", tags=["users"])
logger = logging.getLogger(__name__)


def generate_api_key() -> str:
    """Generate a secure API key."""
    return secrets.token_urlsafe(32)


@router.post(
    "",
    response_model=UserWithAccessKey,
    status_code=status.HTTP_201_CREATED,
    summary="Create a new user",
    description=(
        "Create a new user with a unique username. "
        "Returns the user details including the API key. "
        "Requires a valid API key to be provided."
    ),
)
async def create_user(
    user_data: UserCreate,
    current_user: CurrentUser,
    db: Annotated[Session, Depends(get_db)],
) -> UserWithAccessKey:
    """
    Create a new user and return their details with API key.

    Requires a valid API key (including the bootstrap key).
    """
    logger.info(
        "User '%s' attempting to create user '%s'",
        current_user.username,
        user_data.username,
    )
    # Any authenticated user can create another user
    access_key = generate_api_key()

    # Create new user
    new_user = User(
        username=user_data.username,
        access_key=access_key,
    )

    try:
        db.add(new_user)
        db.commit()
        db.refresh(new_user)
        logger.info(f"Successfully created user '{new_user.username}'")
    except IntegrityError as e:
        db.rollback()
        logger.error(f"Failed to create user '{user_data.username}': Username exists.")
        raise HTTPException(
            status_code=status.HTTP_409_CONFLICT,
            detail="Username already exists",
        ) from e

    # Return user with access key (only time it's shown)
    return UserWithAccessKey(
        user_id=new_user.user_id,
        username=new_user.username,
        access_key=new_user.access_key,
        created_at=new_user.created_at,
        updated_at=new_user.updated_at,
    )

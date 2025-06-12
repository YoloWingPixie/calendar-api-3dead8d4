"""Temporary debugging endpoints."""

from typing import Annotated

from fastapi import APIRouter, Depends
from sqlalchemy.orm import Session

from src.core.database import get_db
from src.models.user import User
from src.schemas.user import UserResponse

router = APIRouter(prefix="/debug", tags=["debug"])


@router.get("/db-check", response_model=list[UserResponse])
async def db_check(
    db: Annotated[Session, Depends(get_db)],
):
    """
    Returns all users from the database. This is a temporary endpoint for
    debugging database connection and seeding issues.
    """
    users = db.query(User).all()
    return users

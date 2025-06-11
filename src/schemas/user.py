"""User-related schemas."""

from datetime import datetime
from uuid import UUID

from pydantic import Field, field_validator

from src.schemas.base import CustomModel


class UserCreate(CustomModel):
    """Schema for creating a new user."""

    username: str = Field(
        ...,
        min_length=1,
        max_length=255,
        pattern="^[A-Za-z0-9-_]+$",
        description="Unique username",
    )

    @field_validator("username", mode="after")
    @classmethod
    def validate_username(cls, v: str) -> str:
        """Validate username format."""
        if not v or v.isspace():
            raise ValueError("Username cannot be empty or whitespace")
        return v.lower()


class UserResponse(CustomModel):
    """Schema for user response."""

    user_id: UUID = Field(..., description="Unique user identifier")
    username: str = Field(..., description="Username")
    owned_calendar_ids: list[UUID] = Field(
        default_factory=list, description="List of calendar IDs owned by the user"
    )
    created_at: datetime = Field(..., description="User creation timestamp")
    updated_at: datetime = Field(..., description="User last update timestamp")


class UserWithAccessKey(UserResponse):
    """Schema for user response with access key (only returned on creation)."""

    access_key: str = Field(..., description="API access key")

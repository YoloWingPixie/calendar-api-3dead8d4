"""Calendar-related schemas."""

from datetime import datetime
from uuid import UUID

from pydantic import Field

from src.schemas.base import CustomModel


class CalendarCreate(CustomModel):
    """Schema for creating a new calendar."""

    calendar_name: str = Field(
        ..., min_length=1, max_length=255, description="Calendar name", alias="name"
    )
    editor_ids: list[UUID] = Field(
        default_factory=list, description="List of user IDs with editor permissions"
    )
    reader_ids: list[UUID] = Field(
        default_factory=list, description="List of user IDs with reader permissions"
    )
    public_read: bool = Field(
        False, description="Whether the calendar is publicly readable"
    )
    public_write: bool = Field(
        False, description="Whether the calendar is publicly writable"
    )


class CalendarUpdate(CustomModel):
    """Schema for updating a calendar."""

    calendar_name: str | None = Field(
        None, min_length=1, max_length=255, description="Calendar name"
    )
    editor_ids: list[UUID] | None = Field(
        None, description="List of user IDs with editor permissions"
    )
    reader_ids: list[UUID] | None = Field(
        None, description="List of user IDs with reader permissions"
    )
    public_read: bool | None = Field(
        None, description="Whether the calendar is publicly readable"
    )
    public_write: bool | None = Field(
        None, description="Whether the calendar is publicly writable"
    )


class CalendarResponse(CustomModel):
    """Schema for calendar response."""

    calendar_id: UUID = Field(..., description="Unique calendar identifier")
    owner_user_id: UUID = Field(..., description="Calendar owner user ID")
    calendar_name: str = Field(..., description="Calendar name", alias="name")
    editor_ids: list[UUID] = Field(
        ..., description="List of user IDs with editor permissions"
    )
    reader_ids: list[UUID] = Field(
        ..., description="List of user IDs with reader permissions"
    )
    public_read: bool = Field(
        ..., description="Whether the calendar is publicly readable"
    )
    public_write: bool = Field(
        ..., description="Whether the calendar is publicly writable"
    )
    created_at: datetime = Field(..., description="Calendar creation timestamp")
    updated_at: datetime = Field(..., description="Calendar last update timestamp")

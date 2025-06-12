"""Calendar event-related schemas."""

from datetime import datetime
from uuid import UUID

from pydantic import Field

from src.schemas.base import CustomModel


class CalendarEventCreate(CustomModel):
    """Schema for creating a new event."""

    title: str = Field(..., min_length=1, max_length=255)
    description: str | None = None
    start_time: datetime
    end_time: datetime


class CalendarEventUpdate(CustomModel):
    """Schema for updating an event."""

    title: str | None = Field(None, min_length=1, max_length=255)
    description: str | None = None
    start_time: datetime | None = None
    end_time: datetime | None = None


class CalendarEventResponse(CustomModel):
    """Schema for an event response."""

    event_id: UUID
    title: str
    description: str | None
    start_time: datetime
    end_time: datetime
    created_at: datetime
    updated_at: datetime

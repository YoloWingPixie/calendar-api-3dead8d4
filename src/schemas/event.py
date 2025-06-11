"""Legacy event-related schemas for backward compatibility."""

from datetime import datetime
from uuid import UUID

from pydantic import Field, model_validator

from src.schemas.base import CustomModel


class EventCreate(CustomModel):
    """Schema for creating a new event (legacy API)."""

    title: str = Field(..., min_length=1, max_length=255, description="Event title")
    description: str | None = Field(None, description="Event description")
    start_time: datetime = Field(..., description="Event start time")
    end_time: datetime = Field(..., description="Event end time")

    @model_validator(mode="after")
    def validate_times(self) -> "EventCreate":
        """Validate event times."""
        if self.end_time <= self.start_time:
            raise ValueError("End time must be after start time")

        duration = self.end_time - self.start_time
        if duration.total_seconds() < 60:
            raise ValueError("Event must be at least 1 minute long")

        return self


class EventUpdate(CustomModel):
    """Schema for updating an event (legacy API)."""

    title: str | None = Field(
        None, min_length=1, max_length=255, description="Event title"
    )
    description: str | None = Field(None, description="Event description")
    start_time: datetime | None = Field(None, description="Event start time")
    end_time: datetime | None = Field(None, description="Event end time")

    @model_validator(mode="after")
    def validate_times(self) -> "EventUpdate":
        """Validate event times if both are provided."""
        if self.start_time is not None and self.end_time is not None:
            if self.end_time <= self.start_time:
                raise ValueError("End time must be after start time")

            duration = self.end_time - self.start_time
            if duration.total_seconds() < 60:
                raise ValueError("Event must be at least 1 minute long")

        return self


class EventResponse(CustomModel):
    """Schema for event response (legacy API)."""

    id: UUID = Field(..., description="Unique event identifier", alias="event_id")
    title: str = Field(..., description="Event title")
    description: str | None = Field(None, description="Event description")
    start_time: datetime = Field(..., description="Event start time")
    end_time: datetime = Field(..., description="Event end time")
    created_at: datetime = Field(..., description="Event creation timestamp")
    updated_at: datetime = Field(..., description="Event last update timestamp")

    model_config = CustomModel.model_config.copy()
    model_config["populate_by_name"] = True

"""Calendar event-related schemas."""

from datetime import datetime
from uuid import UUID

from pydantic import Field, model_validator

from src.schemas.base import CustomModel


class CalendarEventCreate(CustomModel):
    """Schema for creating a new calendar event."""

    title: str = Field(..., min_length=1, max_length=255, description="Event title")
    description: str | None = Field(None, description="Event description")
    start_time: datetime = Field(..., description="Event start time")
    end_time: datetime = Field(..., description="Event end time")
    is_all_day: bool = Field(False, description="Whether the event is all-day")

    @model_validator(mode="after")
    def validate_times(self) -> "CalendarEventCreate":
        """Validate event times."""
        if self.end_time <= self.start_time:
            raise ValueError("End time must be after start time")

        duration = self.end_time - self.start_time
        if duration.total_seconds() < 60:  # Less than 1 minute
            raise ValueError("Event must be at least 1 minute long")

        if self.is_all_day:
            # For all-day events, check if times are at 00:00:00 and 23:59:59
            if (
                self.start_time.time().hour != 0
                or self.start_time.time().minute != 0
                or self.start_time.time().second != 0
            ):
                raise ValueError("All-day events must start at 00:00:00")

            if (
                self.end_time.time().hour != 23
                or self.end_time.time().minute != 59
                or self.end_time.time().second != 59
            ):
                raise ValueError("All-day events must end at 23:59:59")

            if self.start_time.date() != self.end_time.date():
                raise ValueError("All-day events must start and end on the same day")

        return self


class CalendarEventUpdate(CustomModel):
    """Schema for updating a calendar event."""

    title: str | None = Field(
        None, min_length=1, max_length=255, description="Event title"
    )
    description: str | None = Field(None, description="Event description")
    start_time: datetime | None = Field(None, description="Event start time")
    end_time: datetime | None = Field(None, description="Event end time")
    is_all_day: bool | None = Field(None, description="Whether the event is all-day")

    @model_validator(mode="after")
    def validate_times(self) -> "CalendarEventUpdate":
        """Validate event times if both are provided."""
        if self.start_time is not None and self.end_time is not None:
            if self.end_time <= self.start_time:
                raise ValueError("End time must be after start time")

            duration = self.end_time - self.start_time
            if duration.total_seconds() < 60:
                raise ValueError("Event must be at least 1 minute long")

        if self.is_all_day is True and (
            self.start_time is not None or self.end_time is not None
        ):
            # If updating to all-day, validate the times
            start = self.start_time
            end = self.end_time

            if start and (
                start.time().hour != 0
                or start.time().minute != 0
                or start.time().second != 0
            ):
                raise ValueError("All-day events must start at 00:00:00")

            if end and (
                end.time().hour != 23
                or end.time().minute != 59
                or end.time().second != 59
            ):
                raise ValueError("All-day events must end at 23:59:59")

            if start and end and start.date() != end.date():
                raise ValueError("All-day events must start and end on the same day")

        return self


class CalendarEventResponse(CustomModel):
    """Schema for calendar event response."""

    event_id: UUID = Field(..., description="Unique event identifier")
    calendar_id: UUID = Field(..., description="Calendar ID this event belongs to")
    creator_user_id: UUID = Field(..., description="User ID who created the event")
    title: str = Field(..., description="Event title")
    description: str | None = Field(None, description="Event description")
    start_time: datetime = Field(..., description="Event start time")
    end_time: datetime = Field(..., description="Event end time")
    is_all_day: bool = Field(..., description="Whether the event is all-day")
    created_at: datetime = Field(..., description="Event creation timestamp")
    updated_at: datetime = Field(..., description="Event last update timestamp")

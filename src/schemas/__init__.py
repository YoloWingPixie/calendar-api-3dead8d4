"""Expose schemas for easier imports."""

from src.schemas.base import CustomModel
from src.schemas.calendar_event import (
    CalendarEventCreate,
    CalendarEventResponse,
    CalendarEventUpdate,
)

__all__ = [
    "CustomModel",
    "CalendarEventCreate",
    "CalendarEventResponse",
    "CalendarEventUpdate",
]

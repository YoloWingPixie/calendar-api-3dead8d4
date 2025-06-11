"""Pydantic schemas for API validation."""

from src.schemas.base import (
    ErrorResponse,
    HealthResponse,
    PaginatedResponse,
    PaginationParams,
)
from src.schemas.calendar import CalendarCreate, CalendarResponse, CalendarUpdate
from src.schemas.calendar_event import (
    CalendarEventCreate,
    CalendarEventResponse,
    CalendarEventUpdate,
)
from src.schemas.event import EventCreate, EventResponse, EventUpdate
from src.schemas.user import UserCreate, UserResponse, UserWithAccessKey

__all__ = [
    # Base schemas
    "ErrorResponse",
    "HealthResponse",
    "PaginatedResponse",
    "PaginationParams",
    # User schemas
    "UserCreate",
    "UserResponse",
    "UserWithAccessKey",
    # Calendar schemas
    "CalendarCreate",
    "CalendarUpdate",
    "CalendarResponse",
    # Calendar Event schemas
    "CalendarEventCreate",
    "CalendarEventUpdate",
    "CalendarEventResponse",
    # Legacy Event schemas
    "EventCreate",
    "EventUpdate",
    "EventResponse",
]

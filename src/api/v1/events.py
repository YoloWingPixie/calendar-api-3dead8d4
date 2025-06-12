"""Event management endpoints."""

import logging
from typing import Annotated
from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session

from src.core.auth import CurrentUser
from src.core.database import get_db
from src.models.calendar_event import CalendarEvent
from src.schemas.calendar_event import (
    CalendarEventCreate,
    CalendarEventResponse,
    CalendarEventUpdate,
)

router = APIRouter(prefix="/events", tags=["events"])
logger = logging.getLogger(__name__)


@router.post(
    "", response_model=CalendarEventResponse, status_code=status.HTTP_201_CREATED
)
async def create_event(
    event_data: CalendarEventCreate,
    current_user: CurrentUser,
    db: Annotated[Session, Depends(get_db)],
) -> CalendarEvent:
    """Create a new event."""
    logger.info(f"User '{current_user.username}' creating event '{event_data.title}'")
    new_event = CalendarEvent(**event_data.model_dump())
    db.add(new_event)
    db.commit()
    db.refresh(new_event)
    logger.info(f"Successfully created event '{new_event.event_id}'")
    return new_event


@router.get("", response_model=list[CalendarEventResponse])
async def get_events(
    db: Annotated[Session, Depends(get_db)],
) -> list[CalendarEvent]:
    """Get all events."""
    events = db.query(CalendarEvent).all()
    return events


@router.get("/{event_id}", response_model=CalendarEventResponse)
async def get_event(
    event_id: UUID,
    db: Annotated[Session, Depends(get_db)],
) -> CalendarEvent:
    """Get a single event by ID."""
    event = db.query(CalendarEvent).filter(CalendarEvent.event_id == event_id).first()
    if not event:
        raise HTTPException(status_code=404, detail="Event not found")
    return event


@router.put("/{event_id}", response_model=CalendarEventResponse)
async def update_event(
    event_id: UUID,
    event_data: CalendarEventUpdate,
    current_user: CurrentUser,
    db: Annotated[Session, Depends(get_db)],
) -> CalendarEvent:
    """Update an event."""
    event = db.query(CalendarEvent).filter(CalendarEvent.event_id == event_id).first()
    if not event:
        raise HTTPException(status_code=404, detail="Event not found")

    logger.info(f"User '{current_user.username}' updating event '{event_id}'")
    update_data = event_data.model_dump(exclude_unset=True)
    for key, value in update_data.items():
        setattr(event, key, value)

    db.commit()
    db.refresh(event)
    logger.info(f"Successfully updated event '{event.event_id}'")
    return event


@router.delete("/{event_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_event(
    event_id: UUID,
    current_user: CurrentUser,
    db: Annotated[Session, Depends(get_db)],
) -> None:
    """Delete an event."""
    event = db.query(CalendarEvent).filter(CalendarEvent.event_id == event_id).first()
    if not event:
        raise HTTPException(status_code=404, detail="Event not found")

    logger.info(f"User '{current_user.username}' deleting event '{event_id}'")
    db.delete(event)
    db.commit()
    logger.info(f"Successfully deleted event '{event.event_id}'")
    return None

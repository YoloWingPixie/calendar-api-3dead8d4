"""Calendar event management endpoints."""

import logging
from typing import Annotated
from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session

from src.core.auth import CurrentUser
from src.core.database import get_db
from src.models.calendar import Calendar
from src.models.calendar_event import CalendarEvent
from src.schemas.calendar_event import (
    CalendarEventCreate,
    CalendarEventResponse,
    CalendarEventUpdate,
)

# The router will be mounted under /calendars/{calendar_id}, so paths are relative
router = APIRouter(prefix="/events", tags=["events"])
logger = logging.getLogger(__name__)


@router.post(
    "", response_model=CalendarEventResponse, status_code=status.HTTP_201_CREATED
)
async def create_event(
    calendar_id: UUID,
    event_data: CalendarEventCreate,
    current_user: CurrentUser,
    db: Annotated[Session, Depends(get_db)],
):
    """Create a new event in a calendar."""
    logger.info(
        f"User '{current_user.username}' creating event in calendar '{calendar_id}'"
    )
    calendar = db.query(Calendar).filter(Calendar.calendar_id == calendar_id).first()
    if not calendar:
        raise HTTPException(status_code=404, detail="Calendar not found")

    # Simple auth: only calendar owner can add events
    if calendar.owner_user_id != current_user.user_id:
        raise HTTPException(
            status_code=403, detail="Not authorized to create events in this calendar"
        )

    new_event = CalendarEvent(
        **event_data.model_dump(),
        calendar_id=calendar_id,
        creator_user_id=current_user.user_id,
    )
    db.add(new_event)
    db.commit()
    db.refresh(new_event)
    logger.info(f"Successfully created event '{new_event.event_id}'")
    return new_event


@router.get("", response_model=list[CalendarEventResponse])
async def get_events(
    calendar_id: UUID,
    current_user: CurrentUser,
    db: Annotated[Session, Depends(get_db)],
):
    """Get all events for a calendar."""
    calendar = db.query(Calendar).filter(Calendar.calendar_id == calendar_id).first()
    if not calendar:
        raise HTTPException(status_code=404, detail="Calendar not found")

    # Simple auth: only calendar owner can see events
    if calendar.owner_user_id != current_user.user_id:
        raise HTTPException(
            status_code=403, detail="Not authorized to view events in this calendar"
        )

    events = (
        db.query(CalendarEvent).filter(CalendarEvent.calendar_id == calendar_id).all()
    )
    return events


@router.patch("/{event_id}", response_model=CalendarEventResponse)
async def update_event(
    calendar_id: UUID,
    event_id: UUID,
    event_data: CalendarEventUpdate,
    current_user: CurrentUser,
    db: Annotated[Session, Depends(get_db)],
):
    """Update an event."""
    event = db.query(CalendarEvent).filter(CalendarEvent.event_id == event_id).first()
    if not event or event.calendar_id != calendar_id:
        raise HTTPException(status_code=404, detail="Event not found")

    # Simple auth: only event creator can update it
    if event.creator_user_id != current_user.user_id:
        raise HTTPException(
            status_code=403, detail="Not authorized to update this event"
        )

    update_data = event_data.model_dump(exclude_unset=True)
    for key, value in update_data.items():
        setattr(event, key, value)

    db.commit()
    db.refresh(event)
    logger.info(f"Successfully updated event '{event.event_id}'")
    return event


@router.delete("/{event_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_event(
    calendar_id: UUID,
    event_id: UUID,
    current_user: CurrentUser,
    db: Annotated[Session, Depends(get_db)],
):
    """Delete an event."""
    event = db.query(CalendarEvent).filter(CalendarEvent.event_id == event_id).first()
    if not event or event.calendar_id != calendar_id:
        raise HTTPException(status_code=404, detail="Event not found")

    # Simple auth: only event creator can delete it
    if event.creator_user_id != current_user.user_id:
        raise HTTPException(
            status_code=403, detail="Not authorized to delete this event"
        )

    db.delete(event)
    db.commit()
    logger.info(f"Successfully deleted event '{event.event_id}'")
    return None

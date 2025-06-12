"""Event-related API endpoints."""

import logging
from typing import Any
from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session

from src.core.database import get_db
from src.core.security import get_current_user
from src.models import CalendarEvent, User
from src.schemas.event import EventCreate, EventResponse, EventUpdate

router = APIRouter()
logger = logging.getLogger(__name__)


@router.get("/", response_model=list[EventResponse])
def list_events(
    db: Session = Depends(get_db),  # noqa: B008
    current_user: User = Depends(get_current_user),  # noqa: B008
) -> Any:
    """Retrieve all events for the current user."""
    logger.info("User '%s' listing all events.", current_user.username)
    # This is a toy app, so for now we return all events.
    # In a real app, this would be filtered by user/calendar permissions.
    events = db.query(CalendarEvent).all()
    logger.info("Found %d events", len(events))
    return events


@router.post("/", response_model=EventResponse, status_code=status.HTTP_201_CREATED)
def create_event(
    event_in: EventCreate,
    db: Session = Depends(get_db),  # noqa: B008
    current_user: User = Depends(get_current_user),  # noqa: B008
) -> Any:
    """Create a new event."""
    logger.info("User '%s' creating event: '%s'", current_user.username, event_in.title)
    db_event = CalendarEvent(**event_in.model_dump())
    db.add(db_event)
    logger.info("Added event to session, committing...")
    db.commit()
    logger.info("Event committed, refreshing...")
    db.refresh(db_event)
    logger.info("Event created with ID: %s", db_event.event_id)
    return db_event


@router.get("/{event_id}", response_model=EventResponse)
def get_event(
    event_id: UUID,
    db: Session = Depends(get_db),  # noqa: B008
    current_user: User = Depends(get_current_user),  # noqa: B008
) -> Any:
    """Get a single event by ID."""
    logger.info("User '%s' fetching event ID: %s", current_user.username, event_id)
    event = db.query(CalendarEvent).filter(CalendarEvent.event_id == event_id).first()
    if not event:
        logger.warning("Event with ID %s not found", event_id)
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND, detail="Event not found"
        )
    logger.info("Found event: %s", event.title)
    return event


@router.put("/{event_id}", response_model=EventResponse)
def update_event(
    event_id: UUID,
    event_in: EventUpdate,
    db: Session = Depends(get_db),  # noqa: B008
    current_user: User = Depends(get_current_user),  # noqa: B008
) -> Any:
    """Update an existing event."""
    logger.info("User '%s' updating event ID: %s", current_user.username, event_id)
    db_event = (
        db.query(CalendarEvent).filter(CalendarEvent.event_id == event_id).first()
    )
    if not db_event:
        logger.warning("Event with ID %s not found for update", event_id)
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND, detail="Event not found"
        )

    update_data = event_in.model_dump(exclude_unset=True)
    logger.info("Updating event with data: %s", update_data)
    for field, value in update_data.items():
        setattr(db_event, field, value)

    db.commit()
    db.refresh(db_event)
    logger.info("Event updated successfully")
    return db_event


@router.delete("/{event_id}", status_code=status.HTTP_204_NO_CONTENT)
def delete_event(
    event_id: UUID,
    db: Session = Depends(get_db),  # noqa: B008
    current_user: User = Depends(get_current_user),  # noqa: B008
) -> None:
    """Delete an event."""
    logger.info("User '%s' deleting event ID: %s", current_user.username, event_id)
    db_event = (
        db.query(CalendarEvent).filter(CalendarEvent.event_id == event_id).first()
    )
    if not db_event:
        logger.info("Event with ID %s not found for deletion (idempotent)", event_id)
        # No exception on delete if not found, it's idempotent.
        return None
    db.delete(db_event)
    db.commit()
    logger.info("Event deleted successfully")
    return None

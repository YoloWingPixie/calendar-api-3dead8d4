"""Calendar management endpoints."""

import logging
from typing import Annotated
from uuid import UUID

from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.orm import Session

from src.core.auth import CurrentUser
from src.core.database import get_db
from src.models.calendar import Calendar
from src.schemas.calendar import CalendarCreate, CalendarResponse, CalendarUpdate

from . import events as events_router

router = APIRouter(prefix="/calendars", tags=["calendars"])
logger = logging.getLogger(__name__)

# Mount the events router under /calendars/{calendar_id}
router.include_router(events_router.router, prefix="/{calendar_id}")


@router.post(
    "",
    response_model=CalendarResponse,
    status_code=status.HTTP_201_CREATED,
    summary="Create a new calendar",
    description="Create a new calendar owned by the authenticated user.",
)
async def create_calendar(
    calendar_data: CalendarCreate,
    current_user: CurrentUser,
    db: Annotated[Session, Depends(get_db)],
) -> CalendarResponse:
    """
    Create a new calendar for the authenticated user.

    Requires API key authentication via X-API-Key header.
    """
    logger.info(
        "User '%s' creating calendar '%s'",
        current_user.username,
        calendar_data.calendar_name,
    )
    # Create new calendar
    new_calendar = Calendar(
        calendar_name=calendar_data.calendar_name,
        owner_user_id=current_user.user_id,
    )

    db.add(new_calendar)
    db.commit()
    db.refresh(new_calendar)
    logger.info(f"Successfully created calendar '{new_calendar.calendar_name}'")

    return CalendarResponse(
        calendar_id=new_calendar.calendar_id,
        calendar_name=new_calendar.calendar_name,
        owner_user_id=new_calendar.owner_user_id,
        editor_ids=new_calendar.editor_ids or [],
        reader_ids=new_calendar.reader_ids or [],
        public_read=new_calendar.public_read,
        public_write=new_calendar.public_write,
        created_at=new_calendar.created_at,
        updated_at=new_calendar.updated_at,
    )


@router.get("", response_model=list[CalendarResponse])
async def get_calendars(
    current_user: CurrentUser,
    db: Annotated[Session, Depends(get_db)],
):
    """List all calendars owned by the current user."""
    calendars = (
        db.query(Calendar).filter(Calendar.owner_user_id == current_user.user_id).all()
    )
    return calendars


@router.get("/{calendar_id}", response_model=CalendarResponse)
async def get_calendar(
    calendar_id: UUID,
    current_user: CurrentUser,
    db: Annotated[Session, Depends(get_db)],
):
    """Get a single calendar by ID."""
    calendar = db.query(Calendar).filter(Calendar.calendar_id == calendar_id).first()
    if not calendar:
        raise HTTPException(status_code=404, detail="Calendar not found")
    # Simple auth: anyone logged in can read any calendar for now
    return calendar


@router.patch("/{calendar_id}", response_model=CalendarResponse)
async def update_calendar(
    calendar_id: UUID,
    calendar_data: CalendarUpdate,
    current_user: CurrentUser,
    db: Annotated[Session, Depends(get_db)],
):
    """Update a calendar."""
    calendar = db.query(Calendar).filter(Calendar.calendar_id == calendar_id).first()
    if not calendar:
        raise HTTPException(status_code=404, detail="Calendar not found")

    if calendar.owner_user_id != current_user.user_id:
        raise HTTPException(
            status_code=403, detail="Not authorized to update this calendar"
        )

    update_data = calendar_data.model_dump(exclude_unset=True)
    for key, value in update_data.items():
        setattr(calendar, key, value)

    db.commit()
    db.refresh(calendar)
    return calendar


@router.delete("/{calendar_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_calendar(
    calendar_id: UUID,
    current_user: CurrentUser,
    db: Annotated[Session, Depends(get_db)],
):
    """Delete a calendar."""
    calendar = db.query(Calendar).filter(Calendar.calendar_id == calendar_id).first()
    if not calendar:
        raise HTTPException(status_code=404, detail="Calendar not found")

    if calendar.owner_user_id != current_user.user_id:
        raise HTTPException(
            status_code=403, detail="Not authorized to delete this calendar"
        )

    db.delete(calendar)
    db.commit()
    return None

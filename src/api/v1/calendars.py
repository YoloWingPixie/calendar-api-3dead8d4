"""Calendar management endpoints."""

from typing import Annotated

from fastapi import APIRouter, Depends, status
from sqlalchemy.orm import Session

from src.core.auth import CurrentUser
from src.core.database import get_db
from src.models.calendar import Calendar
from src.schemas.calendar import CalendarCreate, CalendarResponse

router = APIRouter(prefix="/calendars", tags=["calendars"])


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
    # Create new calendar
    new_calendar = Calendar(
        calendar_name=calendar_data.calendar_name,
        owner_user_id=current_user.user_id,
    )

    db.add(new_calendar)
    db.commit()
    db.refresh(new_calendar)

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

"""User model."""

from __future__ import annotations

from datetime import datetime
from typing import TYPE_CHECKING
from uuid import UUID as py_UUID

import sqlalchemy as sa
from sqlalchemy import String
from sqlalchemy.dialects.postgresql import TIMESTAMP
from sqlalchemy.dialects.postgresql import UUID as pg_UUID
from sqlalchemy.orm import Mapped, mapped_column, relationship

from src.models.base import Base

if TYPE_CHECKING:
    from src.models.calendar import Calendar


class User(Base):
    """User model representing system users."""

    __tablename__ = "users"

    user_id: Mapped[py_UUID] = mapped_column(
        pg_UUID(as_uuid=True),
        primary_key=True,
        server_default=sa.text("gen_random_uuid()"),
    )
    username: Mapped[str] = mapped_column(String(255), unique=True, nullable=False)
    access_key: Mapped[str] = mapped_column(String(255), unique=True, nullable=False)
    created_at: Mapped[datetime] = mapped_column(
        TIMESTAMP(timezone=True), server_default=sa.text("now()"), nullable=False
    )
    updated_at: Mapped[datetime] = mapped_column(
        TIMESTAMP(timezone=True),
        server_default=sa.text("now()"),
        onupdate=sa.text("now()"),
        nullable=False,
    )

    # Relationships
    calendars: Mapped[list[Calendar]] = relationship(
        "Calendar",
        back_populates="owner",
        foreign_keys="Calendar.owner_user_id",
        cascade="all, delete-orphan",
    )

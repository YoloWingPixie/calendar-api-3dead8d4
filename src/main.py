"""Calendar API main module."""

import logging
from contextlib import asynccontextmanager

from fastapi import FastAPI

from src.core.config import settings
from src.core.database import SessionLocal
from src.core.router import register_routers
from src.models.calendar_event import CalendarEvent
from src.models.user import User


@asynccontextmanager
async def lifespan(app: FastAPI):
    # On startup
    logging.basicConfig(
        level=settings.log_level.upper(),
        format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    )
    logger = logging.getLogger(__name__)
    logger.info("Application startup...")
    logger.info(f"Connecting to database: {settings.database_url}")

    db = SessionLocal()
    try:
        user_count = db.query(User).count()
        event_count = db.query(CalendarEvent).count()
        logger.info(f"Database contains {user_count} users and {event_count} events.")
    finally:
        db.close()

    yield
    # On shutdown
    logger.info("Application shutdown.")


def create_app() -> FastAPI:
    """Create and configure the FastAPI application."""
    app = FastAPI(
        title=settings.app_name,
        version=settings.version,
        debug=settings.debug,
        lifespan=lifespan,
    )

    # Register all routes dynamically
    register_routers(app)

    return app


# Create application instance
app = create_app()


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        "src.main:app",
        host=settings.host,
        port=settings.port,
        reload=settings.debug,
    )

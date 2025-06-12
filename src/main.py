"""Calendar API main module."""

import logging
from collections.abc import AsyncIterator
from contextlib import asynccontextmanager

from fastapi import FastAPI

from src.core.config import settings
from src.core.database import SessionLocal
from src.core.router import register_routers
from src.models import User


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncIterator[None]:
    # On startup, seed the root user if it doesn't exist.
    logger = logging.getLogger(__name__)
    logger.info("Application startup...")

    db = SessionLocal()
    try:
        if settings.bootstrap_admin_key:
            user = db.query(User).filter(User.username == "root").first()
            if not user:
                logger.info("Root user not found, creating it...")
                root_user = User(
                    username="root", access_key=settings.bootstrap_admin_key
                )
                db.add(root_user)
                db.commit()
                logger.info("Root user created successfully.")
    finally:
        db.close()

    yield
    # On shutdown
    logger.info("Application shutdown.")


def create_app(use_lifespan: bool = True) -> FastAPI:
    """Create and configure the FastAPI application."""
    app_lifespan = lifespan if use_lifespan else None
    app = FastAPI(
        title=settings.app_name,
        version=settings.version,
        debug=settings.debug,
        lifespan=app_lifespan,
    )

    # Configure logging
    logging.basicConfig(
        level=settings.log_level.upper(),
        format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    )

    # Register all routes dynamically
    register_routers(app)

    return app


app = create_app()


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(
        "src.main:app",
        host=settings.host,
        port=settings.port,
        reload=settings.debug,
    )

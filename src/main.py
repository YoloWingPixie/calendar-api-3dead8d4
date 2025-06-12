"""Calendar API main module."""

import logging
from collections.abc import AsyncIterator
from contextlib import asynccontextmanager

from fastapi import FastAPI

from src.core.config import settings
from src.core.database import SessionLocal
from src.core.router import register_routers
from src.models import User

# Configure logging at the module level to ensure it's applied before
# Uvicorn or other modules can configure it. Using force=True ensures
# that our configuration replaces any existing handlers.
logging.basicConfig(
    level=settings.log_level.upper(),
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    force=True,
)


@asynccontextmanager
async def lifespan(app: FastAPI) -> AsyncIterator[None]:
    # On startup, seed the root user if it doesn't exist.
    logger = logging.getLogger(__name__)
    logger.info("Application startup: Seeding database...")

    db = None
    try:
        logger.info("Attempting to connect to the database for user seeding...")
        db = SessionLocal()
        logger.info("Database session created for user seeding.")

        if settings.bootstrap_admin_key:
            logger.info("`BOOTSTRAP_ADMIN_KEY` is set, checking for root user...")
            user = db.query(User).filter(User.username == "root").first()
            if not user:
                logger.info("Root user not found, creating it...")
                root_user = User(
                    username="root", access_key=settings.bootstrap_admin_key
                )
                db.add(root_user)
                db.commit()
                logger.info("Root user created and committed successfully.")
            else:
                logger.info("Root user already exists, skipping creation.")
        else:
            logger.warning(
                "`BOOTSTRAP_ADMIN_KEY` not set, skipping root user creation."
            )
    except Exception as e:
        logger.exception("Error during startup user seeding: %s", e)
    finally:
        if db:
            db.close()
            logger.info("Database session for user seeding closed.")

    yield
    # On shutdown
    logger.info("Application shutdown.")


def create_app(use_lifespan: bool = True) -> FastAPI:
    """Create and configure the FastAPI application."""
    # Logging is now configured at the module level.

    app_lifespan = lifespan if use_lifespan else None
    app = FastAPI(
        title=settings.app_name,
        version=settings.version,
        debug=settings.debug,
        lifespan=app_lifespan,
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

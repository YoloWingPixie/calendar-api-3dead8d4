"""Calendar API main module."""

import logging

from fastapi import FastAPI

from src.core.config import settings
from src.core.router import register_routers


def create_app() -> FastAPI:
    """Create and configure the FastAPI application."""
    app = FastAPI(
        title=settings.app_name,
        version=settings.version,
        debug=settings.debug,
    )

    # Configure logging
    logging.basicConfig(
        level=settings.log_level.upper(),
        format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
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

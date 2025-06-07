"""Database connection and session management."""

from collections.abc import Generator

from sqlalchemy import create_engine
from sqlalchemy.orm import Session, sessionmaker
from sqlalchemy.pool import NullPool

from src.core.config import settings

# Create engine with appropriate pooling for production
# For RDS, we use connection pooling with pre-ping to handle network issues
engine = create_engine(
    settings.database_url,
    # Connection pool settings optimized for production
    pool_size=20 if not settings.debug else 5,
    max_overflow=40 if not settings.debug else 10,
    pool_timeout=30,
    pool_recycle=3600,  # Recycle connections after 1 hour
    pool_pre_ping=True,  # Verify connections before use
    # Use NullPool for serverless environments like ECS Fargate
    poolclass=NullPool if settings.database_pool_disabled else None,
    echo=settings.database_echo,
)

# Create session factory
SessionLocal = sessionmaker(
    bind=engine,
    autoflush=False,
    autocommit=False,
    expire_on_commit=False,
)


def get_db() -> Generator[Session]:
    """
    Dependency to get database session.

    Yields:
        Session: Database session
    """
    db = SessionLocal()
    try:
        yield db
        db.commit()
    except Exception:
        db.rollback()
        raise
    finally:
        db.close()

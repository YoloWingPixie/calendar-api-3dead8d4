"""Shared test fixtures and configuration."""

from unittest.mock import MagicMock

import pytest
from fastapi.testclient import TestClient

from src.core.database import get_db
from src.main import create_app


@pytest.fixture
def mock_db():
    """Create a mock database session."""
    # This can be a more sophisticated mock if needed, e.g., an in-memory SQLite
    yield MagicMock()


@pytest.fixture
def app(mock_db):
    """Create a FastAPI app instance with a mocked database."""
    app = create_app(use_lifespan=False)
    app.dependency_overrides[get_db] = lambda: mock_db
    return app


@pytest.fixture
def client(app):
    """Create a test client."""
    with TestClient(app) as client:
        yield client

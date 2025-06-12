"""Shared test fixtures and configuration."""

from unittest.mock import MagicMock

import pytest
from fastapi.testclient import TestClient

from src.core.database import get_db
from src.main import app


@pytest.fixture
def mock_db():
    """Create a mock database session."""
    db = MagicMock()
    return db


@pytest.fixture
def client(mock_db):
    """Create a test client with mocked database."""

    def override_get_db():
        yield mock_db

    app.dependency_overrides[get_db] = override_get_db

    with TestClient(app) as test_client:
        yield test_client

    # Clear overrides after test
    app.dependency_overrides.clear()

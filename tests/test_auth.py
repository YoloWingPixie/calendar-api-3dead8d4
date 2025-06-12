"""Test authentication functionality."""

from unittest.mock import MagicMock
from uuid import uuid4

from fastapi import status
from fastapi.testclient import TestClient

from src.models.user import User


def test_health_check_no_auth_required(client: TestClient) -> None:
    """Test that health check doesn't require authentication."""
    response = client.get("/api/v1/health")
    assert response.status_code == status.HTTP_200_OK
    assert response.json()["status"] == "ok"


def test_create_user_requires_auth(client: TestClient, mock_db, monkeypatch) -> None:
    """Test that user creation requires authentication."""
    # Try to create user without auth - should fail
    response = client.post(
        "/api/v1/users",
        json={"username": "testuser"},
    )
    assert response.status_code == status.HTTP_401_UNAUTHORIZED
    assert "Not authenticated" in response.text


def test_create_user_with_bootstrap_key(
    client: TestClient, mock_db, monkeypatch
) -> None:
    """Test that user creation works with bootstrap admin key."""
    from datetime import UTC, datetime

    # Set bootstrap key in environment
    monkeypatch.setenv("BOOTSTRAP_ADMIN_KEY", "test-bootstrap-key")

    # Mock the root user lookup
    mock_user = User(
        user_id=uuid4(),
        username="root",
        access_key="test-bootstrap-key",
        created_at=datetime.now(UTC),
        updated_at=datetime.now(UTC),
    )
    mock_query = MagicMock()
    mock_query.filter.return_value.first.return_value = mock_user
    mock_db.query.return_value = mock_query

    # Mock the database operations
    mock_db.add = MagicMock()
    mock_db.commit = MagicMock()

    # Mock refresh to populate the user attributes
    def mock_refresh(user):
        user.user_id = uuid4()
        user.created_at = datetime.now(UTC)
        user.updated_at = datetime.now(UTC)

    mock_db.refresh = MagicMock(side_effect=mock_refresh)

    # Use bootstrap key to create first user
    response = client.post(
        "/api/v1/users",
        json={"username": "admin"},
        headers={"X-API-Key": "test-bootstrap-key"},
    )
    assert response.status_code == status.HTTP_201_CREATED
    data = response.json()
    assert data["username"] == "admin"
    assert "access_key" in data
    assert "user_id" in data


def test_create_user_duplicate_username(
    client: TestClient, mock_db, monkeypatch
) -> None:
    """Test that duplicate usernames are rejected."""
    from datetime import UTC, datetime

    from sqlalchemy.exc import IntegrityError

    # Set bootstrap key
    monkeypatch.setenv("BOOTSTRAP_ADMIN_KEY", "test-bootstrap-key")

    # Mock the root user lookup
    mock_user = User(
        user_id=uuid4(),
        username="root",
        access_key="test-bootstrap-key",
        created_at=datetime.now(UTC),
        updated_at=datetime.now(UTC),
    )
    mock_query = MagicMock()
    mock_query.filter.return_value.first.return_value = mock_user
    mock_db.query.return_value = mock_query

    # Mock the database to raise IntegrityError on add
    mock_db.add = MagicMock(side_effect=IntegrityError("duplicate", None, None))
    mock_db.rollback = MagicMock()

    response = client.post(
        "/api/v1/users",
        json={"username": "duplicate"},
        headers={"X-API-Key": "test-bootstrap-key"},
    )
    assert response.status_code == status.HTTP_409_CONFLICT
    assert "Username already exists" in response.text


def test_protected_endpoint_without_auth(client: TestClient) -> None:
    """Test that protected endpoints require authentication."""
    response = client.post(
        "/api/v1/calendars",
        json={"name": "My Calendar"},
    )
    assert response.status_code == status.HTTP_401_UNAUTHORIZED
    assert "Not authenticated" in response.text


def test_protected_endpoint_with_invalid_auth(client: TestClient, mock_db) -> None:
    """Test that invalid API keys are rejected."""
    # Mock query to return no user
    mock_query = MagicMock()
    mock_query.filter.return_value.first.return_value = None
    mock_db.query.return_value = mock_query

    response = client.post(
        "/api/v1/calendars",
        json={"name": "My Calendar"},
        headers={"X-API-Key": "invalid-key"},
    )
    assert response.status_code == status.HTTP_401_UNAUTHORIZED
    assert "Invalid API Key" in response.text


def test_protected_endpoint_with_valid_auth(client: TestClient, mock_db) -> None:
    """Test that valid API keys allow access to protected endpoints."""
    from datetime import UTC, datetime

    # Create a mock user with all required attributes
    mock_user = MagicMock(spec=User)
    mock_user.user_id = uuid4()
    mock_user.username = "authtest"
    mock_user.access_key = "valid-key"
    mock_user.created_at = datetime.now(UTC)
    mock_user.updated_at = datetime.now(UTC)

    # Mock query to return the user
    mock_query = MagicMock()
    mock_query.filter.return_value.first.return_value = mock_user
    mock_db.query.return_value = mock_query

    # Mock database operations for calendar creation
    mock_db.add = MagicMock()
    mock_db.commit = MagicMock()

    # Mock refresh to populate calendar attributes
    def mock_refresh(calendar):
        calendar.calendar_id = uuid4()
        calendar.created_at = datetime.now(UTC)
        calendar.updated_at = datetime.now(UTC)
        calendar.editor_ids = []
        calendar.reader_ids = []
        calendar.public_read = False
        calendar.public_write = False

    mock_db.refresh = MagicMock(side_effect=mock_refresh)

    # Make request with valid API key
    response = client.post(
        "/api/v1/calendars",
        json={"name": "My Protected Calendar"},
        headers={"X-API-Key": "valid-key"},
    )

    # Debug: print response if it fails
    if response.status_code != status.HTTP_201_CREATED:
        print(f"Response status: {response.status_code}")
        print(f"Response body: {response.json()}")

    assert response.status_code == status.HTTP_201_CREATED
    calendar_data = response.json()
    assert calendar_data["name"] == "My Protected Calendar"
    assert str(calendar_data["owner_user_id"]) == str(mock_user.user_id)

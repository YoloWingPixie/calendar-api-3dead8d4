"""Tests for the /api/events endpoint."""

from datetime import datetime
from unittest.mock import MagicMock
from uuid import uuid4

from fastapi import status
from fastapi.testclient import TestClient

from src.models.calendar_event import CalendarEvent
from src.models.user import User


def test_create_event(client: TestClient, mock_db) -> None:
    """Test creating an event."""
    # Mock the user lookup
    mock_user = User(username="root")
    mock_query = MagicMock()
    mock_query.filter.return_value.first.return_value = mock_user
    mock_db.query.return_value = mock_query

    # Mock the refresh to populate the event attributes
    def mock_refresh(event):
        event.event_id = uuid4()
        event.created_at = datetime.now()
        event.updated_at = datetime.now()

    mock_db.refresh = MagicMock(side_effect=mock_refresh)

    response = client.post(
        "/api/v1/events",
        headers={"X-API-Key": "test-key"},
        json={
            "title": "Test Event",
            "start_time": "2025-01-01T10:00:00Z",
            "end_time": "2025-01-01T11:00:00Z",
        },
    )
    assert response.status_code == status.HTTP_201_CREATED
    data = response.json()
    assert data["title"] == "Test Event"


def test_get_events(client: TestClient, mock_db) -> None:
    """Test getting all events."""
    response = client.get("/api/v1/events", headers={"X-API-Key": "test-key"})
    assert response.status_code == status.HTTP_200_OK
    assert isinstance(response.json(), list)


def test_get_event(client: TestClient, mock_db) -> None:
    """Test getting a single event."""
    # Mock the database to return an event
    event_id = uuid4()
    mock_event = CalendarEvent(
        event_id=event_id,
        title="Test Event",
        description="Test Description",
        start_time=datetime.now(),
        end_time=datetime.now(),
        created_at=datetime.now(),
        updated_at=datetime.now(),
    )
    # Mock user lookup first, then event lookup
    mock_user = User(username="root")
    mock_query = MagicMock()
    mock_query.filter.return_value.first.side_effect = [mock_user, mock_event]
    mock_db.query.return_value = mock_query

    response = client.get(
        f"/api/v1/events/{event_id}", headers={"X-API-Key": "test-key"}
    )
    assert response.status_code == status.HTTP_200_OK
    assert response.json()["title"] == "Test Event"


def test_update_event(client: TestClient, mock_db) -> None:
    """Test updating an event."""
    # Mock user and event
    mock_user = User(username="root")
    event_id = uuid4()
    mock_event = CalendarEvent(
        event_id=event_id,
        title="Test Event",
        description="Test Description",
        start_time=datetime.now(),
        end_time=datetime.now(),
        created_at=datetime.now(),
        updated_at=datetime.now(),
    )
    mock_query = MagicMock()
    # The first query is for the user, the second is for the event
    mock_query.filter.return_value.first.side_effect = [mock_user, mock_event]
    mock_db.query.return_value = mock_query

    response = client.put(
        f"/api/v1/events/{event_id}",
        headers={"X-API-Key": "test-key"},
        json={"title": "Updated Event"},
    )
    assert response.status_code == status.HTTP_200_OK
    assert response.json()["title"] == "Updated Event"


def test_delete_event(client: TestClient, mock_db) -> None:
    """Test deleting an event."""
    # Mock user and event
    mock_user = User(username="root")
    event_id = uuid4()
    mock_event = CalendarEvent(
        event_id=event_id,
        title="Test Event",
        description="Test Description",
        start_time=datetime.now(),
        end_time=datetime.now(),
        created_at=datetime.now(),
        updated_at=datetime.now(),
    )
    mock_query = MagicMock()
    # The first query is for the user, the second is for the event
    mock_query.filter.return_value.first.side_effect = [mock_user, mock_event]
    mock_db.query.return_value = mock_query

    response = client.delete(
        f"/api/v1/events/{event_id}",
        headers={"X-API-Key": "test-key"},
    )
    assert response.status_code == status.HTTP_204_NO_CONTENT

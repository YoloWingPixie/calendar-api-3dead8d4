"""Test health check endpoint."""

from datetime import datetime

from fastapi.testclient import TestClient

from src.main import app

client = TestClient(app)


def test_health_check() -> None:
    """Test the health check endpoint returns correct response."""
    response = client.get("/api/v1/health")

    assert response.status_code == 200

    data = response.json()
    assert data["status"] == "ok"
    assert "version" in data
    assert "timestamp" in data

    # Verify timestamp is valid ISO format
    datetime.fromisoformat(data["timestamp"].replace("Z", "+00:00"))

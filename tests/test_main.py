"""Test for main module."""

from fastapi import FastAPI

from src.main import app, create_app


def test_create_app() -> None:
    """Test that create_app returns a FastAPI instance."""
    test_app = create_app()
    assert isinstance(test_app, FastAPI)
    assert test_app.title == "Calendar API"


def test_app_instance() -> None:
    """Test that app is a FastAPI instance."""
    assert isinstance(app, FastAPI)
    assert app.title == "Calendar API"

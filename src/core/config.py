"""Application configuration using pydantic-settings."""

import tomllib
from pathlib import Path

from pydantic_settings import BaseSettings, SettingsConfigDict


def get_project_version() -> str:
    """Get version from pyproject.toml with safe fallback."""
    try:
        pyproject_path = Path(__file__).parent.parent / "pyproject.toml"
        with open(pyproject_path, "rb") as f:
            data = tomllib.load(f)
        return str(data["project"]["version"])
    except (FileNotFoundError, KeyError, tomllib.TOMLDecodeError):
        return "0.0.0"


class Settings(BaseSettings):
    """Application settings."""

    model_config = SettingsConfigDict(
        env_file=".env",
        env_ignore_empty=True,
        extra="ignore",
    )

    # Application settings
    app_name: str = "Calendar API"
    version: str = get_project_version()

    # Environment settings
    debug: bool = False
    environment: str = "development"

    # Server settings
    host: str = "0.0.0.0"
    port: int = 8000

    # Database settings
    database_url: str = "postgresql://user:pass@localhost/calendar_db"

    # Security settings
    api_key_header: str = "X-API-Key"


settings = Settings()

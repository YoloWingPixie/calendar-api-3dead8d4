"""Application configuration using pydantic-settings."""

import json
import os
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

    def __init__(self, **kwargs) -> None:
        """Initialize settings with Doppler JSON parsing support."""
        # Check if running in ECS with Doppler secrets
        doppler_json = os.getenv("DOPPLER_SECRETS_JSON")
        if doppler_json:
            try:
                # Parse the JSON blob from Doppler
                doppler_secrets = json.loads(doppler_json)
                # Merge Doppler secrets with any existing environment variables
                # Environment variables take precedence for local development
                for key, value in doppler_secrets.items():
                    if key not in os.environ:
                        os.environ[key] = str(value)
            except json.JSONDecodeError:
                pass  # Fail silently and fall back to normal env vars

        super().__init__(**kwargs)

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
    database_pool_disabled: bool = False  # Disable pooling for serverless
    database_echo: bool = False  # Log SQL statements (only for debugging)

    # Security settings
    api_key_header: str = "X-API-Key"


settings = Settings()

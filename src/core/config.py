"""
Core application configuration settings.

This module defines the configuration for the application, loading values from
environment variables and Doppler secrets. It uses Pydantic's settings management
to provide a typed, validated configuration object.

The order of precedence for loading settings is:
1.  Initialization arguments (for testing)
2.  Doppler secrets (from DOPPLER_SECRETS_JSON env var)
3.  Environment variables
4.  Values from a .env file
5.  Default values defined in the Settings class
"""

import functools
import json
import logging
import os

from pydantic import Field
from pydantic_settings import BaseSettings


def get_database_url() -> str:
    """
    Gets the database URL, prioritizing the Doppler JSON secret.
    This logic directly mirrors the working implementation in alembic/env.py
    to ensure consistency and reliability.
    """
    doppler_json_str = os.getenv("DOPPLER_SECRETS_JSON")
    if doppler_json_str:
        logging.info("Found DOPPLER_SECRETS_JSON, attempting to parse...")
        try:
            secrets = json.loads(doppler_json_str)
            for key in ["DATABASE_URL", "TF_VAR_database_url", "TF_VAR_DATABASE_URL"]:
                if key in secrets and secrets[key]:
                    logging.info(f"Using database URL from Doppler secret key: {key}")
                    return str(secrets[key])
            raise ValueError(
                "DATABASE_URL not found in DOPPLER_SECRETS_JSON. "
                f"Available keys: {list(secrets.keys())}"
            )
        except json.JSONDecodeError as e:
            raise ValueError(f"Failed to parse DOPPLER_SECRETS_JSON: {e}") from e

    database_url = os.getenv("DATABASE_URL")
    if database_url:
        logging.info("Using database URL from DATABASE_URL environment variable.")
        return database_url

    raise ValueError(
        "Database URL not configured. Set DATABASE_URL or DOPPLER_SECRETS_JSON."
    )


class AppSettings(BaseSettings):
    """
    Application settings.
    """

    # Core settings
    app_name: str = "Calendar API"
    version: str = "0.0.0"
    debug: bool = False
    environment: str = "dev"
    log_level: str = "INFO"

    # Server settings
    host: str = "0.0.0.0"
    port: int = 8000

    # Database settings
    database_url: str = Field(default_factory=get_database_url)
    database_echo: bool = Field(default=False)
    database_pool_disabled: bool = Field(default=False)

    # API settings
    api_v1_str: str = "/api/v1"
    api_key_header: str = "X-API-Key"
    bootstrap_admin_key: str | None = Field(default=None)


@functools.lru_cache
def get_settings() -> AppSettings:
    """
    Get the application settings. This is cached to ensure settings are loaded
    only once.
    """
    return AppSettings()


# Create a single, importable settings instance
settings = get_settings()

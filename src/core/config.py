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
from pydantic_settings import BaseSettings, SettingsConfigDict


class AppSettings(BaseSettings):
    """
    Application settings.

    Settings are loaded from the following sources, in order of precedence:
    1. Doppler secrets (if DOPPLER_SECRETS_JSON env var is set)
    2. Environment variables
    3. Values from a .env file (if present)
    4. Default values defined in this class
    """

    model_config = SettingsConfigDict(
        env_file=".env", env_file_encoding="utf-8", case_sensitive=False
    )

    # Core settings
    app_name: str = "Calendar API"
    version: str = "0.0.0"
    debug: bool = False
    environment: str = "dev"
    log_level: str = "INFO"
    testing: bool = False

    # Server settings
    host: str = "0.0.0.0"
    port: int = 8000

    # Database settings
    database_url: str
    database_echo: bool = Field(default=False)
    database_pool_disabled: bool = Field(default=False)

    # API settings
    api_v1_str: str = "/api/v1"
    api_key_header: str = "X-API-Key"
    bootstrap_admin_key: str | None = Field(default=None)


@functools.lru_cache
def get_settings() -> AppSettings:
    """
    Get the application settings.

    This function is cached to ensure that settings are loaded only once. It
    loads settings from Doppler, if available, and then initializes the main
    AppSettings object.
    """
    init_kwargs = {}
    doppler_json_str = os.getenv("DOPPLER_SECRETS_JSON")

    if doppler_json_str:
        logging.info("Found DOPPLER_SECRETS_JSON, loading secrets...")
        try:
            secrets = json.loads(doppler_json_str)
            # In AWS SM, Doppler secrets are prefixed 'TF_VAR_'. We strip this
            # prefix and convert to lowercase to match Pydantic model fields.
            init_kwargs = {
                key.lower().replace("tf_var_", ""): value
                for key, value in secrets.items()
            }
        except json.JSONDecodeError as e:
            # Raise an error to prevent starting with a bad config.
            raise ValueError(f"Failed to parse DOPPLER_SECRETS_JSON: {e}") from e

    return AppSettings(**init_kwargs)


# Create a single, importable settings instance
settings = get_settings()

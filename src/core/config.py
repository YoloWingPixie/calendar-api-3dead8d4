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
from typing import Any

from pydantic import AliasChoices, Field, PostgresDsn
from pydantic.fields import FieldInfo
from pydantic_settings import (
    BaseSettings,
    PydanticBaseSettingsSource,
    SettingsConfigDict,
)


class DopplerSettingsSource(PydanticBaseSettingsSource):
    """
    A Pydantic settings source that loads configuration from a JSON blob
    stored in the `DOPPLER_SECRETS_JSON` environment variable.

    This source is case-insensitive for secret keys.
    """

    _secrets: dict[str, Any] | None = None

    def _load_secrets(self) -> None:
        if self._secrets is not None:
            return

        doppler_json_str = os.getenv("DOPPLER_SECRETS_JSON")
        if doppler_json_str:
            try:
                self._secrets = {
                    str(k).upper(): v for k, v in json.loads(doppler_json_str).items()
                }
            except json.JSONDecodeError:
                logging.warning("Could not decode DOPPLER_SECRETS_JSON.")
                self._secrets = {}
        else:
            self._secrets = {}

    def get_field_value(
        self, field: FieldInfo, field_name: str
    ) -> tuple[Any, str, bool]:
        """
        Get a value for a field from the loaded Doppler secrets.
        """
        self._load_secrets()
        if not self._secrets:
            return None, field_name, False

        aliases = [field_name]
        if field.validation_alias:
            if isinstance(field.validation_alias, str):
                aliases.append(field.validation_alias)
            elif hasattr(field.validation_alias, "choices"):
                aliases.extend([str(alias) for alias in field.validation_alias.choices])

        field_value = None
        for alias in aliases:
            value = self._secrets.get(str(alias).upper())
            if value is not None:
                field_value = value
                break

        is_complex = self.field_is_complex(field)
        return field_value, field_name, is_complex

    def prepare_field_value(
        self, field_name: str, field: FieldInfo, value: Any, value_is_complex: bool
    ) -> Any:
        """
        No value preparation needed for Doppler, as values are already parsed.
        """
        return value

    def __call__(self) -> dict[str, Any]:
        """
        Load secrets from the Doppler JSON and prepare them for Pydantic.
        """
        d: dict[str, Any] = {}

        for field_name, field in self.settings_cls.model_fields.items():
            field_value, field_key, value_is_complex = self.get_field_value(
                field, field_name
            )
            if field_value is not None:
                d[field_key] = self.prepare_field_value(
                    field_name, field, field_value, value_is_complex
                )

        return d


class AppSettings(BaseSettings):
    """
    Application settings.
    """

    # Core settings
    app_name: str = Field(default="Calendar API", validation_alias="APP_NAME")
    version: str = Field(default="0.4.0", validation_alias="APP_VERSION")
    debug: bool = Field(default=False, validation_alias="DEBUG")
    environment: str = Field(default="dev", validation_alias="ENVIRONMENT")
    log_level: str = Field(default="INFO", validation_alias="LOG_LEVEL")

    # Server settings
    host: str = Field(default="0.0.0.0", validation_alias="HOST")
    port: int = Field(default=8000, validation_alias="PORT")

    # Database settings, with aliases to match Doppler/Terraform variables
    database_url: PostgresDsn = Field(
        default=PostgresDsn("postgresql://user:password@localhost:5432/appdb"),
        validation_alias=AliasChoices(
            "DATABASE_URL", "TF_VAR_database_url", "TF_VAR_DATABASE_URL"
        ),
    )
    database_echo: bool = Field(default=False, validation_alias="DATABASE_ECHO")
    database_pool_disabled: bool = Field(
        default=False, validation_alias="DATABASE_POOL_DISABLED"
    )

    # API settings
    api_v1_str: str = "/api/v1"
    api_key_header: str = "X-API-Key"

    model_config = SettingsConfigDict(
        env_file=".env",
        env_file_encoding="utf-8",
        extra="ignore",
        case_sensitive=False,
    )

    @classmethod
    def settings_customise_sources(
        cls,
        settings_cls: type[BaseSettings],
        init_settings: PydanticBaseSettingsSource,
        env_settings: PydanticBaseSettingsSource,
        dotenv_settings: PydanticBaseSettingsSource,
        file_secret_settings: PydanticBaseSettingsSource,
    ) -> tuple[PydanticBaseSettingsSource, ...]:
        """
        Customise the settings sources to include the Doppler source with
        high priority.
        """
        return (
            init_settings,
            DopplerSettingsSource(settings_cls),
            env_settings,
            dotenv_settings,
            file_secret_settings,
        )


@functools.lru_cache
def get_settings() -> AppSettings:
    """
    Get the application settings.

    This function is cached to ensure that the settings are loaded only once.

    Returns:
        The application settings instance.
    """
    return AppSettings()


# Create a single, importable settings instance
settings = get_settings()

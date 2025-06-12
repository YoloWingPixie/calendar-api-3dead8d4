import json
import os
from logging.config import fileConfig

from sqlalchemy import engine_from_config, pool

from alembic import context
from src.models import CalendarEvent, User  # noqa: F401
from src.models.base import Base

# this is the Alembic Config object, which provides
# access to the values within the .ini file in use.
config = context.config

# Interpret the config file for Python logging.
# This line sets up loggers basically.
if config.config_file_name is not None:
    fileConfig(config.config_file_name)

# add your model's MetaData object here
# for 'autogenerate' support
# from myapp import mymodel
# target_metadata = mymodel.Base.metadata
target_metadata = Base.metadata

# other values from the config, defined by the needs of env.py,
# can be acquired:
# my_important_option = config.get_main_option("my_important_option")
# ... etc.


def get_database_url() -> str:
    """
    Gets the database URL, prioritizing the Doppler JSON secret.

    In the migration ECS task, Doppler injects secrets as a single JSON blob.
    This function parses that secret to find the DATABASE_URL. For local
    development, it falls back to the DATABASE_URL environment variable.
    """
    doppler_json_str = os.getenv("DOPPLER_SECRETS_JSON")
    if doppler_json_str:
        print("Found DOPPLER_SECRETS_JSON, attempting to parse...")
        try:
            secrets = json.loads(doppler_json_str)
            # Doppler can use different keys, check for common ones
            for key in ["DATABASE_URL", "TF_VAR_database_url", "TF_VAR_DATABASE_URL"]:
                if key in secrets and secrets[key]:
                    print(f"Using database URL from Doppler secret key: {key}")
                    return str(secrets[key])

            # If no key was found, raise an error with context
            raise ValueError(
                "DATABASE_URL not found in DOPPLER_SECRETS_JSON. "
                f"Available keys: {list(secrets.keys())}"
            )
        except json.JSONDecodeError as e:
            raise ValueError(f"Failed to parse DOPPLER_SECRETS_JSON: {e}") from e

    # Fallback for local development or other environments
    database_url = os.getenv("DATABASE_URL")
    if database_url:
        print("Using database URL from DATABASE_URL environment variable.")
        return database_url

    # Final fallback to alembic.ini, if configured there.
    ini_url = config.get_main_option("sqlalchemy.url")
    if ini_url:
        print("Using database URL from alembic.ini.")
        return ini_url

    raise ValueError(
        "Database URL not configured. Set DATABASE_URL or DOPPLER_SECRETS_JSON, "
        "or configure sqlalchemy.url in alembic.ini."
    )


# Override the sqlalchemy.url with the correct database URL.
# This is the central point for configuring the database connection for Alembic.
db_url = get_database_url()
if db_url:
    config.set_main_option("sqlalchemy.url", db_url)


def run_migrations_offline() -> None:
    """Run migrations in 'offline' mode.

    This configures the context with just a URL
    and not an Engine, though an Engine is acceptable
    here as well.  By skipping the Engine creation
    we don't even need a DBAPI to be available.

    Calls to context.execute() here emit the given string to the
    script output.

    """
    url = config.get_main_option("sqlalchemy.url")
    context.configure(
        url=url,
        target_metadata=target_metadata,
        literal_binds=True,
        dialect_opts={"paramstyle": "named"},
    )

    with context.begin_transaction():
        context.run_migrations()


def run_migrations_online() -> None:
    """Run migrations in 'online' mode.

    In this scenario we need to create an Engine
    and associate a connection with the context.

    """
    connectable = engine_from_config(
        config.get_section(config.config_ini_section, {}),
        prefix="sqlalchemy.",
        poolclass=pool.NullPool,
    )

    with connectable.connect() as connection:
        context.configure(connection=connection, target_metadata=target_metadata)

        with context.begin_transaction():
            context.run_migrations()


if context.is_offline_mode():
    run_migrations_offline()
else:
    run_migrations_online()

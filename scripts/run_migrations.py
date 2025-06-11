#!/usr/bin/env python
"""Run database migrations with Doppler secrets support."""

import json
import os
import subprocess
import sys


def main():
    """Parse Doppler secrets and run migrations."""
    # Check if running in ECS with Doppler secrets
    doppler_json = os.getenv("DOPPLER_SECRETS_JSON")
    if doppler_json:
        try:
            # Parse the JSON blob from Doppler
            doppler_secrets = json.loads(doppler_json)

            # Look for DATABASE_URL in various possible keys
            database_url = None
            for key in ["DATABASE_URL", "TF_VAR_database_url", "TF_VAR_DATABASE_URL"]:
                if key in doppler_secrets:
                    database_url = doppler_secrets[key]
                    break

            if database_url:
                os.environ["DATABASE_URL"] = database_url
                print("DATABASE_URL set from Doppler secrets")
            else:
                print("WARNING: DATABASE_URL not found in Doppler secrets")
                print("Available keys:", list(doppler_secrets.keys()))

        except json.JSONDecodeError as e:
            print(f"ERROR: Failed to parse DOPPLER_SECRETS_JSON: {e}")
            sys.exit(1)

    # Run alembic upgrade
    print("Running database migrations...")
    result = subprocess.run(["alembic", "upgrade", "head"])
    sys.exit(result.returncode)


if __name__ == "__main__":
    main()

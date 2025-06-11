#!/usr/bin/env python3
"""Check that PR version is higher than base branch version."""

import os
import sys

from packaging.version import InvalidVersion, Version


def main():
    pr_version = os.environ.get("PR_VERSION", "").strip()
    base_version = os.environ.get("BASE_VERSION", "").strip()

    if not pr_version:
        print("❌ PR_VERSION environment variable not set")
        sys.exit(1)

    if not base_version:
        print("❌ BASE_VERSION environment variable not set")
        sys.exit(1)

    print(f"PR version: {pr_version}")
    print(f"Base version: {base_version}")

    try:
        pr_ver = Version(pr_version)
        base_ver = Version(base_version)
    except InvalidVersion as e:
        print(f"❌ Invalid version format: {e}")
        sys.exit(1)

    if pr_ver > base_ver:
        print(f"✅ Version bump validated: {base_version} → {pr_version}")
        sys.exit(0)
    else:
        print(f"❌ Version must be higher than base branch version ({base_version})")
        sys.exit(1)


if __name__ == "__main__":
    main()

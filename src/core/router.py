"""Dynamic router loader for API endpoints."""

import importlib
import pkgutil
from pathlib import Path

from fastapi import APIRouter, FastAPI


def load_routers_from_module(
    module_name: str, prefix: str
) -> list[tuple[APIRouter, str]]:
    """Recursively load all routers from a module.

    Args:
        module_name: The module to search for routers (e.g., 'src.api.v1')
        prefix: The URL prefix for the routes (e.g., '/api/v1')

    Returns:
        List of tuples containing (router, prefix) found in the module
    """
    routers: list[tuple[APIRouter, str]] = []

    try:
        # Import the module
        module = importlib.import_module(module_name)
        if module.__file__ is None:
            return routers
        module_path = Path(module.__file__).parent

        # Iterate through all modules in the package
        for _, name, is_pkg in pkgutil.iter_modules([str(module_path)]):
            full_module_name = f"{module_name}.{name}"

            if is_pkg:
                # Recursively load routers from subpackages
                sub_prefix = f"{prefix}/{name}"
                routers.extend(load_routers_from_module(full_module_name, sub_prefix))
            else:
                # Try to import the module and check for a router
                try:
                    sub_module = importlib.import_module(full_module_name)
                    if hasattr(sub_module, "router"):
                        router = sub_module.router
                        if isinstance(router, APIRouter):
                            # Store the router with its prefix
                            routers.append((router, prefix))
                except ImportError:
                    # Skip modules that can't be imported
                    pass

    except (ImportError, AttributeError):
        # Skip if module can't be imported
        pass

    return routers


def register_routers(app: FastAPI) -> None:
    """Register all API routers with the FastAPI app.

    Args:
        app: The FastAPI application instance
    """
    # Load all v1 API routers
    v1_routers = load_routers_from_module("src.api.v1", "/api/v1")

    # Register each router with the app
    for router, prefix in v1_routers:
        app.include_router(router, prefix=prefix)

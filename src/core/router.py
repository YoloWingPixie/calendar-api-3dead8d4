"""Dynamic router loader for API endpoints."""

import importlib
import pkgutil

from fastapi import APIRouter, FastAPI

from src.core.config import settings


def load_routers_from_module(module_name: str) -> list[APIRouter]:
    """Recursively load all routers from a module.

    Args:
        module_name: The module to search for routers (e.g., 'src.api.v1')

    Returns:
        List of routers found in the module
    """
    routers: list[APIRouter] = []
    module = importlib.import_module(module_name)
    if not hasattr(module, "__path__"):
        return routers

    for _, name, is_pkg in pkgutil.iter_modules(module.__path__):
        full_module_name = f"{module_name}.{name}"
        if is_pkg:
            routers.extend(load_routers_from_module(full_module_name))
        else:
            sub_module = importlib.import_module(full_module_name)
            if hasattr(sub_module, "router") and isinstance(
                sub_module.router, APIRouter
            ):
                routers.append(sub_module.router)

    return routers


def register_routers(app: FastAPI) -> None:
    """Register all API routers with the FastAPI app.

    Args:
        app: The FastAPI application instance
    """
    # Load all v1 API routers
    v1_routers = load_routers_from_module("src.api.v1")

    # Register each router with the app
    for router in v1_routers:
        app.include_router(router, prefix=settings.api_v1_str)

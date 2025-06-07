"""Test for main module."""

import pytest

from src.main import main


def test_main(capsys: pytest.CaptureFixture[str]) -> None:
    """Test main function prints expected output."""
    main()
    captured = capsys.readouterr()
    assert "Calendar API - Starting..." in captured.out

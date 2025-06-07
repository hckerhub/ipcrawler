#!/usr/bin/env python3
"""
Main Application Module
Contains the IPCrawlerApp class
"""

from textual.app import App
from .config import APP_TITLE, APP_SUBTITLE, APP_CSS
from .screens.welcome import WelcomeScreen


class IPCrawlerApp(App):
    """Main application class."""
    
    CSS = APP_CSS
    TITLE = APP_TITLE
    SUB_TITLE = APP_SUBTITLE

    def on_mount(self) -> None:
        """Initialize the application."""
        self.push_screen(WelcomeScreen())

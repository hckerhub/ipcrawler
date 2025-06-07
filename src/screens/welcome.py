#!/usr/bin/env python3
"""
Welcome Screen Module
Contains the WelcomeScreen class for IPCrawler
"""

from textual.app import ComposeResult
from textual.containers import Container
from textual.widgets import Header, Footer, Static
from textual.screen import Screen
from textual.binding import Binding
from rich.text import Text
from rich.align import Align

from ..config import LOGO, APP_TITLE, DEVELOPER


class WelcomeScreen(Screen):
    """Welcome screen with hacker theme."""
    
    CSS = """
    /* Hacker Theme for Welcome Screen */
    Screen {
        background: #000000;
    }
    
    Header {
        background: #0d1117;
        color: #00ff00;
        text-style: bold;
    }
    
    Footer {
        background: #0d1117;
        color: #00ff00;
    }
    
    Footer > .footer--highlight {
        background: #ff0000;
        color: #ffffff;
    }
    
    Footer > .footer--highlight-key {
        background: #ff0000;
        color: #ffffff;
        text-style: bold;
    }
    
    Footer > .footer--key {
        background: #333333;
        color: #00ff00;
        text-style: bold;
    }
    
    #welcome-container {
        background: #000000;
        align: center middle;
        width: 100%;
        height: 100%;
        padding: 2;
    }
    
    #logo {
        background: #000000;
        color: #00ff00;
        text-align: center;
        text-style: bold;
        margin: 1 0;
    }
    
    #welcome-text {
        background: #000000;
        color: #ffff00;
        text-align: center;
        text-style: bold;
        margin: 1 0;
    }
    
    #developer {
        background: #000000;
        color: #ff6600;
        text-align: center;
        margin: 1 0;
    }
    
    #subtitle {
        background: #000000;
        color: #00ff00;
        text-align: center;
        text-style: bold;
        margin: 1 0;
    }
    
    #instructions {
        background: #000000;
        color: #888888;
        text-align: center;
        margin: 2 0;
    }
    
    .matrix-glow {
        background: #000000;
        color: #00ff00;
        text-style: bold;
    }
    """
    
    BINDINGS = [
        Binding("enter", "next_screen", "Continue"),
        Binding("q", "quit", "Quit"),
        Binding("ctrl+c", "quit", "Quit"),
    ]

    def compose(self) -> ComposeResult:
        yield Header()
        yield Container(
            Static(
                Align.center(
                    Text(LOGO, style="bold green")
                ),
                id="logo"
            ),
            Static(
                Align.center(
                    Text("Welcome to IPCrawler – Intelligent Recon Flow Engine", style="bold yellow")
                ),
                id="welcome-text"
            ),
            Static(
                Align.center(
                    Text(f"Developer: {DEVELOPER}", style="bold #ff6600")
                ),
                id="developer"
            ),
            Static(
                Align.center(
                    Text("🔍 Advanced Reconnaissance Automation Tool 🔍", style="bold green")
                ),
                id="subtitle"
            ),
            Static(
                Align.center(
                    Text("⚡ Press [Enter] to continue or [Q] to quit ⚡", style="#888888")
                ),
                id="instructions"
            ),
            id="welcome-container"
        )
        yield Footer()

    def action_next_screen(self) -> None:
        """Navigate to tool selection screen."""
        from .tool_selection import ToolSelectionScreen
        self.app.push_screen(ToolSelectionScreen())

    def action_quit(self) -> None:
        """Quit the application."""
        self.app.exit()

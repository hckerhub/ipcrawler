#!/usr/bin/env python3
"""
Summary Screen Module
Contains the SummaryScreen class for IPCrawler
"""

from textual.app import ComposeResult
from textual.containers import Container
from textual.widgets import Header, Footer, Static, Rule
from textual.screen import Screen
from textual.binding import Binding
from rich.text import Text
from rich.align import Align

from ..config import LOGO


class SummaryScreen(Screen):
    """Summary screen with hacker theme."""
    
    CSS = """
    /* Hacker Theme for Summary */
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
    
    #summary-container {
        background: #000000;
        margin: 1 2;
        padding: 1;
    }
    
    #logo-small {
        background: #000000;
        color: #00ff00;
        text-align: center;
        text-style: bold;
        margin: 1 0;
        height: auto;
    }
    
    #title {
        background: #000000;
        color: #ffff00;
        text-align: center;
        text-style: bold;
        margin: 1;
    }
    
    .summary-content {
        background: #000000;
        margin: 1 0;
    }
    
    .summary-label {
        background: #000000;
        color: #ffff00;
        text-style: bold;
        margin: 1 0;
    }
    
    .summary-value {
        background: #000000;
        color: #00ff00;
        text-style: bold;
        margin: 0 0 0 2;
    }
    
    .target-value {
        background: #000000;
        color: #ff6600;
        text-style: bold;
        margin: 0 0 0 2;
    }
    
    .summary-tool {
        background: #000000;
        color: #00ff00;
        margin: 0 0 0 4;
    }
    
    .note {
        background: #000000;
        color: #ffff00;
        text-style: italic;
        text-align: center;
        margin: 2 0;
    }
    
    .ready-message {
        background: #000000;
        color: #00ff00;
        text-style: bold;
        text-align: center;
        margin: 2 0;
    }
    
    .action-section {
        background: #000000;
        text-align: center;
        margin: 1 0;
    }
    
    .instructions {
        background: #000000;
        color: #888888;
        text-align: center;
        margin: 1 0;
    }
    
    Rule {
        color: #00ff00;
        background: #000000;
    }
    
    .matrix-border {
        background: #000000;
        border: solid #00ff00;
        padding: 1;
        margin: 1 0;
    }
    """
    
    BINDINGS = [
        Binding("enter", "start_scan", "Start Scan"),
        Binding("escape", "back", "Back"),
        Binding("ctrl+c", "quit", "Quit"),
    ]

    def __init__(self, selected_tools, target):
        super().__init__()
        self.selected_tools = selected_tools
        self.target = target

    def compose(self) -> ComposeResult:
        yield Header()
        yield Container(
            Static(
                Align.center(
                    Text(LOGO, style="bold green")
                ),
                id="logo-small"
            ),
            Static(
                Align.center(
                    Text("🎯 Reconnaissance Summary", style="bold yellow")
                ),
                id="title"
            ),
            Rule(),
            Container(
                Static("📋 Configuration Summary:", classes="summary-label"),
                Static("", classes="summary-value"),  # Spacer
                
                Static("🎯 Target:", classes="summary-label"),
                Static(f"   {self.target}", classes="target-value"),
                Static("", classes="summary-value"),  # Spacer
                
                Static("🛠️  Selected Tools:", classes="summary-label"),
                *[Static(f"   • {tool}", classes="summary-tool") for tool in self.selected_tools],
                Static("", classes="summary-value"),  # Spacer
                
                classes="summary-content matrix-border"
            ),
            Rule(),
            Container(
                Static(
                    Align.center(
                        Text("⚠️  Note: This is a simulation. No actual scanning will be performed.", style="bold yellow")
                    ),
                    classes="note"
                ),
                Static(
                    Align.center(
                        Text("🚀 Ready to initiate reconnaissance sequence!", style="bold green")
                    ),
                    classes="ready-message"
                ),
                classes="action-section"
            ),
            Static("", classes="summary-value"),  # Spacer
            Static("", classes="summary-value"),  # Spacer
            Static("", classes="summary-value"),  # Spacer
            Static("", classes="summary-value"),  # Spacer
            Static("", classes="summary-value"),  # Spacer
            Static("", classes="summary-value"),  # Spacer
            Static("", classes="summary-value"),  # Spacer
            Rule(),
            Static(
                Align.center(
                    Text("⚡ Press ENTER to start scan or ESC to go back ⚡", style="#888888")
                ),
                classes="instructions"
            ),
            id="summary-container"
        )
        yield Footer()

    def action_start_scan(self) -> None:
        """Start the reconnaissance scan (simulation)."""
        self.notify("🔥 Reconnaissance sequence initiated! (Simulation mode)", severity="information")
        # In a real implementation, this would start the actual scanning process
        # For now, we'll just show a success message and exit after a delay
        import asyncio
        asyncio.create_task(self._simulate_scan())

    async def _simulate_scan(self):
        """Simulate a scanning process."""
        import asyncio
        await asyncio.sleep(2)
        self.notify("✅ Scan simulation completed successfully!", severity="information")
        await asyncio.sleep(3)
        self.app.exit()

    def action_back(self) -> None:
        """Go back to target input screen."""
        self.app.pop_screen()

    def action_quit(self) -> None:
        """Quit the application."""
        self.app.exit()

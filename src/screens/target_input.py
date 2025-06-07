#!/usr/bin/env python3
"""
Target Input Screen Module
Contains the TargetInputScreen class for IPCrawler
"""

from textual.app import ComposeResult
from textual.containers import Container, Horizontal, Vertical, Grid
from textual.widgets import Header, Footer, Static, Input, Rule
from textual.screen import Screen
from textual.binding import Binding
from rich.text import Text
from rich.align import Align

from ..config import LOGO


class TargetInputScreen(Screen):
    """Target input screen with hacker theme and responsive layout."""
    
    CSS = """
    /* Hacker Theme for Target Input */
    Screen {
        background: #000000;
        layout: vertical;
    }
    
    Header {
        background: #0d1117;
        color: #00ff00;
        text-style: bold;
        height: 3;
    }
    
    Footer {
        background: #0d1117;
        color: #00ff00;
        height: 3;
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
    
    #main-container {
        background: #000000;
        layout: vertical;
        height: 1fr;
        padding: 1 2;
    }
    
    #logo-container {
        background: #000000;
        height: auto;
        margin: 1 0;
    }
    
    #logo-text {
        background: #000000;
        color: #00ff00;
        text-align: center;
        text-style: bold;
        height: auto;
    }
    
    #title-container {
        background: #000000;
        height: auto;
        margin: 0 0 1 0;
    }
    
    #title-text {
        background: #000000;
        color: #ffff00;
        text-align: center;
        text-style: bold;
        height: auto;
    }
    
    #content-area {
        background: #000000;
        layout: vertical;
        height: 1fr;
    }
    
    #tools-section {
        background: #000000;
        height: auto;
        margin: 1 0;
    }
    
    .tools-title {
        background: #000000;
        color: #ffff00;
        text-style: bold;
        height: 1;
        margin: 0 0 1 0;
    }
    
    #tools-grid {
        background: #000000;
        layout: grid;
        grid-size: 4 2;
        grid-gutter: 1 2;
        height: auto;
        margin: 0 0 1 0;
    }
    
    .tool-item {
        background: #000000;
        color: #00ff00;
        text-style: bold;
        height: 1;
        content-align: left middle;
        min-width: 15;
        margin: 0 1 0 0;
    }
    
    #tools-display {
        background: #000000;
        height: auto;
        min-height: 2;
    }
    
    .tools-row {
        background: #000000;
        height: auto;
        margin: 0 0 1 0;
    }
    
    #target-section {
        background: #000000;
        height: auto;
        margin: 1 0;
    }
    
    .target-title {
        background: #000000;
        color: #ffff00;
        text-style: bold;
        height: 1;
        margin: 0 0 1 0;
    }
    
    .input-label {
        background: #000000;
        color: #888888;
        height: 1;
        margin: 0 0 1 0;
    }
    
    #target-input {
        background: #111111;
        color: #00ff00;
        border: solid #00ff00;
        height: 3;
        margin: 0 0 1 0;
    }
    
    #target-input:focus {
        background: #222222;
        border: solid #ffff00;
    }
    
    #examples-section {
        background: #000000;
        height: auto;
        margin: 1 0;
    }
    
    .examples-title {
        background: #000000;
        color: #ffff00;
        text-style: bold;
        height: 1;
        margin: 0 0 1 0;
    }
    
    .example-item {
        background: #000000;
        color: #888888;
        height: 1;
        margin: 0 0 0 2;
    }
    
    #instructions-section {
        background: #000000;
        height: auto;
        margin: 1 0;
    }
    
    .instructions-text {
        background: #000000;
        color: #888888;
        text-align: center;
        height: 1;
    }
    
    Rule {
        color: #00ff00;
        background: #000000;
        height: 1;
    }
    """
    
    BINDINGS = [
        Binding("enter", "next_screen", "Continue"),
        Binding("escape", "back", "Back"),
        Binding("ctrl+c", "quit", "Quit"),
    ]

    def __init__(self, selected_tools):
        super().__init__()
        self.selected_tools = selected_tools
        self.target = ""

    def compose(self) -> ComposeResult:
        yield Header()
        
        yield Container(
            # Logo Section
            Container(
                Static(
                    Align.center(Text(LOGO, style="bold green")),
                    id="logo-text"
                ),
                id="logo-container"
            ),
            
            # Title Section  
            Container(
                Static(
                    Align.center(Text("🎯 Target Configuration", style="bold yellow")),
                    id="title-text"
                ),
                id="title-container"
            ),
            
            Rule(),
            
            # Content Area
            Container(
                # Target Input Section
                Container(
                    Static("Target Specification:", classes="target-title"),
                    Static("Enter IP address or domain name:", classes="input-label"),
                    Input(
                        placeholder="e.g., 192.168.1.1 or example.com",
                        id="target-input"
                    ),
                    id="target-section"
                ),
                
                Rule(),
                
                # Examples Section
                Container(
                    Static("📋 Examples:", classes="examples-title"),
                    Static("• Single IP: 192.168.1.1", classes="example-item"),
                    Static("• Domain: example.com", classes="example-item"),
                    Static("• Subdomain: api.example.com", classes="example-item"),
                    Static("• CIDR: 192.168.1.0/24", classes="example-item"),
                    id="examples-section"
                ),
                
                Rule(),
                
                # Instructions Section
                Container(
                    Static(
                        Align.center(Text("⚡ Press ENTER to continue or ESC to go back ⚡", style="#888888")),
                        classes="instructions-text"
                    ),
                    id="instructions-section"
                ),
                
                id="content-area"
            ),
            
            id="main-container"
        )
        
        yield Footer()

    def _create_tools_display(self) -> Container:
        """Create a responsive layout for selected tools."""
        if not self.selected_tools:
            return Container(
                Static("No tools selected", classes="tool-item"),
                id="tools-display"
            )
        
        # Create tools in rows with multiple columns for space efficiency
        tools_containers = []
        
        # Group tools in rows of 4 for better display
        for i in range(0, len(self.selected_tools), 4):
            row_tools = self.selected_tools[i:i+4]
            row_widgets = []
            
            for tool in row_tools:
                row_widgets.append(
                    Static(f"✅ {tool}", classes="tool-item")
                )
            
            # Create horizontal container for this row
            tools_containers.append(
                Horizontal(*row_widgets, classes="tools-row")
            )
        
        return Container(
            *tools_containers,
            id="tools-display"
        )

    def on_mount(self) -> None:
        """Focus the input field when the screen loads."""
        self.query_one("#target-input", Input).focus()

    def on_input_changed(self, event: Input.Changed) -> None:
        """Handle input field changes."""
        if event.input.id == "target-input":
            self.target = event.value.strip()

    def action_next_screen(self) -> None:
        """Navigate to summary screen."""
        if not self.target:
            self.notify("⚠️  Please enter a target!", severity="warning")
            return
        
        from .summary import SummaryScreen
        self.app.push_screen(SummaryScreen(self.selected_tools, self.target))

    def action_back(self) -> None:
        """Go back to tool selection screen."""
        self.app.pop_screen()

    def action_quit(self) -> None:
        """Quit the application."""
        self.app.exit()

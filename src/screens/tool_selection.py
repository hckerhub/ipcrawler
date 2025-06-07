#!/usr/bin/env python3
"""
Tool Selection Screen Module
Contains the ToolSelectionScreen class for IPCrawler
"""

from textual.app import ComposeResult
from textual.containers import Container
from textual.widgets import Header, Footer, Static, Rule
from textual.screen import Screen
from textual.binding import Binding
from rich.text import Text
from rich.align import Align

from ..config import LOGO, RECON_TOOLS


class ToolSelectionScreen(Screen):
    """Tool selection screen with checkboxes - Hacker Theme."""
    
    CSS = """
    /* Hacker Theme for Tool Selection */
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
    
    #logo-small {
        background: #000000;
        color: #00ff00;
        text-align: center;
        text-style: bold;
        margin: 1 0;
    }
    
    #title {
        background: #000000;
        color: #ffff00;
        text-align: center;
        text-style: bold;
        margin: 1;
    }
    
    #selection-container {
        background: #000000;
        padding: 1;
    }
    
    #tools-container {
        background: #111111;
        border: solid #00ff00;
        padding: 1;
        margin: 1;
        min-height: 15;
    }
    
    .tool-item {
        background: #111111;
        color: #00ff00;
        margin: 0;
        padding: 0 1;
        height: 1;
    }
    
    .current-tool {
        background: #222222;
        color: #ffff00;
        text-style: bold;
        border-left: solid #ff0000;
    }
    
    .normal-tool {
        background: #111111;
        color: #00ff00;
    }
    
    .selected-tool {
        background: #001100;
        color: #00ff00;
        text-style: bold;
    }
    
    .controls {
        background: #000000;
        margin: 1;
    }
    
    .controls-title {
        background: #000000;
        color: #ffff00;
        text-style: bold;
        margin: 0 0 1 0;
    }
    
    .controls-text {
        background: #000000;
        color: #888888;
        margin: 0 0 0 2;
    }
    
    Rule {
        color: #00ff00;
        background: #000000;
    }
    """
    
    BINDINGS = [
        Binding("enter", "next_screen", "Continue"),
        Binding("escape", "back", "Back"),
        Binding("a", "select_all", "Select All"),
        Binding("c", "clear_all", "Clear All"),
        Binding("space", "toggle_current", "Toggle"),
        Binding("q", "quit", "Quit"),
        Binding("ctrl+c", "quit", "Quit"),
    ]

    def __init__(self):
        super().__init__()
        self.selected_tools = set()
        self.current_selection = 0

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
                    Text("🛠️  Select Reconnaissance Tools", style="bold yellow")
                ),
                id="title"
            ),
            Rule(),
            Container(
                *[
                    Static(
                        f"{'► ' if i == 0 else '  '}[ ] {tool_name:<12} | {description}",
                        id=f"tool_{i}",
                        classes="tool-item"
                    )
                    for i, (tool_name, description) in enumerate(RECON_TOOLS)
                ],
                id="tools-container"
            ),
            Rule(),
            Container(
                Static("📋 Controls:", classes="controls-title"),
                Static("• Use ↑↓ arrow keys to navigate tools", classes="controls-text"),
                Static("• Press SPACE to select/deselect current tool", classes="controls-text"),
                Static("• Press A to select all, C to clear all", classes="controls-text"),
                Static("• Press ENTER to continue, ESC to go back", classes="controls-text"),
                Static("• Press Ctrl+C to quit", classes="controls-text"),
                classes="controls"
            ),
            id="selection-container"
        )
        yield Footer()

    def on_mount(self) -> None:
        """Initialize the first tool as selected."""
        self.update_display()

    def update_display(self) -> None:
        """Update the visual display of tools with proper checkboxes."""
        for i, (tool_name, description) in enumerate(RECON_TOOLS):
            tool_widget = self.query_one(f"#tool_{i}", Static)
            is_selected = tool_name in self.selected_tools
            is_current = i == self.current_selection
            
            # Create checkbox display
            if is_selected:
                checkbox = "[X]"
                text_color = "bold green"
            else:
                checkbox = "[ ]"
                text_color = "green"
            
            # Create cursor
            prefix = "► " if is_current else "  "
            
            # Create the display content with proper formatting and colors
            if is_current:
                new_content = Text(f"{prefix}{checkbox} {tool_name:<12} | {description}", style="bold yellow")
            elif is_selected:
                new_content = Text(f"{prefix}{checkbox} {tool_name:<12} | {description}", style="bold green")
            else:
                new_content = Text(f"{prefix}{checkbox} {tool_name:<12} | {description}", style="green")
            
            # Update the widget content
            tool_widget.update(new_content)
            
            # Apply CSS classes for styling
            tool_widget.remove_class("current-tool", "normal-tool", "selected-tool")
            
            if is_current:
                tool_widget.add_class("current-tool")
            elif is_selected:
                tool_widget.add_class("selected-tool")
            else:
                tool_widget.add_class("normal-tool")
            
            # Force refresh
            tool_widget.refresh()

    def on_key(self, event) -> None:
        """Handle key presses."""
        if event.key == "up":
            if self.current_selection == 0:
                self.current_selection = len(RECON_TOOLS) - 1  # Wrap to last item
            else:
                self.current_selection -= 1
            self.update_display()
            event.prevent_default()
        elif event.key == "down":
            if self.current_selection == len(RECON_TOOLS) - 1:
                self.current_selection = 0  # Wrap to first item
            else:
                self.current_selection += 1
            self.update_display()
            event.prevent_default()

    def action_toggle_current(self) -> None:
        """Toggle selection of current tool."""
        if 0 <= self.current_selection < len(RECON_TOOLS):
            tool_name = RECON_TOOLS[self.current_selection][0]
            if tool_name in self.selected_tools:
                self.selected_tools.discard(tool_name)
                self.notify(f"🔥 Deselected {tool_name}", severity="information")
            else:
                self.selected_tools.add(tool_name)
                self.notify(f"⚡ Selected {tool_name}", severity="information")
            self.update_display()

    def action_select_all(self) -> None:
        """Select all tools."""
        for tool_name, _ in RECON_TOOLS:
            self.selected_tools.add(tool_name)
        self.notify("🎯 All tools selected!", severity="information")
        self.update_display()

    def action_clear_all(self) -> None:
        """Clear all selections."""
        self.selected_tools.clear()
        self.notify("💀 All selections cleared!", severity="information")
        self.update_display()

    def action_next_screen(self) -> None:
        """Navigate to target input screen."""
        if not self.selected_tools:
            self.notify("⚠️  Please select at least one tool!", severity="warning")
            return
        from .target_input import TargetInputScreen
        self.app.push_screen(TargetInputScreen(list(self.selected_tools)))

    def action_back(self) -> None:
        """Go back to welcome screen."""
        self.app.pop_screen()

    def action_quit(self) -> None:
        """Quit the application."""
        self.app.exit()

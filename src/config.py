#!/usr/bin/env python3
"""
IPCrawler Configuration
Contains constants, ASCII art, and tool definitions
"""

# ASCII Art Logo
LOGO = """
    ██╗██████╗  ██████╗██████╗  █████╗ ██╗    ██╗██╗     ███████╗██████╗ 
    ██║██╔══██╗██╔════╝██╔══██╗██╔══██╗██║    ██║██║     ██╔════╝██╔══██╗
    ██║██████╔╝██║     ██████╔╝███████║██║ █╗ ██║██║     █████╗  ██████╔╝
    ██║██╔═══╝ ██║     ██╔══██╗██╔══██║██║███╗██║██║     ██╔══╝  ██╔══██╗
    ██║██║     ╚██████╗██║  ██║██║  ██║╚███╔███╔╝███████╗███████╗██║  ██║
    ╚═╝╚═╝      ╚═════╝╚═╝  ╚═╝╚═╝  ╚═╝ ╚══╝╚══╝ ╚══════╝╚══════╝╚═╝  ╚═╝
"""

# Application Information
APP_TITLE = "IPCrawler - Intelligent Recon Flow Engine"
APP_SUBTITLE = "by hckerhub"
DEVELOPER = "hckerhub"

# Available reconnaissance tools
RECON_TOOLS = [
    ("naabu", "Fast port scanner"),
    ("nmap", "Network discovery and security auditing"),
    ("subfinder", "Subdomain discovery tool"),
    ("dnsx", "DNS toolkit with support for multiple DNS queries"),
    ("httpx", "Fast and multi-purpose HTTP toolkit"),
    ("katana", "Next-generation crawling and spidering framework"),
    ("nuclei", "Fast and customizable vulnerability scanner"),
    ("ffuf", "Fast web fuzzer written in Go"),
    ("feroxbuster", "Fast, simple, recursive content discovery tool"),
    ("shuffledns", "Wrapper around massdns for DNS bruteforcing"),
    ("amass", "In-depth attack surface mapping and asset discovery"),
]

# CSS Styles for the application
APP_CSS = """
/* Global Styles */
Screen {
    background: $surface;
}

#welcome-container {
    align: center middle;
    width: 100%;
    height: 100%;
}

#logo {
    margin: 1 0;
    text-align: center;
}

#welcome-text, #developer, #subtitle, #instructions {
    margin: 1 0;
    text-align: center;
}

/* Tool Selection Styles */
#selection-container {
    margin: 1 2;
}

#logo-small {
    margin: 0 0 1 0;
    text-align: center;
    height: auto;
}

#title {
    margin: 1 0;
    text-align: center;
}

#tools-container {
    height: 12;
    border: solid $primary;
    margin: 1 0;
    padding: 1;
}

.tool-item {
    layout: horizontal;
    height: 1;
    margin: 0;
    padding: 0 1;
}

.tool-description {
    margin-left: 2;
    color: $text-muted;
}

.controls {
    margin: 1 0;
}

.controls-title {
    text-style: bold;
    color: $accent;
    margin: 0 0 1 0;
}

.controls-text {
    color: $text-muted;
    margin-left: 2;
    margin: 0 0 0 2;
}

.current-tool {
    background: $primary;
    color: $text;
    text-style: bold;
}

.normal-tool {
    color: $text;
}

/* Target Input Styles */
#target-container {
    margin: 1 2;
}

.target-config {
    margin: 1 0;
}

.section-title {
    text-style: bold;
    color: $accent;
    margin: 1 0;
}

.selected-tools {
    color: $success;
    margin: 0 2;
}

.input-label {
    margin: 1 0 0 0;
}

#target-input {
    margin: 0 0 1 0;
}

.examples-title {
    text-style: bold;
    color: $warning;
}

.example-text {
    color: $text-muted;
    margin-left: 2;
}

.instructions {
    text-align: center;
    color: $text-muted;
    margin: 1 0;
}

/* Summary Screen Styles */
#summary-container {
    margin: 1 2;
}

.summary-content {
    margin: 1 0;
}

.summary-label {
    text-style: bold;
    color: $accent;
}

.summary-value {
    color: $success;
    text-style: bold;
}

.target-value {
    color: $warning;
}

.summary-tool {
    color: $success;
}

.note {
    color: $warning;
    text-style: italic;
    text-align: center;
}

.ready-message {
    text-style: bold;
    color: $success;
    text-align: center;
}

.action-section {
    text-align: center;
}
"""

#!/usr/bin/env python3
"""
Utilities Module
Contains shared utility functions for IPCrawler
"""

import re
import ipaddress
from typing import Tuple, Optional


def validate_target(target: str) -> Tuple[bool, Optional[str]]:
    """
    Validate if the target is a valid IP address, domain, or CIDR.
    
    Args:
        target: The target string to validate
        
    Returns:
        Tuple of (is_valid, error_message)
    """
    if not target or not target.strip():
        return False, "Target cannot be empty"
    
    target = target.strip()
    
    # Check if it's an IP address
    try:
        ipaddress.ip_address(target)
        return True, None
    except ValueError:
        pass
    
    # Check if it's a CIDR network
    try:
        ipaddress.ip_network(target, strict=False)
        return True, None
    except ValueError:
        pass
    
    # Check if it's a valid domain name
    domain_pattern = re.compile(
        r'^(?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)*[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$'
    )
    
    if domain_pattern.match(target):
        return True, None
    
    return False, "Invalid target format. Please enter a valid IP address, domain name, or CIDR notation."


def estimate_scan_time(tools: list) -> Tuple[int, int]:
    """
    Estimate the scan time based on selected tools.
    
    Args:
        tools: List of selected tool names
        
    Returns:
        Tuple of (min_time_minutes, max_time_minutes)
    """
    # Base time estimates per tool (in minutes)
    tool_times = {
        "naabu": (1, 3),
        "nmap": (3, 10),
        "subfinder": (1, 2),
        "dnsx": (1, 2),
        "httpx": (1, 3),
        "katana": (2, 5),
        "nuclei": (5, 15),
        "ffuf": (3, 8),
        "feroxbuster": (3, 8),
        "shuffledns": (2, 4),
        "amass": (5, 20),
    }
    
    min_total = 0
    max_total = 0
    
    for tool in tools:
        if tool in tool_times:
            min_time, max_time = tool_times[tool]
            min_total += min_time
            max_total += max_time
        else:
            # Default estimate for unknown tools
            min_total += 2
            max_total += 5
    
    return min_total, max_total


def format_tool_list(tools: list) -> str:
    """
    Format a list of tools for display.
    
    Args:
        tools: List of tool names
        
    Returns:
        Formatted string of tools
    """
    if not tools:
        return "No tools selected"
    
    if len(tools) == 1:
        return tools[0]
    elif len(tools) == 2:
        return f"{tools[0]} and {tools[1]}"
    else:
        return f"{', '.join(tools[:-1])}, and {tools[-1]}"

#!/usr/bin/env python3
"""
IPCrawler - Intelligent Recon Flow Engine
A modern terminal UI for cybersecurity reconnaissance automation
Developer: hckerhub

Main entry point for the reorganized application
"""

import sys
from src import IPCrawlerApp


def main():
    """Main entry point."""
    try:
        app = IPCrawlerApp()
        app.run()
    except KeyboardInterrupt:
        print("\n👋 Thanks for using IPCrawler!")
        sys.exit(0)
    except Exception as e:
        print(f"❌ Error: {e}")
        sys.exit(1)


if __name__ == "__main__":
    main()

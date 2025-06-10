.PHONY: setup clean setup-docker docker-cmd help

setup:
	@echo "Setting up ipcrawler..."
	python3 -m venv venv
	venv/bin/python3 -m pip install --upgrade pip
	venv/bin/python3 -m pip install -r requirements.txt
	@echo "Creating ipcrawler command..."
	@rm -f ipcrawler-cmd
	@echo '#!/bin/bash' > ipcrawler-cmd
	@echo '# Resolve the real path of the script (follow symlinks)' >> ipcrawler-cmd
	@echo 'SCRIPT_PATH="$$(realpath "$${BASH_SOURCE[0]}")"' >> ipcrawler-cmd
	@echo 'DIR="$$(cd "$$(dirname "$$SCRIPT_PATH")" && pwd)"' >> ipcrawler-cmd
	@echo 'source "$$DIR/venv/bin/activate" && python3 "$$DIR/ipcrawler.py" "$$@"' >> ipcrawler-cmd
	@chmod +x ipcrawler-cmd
	@echo "Installing ipcrawler command to /usr/local/bin..."
	@sudo ln -sf "$$(pwd)/ipcrawler-cmd" /usr/local/bin/ipcrawler
	@echo ""
	@echo "✓ Setup complete!"
	@echo "Now you can run: ipcrawler [options]"

clean:
	@echo "Cleaning up ipcrawler installation..."
	@echo ""
	@echo "🧹 Removing virtual environment and command..."
	@if [ -n "$$VIRTUAL_ENV" ]; then \
		echo "Deactivating virtual environment..."; \
		unset VIRTUAL_ENV; \
		unset PYTHONHOME; \
		export PATH="$$(echo "$$PATH" | sed 's|[^:]*venv[^:]*:||g')"; \
	fi
	rm -rf venv
	rm -f ipcrawler-cmd
	@echo "Removing ipcrawler from /usr/local/bin..."
	@sudo rm -f /usr/local/bin/ipcrawler
	@echo ""
	@echo "🐳 Cleaning up Docker resources..."
	@if [ -n "$$(docker images -q ipcrawler 2>/dev/null)" ]; then \
		echo "Stopping any running ipcrawler containers..."; \
		docker ps -aq --filter ancestor=ipcrawler 2>/dev/null | xargs -r docker stop >/dev/null 2>&1 || true; \
		docker ps -aq --filter ancestor=ipcrawler 2>/dev/null | xargs -r docker rm >/dev/null 2>&1 || true; \
		echo "Removing ipcrawler Docker image..."; \
		docker rmi ipcrawler >/dev/null 2>&1 || true; \
		echo "Docker image removed."; \
	else \
		echo "No ipcrawler Docker image found."; \
	fi
	@echo "Cleaning up results directory..."
	@if [ -d "results" ] && [ -z "$$(ls -A results 2>/dev/null)" ]; then \
		rm -rf results; \
		echo "Empty results directory removed."; \
	elif [ -d "results" ]; then \
		echo "Results directory contains files - keeping it."; \
	else \
		echo "No results directory found."; \
	fi
	@echo ""
	@echo "✓ Clean complete!"
	@echo "Virtual environment, Docker image, and empty directories removed."

setup-docker:
	@echo "Building ipcrawler Docker image..."
	docker build -t ipcrawler .
	@echo ""
	@echo "✓ Docker setup complete!"
	@echo "Now you can run: make docker-cmd"
	@echo "Or manually: docker run -it --rm -v \$$(pwd)/results:/scans ipcrawler"

docker-cmd:
	@echo "Starting ipcrawler Docker container..."
	@echo "Results will be saved to: $$(pwd)/results"
	@echo "Type 'exit' to leave the container"
	@echo ""
	docker run -it --rm -v "$$(pwd)/results:/scans" ipcrawler || true

help:
	@echo "Available make commands:"
	@echo ""
	@echo "  setup         - Set up local Python virtual environment"
	@echo "  clean         - Remove local setup, virtual environment, and Docker resources"
	@echo "  setup-docker  - Build Docker image for ipcrawler"
	@echo "  docker-cmd    - Run interactive Docker container"
	@echo "  help          - Show this help message"
	@echo ""
	@echo "Docker Usage:"
	@echo "  1. make setup-docker    # Build the image"
	@echo "  2. make docker-cmd      # Start interactive session"
	@echo ""
	@echo "Local Usage:"
	@echo "  1. make setup           # Set up locally"
	@echo "  2. ipcrawler --help     # Use the tool"

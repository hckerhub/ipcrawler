.PHONY: setup clean

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
	@echo "Removing virtual environment and command..."
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
	@echo "Clean complete."
	@echo "Virtual environment deactivated (if it was active)."

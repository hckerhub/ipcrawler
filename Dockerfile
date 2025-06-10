FROM python:3.11-slim

# Install system packages including build tools for Python packages
RUN apt-get update && apt-get install -y \
    curl wget git \
    nmap \
    gcc python3-dev \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Copy the local ipcrawler source code
COPY . /app
WORKDIR /app

# Install Python dependencies and ipcrawler
RUN pip install --upgrade pip
RUN pip install -r requirements.txt
RUN pip install .

WORKDIR /scans

# Create a script to help users install additional tools if needed
RUN echo '#!/bin/bash\n\
echo "Installing additional security tools..."\n\
apt-get update\n\
apt-get install -y gobuster nikto smbclient whatweb || true\n\
echo "Done! Some tools may not be available on this minimal image."\n\
echo "ipcrawler will work with basic functionality."\n\
' > /install-tools.sh && chmod +x /install-tools.sh

CMD ["/bin/bash"]

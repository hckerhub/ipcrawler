FROM python:3.11-slim

# Install available security tools + build tools for Debian base
RUN apt-get update && apt-get install -y \
    # Core system tools
    curl wget git gcc python3-dev \
    # Security tools available in Debian repos
    nmap dnsrecon gobuster nbtscan redis-tools \
    smbclient sslscan whatweb \
    # Additional useful tools
    netcat-traditional dnsutils whois \
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

# Create tool management scripts
RUN echo '#!/bin/bash\n\
echo "🔍 Pre-installed security tools:"\n\
echo "✅ Core: nmap, dnsrecon, gobuster, nbtscan"\n\
echo "✅ SMB: smbclient, redis-tools"\n\
echo "✅ SSL/Web: sslscan, whatweb"\n\
echo "✅ Network: netcat, dnsutils, whois"\n\
echo ""\n\
echo "📦 To install additional Kali tools, run: /install-extra-tools.sh"\n\
' > /show-tools.sh && chmod +x /show-tools.sh

# Script to install additional tools (that may not be in Debian repos)
RUN echo '#!/bin/bash\n\
echo "🔧 Installing additional security tools..."\n\
echo "Note: Some tools may not be available in standard Debian repos"\n\
apt-get update\n\
# Try to install additional tools, continue on failure\n\
apt-get install -y nikto enum4linux smbmap snmp sipvicious || echo "Some tools unavailable in Debian repos"\n\
# Install tools via other methods\n\
echo "📥 Installing feroxbuster via GitHub releases..."\n\
curl -sL https://github.com/epi052/feroxbuster/releases/latest/download/x86_64-linux-feroxbuster.tar.gz | tar -xz -C /usr/local/bin 2>/dev/null || echo "feroxbuster install failed"\n\
echo "📥 Installing additional Python tools..."\n\
pip install impacket 2>/dev/null || echo "impacket install failed"\n\
echo ""\n\
echo "✅ Additional tool installation complete!"\n\
echo "Run /show-tools.sh to see what is available"\n\
' > /install-extra-tools.sh && chmod +x /install-extra-tools.sh

CMD ["/bin/bash"]

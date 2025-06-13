from ipcrawler.config import config
from ipcrawler.plugins import ServiceScan
import requests
import urllib3
import ipaddress
from urllib.parse import urlparse
import os

urllib3.disable_warnings()

class VHostRedirectHunter(ServiceScan):

    def __init__(self):
        super().__init__()
        self.name = 'VHost Redirect Hunter'
        self.slug = 'vhost-redirect-hunter' 
        self.tags = ['default', 'http', 'safe', 'quick']
        self.priority = 10  # Higher priority than other vhost plugins

    def configure(self):
        self.match_service_name('^http')
        self.match_service_name('^nacn_http$', negative_match=True)

    async def run(self, service):
        service.info("🕷️  VHost Redirect Hunter activated - searching for hostname redirects...")
        
        # Only process IP addresses, not hostnames
        try:
            ipaddress.ip_address(service.target.address)
        except ValueError:
            service.info("⏭️  Skipping hostname target (only processes IP addresses)")
            return

        scheme = 'https' if service.secure else 'http'
        url = f"{scheme}://{service.target.address}:{service.port}/"
        
        try:
            # Get configuration settings
            vhost_config = config.get('vhost_discovery', {})
            timeout = vhost_config.get('request_timeout', 10)
            user_agent = vhost_config.get('user_agent', 'ipcrawler-vhost-hunter/1.0')
            
            # Use requests for better control and parsing
            resp = requests.get(
                url, 
                verify=False, 
                allow_redirects=False, 
                timeout=timeout,
                headers={'User-Agent': user_agent}
            )
            
            redirect_file = f"{service.target.scandir}/vhost_redirects_{service.port}.txt"
            
            if 'Location' in resp.headers:
                location = resp.headers['Location']
                parsed = urlparse(location)
                redirect_host = parsed.hostname
                
                # Save raw redirect info
                with open(redirect_file, 'w') as f:
                    f.write(f"Status: {resp.status_code}\n")
                    f.write(f"Location: {location}\n")
                    f.write(f"Original URL: {url}\n")
                    if redirect_host:
                        f.write(f"Extracted Hostname: {redirect_host}\n")
                
                if redirect_host and redirect_host != service.target.address:
                    # Store discovered VHost for post-scan processing
                    if not hasattr(service.target, 'discovered_vhosts'):
                        service.target.discovered_vhosts = []
                    
                    vhost_info = {
                        'hostname': redirect_host,
                        'ip': service.target.address,
                        'port': service.port,
                        'scheme': scheme,
                        'original_url': url,
                        'redirect_url': location,
                        'status_code': resp.status_code
                    }
                    
                    service.target.discovered_vhosts.append(vhost_info)
                    service.info(f"✅ VHost discovered: {redirect_host}")
                    service.info(f"   Redirect: {url} → {location}")
                    
                    # Add to manual commands for easy copy-paste
                    manual_cmd = f"echo '{service.target.address} {redirect_host}' | sudo tee -a /etc/hosts"
                    await service.execute(f'echo "# VHost discovered: {redirect_host}" >> "{service.target.scandir}/_manual_commands.txt"')
                    await service.execute(f'echo "{manual_cmd}" >> "{service.target.scandir}/_manual_commands.txt"')
                    await service.execute(f'echo "" >> "{service.target.scandir}/_manual_commands.txt"')
                    
                else:
                    service.info(f"🔍 Redirect found but no useful hostname: {location}")
            else:
                service.info(f"📋 No redirect detected at {url}")
                # Create empty file to show scan was attempted
                with open(redirect_file, 'w') as f:
                    f.write(f"Original URL: {url}\n")
                    f.write(f"Status: {resp.status_code}\n")
                    f.write(f"Result: No redirect detected\n")
                
        except requests.exceptions.RequestException as e:
            service.error(f"❌ Request failed for {url}: {str(e)[:100]}...")
            # Log the error to file
            error_file = f"{service.target.scandir}/vhost_redirects_{service.port}.txt"
            with open(error_file, 'w') as f:
                f.write(f"Original URL: {url}\n")
                f.write(f"Error: {str(e)}\n")
        except Exception as e:
            service.error(f"❌ Unexpected error scanning {url}: {e}")

    def cleanup(self):
        """Called after all scans complete - process discovered VHosts"""
        # This will be called automatically by ipcrawler after scanning 
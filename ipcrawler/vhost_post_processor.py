#!/usr/bin/env python3
"""
VHost Post-Processor for ipcrawler
Handles discovered VHosts and provides interactive /etc/hosts management
"""

import os
import sys
import subprocess
import shutil
from collections import defaultdict
from datetime import datetime

# Try to import config, fallback to defaults if not available
try:
    from ipcrawler.config import config
except ImportError:
    config = {}

class VHostPostProcessor:
    
    def __init__(self, scan_directories):
        # Handle both single directory (string) and multiple directories (list)
        if isinstance(scan_directories, str):
            self.scan_directories = [scan_directories]
        else:
            self.scan_directories = scan_directories
        self.discovered_vhosts = []
        self.existing_hosts = set()
        
    def discover_vhosts_from_files(self):
        """Parse VHost discovery files and extract hostnames"""
        for scan_dir in self.scan_directories:
            for root, dirs, files in os.walk(scan_dir):
                for file in files:
                    if file.startswith('vhost_redirects_') and file.endswith('.txt'):
                        file_path = os.path.join(root, file)
                        try:
                            with open(file_path, 'r') as f:
                                content = f.read()
                                
                            # Extract IP from scan directory name or file content
                            ip = os.path.basename(root)
                            
                            # Extract hostname from file content
                            for line in content.split('\n'):
                                if line.startswith('Extracted Hostname:'):
                                    hostname = line.split(':', 1)[1].strip()
                                    if hostname and hostname != ip:
                                        self.discovered_vhosts.append({
                                            'hostname': hostname,
                                            'ip': ip,
                                            'file': file_path
                                        })
                        except Exception as e:
                            print(f"❌ Error parsing {file_path}: {e}")
                        
    def read_existing_hosts(self):
        """Read existing /etc/hosts entries"""
        try:
            with open('/etc/hosts', 'r') as f:
                for line in f:
                    line = line.strip()
                    if line and not line.startswith('#'):
                        parts = line.split()
                        if len(parts) >= 2:
                            for hostname in parts[1:]:
                                self.existing_hosts.add(hostname)
        except Exception as e:
            print(f"❌ Error reading /etc/hosts: {e}")
            
    def backup_hosts_file(self):
        """Create backup of /etc/hosts"""
        vhost_config = config.get('vhost_discovery', {})
        if not vhost_config.get('backup_hosts_file', True):
            print("⚠️  Backup disabled in config - proceeding without backup")
            return "/etc/hosts"  # Return original path to continue
            
        timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
        backup_path = f"/etc/hosts.backup.{timestamp}"
        try:
            shutil.copy2('/etc/hosts', backup_path)
            print(f"✅ Created backup: {backup_path}")
            return backup_path
        except Exception as e:
            print(f"❌ Failed to create backup: {e}")
            return None
            
    def check_sudo_privileges(self):
        """Check if we have sudo privileges"""
        try:
            result = subprocess.run(['sudo', '-n', 'true'], 
                                  capture_output=True, 
                                  text=True, 
                                  timeout=5)
            return result.returncode == 0
        except:
            return False
            
    def add_hosts_entries(self, entries_to_add):
        """Add entries to /etc/hosts"""
        try:
            # Create temporary file with entries
            temp_entries = []
            for entry in entries_to_add:
                temp_entries.append(f"{entry['ip']} {entry['hostname']}")
                
            # Prepare the command
            entries_text = '\n'.join([
                "\n# Added by ipcrawler VHost discovery",
                f"# Generated: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}"
            ] + temp_entries + [""])
            
            # Write to /etc/hosts using sudo
            process = subprocess.Popen(['sudo', 'tee', '-a', '/etc/hosts'], 
                                     stdin=subprocess.PIPE, 
                                     stdout=subprocess.PIPE, 
                                     stderr=subprocess.PIPE,
                                     text=True)
            
            stdout, stderr = process.communicate(input=entries_text)
            
            if process.returncode == 0:
                print("✅ Successfully added entries to /etc/hosts")
                return True
            else:
                print(f"❌ Failed to add entries: {stderr}")
                return False
                
        except Exception as e:
            print(f"❌ Error adding hosts entries: {e}")
            return False
            
    def display_summary_table(self):
        """Display discovered VHosts in a nice table"""
        if not self.discovered_vhosts:
            print("\n📋 No VHosts discovered during scanning")
            return
            
        print("\n" + "="*70)
        print("🌐 DISCOVERED VIRTUAL HOSTS")
        print("="*70)
        
        # Group by IP for cleaner display
        grouped = defaultdict(list)
        for vhost in self.discovered_vhosts:
            grouped[vhost['ip']].append(vhost['hostname'])
            
        for ip, hostnames in grouped.items():
            print(f"\n📍 Target: {ip}")
            for hostname in hostnames:
                status = "✅ NEW" if hostname not in self.existing_hosts else "⚠️  EXISTS"
                print(f"   {hostname} ({status})")
                
        print("\n" + "="*70)
        
    def run_interactive_session(self):
        """Run the interactive VHost management session"""
        # Check if VHost discovery is enabled
        vhost_config = config.get('vhost_discovery', {})
        if not vhost_config.get('enabled', True):
            return
            
        print("\n🚀 VHost Discovery Post-Processing")
        print("=" * 50)
        
        # Discover VHosts from scan files
        self.discover_vhosts_from_files()
        
        if not self.discovered_vhosts:
            print("📋 No VHosts discovered during scanning")
            return
            
        # Read existing hosts
        self.read_existing_hosts()
        
        # Show summary
        self.display_summary_table()
        
        # Filter out existing entries
        new_vhosts = [v for v in self.discovered_vhosts 
                     if v['hostname'] not in self.existing_hosts]
        
        if not new_vhosts:
            print("\n✅ All discovered VHosts already exist in /etc/hosts")
            return
            
        print(f"\n🎯 Found {len(new_vhosts)} new VHost(s) to add")
        
        # Check privileges
        has_sudo = self.check_sudo_privileges()
        
        if not has_sudo:
            print("\n🔐 Sudo privileges required for /etc/hosts modification")
            print("📝 Manual commands saved to _manual_commands.txt files")
            print("\nManual addition commands:")
            for vhost in new_vhosts:
                print(f"echo '{vhost['ip']} {vhost['hostname']}' | sudo tee -a /etc/hosts")
            return
            
        # Check if auto-add is enabled
        auto_add = vhost_config.get('auto_add_hosts', True)
        
        if not auto_add:
            print("\n📝 Auto-add disabled in config. Manual commands available:")
            for vhost in new_vhosts:
                print(f"echo '{vhost['ip']} {vhost['hostname']}' | sudo tee -a /etc/hosts")
            return
            
        # Interactive prompt
        print(f"\n🤔 Add {len(new_vhosts)} new VHost(s) to /etc/hosts?")
        
        for vhost in new_vhosts:
            print(f"   {vhost['ip']} {vhost['hostname']}")
            
        while True:
            choice = input(f"\n[Y]es / [N]o / [S]how details: ").lower().strip()
            
            if choice in ['y', 'yes']:
                # Create backup first
                backup_path = self.backup_hosts_file()
                if backup_path:
                    # Add entries
                    if self.add_hosts_entries(new_vhosts):
                        print("\n🎉 VHost entries successfully added!")
                        print("💡 Use 'sudo nano /etc/hosts' to edit manually if needed")
                        print(f"🔄 Restore with: sudo cp {backup_path} /etc/hosts")
                    else:
                        print("\n❌ Failed to add VHost entries")
                else:
                    print("\n❌ Cannot proceed without backup")
                break
                
            elif choice in ['n', 'no']:
                print("\n⏭️  Skipping /etc/hosts modification")
                print("📝 Manual commands available in _manual_commands.txt files")
                break
                
            elif choice in ['s', 'show', 'details']:
                print("\n📋 VHost Details:")
                for vhost in new_vhosts:
                    print(f"\n🔗 {vhost['hostname']}")
                    print(f"   IP: {vhost['ip']}")
                    print(f"   Source: {vhost['file']}")
                continue
                
            else:
                print("❌ Please enter Y, N, or S")
                continue

def main():
    """Main entry point"""
    if len(sys.argv) != 2:
        print("Usage: python3 vhost_post_processor.py <scans_directory>")
        sys.exit(1)
        
    scans_dir = sys.argv[1]
    
    if not os.path.exists(scans_dir):
        print(f"❌ Scans directory not found: {scans_dir}")
        sys.exit(1)
        
    processor = VHostPostProcessor(scans_dir)
    processor.run_interactive_session()

if __name__ == "__main__":
    main() 
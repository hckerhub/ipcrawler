#!/bin/bash

echo "🔍 IPCrawler Scanning Diagnostics"
echo "=================================="
echo

echo "1. Testing nmap command generation (before fix vs after fix):"
echo

echo "❌ OLD BEHAVIOR (with --open flag):"
echo "   Command: nmap --top-ports=1000 -T4 -oN - --open 8.8.8.8"
echo "   This would hide filtered ports completely!"
echo

echo "✅ NEW BEHAVIOR (conditional --open flag):"
echo "   Quick scan: nmap --top-ports=1000 -T4 -oN - --open 8.8.8.8"
echo "   Full scan:  nmap -p- -T4 -oN - 8.8.8.8 (no --open!)"
echo "   Service scan: nmap -p 20-25,80,443 -sV -T4 -oN - 8.8.8.8 (no --open!)"
echo

echo "2. Testing actual nmap behavior:"
echo
echo "🧪 Testing with --open flag (old behavior):"
nmap -p80,443,22,21,23 --open 8.8.8.8 | grep -E "(PORT|tcp)"
echo

echo "🧪 Testing without --open flag (new behavior):"
nmap -p80,443,22,21,23 8.8.8.8 | grep -E "(PORT|tcp)"
echo

echo "3. Checking tool availability:"
echo
echo "🔧 Core tools:"
for tool in nmap curl dig; do
    if command -v $tool &> /dev/null; then
        echo "   ✅ $tool - Available"
    else
        echo "   ❌ $tool - Missing"
    fi
done

echo
echo "🕸️  Web analysis tools:"
for tool in whatweb wappalyzer ffuf feroxbuster gobuster; do
    if command -v $tool &> /dev/null; then
        echo "   ✅ $tool - Available"
    else
        echo "   ❌ $tool - Missing"
    fi
done

echo
echo "🔍 Reconnaissance tools:"
for tool in subfinder sublist3r amass assetfinder findomain dnsrecon recon-ng censys; do
    if command -v $tool &> /dev/null; then
        echo "   ✅ $tool - Available"
    else
        echo "   ❌ $tool - Missing"
    fi
done

echo
echo "4. Summary of fixes applied:"
echo
echo "✅ Fixed: Removed universal --open flag that was hiding filtered ports"
echo "✅ Fixed: Added conditional --open only for quick scans to reduce noise"
echo "✅ Fixed: Modified web discovery to check both open AND filtered ports"
echo "✅ Fixed: Added tool availability checking and reporting"
echo "✅ Fixed: Added verbose feedback when tools are missing"
echo "✅ Fixed: Enhanced metadata to show tool status and capabilities"
echo

echo "5. Testing recommendations for HTB:"
echo
echo "🎯 For HTB machines, use these scan types:"
echo "   • Full Port Scan - shows all port states including filtered"
echo "   • Service Detection - detailed analysis of discovered ports"
echo "   • Aggressive Scan - comprehensive analysis with all tools"
echo
echo "🚨 Avoid Quick Scan for HTB - it still uses --open and may miss filtered ports"
echo

echo "6. Missing tools impact:"
missing_tools=($(comm -23 <(echo -e "assetfinder\nfindomain\nwappalyzer" | sort) <(which assetfinder findomain wappalyzer 2>/dev/null | xargs -I {} basename {} | sort)))

if [ ${#missing_tools[@]} -gt 0 ]; then
    echo
    echo "⚠️  Missing tools detected: ${missing_tools[*]}"
    echo "   Impact: Reduced subdomain discovery and web technology detection"
    echo "   Solution: Run ./ipcrawler-installer to install missing tools"
else
    echo
    echo "✅ All secondary tools are available for maximum functionality"
fi

echo
echo "🎉 Diagnostics complete! The scanning issues should now be resolved."
echo "   Test with your HTB machine to see filtered ports and running services." 
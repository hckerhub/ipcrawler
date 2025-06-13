import asyncio, colorama, os, re, string, sys, unidecode
from colorama import Fore, Style
from ipcrawler.config import config

# Rich support for enhanced verbosity output (optional)
try:
	from rich.console import Console
	from rich.text import Text
	from rich.panel import Panel
	from rich.progress import Progress, SpinnerColumn, TextColumn, TimeElapsedColumn
	RICH_AVAILABLE = True
	rich_console = Console()
except ImportError:
	RICH_AVAILABLE = False
	rich_console = None

def slugify(name):
	return re.sub(r'[\W_]+', '-', unidecode.unidecode(name).lower()).strip('-')

def e(*args, frame_index=1, **kvargs):
	frame = sys._getframe(frame_index)

	vals = {}

	vals.update(frame.f_globals)
	vals.update(frame.f_locals)
	vals.update(kvargs)

	return string.Formatter().vformat(' '.join(args), args, vals)

def fformat(s):
	return e(s, frame_index=3)

def cprint(*args, color=Fore.RESET, char='*', sep=' ', end='\n', frame_index=1, file=sys.stdout, printmsg=True, verbosity=0, **kvargs):
	if printmsg and verbosity > config['verbose']:
		return ''
	frame = sys._getframe(frame_index)

	vals = {
		'bgreen':  Fore.GREEN  + Style.BRIGHT,
		'bred':	Fore.RED	+ Style.BRIGHT,
		'bblue':   Fore.BLUE   + Style.BRIGHT,
		'byellow': Fore.YELLOW + Style.BRIGHT,
		'bmagenta': Fore.MAGENTA + Style.BRIGHT,

		'green':  Fore.GREEN,
		'red':	Fore.RED,
		'blue':   Fore.BLUE,
		'yellow': Fore.YELLOW,
		'magenta': Fore.MAGENTA,

		'bright': Style.BRIGHT,
		'srst':   Style.NORMAL,
		'crst':   Fore.RESET,
		'rst':	Style.NORMAL + Fore.RESET
	}

	if config['accessible']:
		 vals = {'bgreen':'', 'bred':'', 'bblue':'', 'byellow':'', 'bmagenta':'', 'green':'', 'red':'', 'blue':'', 'yellow':'', 'magenta':'', 'bright':'', 'srst':'', 'crst':'', 'rst':''}

	vals.update(frame.f_globals)
	vals.update(frame.f_locals)
	vals.update(kvargs)

	unfmt = ''
	if char is not None and not config['accessible']:
		unfmt += color + '[' + Style.BRIGHT + char + Style.NORMAL + ']' + Fore.RESET + sep
	unfmt += sep.join(args)

	fmted = unfmt

	for attempt in range(10):
		try:
			fmted = string.Formatter().vformat(unfmt, args, vals)
			break
		except KeyError as err:
			key = err.args[0]
			unfmt = unfmt.replace('{' + key + '}', '{{' + key + '}}')

	if printmsg:
		print(fmted, sep=sep, end=end, file=file)
	else:
		return fmted

def debug(*args, color=Fore.GREEN, sep=' ', end='\n', file=sys.stdout, **kvargs):
	if config['verbose'] >= 2:
		if config['accessible']:
			args = ('Debug:',) + args
		cprint(*args, color=color, char='-', sep=sep, end=end, file=file, frame_index=2, **kvargs)

def info(*args, sep=' ', end='\n', file=sys.stdout, **kvargs):
	# Enhanced Rich output for certain verbosity messages
	if RICH_AVAILABLE and 'verbosity' in kvargs:
		message = sep.join(str(arg) for arg in args)
		verbosity_level = kvargs['verbosity']
		
		# Level 1 (-v): Plugin starts and discoveries
		if config['verbose'] >= 1 and verbosity_level == 1:
			# Enhanced plugin start messages
			if 'running against' in message and ('Port scan' in message or 'Service scan' in message):
				# Extract plugin info
				if 'Port scan' in message:
					scan_type = "🔍 PORT"
					color = "blue"
				else:
					scan_type = "🔧 SERVICE" 
					color = "green"
				
				# Parse the message for plugin name and target
				import re
				plugin_match = re.search(r'(Port scan|Service scan) ([^{]+?) \(([^)]+)\)', message)
				target_match = re.search(r'against ([^{]+?)$', message.replace('{rst}', '').replace('{byellow}', '').replace('{rst}', ''))
				
				if plugin_match and target_match:
					plugin_name = plugin_match.group(2).strip()
					plugin_slug = plugin_match.group(3).strip()
					target = target_match.group(1).strip()
					
					rich_text = Text()
					rich_text.append(f"[{scan_type}] ", style=f"bold {color}")
					rich_text.append(f"{plugin_name} ", style="bold cyan")
					rich_text.append(f"({plugin_slug}) ", style="dim")
					rich_text.append("→ ", style="bold white")
					rich_text.append(f"{target}", style="bold yellow")
					
					rich_console.print(rich_text)
					return
			
			# Enhanced discovery messages  
			elif 'Discovered open port' in message:
				import re
				port_match = re.search(r'Discovered open port ([^{]+?) on ([^{]+?)$', message.replace('{rst}', '').replace('{bmagenta}', '').replace('{byellow}', ''))
				if port_match:
					port = port_match.group(1).strip()
					target = port_match.group(2).strip()
					
					rich_text = Text()
					rich_text.append("🎯 DISCOVERED ", style="bold green")
					rich_text.append(f"{port} ", style="bold magenta")
					rich_text.append("on ", style="dim")
					rich_text.append(f"{target}", style="bold yellow")
					
					rich_console.print(rich_text)
					return
		
		# Level 2 (-vv): Completion messages with timing
		elif config['verbose'] >= 2 and verbosity_level == 2:
			if 'finished in' in message and ('Port scan' in message or 'Service scan' in message):
				import re
				plugin_match = re.search(r'(Port scan|Service scan) ([^{]+?) \(([^)]+)\)', message)
				target_match = re.search(r'against ([^{]+?) finished in (.+)$', message.replace('{rst}', '').replace('{byellow}', ''))
				
				if plugin_match and target_match:
					scan_type = "✅ COMPLETED" if 'Port scan' in message else "✅ FINISHED"
					plugin_name = plugin_match.group(2).strip()
					plugin_slug = plugin_match.group(3).strip()
					target = target_match.group(1).strip()
					timing = target_match.group(2).strip()
					
					rich_text = Text()
					rich_text.append(f"{scan_type} ", style="bold green")
					rich_text.append(f"{plugin_name} ", style="bold cyan")
					rich_text.append(f"({plugin_slug}) ", style="dim")
					rich_text.append("on ", style="dim")
					rich_text.append(f"{target} ", style="bold yellow")
					rich_text.append("in ", style="dim")
					rich_text.append(f"{timing}", style="bold blue")
					
					rich_console.print(rich_text)
					return
			
			# Enhanced pattern match messages
			elif 'Matched Pattern:' in message or 'pattern' in message.lower():
				rich_text = Text()
				rich_text.append("🔍 PATTERN ", style="bold magenta")
				# Extract the actual pattern content
				pattern_content = message.replace('{rst}', '').replace('{bmagenta}', '').replace('{bright}', '').replace('{yellow}', '').replace('{crst}', '').replace('{bgreen}', '')
				rich_text.append(pattern_content, style="cyan")
				
				rich_console.print(rich_text)
				return
	
	# Fall back to standard cprint
	cprint(*args, color=Fore.BLUE, char='*', sep=sep, end=end, file=file, frame_index=2, **kvargs)

def warn(*args, sep=' ', end='\n', file=sys.stderr,**kvargs):
	if config['accessible']:
		args = ('Warning:',) + args
	cprint(*args, color=Fore.YELLOW, char='!', sep=sep, end=end, file=file, frame_index=2, **kvargs)

def error(*args, sep=' ', end='\n', file=sys.stderr, **kvargs):
	if config['accessible']:
		args = ('Error:',) + args
	cprint(*args, color=Fore.RED, char='!', sep=sep, end=end, file=file, frame_index=2, **kvargs)

def fail(*args, sep=' ', end='\n', file=sys.stderr, **kvargs):
	if config['accessible']:
		args = ('Failure:',) + args
	cprint(*args, color=Fore.RED, char='!', sep=sep, end=end, file=file, frame_index=2, **kvargs)
	exit(-1)

class CommandStreamReader(object):

	def __init__(self, stream, target, tag, patterns=None, outfile=None):
		self.stream = stream
		self.target = target
		self.tag = tag
		self.lines = []
		self.patterns = patterns or []
		self.outfile = outfile
		self.ended = False

		# Empty files that already exist.
		if self.outfile != None:
			with open(self.outfile, 'w'): pass

	# Read lines from the stream until it ends.
	async def _read(self):
		while True:
			if self.stream.at_eof():
				break
			try:
				line = (await self.stream.readline()).decode('utf8').rstrip()
			except ValueError:
				error('{bright}[{yellow}' + self.target.address + '{crst}/{bgreen}' + self.tag + '{crst}]{rst} A line was longer than 64 KiB and cannot be processed. Ignoring.')
				continue

			if line != '':
				# For verbosity 3, enhance only slightly to avoid overwhelming output
				if RICH_AVAILABLE and config['verbose'] >= 3:
					# Very minimal enhancement - just add a subtle indicator
					rich_text = Text()
					rich_text.append("│ ", style="dim blue")
					rich_text.append(f"[{self.target.address}/{self.tag}] ", style="dim")
					rich_text.append(line.strip(), style="white")
					rich_console.print(rich_text)
				else:
					info('{bright}[{yellow}' + self.target.address + '{crst}/{bgreen}' + self.tag + '{crst}]{rst} ' + line.strip().replace('{', '{{').replace('}', '}}'), verbosity=3)

			# Check lines for pattern matches.
			for p in self.patterns:
				description = ''

				# Match and replace entire pattern.
				match = p.pattern.search(line)
				if match:
					if p.description:
						description = p.description.replace('{match}', line[match.start():match.end()])

						# Match and replace substrings.
						matches = p.pattern.findall(line)
						if len(matches) > 0 and isinstance(matches[0], tuple):
							matches = list(matches[0])

						match_count = 1
						for match in matches:
							if p.description:
								description = description.replace('{match' + str(match_count) + '}', match)
							match_count += 1

						async with self.target.lock:
							with open(os.path.join(self.target.scandir, '_patterns.log'), 'a') as file:
								info('{bright}[{yellow}' + self.target.address + '{crst}/{bgreen}' + self.tag + '{crst}]{rst} {bmagenta}' + description + '{rst}', verbosity=2)
								file.writelines(description + '\n\n')
					else:
						info('{bright}[{yellow}' + self.target.address + '{crst}/{bgreen}' + self.tag + '{crst}]{rst} {bmagenta}Matched Pattern: ' + line[match.start():match.end()] + '{rst}', verbosity=2)
						async with self.target.lock:
							with open(os.path.join(self.target.scandir, '_patterns.log'), 'a') as file:
								file.writelines('Matched Pattern: ' + line[match.start():match.end()] + '\n\n')

			if self.outfile is not None:
				with open(self.outfile, 'a') as writer:
					writer.write(line + '\n')
			self.lines.append(line)
		self.ended = True

	# Read a line from the stream cache.
	async def readline(self):
		while True:
			try:
				return self.lines.pop(0)
			except IndexError:
				if self.ended:
					return None
				else:
					await asyncio.sleep(0.1)

	# Read all lines from the stream cache.
	async def readlines(self):
		lines = []
		while True:
			line = await self.readline()
			if line is not None:
				lines.append(line)
			else:
				break
		return lines

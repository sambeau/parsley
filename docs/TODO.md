# TODO List

## IN PLANNING
- Refactor codebase to library: to prepare for more than one command; to prepare for HTTP server
	- use cmd/ for commands
### For V1.0 ALPHA

- Check everything is working:
	- Setup database and sftp environments for integration testing
	- run integration tests
- Build something big to see how it works in practice
- POSTGRES support

### For V1.0 BETA

### For V1.0 RELEASE
- Rewrite README

### For After V1.0 RELEASE

- Parsley Server: Simple, minimal HTTP(S) server that outputs raw HTML files and runs Parsley scripts
	- Ceate Plan
		- Read https://grafana.com/blog/2024/02/09/how-i-write-http-services-in-go-after-13-years/
		- Examine other small, focused, Go HTTP servers, e.g.
		- Comsider useful features from Hugo
	- Import environment
	- Config file? vs Web-based admin
		- HTAccess?
	- HTML/HTTP features for Parsley HTTP Server-to-language API 
		- Investigate in interface/api/environment/context between HTTP and Parsley
		- Request to dictionary
		- dictionary to Response 
		- Cookies
		- Multi-part data

### For After V2.0 RELEASE

	- Website for Parsley hosted with Parsley
	
## CONSIDERING

- Parsley-based static site generator: a small, simple, opinionated site generator
	- Examine Hugo https://gohugo.io/documentation/
- Supabase database support
	- Consider postgrest-go support
		- Examine https://github.com/supabase-community/postgrest-go
	- Examine https://github.com/supabase-community/supabase-go
	- Would not need any realtime features
- MCP support for Parsley
	- Investigate how this would work
- Treesitter grammar
	- Investigate how useful this would be
- Optional basic, type checking: basics + array of basics, no user-defined types, e.g. foo(bar:int){bar}; [int], [[int]]
- Require 'let' to declare variables before use?
- Add 'const' to declare conts to prevent modification?
- Optional Chaining (?.) and maybe (!.)?
- Dictionaries to props in tags: extraProps= {a:"A", b:"B"};<foo {extraProps}/> => <foo a="A" b="B"/>

## DONE

- ~~File type~~ ✅ (v0.8.0)
- ~~Modules, import and export~~ ✅ (v0.9.0)
- ~~Datetime~~ ✅ (v0.6.0)
- ~~Datetime literals with @ syntax~~ ✅ (v0.7.0)
- ~~Duration type~~ ✅ (v0.7.0)
- ~~Duration literals with @ syntax~~ ✅ (v0.7.0)
- ~~URL type~~ ✅ (v0.8.0)
- ~~For loop with indexing~~ ✅ (v0.9.2)
- ~~Open-ended slicing~~ ✅ (v0.9.1)
- ~~Regular expressions~~ ✅ (v0.6.0)
- ~~Regular expression literals with /pattern/ syntax~~ ✅ (v0.6.0)
- ~~i8n/Localisation~~ ✅ (v0.9.7)
- ~~Nullish coalescing operator (??)~~ ✅ (v0.9.9)
- ~~File handle objects (file(), JSON(), CSV(), lines(), text(), bytes())~~ ✅ (v0.9.9)
- ~~Read operator (<==)~~ ✅ (v0.9.9)
- ~~Write operators (==>, ==>>)~~ ✅ (v0.9.9)
- ~~Directory operations (dir(), glob())~~ ✅ (v0.9.9)
- ~~File globbing to dictionary~~ ✅ (v0.9.9)
- ~~Read JSON / write JSON~~ ✅ (v0.9.9)
- ~~File I/O error capture ({data, error} <== file)~~ ✅ (v0.9.9)
- ~~$Decimal type for money?~~ ❌ invalidated by (v0.9.7)
- ~~Markdown support~~ ✅ 0.9.10
- ~~Paths: path and name, display as name~~ ✅ (v0.9.9)
- ~~Sort out what is and isn't a 'let'~~ 🤷‍♂️ turns out they all are so added export instead
- ~~SQL and databases {user} = [$GetUser userID={userId}] <=/=> SQL()~~ ✅ (0.9.15)
- ~~Fetch from URL~~ ✅ (0.9.11)
- ~~1 ++ [2,3,4,5] , [1,2,3,4] ++ 5~~ ✅ (0.9.16)
- ~~File delete() methods: add a delete method to file pseudo-type~~ ✅ (v0.9.17 - implemented as remove())
- ~~File I/O security sandbox: --no-read, --no-write flags~~ ✅ (v0.10.0)
- ~~chroots for write: limit writes to one or more directories and their children~~ ✅ (v0.10.0)
- ~~chroots for read: limit reads to one or more directories and their children~~ ✅ (v0.10.0)
- ~~chroots for execute: limit execution of external scripts/tools to one or more directories and their children, e.g. ./bin~~ ✅ (v0.10.0)
- ~~run scripts: Execute external commands/scripts/tools with command line inputs/options, receive exit status + optional output as result~~ ✅ (v0.11.0)
- ~~pipe scripts: Execute external commands/scripts/tools with command line inputs/options, receive exit status + optional output as result~~ ✅ (v0.11.0)
- ~~Fetch support for format objects~~ ✅ (v0.9.11 - documented in v0.11.0)
- ~~SFTP Support: Read/Write files from FTP server - useful for static site generation~~ ✅ (v0.12.0)
- ~~Directory manipulation methods for file paths: Need methods like .mkdir(), .rmdir(), .remove() for local file paths (currently only available for SFTP in plan)~~ ✅ (v0.12.1)
- ~~Improve REPL~~ ✅ (v0.12.2)
	- ~~Investigate options: what do other Go cli tools do?~~ ✅
	- ~~Better editing Up, down, left, right~~ ✅
- ~~Look at consistency of API one more time~~ ✅ (v0.13.0)
	- ~~Remove deprecated features~~ ✅ (v0.13.0)
- ~~Performance checks~~ ✅ (v0.13.0)
- ~~Code quality checks~~ ✅ (v0.13.0)
- ~~Better errors~~ ✅ (v0.13.1)
	- ~~Human-readable type names in error messages~~ ✅ (v0.13.1)
	- ~~Consistent function name formatting in errors~~ ✅ (v0.13.1)
- ~~run through all code looking for missing tests~~✅ (v0.13.1)
- ~~Support for STDIN/STDOUT/STDERR: Unix pipeline integration with @-, @stdin, @stdout, @stderr~~ ✅ (v0.14.0)
- ~~Require bracket syntax for arrays and array destructuring (consistency for v1.0)~~ ✅ (v0.15.0)
- ~~Datetime intersection operator (`&&`) for combining date and time components~~ ✅ (v0.15.2)
- ~~Tag support for REPL~~ ✅ (v0.15.3)

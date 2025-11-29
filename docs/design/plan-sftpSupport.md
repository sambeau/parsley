# Plan: SFTP Support for Parsley

**Author**: Sam Phillips  
**Status**: ✅ Implemented  
**Version**: 0.12.0  
**Date**: 2025-11-29  
**Approved**: 2025-11-29  
**Implemented**: 2025-11-29

**TL;DR**: Add SFTP (SSH File Transfer Protocol) support using stateful connection objects that integrate with Parsley's network I/O operators (`<=/=`, `=/=>`, `=/=>>`). SFTP connections follow the database pattern (persistent, cached) while file operations use the path-first syntax (`conn(@/path).json`) for consistency with local file I/O.

---

## Design Principles

Following Parsley's core philosophy:

1. **Consistency** — SFTP uses network operators like HTTP fetch, connection pattern like databases
2. **Operators Over Methods** — Network I/O via `<=/=` (read), `=/=>` (write), `=/=>>` (append)
3. **Path-First Syntax** — `conn(@/path).format` matches local `file(@/path).format` pattern
4. **Stateful Connections** — SFTP requires persistent SSH sessions with connection caching
5. **No Special Security Flags** — Like HTTP/Fetch, SFTP is unrestricted (only local file security applies)
6. **Error as Values** — Use `{data, error}` destructuring pattern for robust error handling
7. **Format-Aware** — Works with all format accessors (`.json`, `.text`, `.csv`, `.lines`, `.bytes`)

---

## Core Design Question: SFTP vs HTTP vs Local Files

### Key Differences

| Aspect | Local Files | HTTP Fetch | SFTP |
|--------|-------------|------------|------|
| **I/O Operators** | `<==`, `==>`, `==>>` | `<=/=` only | `<=/=`, `=/=>`, `=/=>>` |
| **State** | Stateless | Stateless | Stateful connection |
| **Auth** | File permissions | Optional headers | SSH keys/passwords |
| **Syntax** | `file(@/path).json` | `JSON(@url)` | `conn(@/path).json` |
| **Connection** | Per-operation | Per-request | Persistent session |
| **Directories** | Yes (`dir()`) | No | Yes (`.dir`, `.mkdir()`, etc.) |

### Design Decision

**SFTP uses connection-based pattern (like databases) with network operators (like HTTP) and path-first syntax (like local files).**

**Rationale**:
- SSH sessions are expensive → persistent connections with caching
- Network transfer → use network operators (`<=/=`, `=/=>`, `=/=>>`)
- Path-first syntax → `conn(@/path).json` mirrors `file(@/path).json`
- Authentication is session-based, not per-request
- Directory operations require maintained connection state

---

## API at a Glance

### Complete SFTP Operations

```parsley
// 1. Create connection (cached by host+user+port)
let conn = SFTP(@sftp://user@example.com, {keyFile: @~/.ssh/id_rsa})

// 2. Read files (network read operator)
let data <=/= conn(@/remote/config.json).json
let lines <=/= conn(@/var/log/app.log).lines
let text <=/= conn(@/home/readme.txt)  // defaults to .text

// 3. Write files (network write operator)
data =/=> conn(@/remote/output.json).json
lines =/=> conn(@/var/log/new.log).lines
text =/=> conn(@/home/readme.txt)  // defaults to .text

// 4. Append files (network append operator)
newEntry =/=>> conn(@/var/log/app.log).lines

// 5. List directories (network read operator)
let {files, error} <=/= conn(@/remote/data/).dir

// 6. Directory operations (methods on handles)
conn(@/remote/newdir/).mkdir()
conn(@/remote/olddir/).rmdir()
conn(@/remote/oldfile.txt).remove()

// 7. Error handling
let {data, error} <=/= conn(@/config.json).json
if (error != null) { log("Failed:", error) }

// 8. Connection management
conn.close()  // Explicit close (auto-closes on script exit)
```

### Operator Patterns

| Operation | Local Files | HTTP | SFTP |
|-----------|-------------|------|------|
| **Read** | `data <== file(@/path).json` | `data <=/= JSON(@url)` | `data <=/= conn(@/path).json` |
| **Write** | `data ==> file(@/path).json` | ❌ | `data =/=> conn(@/path).json` |
| **Append** | `data ==>> file(@/path).lines` | ❌ | `data =/=>> conn(@/path).lines` |

---

## Proposed API Details

### Connection Factory

```parsley
// Basic connection with password
let conn = SFTP(@sftp://user:password@example.com:22)

// With SSH key authentication (preferred)
let conn = SFTP(@sftp://user@example.com, {
    keyFile: @~/.ssh/id_rsa,
    passphrase: "key-passphrase"  // Optional
})

// With explicit options
let conn = SFTP(@sftp://user@example.com, {
    keyFile: @~/.ssh/id_rsa,
    timeout: @30s,
    port: 2222  // Override default 22
})

// Connection cached by host+user+port
let conn1 = SFTP(@sftp://user@host.com)
let conn2 = SFTP(@sftp://user@host.com)  // Returns same connection
```

### Reading Files

```parsley
// Read file through connection using format accessors
let data <=/= conn(@/remote/path/config.json).json
let lines <=/= conn(@/var/log/app.log).lines
let content <=/= conn(@/home/user/readme.txt)  // defaults to .text

// With error handling
let {data, error} <=/= conn(@/config.json).json
if (error != null) {
    log("SFTP read failed:", error)
}

// With fallback
let config <=/= conn(@/config.json).json ?? {default: true}

// Store remote file handle
let remoteFile = conn(@/remote/config.json)
let data <=/= remoteFile.json

// Path interpolation with @(...) template syntax
filename = "data"
let content <=/= conn(@(/remote/path/{filename}.json)).json
```

### Writing Files

```parsley
// Write file through connection
myData =/=> conn(@/remote/path/output.json).json
logLines =/=> conn(@/var/log/app.log).lines
"Hello SFTP" =/=> conn(@/home/user/greeting.txt)  // defaults to .text

// Append to file through connection
newLogEntry =/=>> conn(@/var/log/app.log).lines
moreData =/=>> conn(@/data/output.json).json

// With error handling
let error = myData =/=> conn(@/output.json).json
if (error != null) {
    log("SFTP write failed:", error)
}

// Path interpolation - @ inside parens for templates
dir = "backups"
date = now().format("YYYY-MM-DD")
myData =/=> conn(@(/remote/{dir}/{date}/output.json)).json

// Store remote file handle for reuse
let remoteFile = conn(@/output.json)
myData =/=> remoteFile.json
```

### Directory Operations

```parsley
// List directory contents (network read)
let {files, error} <=/= conn(@/remote/path/).dir

for (f in files ?? []) {
    log(f.name, f.size, f.modified)
}

// Create directory
let error = conn(@/remote/new-directory/).mkdir()

// Remove file
let error = conn(@/remote/file.txt).remove()

// Remove directory
let error = conn(@/remote/empty-directory/).rmdir()
```

### Connection Management

```parsley
// Check connection status
if (conn.connected) {
    data <=/= conn(@/config.json).json
}

// Explicit disconnect (connections auto-close on script exit)
conn.close()

// Connection properties
conn.host       // "example.com"
conn.port       // 22
conn.user       // "username"
conn.connected  // true/false
```

---

## URL Scheme Integration

### SFTP URL Structure

```
sftp://[user[:password]@]host[:port]/path
```

Examples:
- `sftp://example.com/path/to/file.json`
- `sftp://user@example.com:2222/data/file.csv`
- `sftp://user:pass@host.com/remote/path.txt` (discouraged)

### Parsing SFTP URLs

Extend existing `parseUrlString()` to recognize `sftp://` scheme:

```go
func parseUrlString(urlStr string, env *Environment) (*Dictionary, error) {
    // ... existing HTTP/HTTPS parsing ...
    
    // Detect SFTP scheme
    if scheme == "sftp" {
        // SFTP URLs should use connection factory, not direct fetch
        return nil, fmt.Errorf("SFTP URLs require SFTP() connection factory")
    }
    
    // ... rest of parsing ...
}
```

---

## Connection Object Design

### SFTP Connection Structure

```go
type SFTPConnection struct {
    Client        *sftp.Client
    SSHClient     *ssh.Client
    Host          string
    Port          int
    User          string
    Connected     bool
    LastError     string
}

func (sc *SFTPConnection) Type() ObjectType { return SFTP_CONNECTION_OBJ }
func (sc *SFTPConnection) Inspect() string {
    status := "connected"
    if !sc.Connected {
        status = "disconnected"
    }
    return fmt.Sprintf("SFTP(%s@%s:%d) [%s]", sc.User, sc.Host, sc.Port, status)
}
```

### Connection Caching

Like database connections, SFTP connections are expensive to create:

```go
var sftpConnectionsMu sync.Mutex
var sftpConnections = make(map[string]*SFTPConnection)

func getSFTPConnection(host, user string, port int, options *Dictionary, env *Environment) (*SFTPConnection, error) {
    cacheKey := fmt.Sprintf("sftp:%s@%s:%d", user, host, port)
    
    sftpConnectionsMu.Lock()
    defer sftpConnectionsMu.Unlock()
    
    if conn, exists := sftpConnections[cacheKey]; exists && conn.Connected {
        return conn, nil
    }
    
    // Create new connection
    // ... (authentication, SSH setup) ...
    
    sftpConnections[cacheKey] = conn
    return conn, nil
}
```

---

### Format Factory Integration

### Connection Call Syntax

SFTP connections are callable with a path to create remote file handles:

```parsley
conn(@/path)              // Returns SFTP file handle (format unspecified)

// Format accessed via properties (lowercase)
conn(@/path).json         // SFTP file handle for JSON format
conn(@/path).text         // SFTP file handle for text format
conn(@/path).csv          // SFTP file handle for CSV format
conn(@/path).lines        // SFTP file handle for lines format
conn(@/path).bytes        // SFTP file handle for bytes format
conn(@/path).file         // Auto-detect format from extension
```

**Pattern Consistency**:
```parsley
// Local files
let data <== file(@/local/path.json).json

// Remote SFTP files  
let data <=/= conn(@/remote/path.json).json
```

This provides clear distinction between local and network I/O:
- Local files: `data <== file(@/local/path).json` (file system)
- Network (HTTP): `data <=/= JSON(@https://api.example.com/data.json)` (stateless, no append)
- Network (SFTP): `data <=/= conn(@/remote/path).json` (stateful, append supported via SSH_FXF_APPEND)

### SFTP File Handle Structure

```go
type SFTPFileHandle struct {
    Connection *SFTPConnection
    Path       string
    Format     string  // "json", "csv", "text", "lines", "bytes" (defaults to "text")
    Options    *Dictionary
}
```

### Read/Write Operations

```go
// Reading SFTP file
func evalSFTPRead(handle *SFTPFileHandle, env *Environment) (Object, error) {
    // Open remote file via SFTP
    file, err := handle.Connection.Client.Open(handle.Path)
    if err != nil {
        return nil, fmt.Errorf("SFTP read failed: %s", err)
    }
    defer file.Close()
    
    // Read content
    data, err := io.ReadAll(file)
    if err != nil {
        return nil, fmt.Errorf("SFTP read failed: %s", err)
    }
    
    // Parse based on format
    switch handle.Format {
    case "json":
        return parseJSON(string(data))
    case "text":
        return &String{Value: string(data)}, nil
    // ... other formats ...
    }
}

// Writing SFTP file
func evalSFTPWrite(handle *SFTPFileHandle, data Object, append bool, env *Environment) error {
    // Determine open flags
    flags := os.O_WRONLY | os.O_CREATE
    if append {
        flags |= os.O_APPEND  // SSH_FXF_APPEND (0x00000004)
    } else {
        flags |= os.O_TRUNC
    }
    
    // Encode based on format
    var content string
    switch handle.Format {
    case "json":
        content = objectToJSON(data)
    case "text":
        if str, ok := data.(*String); ok {
            content = str.Value
        }
    // ... other formats ...
    }
    
    // Open remote file via SFTP with appropriate flags
    file, err := handle.Connection.Client.OpenFile(handle.Path, flags)
    if err != nil {
        return fmt.Errorf("SFTP write failed: %s", err)
    }
    defer file.Close()
    
    // Write content
    _, err = file.Write([]byte(content))
    return err
}
```

**Note**: The `os.O_APPEND` flag is converted to SFTP's `SSH_FXF_APPEND` (0x00000004), fully supported by github.com/pkg/sftp.

---

## Authentication

### SSH Key Authentication (Recommended)

```parsley
// Default key location
let conn = SFTP(@sftp://user@example.com, {
    keyFile: @~/.ssh/id_rsa
})

// Encrypted key with passphrase
let conn = SFTP(@sftp://user@example.com, {
    keyFile: @~/.ssh/id_rsa_encrypted,
    passphrase: "my-secret-passphrase"
})

// Multiple key attempts
let conn = SFTP(@sftp://user@example.com, {
    keyFiles: [@~/.ssh/id_rsa, @~/.ssh/id_ed25519]
})
```

### Password Authentication (Discouraged)

```parsley
// In URL (not recommended - visible in code)
let conn = SFTP(@sftp://user:password@example.com)

// In options (slightly better)
let conn = SFTP(@sftp://user@example.com, {
    password: "my-password"
})

// From environment (best for passwords)
let password = env("SFTP_PASSWORD")
let conn = SFTP(@sftp://user@example.com, {
    password: password
})
```

### Host Key Verification

```parsley
// Strict verification (default - safest)
let conn = SFTP(@sftp://user@example.com, {
    keyFile: @~/.ssh/id_rsa,
    hostKeyVerify: "strict"  // Must match known_hosts
})

// Accept any (insecure - dev only)
let conn = SFTP(@sftp://user@example.com, {
    keyFile: @~/.ssh/id_rsa,
    hostKeyVerify: "accept-any"  // INSECURE!
})

// Known hosts file
let conn = SFTP(@sftp://user@example.com, {
    keyFile: @~/.ssh/id_rsa,
    knownHostsFile: @~/.ssh/known_hosts  // Default
})
```

---

## Directory Operations

### Listing Directories

```parsley
// List directory (network read)
let {files, error} <=/= conn(@/remote/path/).dir

if (error == null) {
    for (file in files) {
        log(file.name, "-", file.size, "bytes")
        log("  Modified:", file.modified)
        log("  Is dir:", file.isDir)
    }
}

// Recursive listing (network read)
let {files, error} <=/= conn(@/remote/path/).dir({recursive: true})
```

### File Info Object

```parsley
{
    name: "config.json",
    path: "/remote/path/config.json",
    size: 1024,
    modified: @2025-11-29T10:30:00,
    isDir: false,
    isFile: true,
    mode: "rw-r--r--",
    owner: "username",
    group: "staff"
}
```

### Creating/Removing Directories

```parsley
// Create directory
let error = conn(@/remote/new-dir/).mkdir()
if (error != null) {
    log("Failed to create directory:", error)
}

// Create nested directories
let error = conn(@/remote/path/to/nested/dir/).mkdir({parents: true})

// Remove empty directory
let error = conn(@/remote/empty-dir/).rmdir()

// Remove directory recursively (dangerous!)
let error = conn(@/remote/dir/).rmdir({recursive: true})
```

---

## Security Integration

### Consistency with HTTP/Fetch

**SFTP follows the same security model as HTTP/Fetch**: No special restrictions.

- ✅ **HTTP/Fetch**: No flags required, unrestricted by default
- ✅ **SFTP**: No flags required, unrestricted by default
- ⚠️ **Local file writes**: Denied by default (`--allow-write` required)

**Rationale**: Network operations (HTTP, SFTP) access *remote* resources, not local file system. Security model focuses on protecting the *local* machine.

### File System Security for SSH Keys

SFTP respects existing file security when reading local SSH key files:

```bash
# Restrict reading sensitive local directories (including /etc)
./pars --restrict-read=/etc deploy.pars  # Error: Can't read /etc/ssh/id_rsa

# Normally SSH keys in ~/.ssh are fine (read allowed by default)
./pars deploy.pars  # ✅ Can read ~/.ssh/id_rsa
```

### Best Practices

**Credential Safety:**
- Store SSH keys in `~/.ssh/` (standard location)
- Use SSH agent for passphrase-protected keys
- Don't hardcode passwords in scripts (use key auth)
- Use environment variables for sensitive config:
  ```parsley
  let keyPath = ENV.SSH_KEY_PATH ?? @~/.ssh/id_rsa
  let conn = SFTP(@sftp://user@host, {keyFile: keyPath})
  ```

**Host Key Verification:**
- Rely on `~/.ssh/known_hosts` (SSH's standard mechanism)
- First connection requires manual verification
- Subsequent connections validated automatically

---

## Practical Examples

### Example 1: Static Site Deployment

```parsley
// Build and deploy static site to SFTP server

// Build site locally
let pages <== dir(@./src/pages/)
let output = @./build/

for (page in pages) {
    let html = renderPage(page)
    html ==> text(output + page.basename)
}

// Deploy to SFTP
let conn = SFTP(@sftp://deploy@example.com, {
    keyFile: @~/.ssh/deploy_key
})

let {files} <== dir(output)
for (file in files) {
    let content <== text(file.path)
    let remotePath = @/var/www/html/ + file.basename
    content =/=> conn(remotePath).text
    log("Deployed:", file.basename)
}

conn.close()
log("✓ Deployment complete")
```

### Example 2: Backup to Remote Server

```parsley
// Backup local data to SFTP server

let conn = SFTP(@sftp://backup@backup-server.com, {
    keyFile: @~/.ssh/backup_key
})

// Create backup directory with timestamp
let backupDir = @/backups/ + now().format("YYYY-MM-DD")
conn(backupDir).mkdir({parents: true})

// Backup database export
let dbData <== JSON(@./local/database.json)
dbData =/=> conn(backupDir + "/database.json").json

// Backup config files
let config <== JSON(@./config/app.json)
config =/=> conn(backupDir + "/app.json").json

log("✓ Backup complete:", backupDir)
```

### Example 3: Download Remote Logs

```parsley
// Download and analyze remote logs via SFTP

let conn = SFTP(@sftp://logger@log-server.com, {
    keyFile: @~/.ssh/logger_key
})

// Read remote log file
let {data: logLines, error} <=/= conn(@/var/log/app.log).lines

if (error != null) {
    log("Failed to read logs:", error)
} else {
    // Analyze logs locally
    let errors = logLines.filter(fn(line) {
        line.contains("ERROR")
    })
    
    log("Found", errors.length(), "errors")
    
    // Save error summary locally
    errors ==> lines(@./logs/error-summary.txt)
}
```

### Example 4: Sync Configuration

```parsley
// Sync configuration between local and remote

let conn = SFTP(@sftp://admin@server.com, {
    keyFile: @~/.ssh/admin_key
})

// Read local config
let localConfig <== JSON(@./config/local.json)

// Read remote config
let {data: remoteConfig, error} <=/= conn(@/etc/app/config.json).json

if (error != null) {
    // No remote config, upload local
    localConfig =/=> conn(@/etc/app/config.json).json
    log("✓ Uploaded initial config")
} else {
    // Merge configs
    let merged = {
        ...remoteConfig,
        ...localConfig  // Local overrides remote
    }
    
    // Write back to remote
    merged =/=> conn(@/etc/app/config.json).json
    log("✓ Synced configuration")
}
```

---

## Implementation Checklist

### Phase 1: Core SFTP Infrastructure (v0.12.0) ✅

- [x] Add `github.com/pkg/sftp` and `golang.org/x/crypto/ssh` dependencies
- [x] Create `SFTPConnection` struct
- [x] Create `SFTPFileHandle` struct  
- [x] Implement `SFTP()` builtin factory
- [x] Implement connection caching
- [x] Implement SSH key authentication
- [x] Implement password authentication (with warnings)
- [x] Implement host key verification
- [x] Add SFTP connection object type
- [x] Add connection property accessors (`.host`, `.user`, `.connected`)

### Phase 2: File Operations (v0.12.0) ✅

- [x] Implement `conn(path)` call syntax returning SFTP file handle
- [x] Implement `.json` format accessor on SFTP file handles
- [x] Implement `.text` format accessor on SFTP file handles
- [x] Implement `.csv` format accessor on SFTP file handles
- [x] Implement `.lines` format accessor on SFTP file handles
- [x] Implement `.bytes` format accessor on SFTP file handles
- [x] Implement `.file` format accessor (auto-detect from extension)
- [x] Implement SFTP read operation via `<=/=` (network read)
- [x] Implement SFTP write operation via `=/=>` (network write)
- [x] Implement SFTP append operation via `=/=>>` (network append with SSH_FXF_APPEND)
- [x] Add error capture pattern support `{data, error} <=/=`
- [x] Add fallback support with `??`

### Phase 3: Directory Operations (v0.12.0) ✅

- [x] Implement `.dir` accessor for directory listing
- [x] Implement `.dir(options)` with recursive option
- [x] Implement `.mkdir()` method for directory creation
- [x] Implement `.mkdir(options)` with parents option
- [x] Implement `.rmdir()` method for directory removal
- [x] Implement `.rmdir(options)` with recursive option
- [x] Implement `.remove()` method for file deletion
- [x] Support error capture for directory operations
- [x] Add recursive directory listing
- [x] Add recursive directory creation (`parents: true`)
- [x] Return file info objects with metadata

### Phase 4: Testing (v0.12.0) ✅

- [x] Add unit tests for SFTP connection creation
- [x] Add unit tests for authentication methods
- [x] Add unit tests for read/write operations
- [x] Add unit tests for directory operations
- [x] Add integration tests with local SFTP server (syntax tests - real server optional)
- [x] Test error handling and edge cases
- [x] Test connection caching
- [x] Test security policy enforcement

### Phase 5: Documentation (v0.12.0) ✅

- [x] Add SFTP section to docs/reference.md
- [x] Add SFTP examples to README.md (covered in reference.md)
- [x] Create examples/sftp_demo.pars
- [x] Update CHANGELOG.md
- [x] Document authentication best practices
- [x] Add troubleshooting guide (covered in reference.md)

---

## Design Decisions (Resolved)

### 1. Connection Lifecycle ✅

**Decision**: Auto-close on script exit with optional explicit `.close()` for immediate cleanup.

**Rationale**: Simpler, prevents leaks, matches database connection pattern. Long-running scripts can call `.close()` when needed.

---

### 2. Glob Pattern Support ✅

**Decision**: Phase 2 feature — use `dir()` + filter for now.

**Rationale**: Implement glob locally after listing (simple). Can add dedicated glob method in future enhancement without breaking changes.

---

### 3. File Permissions ✅

**Decision**: Phase 2 feature — not needed for MVP.

**Rationale**: Most use cases don't need permission control. Can add `{mode: "644"}` option later if needed.

---

### 4. Append Operations ✅

**Decision**: ✅ **Yes, included in MVP.**

**Rationale**: SFTP protocol natively supports append via `SSH_FXF_APPEND` flag (0x00000004). github.com/pkg/sftp library fully supports this. Using `=/=>>` operator (network append) to match network I/O pattern.

**Example**: `logEntry =/=>> conn(@/var/log/app.log).lines`

---

### 5. Symbolic Links ✅

**Decision**: Follow symlinks by default, add option for control in Phase 2.

**Rationale**: Matches standard behavior. Option for `followSymlinks: false` can be added later if needed.

---

### 6. Error Recovery ✅

**Decision**: Phase 2 feature — explicit error handling for MVP.

**Rationale**: Auto-reconnect adds complexity. Users can implement retry logic themselves using error capture pattern. Can add `autoReconnect: true` option in future if there's demand.

---

## Alternative Designs Considered

### Alternative 1: Direct URL Fetch (like HTTP)

```parsley
// Not chosen: Stateless SFTP (inefficient)
let data <=/= JSON(@sftp://user@host.com/path/file.json)
```

**Rejected because**:
- Creates new SSH session per operation (very expensive)
- No way to reuse connections efficiently
- Authentication complexity per request
- Doesn't match SFTP's stateful nature

---

### Alternative 2: Unified Remote File API

```parsley
// Not chosen: Generic remote file abstraction
let remote = REMOTE(@sftp://user@host.com)
let data <== remote.JSON(@/path/file.json)

// Also works for HTTP
let http = REMOTE(@https://api.example.com)
let data <== http.JSON(@/data)
```

**Rejected because**:
- Overabstraction — HTTP and SFTP are fundamentally different
- Loses protocol-specific features
- Harder to extend independently
- Not aligned with Parsley's explicit design

---

### Alternative 3: Format Objects with Remote Support

```parsley
// Not chosen: Extend format factories for remote
let data <== JSON(@sftp://user@host.com/file.json, {
    auth: {keyFile: @~/.ssh/id_rsa}
})
```

**Rejected because**:
- Mixes format concerns with network protocol
- No connection reuse
- Duplicate authentication for every file
- Inconsistent with database pattern

---

## Migration from HTTP Fetch

Users familiar with HTTP fetch will need to understand:

| HTTP Fetch | SFTP |
|------------|------|
| `<=/=` operator directly | Connection + `<==` operator |
| Stateless requests | Stateful connection |
| URL in format factory | Connection method |
| No auth management | SSH keys/passwords |
**Example comparison**:

```parsley
// HTTP fetch (stateless, network I/O)
let data <=/= JSON(@https://api.example.com/data.json)

// SFTP (stateful connection, network I/O)
let conn = SFTP(@sftp://user@example.com, {keyFile: @~/.ssh/id_rsa})
let data <=/= conn(@/path/data.json).json
data =/=>> conn(@/path/log.txt).lines  // Append supported

// Local file (file I/O)
let data <== JSON(@./local/data.json)
data ==>> lines(@./log.txt)  // Append to local file
```

---

## Future Enhancements (v0.13.0+)

### Connection Pooling

```parsley
// Multiple concurrent SFTP operations
let pool = SFTPPool(@sftp://user@host.com, {
    keyFile: @~/.ssh/id_rsa,
    maxConnections: 5
})
```

### Async Operations (Phase 2/v0.13.0)

```parsley
// Non-blocking SFTP operations
let upload = async(data =/=> conn(@/remote/big-file.json).json)
let download = async(data <=/= conn(@/remote/other.json).json)

await(upload)
await(download)
```

### Progress Callbacks (Phase 2/v0.13.0)

```parsley
// Large file upload with progress
let onProgress = fn(bytes, total) {
    log("Uploaded:", bytes, "/", total)
}

bigData =/=> conn(@/remote/large-file.bin).bytes({onProgress})
```

### SFTP Server Support (Future/v0.14.0+)

```parsley
// Run SFTP server from Parsley
let server = SFTPServer({
    port: 2222,
    host: "0.0.0.0",
    keyFile: @./server_key,
    root: @./sftp-root/,
    onConnect: fn(user) { log("Connected:", user) },
    onDisconnect: fn(user) { log("Disconnected:", user) }
})

server.start()
```

---

## Summary

### What We're Building (v0.12.0)

- `SFTP(url, options?)` builtin creating stateful connection
- Connection callable: `conn(@/path)` returns SFTP file handle
- Format accessors (lowercase): `.json`, `.text`, `.csv`, `.lines`, `.bytes`, `.file`
- Directory accessor: `.dir`, `.dir(options)` for listing
- Network read operator: `data <=/= conn(@/path).json`
- Network write operator: `data =/=> conn(@/path).json`
- Network append operator: `data =/=>> conn(@/path).lines` (SSH_FXF_APPEND)
- Directory operations: `.mkdir()`, `.rmdir()`, `.remove()` methods on handles
- SSH key + password authentication
- Host key verification
- Connection caching and reuse
- No special security flags (follows HTTP/Fetch pattern)
- Error capture: `{data, error} <=/= conn(@/path).json`

### Why This Design

✅ **Network operators signal network I/O** — `<=/=`, `=/=>`, `=/=>>` vs local `<==`, `==>`, `==>>`  
✅ **Path-first syntax matches local files** — `conn(@/path).json` mirrors `file(@/path).json`  
✅ **Connection pattern matches databases** — Stateful, cached, expensive to create  
✅ **Efficient connection reuse** — SSH sessions cached by host+user+port  
✅ **Append support via protocol** — `SSH_FXF_APPEND` (0x00000004) fully supported  
✅ **Works with all format accessors** — `.json`, `.text`, `.csv`, `.lines`, `.bytes`  
✅ **Clear authentication model** — SSH keys, passwords, host verification  
✅ **No network restrictions** — Like HTTP/Fetch, SFTP unrestricted (only local file security)  
✅ **Minimal API surface** — One factory, network operators, format accessors  
✅ **Composable and consistent** — Integrates with existing I/O patterns  

### Why Connection-Based (Not URL-Based Like HTTP)

SFTP differs fundamentally from HTTP:
- **Stateful protocol** vs stateless requests
- **SSH session lifecycle** vs ephemeral connections  
- **Expensive to establish** vs cheap HTTP requests
- **Session-based auth** vs per-request headers
- **Directory operations** require maintained state

Connection objects match SFTP's nature; stateless fetch does not.

### Why Network Operators (Not File Operators)

SFTP is network I/O, not local file I/O:
- **Data crosses network** → latency, bandwidth, failures
- **Clear distinction** from local operations in code
- **Matches HTTP pattern** → `<=/=` for network reads
- **Enables append** → `=/=>>` parallel to `==>>` for local files

---

## Dependencies

### Go Libraries

```go
import (
    "github.com/pkg/sftp"           // SFTP protocol implementation
    "golang.org/x/crypto/ssh"        // SSH client
    "golang.org/x/crypto/ssh/knownhosts"  // Host key verification
)
```

### License Compatibility

- `github.com/pkg/sftp`: BSD-2-Clause (✅ Compatible with MIT)
- `golang.org/x/crypto/ssh`: BSD-3-Clause (✅ Compatible with MIT)

---

## Timeline

**v0.12.0 (Target: ~2 weeks)**
- Core SFTP connection and file operations
- SSH key authentication
- Basic directory operations
- Security integration
- Documentation and examples

**v0.13.0 (Future)**
- Connection pooling
- Advanced directory operations (glob, recursive)
- File permissions support
- Append operations
- Performance optimizations

---

## Conclusion

This design provides SFTP support that:
1. Feels natural to Parsley users
2. Handles SFTP's stateful nature correctly
3. Integrates with existing patterns (connections, operators, formats)
4. Maintains security and error handling conventions
5. Leaves room for future enhancements

The connection-based approach is the right choice for SFTP, even though it differs from HTTP fetch, because it matches the protocol's fundamental characteristics.

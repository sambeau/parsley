# Plan: File System Security (Chroots)

**TL;DR**: Add command-line flags to restrict file system access to specific directories, preventing scripts from reading/writing outside allowed paths. Similar in purpose to Unix chroot but implemented via path validation rather than OS-level isolation.

---

## Design Principles

1. **Secure by Default** — Write and execute operations are restricted by default
2. **Explicit Allow Lists** — Scripts must explicitly be granted write/execute access
3. **Read Freedom** — Read operations unrestricted by default, with opt-in blacklisting
4. **Minimal Performance Impact** — Path validation should be fast and efficient
5. **Clear Error Messages** — Users should understand why access was denied
6. **Composable Restrictions** — Read, write, and execute can be controlled independently

---

## Command-Line API

### Flag Specifications

| Flag | Shorthand | Description | Example |
|------|-----------|-------------|---------|
| `--restrict-read=PATHS` | | Deny reading from comma-separated paths | `--restrict-read=/etc,/var` |
| `--no-read` | | Deny all file reads | `--no-read` |
| `--allow-write=PATHS` | | Allow writing to comma-separated paths | `--allow-write=./output,./cache` |
| `--allow-write-all` | `-w` | Allow unrestricted writes | `-w` |
| `--allow-execute=PATHS` | | Allow executing scripts from paths | `--allow-execute=./bin,./tools` |
| `--allow-execute-all` | `-x` | Allow unrestricted script execution | `-x` |

### Basic Usage

```bash
# Default: unrestricted reads, NO writes, NO executes
./pars script.pars

# Allow writes to specific directory
./pars --allow-write=./output script.pars

# Allow writes anywhere (shorthand)
./pars -w script.pars

# Allow execute anywhere (shorthand)
./pars -x script.pars

# Restrict reads from sensitive directories
./pars --restrict-read=/etc,/var script.pars

# Deny all reads (for stdin-only scripts)
./pars --no-read < data.json

# Combined: allow all writes and executes
./pars -w -x script.pars

# Specific write paths with unrestricted execute
./pars --allow-write=./output -x script.pars
```

### Path Resolution

All paths in allow/restrict lists are:
- Resolved to absolute paths at startup
- Cleaned using Rob Pike's cleanname algorithm
- Allow/deny access to the directory and all subdirectories
- Support `~` for home directory expansion

```bash
# These are equivalent
./pars --allow-write=./output script.pars
./pars --allow-write=$(pwd)/output script.pars

# Home directory expansion
./pars --allow-write=~/Documents/output script.pars
```

### Behavior Matrix

| Scenario | Read | Write | Execute |
|----------|------|-------|---------|
| Default | ✅ All allowed | ❌ All denied | ❌ All denied |
| `--allow-write=./output` | ✅ All allowed | ✅ `./output/*` only | ❌ All denied |
| `-w` | ✅ All allowed | ✅ All allowed | ❌ All denied |
| `--restrict-read=/etc` | ✅ Except `/etc/*` | ❌ All denied | ❌ All denied |
| `--no-read` | ❌ All denied | ❌ All denied | ❌ All denied |
| `--allow-execute=./bin` | ✅ All allowed | ❌ All denied | ✅ `./bin/*` only |
| `-x` | ✅ All allowed | ❌ All denied | ✅ All allowed |

---

## Implementation Strategy

### Phase 1: Core Infrastructure

**Add Security Context to Environment**

```go
type SecurityPolicy struct {
    RestrictRead    []string  // Denied read directories (blacklist)
    NoRead          bool      // Deny all reads
    AllowWrite      []string  // Allowed write directories (whitelist)
    AllowWriteAll   bool      // Allow all writes
    AllowExecute    []string  // Allowed execute directories (whitelist)
    AllowExecuteAll bool      // Allow all executes
}

type Environment struct {
    // ... existing fields ...
    Security *SecurityPolicy
}
```

**Path Validation Function**

```go
func (env *Environment) checkPathAccess(path string, operation string) error {
    if env.Security == nil {
        // No policy = default behavior
        // Read: allowed
        // Write: denied
        // Execute: denied
        if operation == "write" {
            return fmt.Errorf("write access denied (use --allow-write or -w)")
        }
        if operation == "execute" {
            return fmt.Errorf("execute access denied (use --allow-execute or -x)")
        }
        return nil
    }
    
    // Convert to absolute path
    absPath, err := filepath.Abs(path)
    if err != nil {
        return fmt.Errorf("invalid path: %s", err)
    }
    absPath = filepath.Clean(absPath)
    
    switch operation {
    case "read":
        if env.Security.NoRead {
            return fmt.Errorf("file read access denied: %s", path)
        }
        // Check blacklist
        if isPathRestricted(absPath, env.Security.RestrictRead) {
            return fmt.Errorf("file read restricted: %s", path)
        }
        
    case "write":
        if env.Security.AllowWriteAll {
            return nil // Unrestricted
        }
        if !isPathAllowed(absPath, env.Security.AllowWrite) {
            return fmt.Errorf("file write not allowed: %s (use --allow-write or -w)", path)
        }
        
    case "execute":
        if env.Security.AllowExecuteAll {
            return nil // Unrestricted
        }
        if !isPathAllowed(absPath, env.Security.AllowExecute) {
            return fmt.Errorf("script execution not allowed: %s (use --allow-execute or -x)", path)
        }
    }
    
    return nil
}

func isPathAllowed(path string, allowList []string) bool {
    // Empty allow list means nothing is allowed
    if len(allowList) == 0 {
        return false
    }
    
    // Check if path is within any allowed directory
    for _, allowed := range allowList {
        if path == allowed || strings.HasPrefix(path, allowed + string(filepath.Separator)) {
            return true
        }
    }
    
    return false
}

func isPathRestricted(path string, restrictList []string) bool {
    // Empty restrict list = no restrictions
    if len(restrictList) == 0 {
        return false
    }
    
    // Check if path is within any restricted directory
    for _, restricted := range restrictList {
        if path == restricted || strings.HasPrefix(path, restricted + string(filepath.Separator)) {
            return true
        }
    }
    
    return false
}
```

### Phase 2: Integration Points

**1. File Reading (`readFileContent`)**

```go
func readFileContent(fileDict *Dictionary, env *Environment) (Object, *Error) {
    pathStr := getFilePathString(fileDict, env)
    
    // Security check
    if err := env.checkPathAccess(pathStr, "read"); err != nil {
        return nil, newError("security: %s", err.Error())
    }
    
    // ... existing read logic ...
}
```

**2. File Writing (`writeFileContent`)**

```go
func writeFileContent(fileDict *Dictionary, value Object, appendMode bool, env *Environment) *Error {
    pathStr := getFilePathString(fileDict, env)
    
    // Security check
    if err := env.checkPathAccess(pathStr, "write"); err != nil {
        return newError("security: %s", err.Error())
    }
    
    // ... existing write logic ...
}
```

**3. File Deletion (`evalFileRemove`)**

```go
func evalFileRemove(fileDict *Dictionary, env *Environment) Object {
    pathStr := getFilePathString(fileDict, env)
    
    // Security check (treat as write operation)
    if err := env.checkPathAccess(pathStr, "write"); err != nil {
        return newError("security: %s", err.Error())
    }
    
    // ... existing delete logic ...
}
```

**4. Directory Listing (`dir`, `files`)**

```go
// In dir() and files() builtins
if err := env.checkPathAccess(pathStr, "read"); err != nil {
    return newError("security: %s", err.Error())
}
```

**5. Module Imports**

```go
// When importing modules
if err := env.checkPathAccess(modulePath, "execute"); err != nil {
    return newError("security: %s", err.Error())
}
```

### Phase 3: CLI Flag Parsing

**Update main.go**

```go
import (
    "flag"
    "strings"
)

var (
    restrictReadFlag     = flag.String("restrict-read", "", "Comma-separated read blacklist paths")
    noReadFlag           = flag.Bool("no-read", false, "Deny all file reads")
    allowWriteFlag       = flag.String("allow-write", "", "Comma-separated write whitelist paths")
    allowWriteAllFlag    = flag.Bool("allow-write-all", false, "Allow unrestricted writes")
    allowWriteAllShort   = flag.Bool("w", false, "Shorthand for --allow-write-all")
    allowExecuteFlag     = flag.String("allow-execute", "", "Comma-separated execute whitelist paths")
    allowExecuteAllFlag  = flag.Bool("allow-execute-all", false, "Allow unrestricted executes")
    allowExecuteAllShort = flag.Bool("x", false, "Shorthand for --allow-execute-all")
)

func buildSecurityPolicy() (*evaluator.SecurityPolicy, error) {
    policy := &evaluator.SecurityPolicy{
        NoRead:          *noReadFlag,
        AllowWriteAll:   *allowWriteAllFlag || *allowWriteAllShort,
        AllowExecuteAll: *allowExecuteAllFlag || *allowExecuteAllShort,
    }
    
    // Parse restrict list
    if *restrictReadFlag != "" {
        paths, err := parseAndResolvePaths(*restrictReadFlag)
        if err != nil {
            return nil, fmt.Errorf("invalid --restrict-read: %s", err)
        }
        policy.RestrictRead = paths
    }
    
    // Parse allow lists
    if *allowWriteFlag != "" {
        paths, err := parseAndResolvePaths(*allowWriteFlag)
        if err != nil {
            return nil, fmt.Errorf("invalid --allow-write: %s", err)
        }
        policy.AllowWrite = paths
    }
    
    if *allowExecuteFlag != "" {
        paths, err := parseAndResolvePaths(*allowExecuteFlag)
        if err != nil {
            return nil, fmt.Errorf("invalid --allow-execute: %s", err)
        }
        policy.AllowExecute = paths
    }
    
    return policy, nil
}

func parseAndResolvePaths(pathList string) ([]string, error) {
    parts := strings.Split(pathList, ",")
    resolved := make([]string, 0, len(parts))
    
    for _, p := range parts {
        p = strings.TrimSpace(p)
        if p == "" {
            continue
        }
        
        // Expand home directory
        if strings.HasPrefix(p, "~/") {
            home, err := os.UserHomeDir()
            if err != nil {
                return nil, fmt.Errorf("cannot expand ~: %s", err)
            }
            p = filepath.Join(home, p[2:])
        }
        
        // Convert to absolute path
        absPath, err := filepath.Abs(p)
        if err != nil {
            return nil, fmt.Errorf("invalid path %s: %s", p, err)
        }
        
        // Clean path
        absPath = filepath.Clean(absPath)
        
        resolved = append(resolved, absPath)
    }
    
    return resolved, nil
}

func main() {
    flag.Parse()
    
    // Build security policy (always create one to enable default restrictions)
    policy, err := buildSecurityPolicy()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %s\n", err)
        os.Exit(1)
    }
    
    // Create environment with security policy
    env := evaluator.NewEnvironment()
    env.Security = policy
    
    // ... rest of main ...
}
```

---

## Error Messages

Security errors should be clear and actionable:

```
Error: security: write access denied (use --allow-write or -w)

Error: security: file read restricted: /etc/passwd
Hint: Use --restrict-read to modify read restrictions

Error: security: file write not allowed: ./output/result.json (use --allow-write or -w)
Hint: Use --allow-write=./output to allow writing to ./output directory

Error: security: script execution not allowed: ../tools/converter.pars (use --allow-execute or -x)
Hint: Use --allow-execute=../tools to allow executing scripts from ../tools directory
```

---

## Use Cases

### Static Site Generator (Read Freely, Write to Output)

```bash
./pars --allow-write=./public build.pars
```

### API Data Processor (Restrict Sensitive Reads, Write Results)

```bash
./pars --restrict-read=/etc --allow-write=./output process.pars
```

### Template Renderer (No File Access)

```bash
./pars --no-read < data.json
```

### Development Mode (Unrestricted)

```bash
./pars -w -x dev-script.pars
```

### Script with External Tools

```bash
./pars -x --allow-write=./output process.pars
```

---

## Testing Strategy

### Unit Tests

```go
func TestSecurityPolicyReadRestriction(t *testing.T) {
    policy := &SecurityPolicy{
        RestrictRead: []string{"/tmp/restricted"},
    }
    
    env := NewEnvironment()
    env.Security = policy
    
    // Should deny
    err := env.checkPathAccess("/tmp/restricted/file.txt", "read")
    if err == nil {
        t.Errorf("Should deny read from /tmp/restricted/file.txt")
    }
    
    // Should allow
    err = env.checkPathAccess("/tmp/allowed/file.txt", "read")
    if err != nil {
        t.Errorf("Should allow read from /tmp/allowed/file.txt")
    }
}

func TestSecurityPolicyWriteDefault(t *testing.T) {
    policy := &SecurityPolicy{}
    
    env := NewEnvironment()
    env.Security = policy
    
    // Should deny by default
    err := env.checkPathAccess("/tmp/file.txt", "write")
    if err == nil {
        t.Errorf("Should deny write by default")
    }
}

func TestSecurityPolicyWriteAllowed(t *testing.T) {
    policy := &SecurityPolicy{
        AllowWrite: []string{"/tmp/output"},
    }
    
    env := NewEnvironment()
    env.Security = policy
    
    // Should allow
    err := env.checkPathAccess("/tmp/output/file.txt", "write")
    if err != nil {
        t.Errorf("Should allow write to /tmp/output/file.txt")
    }
    
    // Should deny
    err = env.checkPathAccess("/tmp/other/file.txt", "write")
    if err == nil {
        t.Errorf("Should deny write to /tmp/other/file.txt")
    }
}
```

### Integration Tests

Create test scripts that attempt to access files and verify they're blocked:

```parsley
// test-security-write.pars
// Should fail when run without --allow-write

"test" ==> text(@./output/test.txt)  // Should error
```

```bash
./pars test-security-write.pars
# Expected: Error: security: write access denied (use --allow-write or -w)

./pars --allow-write=./output test-security-write.pars
# Expected: Success
```

---

## Migration Path

### Current State (v0.9.18)
- No security restrictions
- All file operations allowed

### Phase 1 (v0.10.0) - Add Security Framework
- Add `SecurityPolicy` struct and path validation
- Add command-line flags
- **Default behavior change**: Write and execute now denied by default
- Add `--allow-write`, `-w`, `--allow-execute`, `-x` flags
- Add `--restrict-read`, `--no-read` flags
- Document security features with migration guide

### Phase 2 (v0.11.0) - Refine
- Add better error messages with hints
- Add config file support for security policies
- Add environment variable support (`PARSLEY_ALLOW_WRITE`, etc.)
- Add security audit logging (optional flag)

### Phase 3 (v1.0.0+) - Additional Features
- Network access restrictions
- Resource limits (memory, file size, execution time)
- More sophisticated path matching (globs, regex)

---

## Configuration File Support (Future)

Allow security policies in `.parsleyrc`:

```json
{
  "security": {
    "restrictRead": ["/etc", "/var"],
    "allowWrite": ["./output", "./cache"],
    "allowExecute": ["./bin"]
  }
}
```

Or environment variables:

```bash
export PARSLEY_RESTRICT_READ="/etc,/var"
export PARSLEY_ALLOW_WRITE="./output"
./pars script.pars
```

---

## Implementation Checklist

### Phase 1: Core Infrastructure (v0.10.0)
- [ ] Add `SecurityPolicy` struct to evaluator
- [ ] Add `Security` field to `Environment`
- [ ] Implement `checkPathAccess()` function
- [ ] Implement `isPathAllowed()` helper
- [ ] Implement `isPathRestricted()` helper
- [ ] Add security checks to `readFileContent()`
- [ ] Add security checks to `writeFileContent()`
- [ ] Add security checks to `evalFileRemove()`
- [ ] Add security checks to `dir()` and `files()` builtins
- [ ] Add security checks to module imports

### Phase 2: CLI Integration (v0.10.0)
- [ ] Add command-line flags to `main.go`
- [ ] Implement `buildSecurityPolicy()` function
- [ ] Implement `parseAndResolvePaths()` function
- [ ] Initialize environment with security policy
- [ ] Handle shorthand flags properly

### Phase 3: Testing (v0.10.0)
- [ ] Add unit tests for path validation (read, write, execute)
- [ ] Add unit tests for default deny behavior
- [ ] Add integration tests for read restrictions
- [ ] Add integration tests for write allowances
- [ ] Add integration tests for execute allowances
- [ ] Test edge cases (symlinks, relative paths, etc.)
- [ ] Test shorthand flags

### Phase 4: Documentation (v0.10.0)
- [ ] Update README with security examples
- [ ] Add security section to reference.md
- [ ] Create migration guide for v0.10.0 breaking change
- [ ] Create security best practices guide
- [ ] Add examples for common use cases
- [ ] Update CHANGELOG with breaking change notice

### Phase 5: Future Enhancements (v0.11.0+)
- [ ] Config file support (`.parsleyrc`)
- [ ] Environment variable support
- [ ] Improved error messages with suggestions
- [ ] Security audit logging
- [ ] Symlink handling policy
- [ ] Glob pattern support for paths

---

## Security Considerations

### What This Protects Against

✅ Accidental file writes outside intended directories  
✅ Scripts writing to system directories  
✅ Untrusted scripts accessing user data for writing  
✅ Scripts reading from explicitly restricted paths  
✅ Unauthorized execution of external scripts  

### What This Does NOT Protect Against

❌ Malicious scripts with access to allowed directories  
❌ Resource exhaustion (CPU, memory, disk space)  
❌ Network access (not yet implemented)  
❌ Process spawning (when that feature is added)  
❌ Symlink attacks (basic validation only)  

### Additional Security Features Needed

- **Network restrictions**: `--allow-net=example.com` for HTTP/database access
- **Resource limits**: `--max-memory`, `--max-file-size`, `--timeout`
- **Process isolation**: Consider running scripts in separate processes
- **Audit logging**: Log all file operations when security is enabled

---

## Breaking Changes in v0.10.0

### Default Behavior Change

**Before v0.10.0:**
- All file operations unrestricted

**After v0.10.0:**
- **Reads**: Unrestricted (no change)
- **Writes**: Denied by default (breaking change)
- **Executes**: Denied by default (breaking change)

### Migration Guide

Scripts that write files will need to:

1. Add `--allow-write=PATHS` for specific directories
2. Or add `-w` for unrestricted writes (old behavior)

```bash
# Old (v0.9.x)
./pars build-site.pars

# New (v0.10.0+) - specific directory
./pars --allow-write=./public build-site.pars

# New (v0.10.0+) - unrestricted (old behavior)
./pars -w build-site.pars
```

Scripts that import modules will need:

1. Add `--allow-execute=PATHS` for module directories
2. Or add `-x` for unrestricted execution (old behavior)

```bash
# Old (v0.9.x)
./pars app.pars

# New (v0.10.0+) - specific directory
./pars --allow-execute=./lib app.pars

# New (v0.10.0+) - unrestricted (old behavior)
./pars -x app.pars
```

---

## Open Questions

1. **Should we allow multiple allow lists?**
   - Currently: `--allow-write=./data,./config`
   - Alternative: `--allow-write=./data --allow-write=./config`
   - **Decision**: Support both for flexibility

2. **How to handle symlinks?**
   - Follow symlinks and check target path?
   - Or: deny all symlinks in restricted mode?
   - **Decision**: Follow symlinks but validate target path (Phase 1)

3. **Should path restrictions apply to the script itself?**
   - Should `pars script.pars` allow reading script.pars even with `--no-read`?
   - **Decision**: Yes, automatically allow reading the script being executed

4. **Config file precedence?**
   - CLI flags override config file?
   - Or: merge them?
   - **Decision**: CLI flags override config file (Phase 2)

5. **Should we make write restrictions opt-out instead of opt-in?**
   - Current: Denied by default, opt-in with flags
   - Alternative: Allowed by default, opt-out with `--restrict-write`
   - **Decision**: Keep opt-in (more secure)

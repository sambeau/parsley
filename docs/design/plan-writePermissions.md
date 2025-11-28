# Plan: Write Permission Control for Parsley

**TL;DR**: Add `--allow-write` CLI flag to control which directories/files Parsley can write to. By default, writing is unrestricted (current behavior), but users can opt-in to restricted mode for safer execution of untrusted templates.

---

## Motivation

While reading files is generally safe, writing files can be dangerous:
- Overwriting important configuration files
- Creating files in system directories
- Filling up disk space
- Modifying files outside the intended output directory

A `--allow-write` flag gives users control over write permissions when running templates from untrusted sources or in production environments.

---

## Design Principles

1. **Opt-in restriction** — Default behavior remains unrestricted (backward compatible)
2. **Simple mental model** — If you specify `--allow-write`, only those paths are writable
3. **Fail clearly** — Denied writes produce clear error messages
4. **Path-based** — Permissions are directory/file based, not format-based

---

## CLI Interface

### Default (unrestricted)
```bash
# Current behavior - can write anywhere
pars template.pars
```

### Restricted Mode
```bash
# Only allow writes to ./output directory
pars --allow-write=./output template.pars

# Allow writes to multiple directories
pars --allow-write=./output,./cache,./logs template.pars

# Allow writes to specific file only
pars --allow-write=./config.json template.pars

# No writes allowed at all
pars --allow-write= template.pars
# or
pars --deny-write template.pars
```

### Path Resolution

- Relative paths are resolved relative to current working directory
- Paths are normalized (no `..` escapes)
- Symbolic links are resolved before checking

```bash
# These are equivalent if CWD is /home/user/project
pars --allow-write=./output template.pars
pars --allow-write=/home/user/project/output template.pars
```

---

## Behavior

### When `--allow-write` is NOT specified
- All writes are permitted (current behavior)
- No permission checking overhead

### When `--allow-write` is specified
- Only paths under the specified directories/files are writable
- Attempting to write elsewhere returns an error (not a runtime panic)
- The error can be captured with `{data, error}` pattern

### Error Messages

```parsley
// Attempting to write to /etc/passwd with --allow-write=./output
data ==> JSON(@/etc/passwd)
// ERROR: write permission denied: '/etc/passwd' is not within allowed paths [./output]

// With error capture
let {data, error} = data ==> JSON(@/etc/passwd)
// error = "write permission denied: '/etc/passwd' is not within allowed paths [./output]"
```

---

## Implementation Plan

### Phase 1: Environment Flag Storage

Add write permission state to the evaluator environment:

```go
// In Environment struct
type Environment struct {
    // ... existing fields ...
    
    // Write permissions (nil = unrestricted, empty = deny all, populated = allow listed)
    AllowedWritePaths []string
}

// Helper method
func (e *Environment) CanWriteTo(path string) bool {
    if e.AllowedWritePaths == nil {
        return true // Unrestricted mode
    }
    if len(e.AllowedWritePaths) == 0 {
        return false // Deny all mode
    }
    // Check if path is under any allowed path
    absPath, _ := filepath.Abs(path)
    for _, allowed := range e.AllowedWritePaths {
        absAllowed, _ := filepath.Abs(allowed)
        if strings.HasPrefix(absPath, absAllowed) {
            return true
        }
    }
    return false
}
```

### Phase 2: CLI Flag Parsing

Update `main.go` to parse the flag:

```go
var allowWrite = flag.String("allow-write", "", "Comma-separated paths where writing is allowed (empty = deny all, unset = allow all)")
var denyWrite = flag.Bool("deny-write", false, "Deny all file writes")

// In main()
if *denyWrite {
    env.AllowedWritePaths = []string{} // Empty slice = deny all
} else if *allowWrite != "" {
    paths := strings.Split(*allowWrite, ",")
    for _, p := range paths {
        p = strings.TrimSpace(p)
        if p != "" {
            absPath, err := filepath.Abs(p)
            if err != nil {
                log.Fatalf("Invalid path in --allow-write: %s", p)
            }
            env.AllowedWritePaths = append(env.AllowedWritePaths, absPath)
        }
    }
    if len(env.AllowedWritePaths) == 0 {
        env.AllowedWritePaths = []string{} // Explicit empty = deny all
    }
} else {
    env.AllowedWritePaths = nil // nil = unrestricted
}
```

### Phase 3: Write Permission Checks

Add permission check to `evalWriteStatement`:

```go
func evalWriteStatement(node *ast.WriteStatement, env *Environment) Object {
    // ... existing code to get file path ...
    
    pathStr := getFilePathString(fileDict, env)
    
    // Check write permission
    if !env.CanWriteTo(pathStr) {
        return newError("write permission denied: '%s' is not within allowed paths %v", 
            pathStr, env.AllowedWritePaths)
    }
    
    // ... rest of existing write logic ...
}
```

### Phase 4: Directory Creation Check

Also check permissions when creating directories (if we add that feature):

```go
func evalMkdir(path string, env *Environment) Object {
    if !env.CanWriteTo(path) {
        return newError("write permission denied: cannot create directory '%s'", path)
    }
    // ... create directory ...
}
```

---

## Security Considerations

### Path Traversal Prevention

```go
func (e *Environment) CanWriteTo(path string) bool {
    // Resolve to absolute path
    absPath, err := filepath.Abs(path)
    if err != nil {
        return false
    }
    
    // Resolve any symlinks to prevent escaping via symlink
    realPath, err := filepath.EvalSymlinks(absPath)
    if err != nil {
        // File doesn't exist yet - check parent directory
        parentDir := filepath.Dir(absPath)
        realPath, err = filepath.EvalSymlinks(parentDir)
        if err != nil {
            return false
        }
        realPath = filepath.Join(realPath, filepath.Base(absPath))
    }
    
    // Clean the path to remove any .. components
    realPath = filepath.Clean(realPath)
    
    // Check against allowed paths
    for _, allowed := range e.AllowedWritePaths {
        allowedReal, _ := filepath.EvalSymlinks(allowed)
        allowedReal = filepath.Clean(allowedReal)
        
        if realPath == allowedReal || strings.HasPrefix(realPath, allowedReal+string(filepath.Separator)) {
            return true
        }
    }
    return false
}
```

### Race Conditions

- TOCTOU (time-of-check-time-of-use) attacks are possible but mitigated by:
  - Checking permissions immediately before write
  - Using absolute resolved paths
  - This is defense-in-depth, not a security sandbox

### Limitations

This is **not** a security sandbox. It's a safety feature to prevent accidental writes. A determined attacker with code execution could potentially bypass it. For true sandboxing, use OS-level isolation (containers, VMs, etc.).

---

## Testing Plan

### Unit Tests

```go
func TestWritePermissions(t *testing.T) {
    tests := []struct {
        name        string
        allowPaths  []string
        writePath   string
        shouldAllow bool
    }{
        {"nil allows all", nil, "/any/path", true},
        {"empty denies all", []string{}, "/any/path", false},
        {"exact match", []string{"/tmp/out"}, "/tmp/out/file.txt", true},
        {"outside allowed", []string{"/tmp/out"}, "/etc/passwd", false},
        {"parent traversal blocked", []string{"/tmp/out"}, "/tmp/out/../etc/passwd", false},
    }
    // ...
}
```

### Integration Tests

```go
func TestWritePermissionCLI(t *testing.T) {
    // Test with --allow-write flag
    // Test with --deny-write flag
    // Test error messages
}
```

---

## Documentation

### README Addition

```markdown
## Write Permissions

By default, Parsley can write to any location. For safer execution, use `--allow-write`:

\`\`\`bash
# Only allow writes to the output directory
pars --allow-write=./output template.pars

# Allow multiple directories
pars --allow-write=./output,./logs template.pars

# Deny all writes
pars --deny-write template.pars
\`\`\`

Attempting to write outside allowed paths produces an error that can be captured:

\`\`\`parsley
let {data, error} = myData ==> JSON(@/etc/passwd)
if (error) {
    // error = "write permission denied: ..."
}
\`\`\`
```

---

## Open Questions

1. **Should there be a config file option?**
   ```yaml
   # .parsleyrc
   allow-write:
     - ./output
     - ./cache
   ```
   *Recommendation: Not initially. CLI flags are simpler and more explicit.*

2. **Should append (`==>>`) have separate permissions?**
   *Recommendation: No. Append is a form of write, same permissions apply.*

3. **What about file deletion (if added later)?**
   *Recommendation: Same `--allow-write` controls deletion. Writing implies destructive operations.*

4. **Environment variable alternative?**
   ```bash
   PARSLEY_ALLOW_WRITE=./output pars template.pars
   ```
   *Recommendation: Yes, add as alternative. Useful for CI/CD.*

---

## Implementation Priority

This feature is **not urgent** but would be valuable for:
- Running templates in CI/CD pipelines
- Processing user-uploaded templates
- Production deployments where accidental writes could be costly

Estimated effort: ~2-3 hours for basic implementation, ~1 day with full testing.

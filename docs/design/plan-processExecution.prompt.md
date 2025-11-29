# Plan: External Process Execution

**TL;DR**: Add `COMMAND()` factory creating command handle objects, use `<=#=>` operator for execution (matching database mutation pattern), return result dictionary with stdout/stderr/exit code. Leverage existing `--allow-execute` security flag. Pattern mirrors database and file I/O APIs for consistency.

---

## Design Principles

1. **Commands are objects** — Like files and databases, commands are special dictionaries
2. **Execution is explicit** — The `<=#=>` operator is when execution happens
3. **Safe by default** — No shell interpretation, explicit arguments, security checks
4. **Results are values** — Return `{stdout, stderr, exitCode, error}` dict
5. **Composable** — Command handles can be stored, passed, inspected
6. **Synchronous** — Pragmatic blocking execution (async is future consideration)

---

## Command Handle Objects

### Creating Command Handles

```parsley
// Basic execution
let cmd = COMMAND("ls", ["-la", "/tmp"])

// Simple command (no arguments)
let cmd = COMMAND("pwd")

// With options
let cmd = COMMAND("cat", [], {
    env: {PATH: "/usr/bin"},
    dir: @./workspace,
    timeout: @5s
})
```

### Command Handle Structure

```parsley
{
    __type: "command",
    binary: "ls",           // String: executable path/name
    args: ["-la", "/tmp"],  // Array: arguments (default: [])
    options: {              // Dict: execution options
        env: null,          // Dict or null (inherits parent if null)
        dir: null,          // Path or null (uses current if null)
        timeout: null       // Duration or null (no timeout if null)
    }
}
```

### Command Handle Properties

All properties are eager (no filesystem access needed):

| Property | Type | Description |
|----------|------|-------------|
| `cmd.binary` | String | Executable name or path |
| `cmd.args` | Array | Argument list |
| `cmd.options` | Dict | Execution options |
| `cmd.options.env` | Dict/Null | Environment variables |
| `cmd.options.dir` | Path/Null | Working directory |
| `cmd.options.timeout` | Duration/Null | Execution timeout |

```parsley
let cmd = COMMAND("ls", ["-la"])

cmd.binary     // "ls"
cmd.args       // ["-la"]
cmd.options    // {env: null, dir: null, timeout: null}
```

---

## Execution Operator: `<=#=>`

### Pattern Consistency

Following Parsley's existing patterns:

| Feature | File I/O | Database | Process Execution |
|---------|----------|----------|-------------------|
| **Factory** | `JSON(@./file)` | `SQLITE(@./db)` | `COMMAND("ls", ["-la"])` |
| **Handle Type** | `{__type: "file"}` | `{__type: "connection"}` | `{__type: "command"}` |
| **Operator** | `<==` | `<=!=>` | `<=#=>` |
| **Usage** | `data <== file` | `result = db <=!=> query` | `result = cmd <=#=> input` |
| **Security** | `--allow-write` | N/A | `--allow-execute` |

### Execution Semantics

```parsley
// Execute with input
let result = COMMAND("cat") <=#=> "hello world"

// Execute without input (use null)
let result = COMMAND("ls", ["-la"]) <=#=> null

// result structure:
// {
//     stdout: "hello world",
//     stderr: "",
//     exitCode: 0,
//     error: null
// }
```

**Execution Steps:**

1. **Security check** — Validates binary path against `--allow-execute` policy
2. **Start process** — Launches with provided args/options
3. **Send input** — Pipes input to stdin (if not null)
4. **Wait synchronously** — Blocks until completion or timeout
5. **Capture output** — Buffers stdout and stderr separately
6. **Return result** — Dict with all execution information

---

## Result Object

### Result Structure

```parsley
{
    stdout: "...",      // String: standard output (empty string if none)
    stderr: "...",      // String: standard error (empty string if none)
    exitCode: 0,        // Integer: process exit code
    error: null         // String or null: error message if execution failed
}
```

### Result Patterns

**Destructure what you need:**

```parsley
// All fields
let {stdout, stderr, exitCode, error} = COMMAND("ls") <=#=> null

// Just output
let {stdout} = COMMAND("echo", ["hello"]) <=#=> null

// Just exit code
let {exitCode} = COMMAND("test", ["-f", "file.txt"]) <=#=> null

// Check for errors
let {error} = COMMAND("nonexistent") <=#=> null
if (error) {
    log("Execution failed:", error)
}
```

**Simple cases:**

```parsley
// If you only care about stdout
let result = COMMAND("cat", ["file.txt"]) <=#=> null
log(result.stdout)

// With fallback
let output = COMMAND("cat", ["file.txt"]) <=#=> null ?? {stdout: "default", stderr: "", exitCode: 1, error: "not found"}
log(output.stdout)
```

### Exit Codes and Errors

**Exit Code Semantics:**
- `0` = success (standard Unix convention)
- Non-zero = failure (specific meaning depends on program)
- Available even if `error` is present

**Error Field:**
- `null` if process executed (even if exit code non-zero)
- String message if process couldn't start
- Examples: "executable not found", "permission denied", "timeout exceeded"

**Distinction:**
```parsley
// Process executed but failed (exit code 1)
let {exitCode, error} = COMMAND("false") <=#=> null
// exitCode = 1, error = null

// Process couldn't execute
let {exitCode, error} = COMMAND("/nonexistent") <=#=> null
// exitCode = -1, error = "executable not found"

// Process killed by timeout
let {exitCode, error} = COMMAND("sleep", ["60"], {timeout: @1s}) <=#=> null
// exitCode = -1, error = "timeout exceeded"
```

---

## Options Dictionary

### env Option

Override environment variables:

```parsley
// Set specific variables
let result = COMMAND("./script.sh", [], {
    env: {
        PATH: "/usr/bin:/bin",
        HOME: "/tmp",
        MY_VAR: "value"
    }
}) <=#=> null

// Empty env (minimal environment)
let result = COMMAND("env", [], {env: {}}) <=#=> null
```

**Behavior:**
- `null` (default) = inherit parent process environment
- `{}` (empty dict) = minimal environment + system-required vars
- `{KEY: "val"}` = specified variables only + system-required vars

### dir Option

Set working directory:

```parsley
// Execute in specific directory
let result = COMMAND("ls", [], {dir: @./subfolder}) <=#=> null

// Relative to path
let workspace = @./projects/web
let result = COMMAND("npm", ["install"], {dir: workspace}) <=#=> null

// Absolute path
let result = COMMAND("make", ["build"], {dir: @/home/user/project}) <=#=> null
```

**Security note:** Directory must be readable (subject to `--restrict-read` if set)

### timeout Option

Kill process after duration:

```parsley
// 5 second timeout
let result = COMMAND("./slow-script.sh", [], {timeout: @5s}) <=#=> null

// If timeout exceeded:
// exitCode = -1, error = "timeout exceeded"

// Typical timeouts
let quick = COMMAND("fast-cmd", [], {timeout: @1s}) <=#=> null
let medium = COMMAND("build.sh", [], {timeout: @5m}) <=#=> null
let long = COMMAND("deploy.sh", [], {timeout: @1h}) <=#=> null
```

**Behavior:**
- `null` (default) = no timeout, waits indefinitely
- Duration = kills process if exceeds limit
- Killed processes return error and exit code -1

---

## Format Object Pattern (NEW)

### Encoding and Decoding

Formats (JSON, CSV, etc.) are objects with `encode` and `decode` methods:

```parsley
// JSON
let data = JSON.decode('{"name": "Alice", "age": 30}')
// data = {name: "Alice", age: 30}

let text = JSON.encode({name: "Bob", age: 25})
// text = '{"name":"Bob","age":25}'

// CSV
let rows = CSV.decode("a,b,c\n1,2,3\n4,5,6")
// rows = [["a","b","c"], ["1","2","3"], ["4","5","6"]]

let csv = CSV.encode([["x","y"], ["1","2"]])
// csv = "x,y\n1,2"
```

### Format Namespaces

| Format | decode(string) | encode(data) | Factory (for files) |
|--------|----------------|--------------|---------------------|
| `JSON` | String → Dict/Array | Dict/Array → String | `JSON(@./file.json)` |
| `CSV` | String → Array | Array → String | `CSV(@./file.csv)` |
| `lines` | String → Array | Array → String | `lines(@./file.txt)` |

### Separation of Concerns

```parsley
// File I/O: JSON() creates file handle
let fileHandle = JSON(@./data.json)
let data <== fileHandle

// String conversion: JSON.decode/encode
let data = JSON.decode('{"key": "value"}')
let text = JSON.encode({key: "value"})

// Combined: read file, parse JSON
let jsonText <== text(@./data.json)
let data = JSON.decode(jsonText)

// Or use file handle directly (automatic decode)
let data <== JSON(@./data.json)
```

---

## Command Handles as Variables with Dynamic Input

### Example 1: Stored Command + Format Encoding

```parsley
// Set up command handles once
let jq = COMMAND("jq", [".users[].name"])
let grep = COMMAND("grep", ["-i"])
let wc = COMMAND("wc", ["-l"])

// Function to format JSON input
let formatUserQuery = fn(filters) {
    JSON.encode({
        users: filters.users,
        timestamp: @now.format(),
        query: filters.search
    })
}

// Use the command with dynamic input
let input = formatUserQuery({users: userList, search: "alice"})
let result = jq <=#=> input

log(result.stdout)
```

### Example 2: Command Component Pattern

```parsley
// Like SQL components, but for shell commands
let FormatCSV = fn(props) {
    let rows = props.data.map(fn(row) {
        [row.name, row.email, row.age]
    })
    CSV.encode([["name","email","age"]].concat(rows))
}

// Set up CSV processor command
let csvjson = COMMAND("csvjson")

// Use component to format input
let users = [{name: "Alice", email: "a@ex.com", age: 30}]
let csvInput = <FormatCSV data={users} />
let jsonResult = csvjson <=#=> csvInput

// Parse the JSON output
let data = JSON.decode(jsonResult.stdout)
log(data)
```

### Example 3: Git Command with Formatted Commit Message

```parsley
// Store git command
let gitCommit = COMMAND("git", ["commit", "-F", "-"], {
    dir: @./myproject
})

// Component formats commit message
let CommitMessage = fn(props) {
    let lines = [
        props.type + ": " + props.subject,
        "",
        props.body,
        "",
        "Refs: " + props.issue
    ]
    lines.join("\n")
}

// Use it
let msg = <CommitMessage 
    type="feat"
    subject="Add user authentication"
    body="Implements OAuth2 flow with JWT tokens"
    issue="#123"
/>

let result = gitCommit <=#=> msg

if (result.exitCode == 0) {
    log("Committed successfully!")
}
```

### Example 4: Docker Commands with Configuration

```parsley
// Store docker commands
let dockerRun = COMMAND("docker", ["run", "-i", "myimage"])
let dockerBuild = COMMAND("docker", ["build", "-f", "-", "."], {
    dir: @./app
})

// Function generates Dockerfile
let GenerateDockerfile = fn(config) {
    let lines = [
        "FROM " + config.base,
        "WORKDIR /app",
        "COPY . .",
        "RUN " + config.buildCmd,
        "CMD [\"" + config.cmd + "\"]"
    ]
    lines.join("\n")
}

// Build with generated Dockerfile
let dockerfile = GenerateDockerfile({
    base: "node:18-alpine",
    buildCmd: "npm install && npm run build",
    cmd: "npm start"
})

let buildResult = dockerBuild <=#=> dockerfile

if (buildResult.exitCode == 0) {
    log("Docker image built successfully!")
}
```

### Example 5: SQL Query via CLI with Template

```parsley
// Store psql command
let psql = COMMAND("psql", [
    "-h", "localhost",
    "-U", "admin",
    "-d", "mydb",
    "-t",  // tuples only
    "-A"   // unaligned
], {
    env: {PGPASSWORD: "secret"}
})

// SQL template component (like <SQL> but for CLI)
let UserQuery = fn(props) {
    let sql = "SELECT name, email FROM users"
    let where = []
    
    if (props.active) {
        where = where.concat(["active = true"])
    }
    if (props.role) {
        where = where.concat(["role = '" + props.role + "'"])
    }
    
    if (where.length() > 0) {
        sql + " WHERE " + where.join(" AND ")
    } else {
        sql
    }
}

// Execute query
let query = <UserQuery active={true} role="admin" />
let result = psql <=#=> query

let users = result.stdout.trim().split("\n")
log("Found", users.length(), "users")
```

### Example 6: API Request via curl with JSON Body

```parsley
// Store curl command for API
let apiPost = COMMAND("curl", [
    "-X", "POST",
    "-H", "Content-Type: application/json",
    "-d", "@-",  // read from stdin
    "https://api.example.com/users"
])

// Component formats request body
let CreateUserRequest = fn(props) {
    JSON.encode({
        name: props.name,
        email: props.email,
        metadata: {
            source: "parsley-script",
            timestamp: @now.iso()
        }
    })
}

// Make API call
let reqBody = <CreateUserRequest name="Alice" email="alice@example.com" />
let response = apiPost <=#=> reqBody

// Parse response
let data = JSON.decode(response.stdout)
log("Created user:", data.id)
```

### Example 7: Multi-Stage Pipeline with Format Conversion

```parsley
// Set up pipeline commands
let psAux = COMMAND("ps", ["aux"])
let grepNode = COMMAND("grep", ["node"])
let wcLines = COMMAND("wc", ["-l"])

// Execute pipeline
let ps = psAux <=#=> null
let filtered = grepNode <=#=> ps.stdout
let count = wcLines <=#=> filtered.stdout

log("Node processes:", count.stdout.trim())

// Alternative: process data in Parsley
let ps = psAux <=#=> null
let lines = ps.stdout.split("\n")
    .filter(fn(line) { line.contains("node") })
    
log("Node processes:", lines.length())
```

### Example 8: JSON API to CSV Export

```parsley
// Fetch JSON from API
let apiGet = COMMAND("curl", [
    "-H", "Accept: application/json",
    "https://api.example.com/users"
])

let response = apiGet <=#=> null

// Parse JSON response
let users = JSON.decode(response.stdout)

// Convert to CSV
let csvData = CSV.encode(
    [["name", "email", "role"]].concat(
        users.map(fn(u) { [u.name, u.email, u.role] })
    )
)

// Write to file
csvData ==> text(@./users.csv)

log("Exported", users.length(), "users to CSV")
```

### Example 9: Configuration-Driven Execution

```parsley
// Config with command definitions
let config = JSON.decode(<== text(@./commands.json))

// Execute based on config
for (cmdDef in config.commands) {
    log("Running:", cmdDef.name)
    
    let cmd = COMMAND(cmdDef.binary, cmdDef.args, {
        dir: cmdDef.workdir ? @{cmdDef.workdir} : null,
        timeout: @{cmdDef.timeout ?? "30s"}
    })
    
    // Prepare input from config
    let input = cmdDef.stdin ? JSON.encode(cmdDef.stdin) : null
    
    let result = cmd <=#=> input
    
    if (result.exitCode != 0) {
        log("FAILED:", cmdDef.name)
        log(result.stderr)
        break
    }
    
    // Parse output if JSON expected
    if (cmdDef.outputFormat == "json") {
        let data = JSON.decode(result.stdout)
        log("Result:", data)
    }
}
```

---

## Security Integration

### Leverage Existing `--allow-execute` Flag

The security model from v0.10.0 already provides process execution control:

```bash
# Default: execution denied
./pars script.pars
# ERROR: command execution not allowed (use --allow-execute or -x)

# Allow specific directory
./pars --allow-execute=./scripts script.pars

# Allow multiple directories
./pars --allow-execute=./bin,./tools script.pars

# Unrestricted (development)
./pars -x script.pars
```

### Path Resolution and Validation

**Security checks before execution:**

1. **Resolve command path**
   - Absolute paths used as-is: `/usr/bin/ls`
   - Relative paths prefixed with `./`: `./script.sh`
   - Simple names searched in PATH: `ls` → `/usr/bin/ls`

2. **Apply security policy**
   - Check resolved path against `checkPathAccess(path, "execute")`
   - Deny if path not in allow-list (when policy active)
   - Same validation as module `import()`

3. **Additional safety**
   - Reject paths containing `..` (parent directory traversal)
   - Reject paths with null bytes
   - Validate working directory (if specified)

### Command Argument Safety

**No shell interpretation:**
- Arguments passed directly to executable (like Go's `exec.Command`)
- No glob expansion: `["*.txt"]` is literal asterisk-star-dot-txt
- No variable substitution: `["$HOME"]` is literal dollar-HOME
- No pipe/redirect: `["cat", "file", "|", "grep", "x"]` fails (pipe is argument)

**This prevents injection attacks:**
```parsley
// SAFE - semicolon is literal argument
let userInput = "; rm -rf /"
let result = COMMAND("echo", [userInput]) <=#=> null
// Outputs: ; rm -rf /

// SAFE - even with dangerous input
let result = COMMAND("cat", [userInput]) <=#=> null
// Error: file "; rm -rf /" not found
```

### Best Practices

**1. Use explicit paths:**
```parsley
// Good - explicit location
let result = COMMAND("./scripts/build.sh") <=#=> null

// Risky - depends on PATH
let result = COMMAND("build.sh") <=#=> null
```

**2. Validate user input:**
```parsley
// Check input before using
let sanitize = fn(input) {
    if (input.startsWith("-")) {
        return null  // Reject options
    }
    if (input.contains("..")) {
        return null  // Reject traversal
    }
    input
}

let safe = sanitize(userInput)
if (safe) {
    let result = COMMAND("cat", [safe]) <=#=> null
}
```

**3. Use timeouts for untrusted commands:**
```parsley
// Prevent hanging
let result = COMMAND(userCommand, [], {timeout: @5s}) <=#=> null
```

**4. Minimal environment:**
```parsley
// Don't leak secrets
let result = COMMAND("./script.sh", [], {
    env: {PATH: "/usr/bin", HOME: "/tmp"}
}) <=#=> null
```

---

## Implementation Strategy

### Phase 1: Core Execution (v0.11.0)

**Evaluator Changes:**

```go
// Add to builtins
"COMMAND": {
    Fn: func(args []Object, env *Environment) Object {
        // Validate arguments
        if len(args) < 1 {
            return newError("COMMAND() requires at least 1 argument")
        }
        
        binary, ok := args[0].(*String)
        if !ok {
            return newError("COMMAND() first argument must be string")
        }
        
        // Parse args array
        var cmdArgs []string
        if len(args) >= 2 {
            argsArray, ok := args[1].(*Array)
            if !ok {
                return newError("COMMAND() second argument must be array")
            }
            for _, arg := range argsArray.Elements {
                str, ok := arg.(*String)
                if !ok {
                    return newError("COMMAND() arguments must be strings")
                }
                cmdArgs = append(cmdArgs, str.Value)
            }
        }
        
        // Parse options dict (if provided)
        var opts *Dictionary
        if len(args) >= 3 {
            opts, ok = args[2].(*Dictionary)
            if !ok {
                return newError("COMMAND() third argument must be dict")
            }
        }
        
        // Create command handle
        return createCommandHandle(binary.Value, cmdArgs, opts, env)
    },
},
```

**Command Handle Creation:**

```go
func createCommandHandle(binary string, args []string, opts *Dictionary, env *Environment) Object {
    handle := &Dictionary{
        Pairs: make(map[HashKey]HashPair),
    }
    
    // Set __type
    typeKey := (&String{Value: "__type"}).HashKey()
    handle.Pairs[typeKey] = HashPair{
        Key:   &String{Value: "__type"},
        Value: &String{Value: "command"},
    }
    
    // Set binary
    binaryKey := (&String{Value: "binary"}).HashKey()
    handle.Pairs[binaryKey] = HashPair{
        Key:   &String{Value: "binary"},
        Value: &String{Value: binary},
    }
    
    // Set args
    argElements := make([]Object, len(args))
    for i, arg := range args {
        argElements[i] = &String{Value: arg}
    }
    argsKey := (&String{Value: "args"}).HashKey()
    handle.Pairs[argsKey] = HashPair{
        Key:   &String{Value: "args"},
        Value: &Array{Elements: argElements},
    }
    
    // Parse and set options
    optionsDict := parseCommandOptions(opts, env)
    optsKey := (&String{Value: "options"}).HashKey()
    handle.Pairs[optsKey] = HashPair{
        Key:   &String{Value: "options"},
        Value: optionsDict,
    }
    
    return handle
}
```

**Add `<=#=>` Token and Parser:**

```go
// In lexer
EXECUTE_WITH  // <=#=>

// In parser (infix expression)
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
    // ... existing code ...
    
    if p.curTokenIs(token.EXECUTE_WITH) {
        return p.parseExecuteExpression(left)
    }
}

func (p *Parser) parseExecuteExpression(left ast.Expression) ast.Expression {
    exp := &ast.ExecuteExpression{
        Token: p.curToken,
        Command: left,
    }
    
    precedence := p.curPrecedence()
    p.nextToken()
    exp.Input = p.parseExpression(precedence)
    
    return exp
}
```

**Execution (in evalExecuteExpression):**

```go
func evalExecuteExpression(node *ast.ExecuteExpression, env *Environment) Object {
    // Evaluate command handle
    cmdObj := Eval(node.Command, env)
    if isError(cmdObj) {
        return cmdObj
    }
    
    // Verify it's a command handle
    cmdDict, ok := cmdObj.(*Dictionary)
    if !ok {
        return newError("left operand of <=#=> must be command handle")
    }
    
    typeKey := (&String{Value: "__type"}).HashKey()
    typePair, ok := cmdDict.Pairs[typeKey]
    if !ok || typePair.Value.(*String).Value != "command" {
        return newError("left operand of <=#=> must be command handle")
    }
    
    // Evaluate input
    var input *String
    inputObj := Eval(node.Input, env)
    if !isNull(inputObj) {
        var ok bool
        input, ok = inputObj.(*String)
        if !ok {
            return newError("right operand of <=#=> must be string or null")
        }
    }
    
    // Execute the command
    return executeCommand(cmdDict, input, env)
}

func executeCommand(cmdDict *Dictionary, input *String, env *Environment) Object {
    // Extract binary
    binaryKey := (&String{Value: "binary"}).HashKey()
    binary := cmdDict.Pairs[binaryKey].Value.(*String).Value
    
    // Resolve command path
    resolvedPath, err := exec.LookPath(binary)
    if err != nil {
        if strings.Contains(binary, "/") {
            resolvedPath = binary
        } else {
            return createErrorResult("command not found: " + binary)
        }
    }
    
    // Security check
    if err := env.Security.checkPathAccess(resolvedPath, "execute"); err != nil {
        return createErrorResult("security: " + err.Error())
    }
    
    // Extract args
    argsKey := (&String{Value: "args"}).HashKey()
    argsArray := cmdDict.Pairs[argsKey].Value.(*Array)
    args := make([]string, len(argsArray.Elements))
    for i, arg := range argsArray.Elements {
        args[i] = arg.(*String).Value
    }
    
    // Extract options
    optsKey := (&String{Value: "options"}).HashKey()
    opts := cmdDict.Pairs[optsKey].Value.(*Dictionary)
    
    // Build exec.Command
    cmd := exec.Command(resolvedPath, args...)
    
    // Apply options
    applyCommandOptions(cmd, opts, env)
    
    // Set stdin if provided
    if input != nil {
        cmd.Stdin = strings.NewReader(input.Value)
    }
    
    // Execute and capture
    var stdout, stderr bytes.Buffer
    cmd.Stdout = &stdout
    cmd.Stderr = &stderr
    
    err = cmd.Run()
    
    // Build result dict
    return createResultDict(stdout.String(), stderr.String(), err)
}

func createResultDict(stdout, stderr string, err error) Object {
    result := &Dictionary{Pairs: make(map[HashKey]HashPair)}
    
    // stdout
    stdoutKey := (&String{Value: "stdout"}).HashKey()
    result.Pairs[stdoutKey] = HashPair{
        Key:   &String{Value: "stdout"},
        Value: &String{Value: stdout},
    }
    
    // stderr
    stderrKey := (&String{Value: "stderr"}).HashKey()
    result.Pairs[stderrKey] = HashPair{
        Key:   &String{Value: "stderr"},
        Value: &String{Value: stderr},
    }
    
    // exitCode and error
    if err != nil {
        if exitErr, ok := err.(*exec.ExitError); ok {
            // Non-zero exit
            exitCodeKey := (&String{Value: "exitCode"}).HashKey()
            result.Pairs[exitCodeKey] = HashPair{
                Key:   &String{Value: "exitCode"},
                Value: &Integer{Value: int64(exitErr.ExitCode())},
            }
            errorKey := (&String{Value: "error"}).HashKey()
            result.Pairs[errorKey] = HashPair{
                Key:   &String{Value: "error"},
                Value: NULL,
            }
        } else {
            // Execution failed
            exitCodeKey := (&String{Value: "exitCode"}).HashKey()
            result.Pairs[exitCodeKey] = HashPair{
                Key:   &String{Value: "exitCode"},
                Value: &Integer{Value: -1},
            }
            errorKey := (&String{Value: "error"}).HashKey()
            result.Pairs[errorKey] = HashPair{
                Key:   &String{Value: "error"},
                Value: &String{Value: err.Error()},
            }
        }
    } else {
        // Success
        exitCodeKey := (&String{Value: "exitCode"}).HashKey()
        result.Pairs[exitCodeKey] = HashPair{
            Key:   &String{Value: "exitCode"},
            Value: &Integer{Value: 0},
        }
        errorKey := (&String{Value: "error"}).HashKey()
        result.Pairs[errorKey] = HashPair{
            Key:   &String{Value: "error"},
            Value: NULL,
        }
    }
    
    return result
}
```

### Phase 2: Options Support (v0.11.0)

**stdin, env, dir, timeout implementation:**

```go
func applyCommandOptions(cmd *exec.Cmd, opts *Dictionary, env *Environment) {
    // env
    envKey := (&String{Value: "env"}).HashKey()
    if envPair, ok := opts.Pairs[envKey]; ok && !isNull(envPair.Value) {
        if envDict, ok := envPair.Value.(*Dictionary); ok {
            var envVars []string
            for _, pair := range envDict.Pairs {
                key := pair.Key.(*String).Value
                val := pair.Value.(*String).Value
                envVars = append(envVars, key+"="+val)
            }
            cmd.Env = envVars
        }
    }
    
    // dir
    dirKey := (&String{Value: "dir"}).HashKey()
    if dirPair, ok := opts.Pairs[dirKey]; ok && !isNull(dirPair.Value) {
        if pathDict, ok := dirPair.Value.(*Dictionary); ok {
            pathStr := pathDictToString(pathDict)
            cmd.Dir = pathStr
        }
    }
    
    // timeout (using context)
    timeoutKey := (&String{Value: "timeout"}).HashKey()
    if timeoutPair, ok := opts.Pairs[timeoutKey]; ok && !isNull(timeoutPair.Value) {
        if durDict, ok := timeoutPair.Value.(*Dictionary); ok {
            dur := durationDictToGoDuration(durDict)
            ctx, cancel := context.WithTimeout(context.Background(), dur)
            defer cancel()
            // Note: use CommandContext instead
        }
    }
}
```

### Phase 3: Format Object Pattern (v0.11.0)

**Add JSON.decode and JSON.encode:**

```go
// Make JSON an object with methods
"JSON": {
    Fn: func(args []Object, env *Environment) Object {
        // If called as JSON(path), create file handle (existing behavior)
        if len(args) == 1 {
            if pathDict, ok := args[0].(*Dictionary); ok {
                // Check if it's a path
                if isPathDict(pathDict) {
                    return createJSONFileHandle(pathDict, nil)
                }
            }
        }
        return newError("JSON() expects path argument for file handle")
    },
    // Add namespace methods
    Namespace: &Dictionary{
        Pairs: map[HashKey]HashPair{
            (&String{Value: "decode"}).HashKey(): {
                Key: &String{Value: "decode"},
                Value: &Builtin{
                    Fn: func(args []Object, env *Environment) Object {
                        if len(args) != 1 {
                            return newError("JSON.decode() expects 1 argument")
                        }
                        str, ok := args[0].(*String)
                        if !ok {
                            return newError("JSON.decode() expects string")
                        }
                        return parseJSON(str.Value)
                    },
                },
            },
            (&String{Value: "encode"}).HashKey(): {
                Key: &String{Value: "encode"},
                Value: &Builtin{
                    Fn: func(args []Object, env *Environment) Object {
                        if len(args) != 1 {
                            return newError("JSON.encode() expects 1 argument")
                        }
                        return stringifyJSON(args[0])
                    },
                },
            },
        },
    },
},

// Similar for CSV
"CSV": {
    // ... similar structure with decode/encode
},
```

### Phase 4: Testing & Documentation (v0.11.0)

**Unit Tests:**
- `TestCommandBasic` - Simple command execution
- `TestCommandArgs` - Command with arguments
- `TestCommandWithInput` - Input piping via `<=#=>`
- `TestCommandEnv` - Environment variables
- `TestCommandDir` - Working directory
- `TestCommandTimeout` - Timeout handling
- `TestCommandSecurity` - Security policy enforcement
- `TestCommandNonZeroExit` - Exit code handling
- `TestCommandNotFound` - Missing command errors
- `TestFormatDecode` - JSON.decode, CSV.decode
- `TestFormatEncode` - JSON.encode, CSV.encode

**Integration Tests:**
- Real commands (`ls`, `echo`, `cat`, etc.)
- Pipeline simulations
- Error scenarios
- Format conversions

**Documentation:**
- Update README.md with examples
- Add Process Execution section to reference.md
- Update CHANGELOG.md
- Add examples/process_demo.pars
- Document Format Object Pattern in reference.md

### Phase 5: Advanced Features (Future - v0.12.0)

- **Streaming output** - Real-time stdout/stderr access
- **Async execution** - Non-blocking with callbacks/promises
- **Bidirectional pipes** - Interactive processes
- **Signal handling** - Send signals to processes
- **Process groups** - Manage multiple related processes

---

## Open Questions

1. **Should we allow shell mode?**
   - Add `{shell: true}` option to enable shell interpretation?
   - **Recommendation: NO** - too dangerous, keep safe by default (✅ - agree, SP)
   
2. **Should we provide helper builtins?**
   - `shell(command)` that automatically uses `sh -c`?
   - **Recommendation: NO** - keep API minimal, users can do `COMMAND("sh", ["-c", cmd]) <=#=> null` (✅ - agree, SP)
   
3. **Should we allow background execution?**
   - Return process handle before completion, poll for status?
   - **Recommendation: FUTURE** - start synchronous, add async later if needed (✅ - agree, SP)
   
4. **Should we expose process PID?**
   - Useful for signal sending, process management
   - **Recommendation: FUTURE** - not needed for initial version (✅ - agree, SP)
   
5. **Should result include timing information?**
   - Add `duration` field showing execution time?
   - **Recommendation: YES** - useful for benchmarking, minimal cost (✅ - agree, SP)

6. **Should Format.decode/encode support options?**
   - `CSV.decode(text, {header: true})`?
   - **Recommendation: YES** - second parameter for options dict (✅ - agree, SP) 

---

## Summary

### What We're Building

- `COMMAND(binary, args?, options?)` builtin creating command handle
- `<=#=>` operator executes: `result = cmd <=#=> input`
- Returns `{stdout, stderr, exitCode, error}`
- Options: `{env, dir, timeout}`
- Security via existing `--allow-execute` flag
- Format objects: `JSON.decode/encode`, `CSV.decode/encode`
- Synchronous execution (async is future work)

### Why This Design

- ✅ Matches database pattern: `result = connection <=!=> query`
- ✅ Consistent with file I/O operator style
- ✅ Safe by default (no shell, path validation)
- ✅ Type-complete (returns useful objects)
- ✅ Composable (handles are first-class values)
- ✅ Minimal API (one builtin, one operator)
- ✅ Leverages existing security infrastructure
- ✅ Format objects separate encoding from I/O

### Implementation Phases

1. **v0.11.0** - Core execution + options + format objects + tests + docs
2. **v0.12.0** - Advanced features (streaming, async, signals)

This design brings external process execution to Parsley while maintaining the language's core principles of simplicity, safety, and composability.

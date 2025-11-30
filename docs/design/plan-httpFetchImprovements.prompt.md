# Plan: HTTP/Fetch Improvements Implementation

## Current State

- `<=/=` (fetch) works for HTTP GET requests via `JSON()`, `YAML()`, etc.
- `=/=>` (write) **does NOT support HTTP** - only local files and SFTP
- HTTP method is set via options dictionary: `{method: "POST", body: {...}}`
- SFTP has `.json`, `.text`, `.csv` format accessors we can follow as a pattern

## Design Goals

1. **Arrow direction matches data flow** - The payload (significant data) determines the arrow direction
2. **Sensible defaults** - `<=/=` defaults to GET, `=/=>` defaults to POST
3. **Method shortcuts** - `.put`, `.patch`, `.delete` accessors for HTTP methods
4. **Reusable formatters** - Create configured formatters once, use multiple times
5. **Full response metadata** - Access status, headers, URL, errors when needed

---

## Response Structure

Responses follow Parsley's typed dictionary pattern (like datetime, duration, regex).

**All internal fields use `__` prefix** to avoid collisions with JSON data:

```parsley
let users <=/= JSON(@https://api.example.com/users)

// OK response:
{
  __type: "response",
  __format: "json",
  __data: [...],              // The fetched data (auto-unwraps)
  __response: {               // Returned by .response() method
    status: 200,
    statusText: "OK",
    ok: true,
    url: @https://api.example.com/users,
    headers: {"content-type": "application/json", ...},
    error: null
  }
}

// Error response:
{
  __type: "response",
  __format: "json",
  __data: null,               // No data due to error
  __response: {
    status: 404,
    statusText: "Not Found",
    ok: false,
    url: @https://api.example.com/users,
    headers: {"content-type": "application/json", ...},
    error: "Not Found"
  }
}
```

### Usage

```parsley
let users <=/= JSON(@https://api.example.com/users)

// Direct data access - auto-unwraps to .__data
for user in users { ... }
print(users[0].name)
print(users.length)

// Metadata via .response() method - returns dictionary for destructuring
let {status, ok, error} = users.response()

if !users.response().ok {
    print("Error: " + users.response().error)
}

// Check before accessing data
if users.response().ok {
    for user in users {
        print(user.name)
    }
}
```

---

## Implementation Steps

### Phase 1: HTTP Method Accessors (`.get`, `.post`, `.put`, `.patch`, `.delete`)

**File:** `pkg/evaluator/evaluator.go` - modify `evalDotExpression()`

Add handling for request dictionaries similar to SFTP format accessors:

```go
if isRequestDict(dict) {
    methods := map[string]string{
        "get": "GET", "post": "POST", "put": "PUT", 
        "patch": "PATCH", "del": "DELETE"
    }
    if method, ok := methods[node.Key]; ok {
        // Clone request dict with new method
        return setRequestMethod(dict, method, env)
    }
}
```

**Result:**

```parsley
let response <=/= JSON(@https://api.example.com/users).post  // Sets method to POST
```

---

### Phase 2: Enable Write Operator for HTTP

**File:** `pkg/evaluator/evaluator.go` - modify `evalWriteStatement()`

Currently rejects request dictionaries. Change to:

1. Detect if target is a request dictionary
2. Set the value being written as the request body
3. Default method to `POST` (unless already set via accessor)
4. Execute HTTP request and return response

**Result:**

```parsley
let response = {name: "Alice"} =/=> JSON(@https://api.example.com/users)
// Sends POST with body: {"name": "Alice"}
```

---

### Phase 3: Combine Accessors with Write Operator

Once Phase 1 and 2 are complete, this automatically works:

```parsley
let response = {name: "Alice"} =/=> JSON(@https://api.example.com/users/1).put
// Sends PUT with body: {"name": "Alice"}

let response = {name: "Updated"} =/=> JSON(@https://api.example.com/users/1).patch
// Sends PATCH with body: {"name": "Updated"}
```

**DELETE uses fetch operator** (no payload - the request itself is the significant message):

```parsley
let result <=/= JSON(@https://api.example.com/users/1).delete
// Sends DELETE request, returns response
```

### Arrow Direction Summary

| Method | Operator | Accessor | Reason |
|--------|----------|----------|--------|
| GET | `<=/=` | `.get` | Server sends data to client |
| POST | `=/=>` | `.post` | Client sends payload to server |
| PUT | `=/=>` | `.put` | Client sends payload to server |
| PATCH | `=/=>` | `.patch` | Client sends payload to server |
| DELETE | `<=/=` | `.delete` | Request is the message, no payload |

---

### Phase 4: Reusable Formatters (Future Enhancement)

**New feature:** Allow format factories without a path argument

```parsley
// Create formatter with options - these persist across requests
let api = JSON({
    headers: {
        "Authorization": "Bearer " + token,
        "X-API-Key": apiKey
    },
    timeout: 30
})

// Reuse formatter - options remembered
let users <=/= api(@https://api.example.com/users)
let posts <=/= api(@https://api.example.com/posts)
let result = {name: "New"} =/=> api(@https://api.example.com/users).post
```

This requires:

1. `JSON()` with no path returns a "partial formatter" object
2. Partial formatter is callable with a path
3. Options are preserved and merged

---

## Implementation Order

| Phase | Description | Complexity | Files |
|-------|-------------|------------|-------|
| 1 | HTTP method accessors | Medium | evaluator.go |
| 2 | Write operator for HTTP | Medium | evaluator.go |
| 3 | Combined (automatic) | None | - |
| 4 | Reusable formatters | High | evaluator.go |

---

## Design Decisions Made

1. **DELETE uses `<=/=` with `.delete` accessor** - The request itself is the message, no payload needed:
   ```parsley
   let result <=/= JSON(@https://api.example.com/users/1).delete
   ```

2. **Response is a typed dictionary** - Follows pattern of datetime, duration, regex with `__type: "response"`

3. **Data auto-unwraps** - Iteration, indexing, length all work on `.__data` automatically

4. **`.response()` method returns `__response` dict** - Clean destructuring: `let {status, error} = users.response()`

5. **`url` follows fetch spec** - Returns final URL after any redirects

6. **`error` is always text** - Human-readable error message when request fails

7. **`data` is null on error** - Clear signal that no data was retrieved

---

## Future Pattern: Auto-Unwrapping Typed Dictionaries

The response type establishes a useful pattern that could be generalized:

**Any `__type` with a `__data` field could auto-unwrap for iteration/indexing:**

```parsley
// HTTP Response
{__type: "response", __format: "json", __data: [...], __response: {...}}

// Paginated results
{__type: "paginated", __data: [...], __page: 1, __totalPages: 10, __nextUrl: @...}

// Cached data
{__type: "cached", __data: {...}, __cachedAt: @2024-01-01T12:00:00Z, __ttl: 3600}

// Form data (for multipart/form-data requests)
{__type: "formdata", __data: {...}, __boundary: "----FormBoundary", __encoding: "multipart/form-data"}
```

**Benefits:**
- `__` prefix prevents collisions with user data
- Consistent pattern across the language
- Metadata travels with data without getting in the way
- Easy to add new wrapper types without special syntax
- Destructuring and methods provide clean access to metadata

### Other Response Formats to Consider

- **Form Data** - `multipart/form-data` for file uploads
- **URL Encoded** - `application/x-www-form-urlencoded` 
- **Binary/Blob** - Raw binary responses
- **Stream** - For chunked/streaming responses (future)


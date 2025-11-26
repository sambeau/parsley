# Dictionary Destructuring in Parsley

Dictionary destructuring allows you to extract multiple values from a dictionary into individual variables in a single statement.

## Basic Syntax

```parsley
let {key1, key2, key3} = dictionary;
```

This extracts the values associated with `key1`, `key2`, and `key3` from the dictionary and assigns them to variables with the same names.

## Features

### 1. Basic Destructuring

Extract values from a dictionary:

```parsley
let person = {name: "Alice", age: 30, city: "NYC"};
let {name, age} = person;
// name = "Alice", age = 30
```

### 2. Aliasing with `as`

Rename variables during destructuring:

```parsley
let coords = {x: 100, y: 200};
let {x as posX, y as posY} = coords;
// posX = 100, posY = 200
```

### 3. Rest Operator (`...`)

Collect remaining properties into a new dictionary:

```parsley
let data = {a: 1, b: 2, c: 3, d: 4};
let {a, b, ...rest} = data;
// a = 1, b = 2, rest = {c: 3, d: 4}
```

**Note:** The rest operator must appear last in the pattern.

### 4. Missing Keys

Keys that don't exist in the source dictionary are set to `null`:

```parsley
let obj = {x: 10};
let {x, y, z} = obj;
// x = 10, y = null, z = null
```

### 5. Nested Destructuring

Extract values from nested dictionaries:

```parsley
let user = {
    profile: {
        username: "bob",
        email: "bob@example.com"
    }
};

let {profile: {username, email}} = user;
// username = "bob", email = "bob@example.com"
```

You can combine nesting with aliasing:

```parsley
let {profile: {username as name}} = user;
// name = "bob"
```

### 6. Destructuring in Assignments

Update existing variables using destructuring:

```parsley
let a = 0;
let b = 0;
{a, b} = {a: 42, b: 99};
// a = 42, b = 99
```

## Complete Example

```parsley
// API response simulation
let apiResponse = {
    status: 200,
    data: {
        users: [{id: 1, name: "Alice"}, {id: 2, name: "Bob"}],
        total: 2
    },
    meta: {
        timestamp: 1234567890,
        version: "1.0"
    }
};

// Extract specific fields
let {
    status,
    data: {total},
    meta: {version}
} = apiResponse;

log("Status:", status);           // 200
log("Total users:", total);       // 2
log("API version:", version);     // "1.0"

// Extract with rest
let config = {host: "localhost", port: 8080, debug: true, timeout: 5000};
let {host, port, ...options} = config;
log("Server:", host + ":" + port);  // "localhost:8080"
log("Options:", options);           // {debug: true, timeout: 5000}
```

## Type Safety

Attempting to destructure a non-dictionary value results in an error:

```parsley
let {a} = "not a dict";  // ERROR: dictionary destructuring requires a dictionary value
let {b} = 42;            // ERROR: dictionary destructuring requires a dictionary value
let {c} = [1, 2, 3];     // ERROR: dictionary destructuring requires a dictionary value
```

## Limitations

- **Function parameters**: Dictionary destructuring is currently only supported in `let` statements and assignments, not in function parameters.
- **Rest operator position**: The rest operator (`...`) must be the last element in the destructuring pattern.
- **Array destructuring in nested patterns**: Nested array destructuring within dictionary patterns is not yet supported.

## Comparison with Array Destructuring

Parsley also supports array destructuring:

```parsley
// Array destructuring
let [a, b, c] = [1, 2, 3];

// Dictionary destructuring
let {a, b, c} = {a: 1, b: 2, c: 3};
```

The key difference is that array destructuring assigns by position, while dictionary destructuring assigns by key name.

# Trigonometric Functions Demo

This file demonstrates the trigonometric functions available in the Pars programming language.

## Example Usage

Start the REPL with `go run main.go` and try these examples:

### Basic Trigonometric Functions
```
sin(0)                    # Output: 0
cos(0)                    # Output: 1
tan(0)                    # Output: 0
sin(pi() / 2)            # Output: 1
cos(pi())                # Output: -1
tan(pi() / 4)            # Output: 1
```

### Inverse Trigonometric Functions
```
asin(0)                  # Output: 0
acos(1)                  # Output: 0
atan(1)                  # Output: 0.7853981633974483 (π/4)
```

### Mathematical Functions
```
sqrt(4)                  # Output: 2
sqrt(16)                 # Output: 4
pow(2, 3)               # Output: 8
pow(3, 2)               # Output: 9
pi()                     # Output: 3.141592653589793
```

### Variable Assignment and Updates
```
# Assign custom values
pi_custom = 3.1415926
e_custom = 2.71828

# Use variables in calculations
radius = 5
area = pi_custom * pow(radius, 2)
area                     # Output: 78.539815

# Update variables
radius = 10
area = pi_custom * pow(radius, 2)  
area                     # Output: 314.15926

# Variable assignment with trigonometric functions
x = pi() / 4
sin_val = sin(x)         # Output: 0.7071067811865476
cos_val = cos(x)         # Output: 0.7071067811865475

# Update and recompute
x = pi() / 2
sin_val = sin(x)         # Output: 1
cos_val = cos(x)         # Output: 6.123233995736766e-17 (≈ 0)
```

### Complex Calculations
```
# Calculate the area of a circle with radius 5
let radius = 5
let area = pi() * pow(radius, 2)
area                     # Output: 78.53981633974483

# Calculate hypotenuse using Pythagorean theorem
let a = 3
let b = 4
let c = sqrt(pow(a, 2) + pow(b, 2))
c                        # Output: 5

# Convert degrees to radians and calculate sin
let degrees = 30
let radians = degrees * pi() / 180
sin(radians)             # Output: 0.49999999999999994 (≈ 0.5)
```

### Function Definitions
```
# Define a function to calculate distance between two points
let distance = fn(x1, y1, x2, y2) {
    let dx = x2 - x1
    let dy = y2 - y1
    sqrt(pow(dx, 2) + pow(dy, 2))
}

distance(0, 0, 3, 4)     # Output: 5

# Define a function to convert degrees to radians
let toRadians = fn(degrees) {
    degrees * pi() / 180
}

toRadians(90)            # Output: 1.5707963267948966 (π/2)
```

## Supported Functions

| Function | Parameters | Description | Example |
|----------|------------|-------------|---------|
| `sin(x)` | x: angle in radians | Sine function | `sin(pi()/2)` → 1 |
| `cos(x)` | x: angle in radians | Cosine function | `cos(0)` → 1 |
| `tan(x)` | x: angle in radians | Tangent function | `tan(pi()/4)` → 1 |
| `asin(x)` | x: value between -1 and 1 | Arcsine function | `asin(1)` → π/2 |
| `acos(x)` | x: value between -1 and 1 | Arccosine function | `acos(0)` → π/2 |
| `atan(x)` | x: any real number | Arctangent function | `atan(1)` → π/4 |
| `sqrt(x)` | x: non-negative number | Square root | `sqrt(25)` → 5 |
| `pow(base, exp)` | base, exp: real numbers | Power function | `pow(2, 3)` → 8 |
| `pi()` | none | Returns π | `pi()` → 3.141592... |

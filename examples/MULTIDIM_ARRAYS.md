# Multi-dimensional Arrays Demo

This file demonstrates multi-dimensional array support in pars using square bracket notation.

## Square Bracket Notation

Square brackets `[...]` create arrays and prevent the comma operator from flattening nested structures.

### Basic Arrays
```pars
[1,2,3]              // Array with 3 elements
1,2,3                // Same as above (comma creates arrays)
```

### Nested Arrays
```pars
[[1,2,3],[4,5,6]]    // 2D array (array of arrays)
[1,2,3],[4,5,6]      // Same as above
```

### Empty Arrays
```pars
[]                   // Empty array
[[]]                 // Array containing one empty array
```

## Indexing Multi-dimensional Arrays

Use chained `[index]` to access nested elements:

### 2D Arrays
```pars
xs = [[1,2,3],[4,5,6]]
xs[0]                // Returns [1,2,3]
xs[1]                // Returns [4,5,6]
xs[1][2]             // Returns 6
```

### 3D Arrays
```pars
ys = [[1,2],[3,4],[5,6]]
ys[2][0]             // Returns 5

zs = [[[1],[2]],[[3],[4]],[[5],[6]]]
zs[1][0][0]          // Returns 3
```

### String Arrays
```pars
ws = [[["a"],["b"]],[["c"],["d"]],[["e"],["f"]]]
ws[0][1][0]          // Returns "b"
```

## Practical Examples

### Matrix Operations
```pars
matrix = [[1,2,3],[4,5,6],[7,8,9]]
matrix[0]            // First row: [1,2,3]
matrix[2][1]         // Element at row 2, column 1: 8
```

### 3D Tensor
```pars
tensor = [[[1,2],[3,4]],[[5,6],[7,8]],[[9,10],[11,12]]]
tensor[1][0][1]      // Returns 6
```

### Phone Number Lookup
```pars
numbers = ["zero","one","two","three","four","five","six","seven","eight","nine"]
phone = "07595954919"
main = fn(){
    for (x in phone) {
        numbers[toNumber(x)]
    }
}
main()
// Returns: zero, seven, five, nine, five, nine, five, four, nine, one, nine
```

## Key Features

1. **Nested structures**: Arrays can contain other arrays to any depth
2. **Chained indexing**: Use `arr[i][j][k]` to access deeply nested elements
3. **Mixed types**: Arrays can contain different types including nested arrays
4. **Empty arrays**: Both `[]` and `[[]]` are supported
5. **Backward compatible**: Comma notation still works for creating arrays

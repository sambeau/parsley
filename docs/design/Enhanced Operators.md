
## Performance Considerations

1. **Hash Maps:** Use `map[string]bool` with `Inspect()` output as keys for O(n) lookups in set operations

2. **Pre-allocation:** Pre-allocate result slices when size is known (repetition, chunking)

3. **String Builder:** Use `strings.Builder` for efficient string repetition

4. **Equality Semantics:** Use existing `Inspect()` method for consistent equality across mixed-type arrays

## Design Decisions

### Why These Operators?

- **Familiar:** Leverage mathematical operators for intuitive meaning
- **Composable:** Operators chain naturally with existing syntax
- **Concise:** Replace verbose filter/map patterns with single operators
- **Consistent:** Same operators work across arrays and dictionaries where applicable

### Why Not Dictionary Union (`||`)?

Dictionary merge already exists with `++` (right wins). For left-wins behavior, just reverse the operands:
- `a ++ b` (right wins) is equivalent to reversing for left-wins
- Adding `||` would be redundant and confusing

### Shallow vs Deep Operations

All dictionary operations are **shallow only**:
- Simpler to reason about
- Predictable performance
- Avoids ambiguity with nested structures
- Deep operations can be built with functions if needed

### Equality Semantics

Use `Inspect()` output for comparing values:
- Consistent with existing Parsley behavior
- Works with mixed-type arrays
- Simple implementation
- Predictable results

## Backward Compatibility

All changes extend existing operators with new type combinations. No breaking changes to existing behavior:

- `&&` and `||` currently only work with booleans → extend to arrays/dicts
- `++` currently works with arrays and dicts → extend to scalar + array
- `-` currently only arithmetic → extend to arrays/dicts
- `/` currently only arithmetic → extend to arrays
- `*` currently only arithmetic → extend to strings/arrays

## Future Enhancements

Potential additions not in this design:

1. **Array modulo (`%`):** Every nth element - `[1,2,3,4,5,6] % 2` → `[2,4,6]`
2. **Deep dictionary operations:** Recursive merge/intersection
3. **Set type:** Dedicated set data structure with guaranteed uniqueness
4. **Performance optimizations:** Specialized hash functions for common types

## Implementation Checklist

- [x] Implement scalar concatenation (`++`)
- [x] Implement array/dict intersection (`&&`)
- [x] Implement array union (`||`)
- [x] Implement array/dict subtraction (`--`)
- [x] Implement array chunking (`/`)
- [x] Implement string/array repetition (`*`)
- [x] Add comprehensive tests for array operators
- [x] Add comprehensive tests for dictionary operators
- [x] Update `docs/reference.md` with new operators
- [x] Update `README.md` with examples
- [ ] Performance testing with large collections
- [x] Update version number in `VERSION` file

---

Please save this content as `/Users/samphillips/Dev/parsley/docs/design/Enhanced Operators.md`## Performance Considerations

1. **Hash Maps:** Use `map[string]bool` with `Inspect()` output as keys for O(n) lookups in set operations

2. **Pre-allocation:** Pre-allocate result slices when size is known (repetition, chunking)

3. **String Builder:** Use `strings.Builder` for efficient string repetition

4. **Equality Semantics:** Use existing `Inspect()` method for consistent equality across mixed-type arrays

## Design Decisions

### Why These Operators?

- **Familiar:** Leverage mathematical operators for intuitive meaning
- **Composable:** Operators chain naturally with existing syntax
- **Concise:** Replace verbose filter/map patterns with single operators
- **Consistent:** Same operators work across arrays and dictionaries where applicable

### Why Not Dictionary Union (`||`)?

Dictionary merge already exists with `++` (right wins). For left-wins behavior, just reverse the operands:
- `a ++ b` (right wins) is equivalent to reversing for left-wins
- Adding `||` would be redundant and confusing

### Shallow vs Deep Operations

All dictionary operations are **shallow only**:
- Simpler to reason about
- Predictable performance
- Avoids ambiguity with nested structures
- Deep operations can be built with functions if needed

### Equality Semantics

Use `Inspect()` output for comparing values:
- Consistent with existing Parsley behavior
- Works with mixed-type arrays
- Simple implementation
- Predictable results

## Backward Compatibility

All changes extend existing operators with new type combinations. No breaking changes to existing behavior:

- `&&` and `||` currently only work with booleans → extend to arrays/dicts
- `++` currently works with arrays and dicts → extend to scalar + array
- `-` currently only arithmetic → extend to arrays/dicts
- `/` currently only arithmetic → extend to arrays
- `*` currently only arithmetic → extend to strings/arrays

## Future Enhancements

Potential additions not in this design:

1. **Array modulo (`%`):** Every nth element - `[1,2,3,4,5,6] % 2` → `[2,4,6]`
2. **Deep dictionary operations:** Recursive merge/intersection
3. **Set type:** Dedicated set data structure with guaranteed uniqueness
4. **Performance optimizations:** Specialized hash functions for common types

## Implementation Checklist

- [x] Implement scalar concatenation (`++`)
- [x] Implement array/dict intersection (`&&`)
- [x] Implement array union (`||`)
- [x] Implement array/dict subtraction (`--`)
- [x] Implement array chunking (`/`)
- [x] Implement string/array repetition (`*`)
- [x] Add comprehensive tests for array operators
- [x] Add comprehensive tests for dictionary operators
- [x] Update `docs/reference.md` with new operators
- [x] Update `README.md` with examples
- [ ] Performance testing with large collections
- [x] Update version number in `VERSION` file

---

Please save this content as `/Users/samphillips/Dev/parsley/docs/design/Enhanced Operators.md`
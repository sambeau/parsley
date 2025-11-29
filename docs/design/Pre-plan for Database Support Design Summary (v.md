Plan: Database Support Design Summary (v1)
Consolidate all design decisions into a complete, implementable specification.

Complete Database API Design for Parsley v1
Connection Factories
Query Operators
Query Components with SQL Tags
Transaction API
Connection Dictionary Properties
Default Options per Database
This design is complete, implementable, and ready for v1. All complexity deferred to future versions!

Plan: Flexible Transaction Error Handling
Design transaction API that supports both simple (check-commit-only) and detailed (check-critical-steps) patterns.

Implementation
Pattern 1: Simple - Check Commit Only
Pattern 2: Hybrid - Check When Needed
Pattern 3: Conditional Logic
Pattern 4: Early Exit on Specific Errors
Key Design Points
Operators always return {data, error} pattern - Consistent whether in transaction or not:
Database auto-aborts on error - Once any query fails in a transaction, database enters aborted state. Subsequent queries are ignored until rollback.

commit() returns false if transaction aborted - Even if you don't check individual queries, commit() will fail and db.lastError contains the first error.

Developers choose their pattern - Simple transactions: check commit. Complex logic: check critical steps. Best of both worlds.

// Simple pattern - recommended for most use cases
db.begin()
db <=!=> <CreateUser name="Alice" />
db <=!=> <CreatePost authorId={1} />
if (!db.commit()) {
    log("Failed:", db.lastError)
}

// Detailed pattern - when you need the data or specific error handling
db.begin()
let {user, error} = db <=!=> <CreateUser name="Alice" />
if (error) {
    db.rollback()
    handleUserCreationError(error)
} else {
    db <=!=> <CreatePost authorId={user.id} />
    db.commit()
}

This gives maximum flexibility while keeping the simple case really simple!

Plan: Add Database Support to Parsley
Add database connectivity with three query operators, component-based queries, and flexible transaction handling.

Steps
Add database connection factories to evaluator.go - Implement SQLITE(), POSTGRES(), MYSQL() builtin functions that accept path/URL literals and optional options dict. Return connection dictionaries with properties (type, host, database, connected, inTransaction, lastError) and methods (begin(), commit(), rollback(), close(), ping()). Cache connections by normalized path/URL + options.

Add query operator tokens to lexer.go - Define three new operators:

QUERY_ONE for <=?=> (expect single row → dict/null)
QUERY_MANY for <=??=> (expect multiple rows → array)
EXECUTE for <=!=> (mutation → {affected, lastId} or RETURNING results)
Use multi-character lookahead similar to existing READ_FROM (<==), FETCH_FROM (<=/=).

Create query statement AST nodes in ast.go - Add QueryStatement, QueryManyStatement, ExecuteStatement similar to existing ReadStatement/FetchStatement. Each contains left expression (connection), right expression (query component), and operator token.

Implement <SQL> built-in component in evaluator.go - Parse <SQL params={...}>SQL CODE</SQL> and return query handle dict: {__type: "query", sql: "...", params: {...}}. Process {expr} interpolations in SQL content by replacing with :name placeholders and extracting params. Support both named params (:name) and spread operator (params={...props}).

Add query operator evaluation in evaluator.go - Implement evalQueryStatement(), evalQueryManyStatement(), evalExecuteStatement(). Extract SQL and params from query handle, send to appropriate database driver (use modernc.org/sqlite for pure Go SQLite, github.com/lib/pq for PostgreSQL, github.com/go-sql-driver/mysql for MySQL). Return {data, error} pattern consistently. Handle PostgreSQL RETURNING clause (returns rows instead of metadata).

Implement transaction API methods - Add begin(), commit(), rollback() to connection dictionaries. begin() returns true on success, commit() returns true if commit succeeds (false if fails + auto-rolled back), rollback() always succeeds. Track transaction state in connection, populate lastError on failures. Database auto-aborts transactions on query errors.

Further Considerations
Transaction pattern flexibility - Simple pattern: check only commit() for straightforward transactions. Detailed pattern: check critical query errors when you need result data or specific error handling. Both patterns coexist—developers choose based on needs.

Placeholder style - Use :name (named parameters) for readability and compatibility with spread operator. Each database driver translates to native format (PostgreSQL: $1, $2, SQLite/MySQL: ?).

Connection options and defaults - SQLite defaults: {timeout: @5s, cache_size: 2000, journal_mode: "WAL"}. PostgreSQL: {pool_size: 5, pool_timeout: @30s, ssl_mode: "prefer"}. MySQL: {pool_size: 5, timeout: @30s, charset: "utf8mb4"}. Support URL query params and explicit options dict (dict takes precedence).

## Connection Factories

```parsely
// Explicit database types with path/URL literals
let db = SQLITE(@./data.db)
let db = POSTGRES(@postgres://user:pass@localhost:5432/mydb)
let db = MYSQL(@mysql://user:pass@localhost:3306/mydb)

// With options (second parameter)
let db = SQLITE(@./data.db, {timeout: @5s, readonly: true})
let db = POSTGRES(@postgres://localhost/mydb, {
    pool_size: 10,
    timeout: @30s
})

// Connections cached by path/URL + options
let db1 = SQLITE(@./data.db)
let db2 = SQLITE(@./data.db)  // Returns same connection
```

## Query Operators

```parsely
// <=?=> expects single row (returns dict or null)
let user = db <=?=> <GetUser id={1} />

// <=??=> expects multiple rows (returns array)
let users = db <=??=> <SearchUsers name="Alice" />

// <=!=> for mutations (returns {affected, lastId} or RETURNING results)
let { affected } = db <=!=> <CreateUser name="Alice" email="a@example.com" />

// With PostgreSQL RETURNING
let user = db <=!=> <CreateUser name="Alice" /> // Returns full row
let users = db <=!=> <CreateUsers batch={data} /> // Returns array

## Query Components with SQL Tags

    ```parsely
// Component returns <SQL> tag with params
let GetUser = fn( props ) {
    <SQL params={...props}>
    SELECT * FROM users WHERE id = :id
</SQL>
}

// Spread operator when columns match props
let CreateUser = fn( props ) {
    <SQL params={...props}>
    INSERT INTO users (name, email, age)
    VALUES (:name, :email, :age)
</SQL>
}

// Explicit mapping when names differ
let CreatePost = fn( props ) {
    <SQL params={
    post_title: props.title,
    post_body:props.body,
    author_id:props.authorId
    }>
    INSERT INTO posts (post_title, post_body, author_id)
    VALUES (:post_title, :post_body, :author_id)
</SQL>
}

// With validation
let { isEmail } = import( @./ validators.pars )

let RegisterUser = fn( props ) {
    if ( !isEmail( props.email ) ) {
    return { error: "Invalid email" }
}

<SQL params={...props, created_at: @now, verified: false}>
    INSERT INTO users( name, email, password_hash, created_at, verified )
VALUES(: name, : email, : password_hash, : created_at, : verified )
    </SQL >
}
```

## Flexible Transaction Error Handling

```parsely
// For straightforward transactions where any failure = rollback everything

if ( db.begin() ) {
    db <=!=> <CreateUser name="Alice" email="alice@example.com" />
    db <=!=> <CreatePost authorId={1} title="First Post" />
    db <=!=> <AddTag postId={1} tag="announcement" />

    if ( db.commit() ) {
        log( "User, post, and tag created successfully!" )
    } else {
        log( "Transaction failed:", db.lastError )
    }
}
```

## Connection Dictionary Properties

    ```parsely
db.type           // "postgres", "mysql", "sqlite"
db.host           // "localhost"
db.port           // 5432
db.database       // "mydb"
db.user           // "user"
db.connected      // true/false
db.inTransaction  // true/false
db.lastError      // {code, message} or null

// Methods
db.begin()        // Start transaction
db.commit()       // Commit (returns boolean)
db.rollback()     // Rollback
db.close()        // Close connection
db.ping()         // Test connection
```

## Default Options per Database

```parsely
// SQLite defaults
{ timeout: @5s, cache_size: 2000, journal_mode: "WAL" }

// PostgreSQL defaults
{ pool_size: 5, pool_timeout: @30s, ssl_mode: "prefer" }

// MySQL defaults
{ pool_size: 5, timeout: @30s, charset: "utf8mb4" }
```
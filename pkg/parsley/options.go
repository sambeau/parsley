package parsley

import (
	"database/sql"

	"github.com/sambeau/parsley/pkg/evaluator"
)

// DBConnectionConfig holds an injected database connection.
// Used by WithDB() to pass server-managed connections to Parsley scripts.
type DBConnectionConfig struct {
	DB     *sql.DB
	Driver string // "sqlite", "postgres", "mysql"
}

// Config holds evaluation configuration
type Config struct {
	Env           *evaluator.Environment
	Security      *evaluator.SecurityPolicy
	Logger        evaluator.Logger
	Filename      string
	Vars          map[string]interface{}
	DBConnections map[string]*DBConnectionConfig // Injected database connections
}

// Option configures evaluation
type Option func(*Config)

// WithEnv uses a pre-configured environment
func WithEnv(env *evaluator.Environment) Option {
	return func(c *Config) {
		c.Env = env
	}
}

// WithSecurity sets the security policy for file system access
func WithSecurity(policy *evaluator.SecurityPolicy) Option {
	return func(c *Config) {
		c.Security = policy
	}
}

// WithLogger sets the logger for log()/logLine() output
func WithLogger(logger evaluator.Logger) Option {
	return func(c *Config) {
		c.Logger = logger
	}
}

// WithFilename sets the filename for error messages
func WithFilename(name string) Option {
	return func(c *Config) {
		c.Filename = name
	}
}

// WithVar pre-populates a variable in the environment.
// The value is converted from Go types to Parsley types using ToParsley().
func WithVar(name string, value interface{}) Option {
	return func(c *Config) {
		if c.Vars == nil {
			c.Vars = make(map[string]interface{})
		}
		c.Vars[name] = value
	}
}

// WithDB injects a database connection into the Parsley environment.
// The connection is available to scripts as a variable with the given name.
// The host application is responsible for managing the connection lifecycle
// (opening, closing, pooling). Parsley will NOT close this connection.
//
// Example:
//
//	db, _ := sql.Open("sqlite", "./app.db")
//	defer db.Close()
//
//	result, err := parsley.EvalFile("handler.pars",
//	    parsley.WithDB("db", db, "sqlite"),
//	)
//
// In Parsley script:
//
//	let user = db <=?=> "SELECT * FROM users WHERE id = 1"
func WithDB(name string, db *sql.DB, driver string) Option {
	return func(c *Config) {
		if c.DBConnections == nil {
			c.DBConnections = make(map[string]*DBConnectionConfig)
		}
		c.DBConnections[name] = &DBConnectionConfig{
			DB:     db,
			Driver: driver,
		}
	}
}

// newConfig creates a new Config with defaults and applies options
func newConfig(opts ...Option) *Config {
	c := &Config{
		Logger: evaluator.DefaultLogger,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// applyConfig applies the config to an environment
func applyConfig(env *evaluator.Environment, c *Config) error {
	// Apply security policy
	if c.Security != nil {
		env.Security = c.Security
	}

	// Apply filename
	if c.Filename != "" {
		env.Filename = c.Filename
	}

	// Apply logger
	if c.Logger != nil {
		env.Logger = c.Logger
	}

	// Apply variables
	for name, value := range c.Vars {
		obj, err := ToParsley(value)
		if err != nil {
			return err
		}
		env.Set(name, obj)
	}

	// Apply database connections
	for name, dbConfig := range c.DBConnections {
		conn := evaluator.NewManagedDBConnection(dbConfig.DB, dbConfig.Driver)
		env.Set(name, conn)
	}

	return nil
}

package schema

import "time"

// Column represents a single column in a PostgreSQL table
type Column struct {
	Name string `json:"name"`
	Type string `json:"type"` // PostgreSQL data type
	// PGTypeOID          int
	GoType          string      `json:"go_type"` // Corresponding Go type
	IsPrimaryKey    bool        `json:"is_primary_key"`
	IsNullable      bool        `json:"is_nullable"`
	IsUnique        bool        `json:"is_unique"`
	DefaultValue    string      `json:"default_value"` // SQL default value
	MaxLength       int         `json:"max_length"`    // For VARCHAR types
	Precision       int         `json:"precision"`     // For NUMERIC types
	Scale           int         `json:"scale"`         // For NUMERIC types
	CheckConstraint string      `json:"check_constraint"`
	Comment         string      `json:"comment"`          // Column comment
	OrdinalPosition int         `json:"ordinal_position"` // Column order in table
	ForeignKey      *ForeignKey `json:"foreign_key,omitempty"`
}

// ForeignKey represents foreign key constraints
type ForeignKey struct {
	ReferenceTable  string `json:"reference_table"`
	ReferenceColumn string `json:"reference_column"`
	OnDelete        string `json:"on_delete"` // CASCADE, SET NULL, etc.
	OnUpdate        string `json:"on_update"`
}

// Index represents a database index
type Index struct {
	Name     string   `json:"name"`
	Columns  []string `json:"columns"`
	IsUnique bool     `json:"is_unique"`
	Method   string   `json:"method"` // BTREE, GIN, etc.
}

// Table represents a complete PostgreSQL table definition
type Table struct {
	Name        string    `json:"name"`
	Schema      string    `json:"schema"` // public, etc.
	Columns     []Column  `json:"columns"`
	Indexes     []Index   `json:"indexes"`
	PrimaryKey  []string  `json:"primary_key"`  // For composite primary keys
	Comment     string    `json:"comment"`      // Table comment
	Owner       string    `json:"owner"`        // Table owner
	Tablespace  string    `json:"tablespace"`   // Tablespace name
	IsTemporary bool      `json:"is_temporary"` // Is temporary table
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

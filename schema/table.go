package schema

import "time"

// Column represents a single column in a PostgreSQL table
type Column struct {
	Name string `json:"name"`
	Type string `json:"type"` // PostgreSQL data type
	// PGTypeOID          int
	GoType string `json:"go_type"` // Corresponding Go type
	// IsPrimaryKey    bool        `json:"is_primary_key"` // primary key constraint
	IsNullable bool `json:"is_nullable"`
	// IsUnique        bool        `json:"is_unique"`     // unique constraint
	DefaultValue string `json:"default_value"` // SQL default value
	MaxLength    int    `json:"max_length"`    // For VARCHAR types
	// CheckConstraint string      `json:"check_constraint"`
	Comment         string `json:"comment"`          // Column comment
	OrdinalPosition int    `json:"ordinal_position"` // Column order in table

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
	Name           string   `json:"name"`
	IndexedColumns []string `json:"columns"`
	IsUnique       bool     `json:"is_unique"`
	IsPrimaryKey   bool
	IsExclusion    bool
	Method         string `json:"method"` // BTREE, GIN, etc.
	Definition     string
}

// Table represents a complete PostgreSQL table definition
type Table struct {
	Name       string `json:"name"`
	Schema     string `json:"schema"` // public, etc.
	Info       *TableInfo
	Columns    []Column `json:"columns"`
	Indexes    []Index  `json:"indexes"`
	Constraint *Constraint
}

type Numeric struct {
	Precision int `json:"precision"` // For NUMERIC types
	Scale     int `json:"scale"`     // For NUMERIC types
}

type TableInfo struct {
	Comment     string    `json:"comment"`      // Table comment
	IsTemporary bool      `json:"is_temporary"` // Is temporary table
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Tablespace  string    `json:"tablespace"` // Tablespace name
	Owner       string    `json:"owner"`      // Table owner

}

type Constraint struct {
	PrimaryKey []PrimaryKey
	Unique     []Unique
	Check      []Check
	ForeignKey []ForeignKey `json:"foreign_key,omitempty"`
}

type Unique struct {
	Name      string
	Defferred int
}

type PrimaryKey struct {
	Name       string
	Column     string
	Definition string
}

type Check struct {
	Name       string
	Unique     bool
	Definition string
	Column     string
	Defferred  int
}

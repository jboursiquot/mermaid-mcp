package erd

type Column struct {
	Name  string `db:"column_name"`
	Type  string `db:"data_type"`
	IsKey bool   // Indicates if this column is used in a relationship
}

type Table struct {
	Name    string
	Columns []Column
}

type ForeignKey struct {
	Table         string `db:"table"`
	Column        string `db:"column"`
	ForeignTable  string `db:"foreign_table"`
	ForeignColumn string `db:"foreign_column"`
}

type metadata struct {
	Direction   string
	Tables      []Table
	ForeignKeys []ForeignKey
}

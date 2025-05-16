package erd

// SQL queries for schema introspection
var queries = map[string]string{
	// Get tables with filtering condition
	"tables": `
    SELECT table_name
    FROM information_schema.tables
    WHERE table_schema = 'public'`,

	// Get columns for a specific table
	"all_columns": `
    SELECT column_name, data_type
    FROM information_schema.columns
    WHERE table_schema = 'public'
      AND table_name = $1
    ORDER BY ordinal_position`,

	// Get primary and foreign key columns for a specific table
	"key_columns": `
    SELECT c.column_name, c.data_type
    FROM information_schema.columns c
    JOIN information_schema.key_column_usage kcu
      ON c.table_name = kcu.table_name
     AND c.column_name = kcu.column_name
     AND c.table_schema = kcu.table_schema
    JOIN information_schema.table_constraints tc
      ON kcu.constraint_name = tc.constraint_name
     AND kcu.table_schema = tc.table_schema
    WHERE c.table_schema = 'public'
      AND c.table_name = $1
      AND tc.constraint_type IN ('PRIMARY KEY', 'FOREIGN KEY')
    ORDER BY c.ordinal_position`,

	// Get foreign keys filtered by tables
	"foreign_keys": `
    SELECT DISTINCT
      tc.table_name AS "table",
      kcu.column_name AS "column",
      ccu.table_name AS foreign_table,
      ccu.column_name AS foreign_column
    FROM information_schema.table_constraints tc
    JOIN information_schema.key_column_usage kcu
      ON tc.constraint_name = kcu.constraint_name
    JOIN information_schema.constraint_column_usage ccu
      ON ccu.constraint_name = tc.constraint_name
    WHERE tc.constraint_type = 'FOREIGN KEY'
      AND tc.table_schema = 'public'
      AND (tc.table_name IN (?) OR ccu.table_name IN (?))`,
}

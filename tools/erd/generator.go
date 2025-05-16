package erd

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	mcp "github.com/metoro-io/mcp-golang"
)

// Generator handles Mermaid ERD diagram generation
type Generator struct {
	db          *sqlx.DB
	Name        string
	Description string
}

func NewGenerator(db *sqlx.DB) (*Generator, error) {
	if db == nil {
		return nil, errors.New("database connection is nil")
	}
	return &Generator{
		db:          db,
		Name:        "ERD Generator",
		Description: "Generate Mermaid erDiagram markup from a database schema",
	}, nil
}

// Arguments represents the input parameters for diagram generation
type Arguments struct {
	TableNames        []string `json:"table_names" jsonschema:"required,description=List of tables to generate mermaid diagram for."`
	Direction         string   `json:"direction" jsonschema:"default=LR,description=Direction of the diagram (LR, TB, RL, BT)."`
	IncludeAllColumns bool     `json:"include_all_columns" jsonschema:"default=false,description=Include all columns in the diagram. By default, only key columns are included."`
}

// Generate processes the generation request
func (g *Generator) Generate(arguments Arguments) (*mcp.ToolResponse, error) {
	ctx := context.Background()

	md, err := g.generateMetadata(ctx, arguments.TableNames, arguments.Direction, arguments.IncludeAllColumns)
	if err != nil {
		return nil, fmt.Errorf("failed to generate metadata: %w", err)
	}

	buf := &bytes.Buffer{}
	if err := mermaidERDTemplate.Execute(buf, md); err != nil {
		return nil, fmt.Errorf("failed to render ER diagram: %w", err)
	}

	return mcp.NewToolResponse(
		mcp.NewTextContent(buf.String()),
	), nil
}

func (g *Generator) generateMetadata(ctx context.Context, tableNames []string, direction string, includeAllColumns bool) (metadata, error) {
	var tablesFound []string

	// If tableNames is empty, get all tables
	if len(tableNames) == 0 {
		// Query to get all tables
		err := g.db.SelectContext(ctx, &tablesFound, queries["tables"])
		if err != nil {
			return metadata{}, fmt.Errorf("failed to load all tables: %w", err)
		}
	} else {
		// Filter tables by the provided names
		query, args, err := sqlx.In(queries["tables"]+" AND table_name IN (?)", tableNames)
		if err != nil {
			return metadata{}, fmt.Errorf("failed to build query for tables: %w", err)
		}
		query = g.db.Rebind(query)
		err = g.db.SelectContext(ctx, &tablesFound, query, args...)
		if err != nil {
			return metadata{}, fmt.Errorf("failed to load tables: %w", err)
		}
	}

	tables := make([]Table, 0, len(tablesFound))

	// Load relevant columns
	for _, tn := range tablesFound {
		cols := []Column{}
		if includeAllColumns {
			err := g.db.SelectContext(ctx, &cols, queries["all_columns"], tn)
			if err != nil {
				return metadata{}, fmt.Errorf("failed to load all columns for table %s: %w", tn, err)
			}
		} else {
			err := g.db.SelectContext(ctx, &cols, queries["key_columns"], tn)
			if err != nil {
				return metadata{}, fmt.Errorf("failed to load key columns for table %s: %w", tn, err)
			}
		}
		tables = append(tables, Table{Name: tn, Columns: cols})
	}

	// Load relevant foreign keys - filter by the tables we're interested in
	fks := []ForeignKey{}

	if len(tablesFound) != 0 {
		query, args, err := sqlx.In(queries["foreign_keys"], tablesFound, tablesFound)
		if err != nil {
			return metadata{}, fmt.Errorf("failed to build query for foreign keys: %w", err)
		}
		query = g.db.Rebind(query)
		err = g.db.SelectContext(ctx, &fks, query, args...)
		if err != nil {
			return metadata{}, fmt.Errorf("failed to load foreign keys: %w", err)
		}
	}

	return metadata{Direction: direction, Tables: tables, ForeignKeys: fks}, nil
}

package erd_test

import (
	"testing"

	erd "github.com/jboursiquot/mermaid-mcp/tools/erd"
	mcp "github.com/metoro-io/mcp-golang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerator_Generate(t *testing.T) {
	// Get the test DB
	db := GetTestDB()
	require.NotNil(t, db, "test database connection is nil")

	// Test cases as a map with descriptive keys
	testCases := map[string]struct {
		tableNames        []string
		direction         string
		includeAllColumns bool
		wantErr           bool
		checkFunc         func(t *testing.T, resp *mcp.ToolResponse)
	}{
		"All tables with key columns only (default)": {
			tableNames: []string{},
			direction:  "LR",
			wantErr:    false,
			checkFunc: func(t *testing.T, resp *mcp.ToolResponse) {
				content := resp.Content[0].TextContent.Text

				// Verify we have an ER diagram
				assert.Contains(t, content, "erDiagram")

				// Check for expected tables
				assert.Contains(t, content, "customers {")
				assert.Contains(t, content, "products {")
				assert.Contains(t, content, "orders {")
				assert.Contains(t, content, "order_items {")
				assert.Contains(t, content, "inventory {")
				assert.Contains(t, content, "suppliers {")
				assert.Contains(t, content, "shipments {")
				assert.Contains(t, content, "payments {")
				assert.Contains(t, content, "reviews {")

				// Check for expected relationships
				assert.Contains(t, content, "orders }o--|| customers")
			},
		},
		"Specific tables": {
			tableNames: []string{"customers", "orders", "products", "order_items"},
			direction:  "TB",
			wantErr:    false,
			checkFunc: func(t *testing.T, resp *mcp.ToolResponse) {
				content := resp.Content[0].TextContent.Text

				// Verify we have an ER diagram
				assert.Contains(t, content, "erDiagram")

				// Check for expected tables
				assert.Contains(t, content, "customers {")
				assert.Contains(t, content, "orders {")
				assert.Contains(t, content, "products {")
				assert.Contains(t, content, "order_items {")

				// Check for unexpected tables
				assert.NotContains(t, content, "inventory {")
				assert.NotContains(t, content, "suppliers {")
				assert.NotContains(t, content, "shipments {")
				assert.NotContains(t, content, "payments {")
				assert.NotContains(t, content, "reviews {")

				// Check for expected relationships
				assert.Contains(t, content, "orders }o--|| customers")
				assert.Contains(t, content, "order_items }o--|| orders")
				assert.Contains(t, content, "payments }o--|| orders")
				assert.Contains(t, content, "shipments }o--|| orders")
				assert.Contains(t, content, "inventory }o--|| products")
				assert.Contains(t, content, "order_items }o--|| products")
				assert.Contains(t, content, "reviews }o--|| products")
				assert.Contains(t, content, "reviews }o--|| customers")

				// Check for unexpected relationships
				assert.NotContains(t, content, "inventory }o--|| suppliers")
			},
		},
		"All tables with all columns": {
			tableNames:        []string{"customers", "orders", "products", "order_items"},
			direction:         "TB",
			includeAllColumns: true,
			wantErr:           false,
			checkFunc: func(t *testing.T, resp *mcp.ToolResponse) {
				content := resp.Content[0].TextContent.Text

				// Verify we have an ER diagram
				assert.Contains(t, content, "erDiagram")

				// Check for expected tables
				assert.Contains(t, content, "customers {")
				assert.Contains(t, content, "orders {")
				assert.Contains(t, content, "products {")
				assert.Contains(t, content, "order_items {")

				// Check for all columns in customers table
				assert.Regexp(t, `customers \{[^}]*id`, content, "customers table should include 'id' column")
				assert.Regexp(t, `customers \{[^}]*name`, content, "customers table should include 'name' column")
				assert.Regexp(t, `customers \{[^}]*email`, content, "customers table should include 'email' column")

				// Check for all columns in products table
				assert.Regexp(t, `products \{[^}]*id`, content, "products table should include 'id' column")
				assert.Regexp(t, `products \{[^}]*name`, content, "products table should include 'name' column")

				// Check for all columns in orders table
				assert.Regexp(t, `orders \{[^}]*id`, content, "orders table should include 'id' column")
				assert.Regexp(t, `orders \{[^}]*customer_id`, content, "orders table should include 'customer_id' column")
				assert.Regexp(t, `orders \{[^}]*order_date`, content, "orders table should include 'order_date' column")

				// Check for all columns in order_items table
				assert.Regexp(t, `order_items \{[^}]*id`, content, "order_items table should include 'id' column")
				assert.Regexp(t, `order_items \{[^}]*order_id`, content, "order_items table should include 'order_id' column")
				assert.Regexp(t, `order_items \{[^}]*product_id`, content, "order_items table should include 'product_id' column")
				assert.Regexp(t, `order_items \{[^}]*quantity`, content, "order_items table should include 'quantity' column")

				// Check for expected relationships
				assert.Contains(t, content, "orders }o--|| customers")
				assert.Contains(t, content, "order_items }o--|| orders")
				assert.Contains(t, content, "order_items }o--|| products")

				// Should not contain unrelated tables
				assert.NotContains(t, content, "inventory {")
				assert.NotContains(t, content, "suppliers {")
				assert.NotContains(t, content, "shipments {")
				assert.NotContains(t, content, "payments {")
				assert.NotContains(t, content, "reviews {")
			},
		},
		"Non-existent table": {
			tableNames: []string{"table_that_does_not_exist"},
			direction:  "LR",
			wantErr:    false,
			checkFunc: func(t *testing.T, resp *mcp.ToolResponse) {
				content := resp.Content[0].TextContent.Text

				// Should still have the erDiagram header
				assert.Contains(t, content, "erDiagram")

				// Output should be mostly empty (no tables)
				assert.True(t, assert.NotContains(t, content, "{"), "Expected no table definitions")
			},
		},
	}

	// Run the tests
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tool, err := erd.NewGenerator(db)
			require.NoError(t, err, "failed to create generator")

			resp, err := tool.Generate(erd.Arguments{
				TableNames:        tc.tableNames,
				Direction:         tc.direction,
				IncludeAllColumns: tc.includeAllColumns,
			})

			// Check for errors
			if tc.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			// Run the specific checks for this test case
			tc.checkFunc(t, resp)
		})
	}
}

package erd_test

import (
	_ "github.com/lib/pq"
)

// func TestGenerateMetadata(t *testing.T) {
// 	ctx := context.Background()

// 	t.Run("Test GenerateMetadata", func(t *testing.T) {
// 		md, err := mermaidmcp.GenerateMetadata(ctx, testDB, []string{"customers", "orders"}, "LR")
// 		assert.NoError(t, err)
// 		assert.NotEmpty(t, md.Tables)
// 		assert.NotEmpty(t, md.ForeignKeys)
// 		assert.Equal(t, 2, len(md.Tables))
// 		assert.Equal(t, 5, len(md.ForeignKeys))
// 	})
// }

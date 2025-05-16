package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jboursiquot/mermaid-mcp/tools/erd"
	mcp "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/http"
)

func main() {
	// Set up channel to listen for interrupt signals
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	erdGenTool, err := erd.NewGenerator(db)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create ERD generator: %v\n", err)
		os.Exit(1)
	}

	// Use HTTP transport with default options (port 8080)
	transport := http.NewHTTPTransport("/mcp").WithAddr(":8080")

	slog.Info("Starting server...")
	server := mcp.NewServer(
		transport,
		mcp.WithName("mermaid-mcp"),
		mcp.WithInstructions("MCP server for generating Mermaid diagrams"),
		mcp.WithVersion("0.0.1"),
	)

	slog.Info("Registering tools")
	if err := server.RegisterTool(erdGenTool.Name, erdGenTool.Description, erdGenTool.Generate); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to register tool: %v\n", err)
		os.Exit(1)
	}

	// Start server in goroutine
	go func() {
		if err := server.Serve(); err != nil {
			fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
			os.Exit(1)
		}
	}()

	slog.Info("Server started at http://localhost:8080...")

	// Wait for termination signal
	<-done
	slog.Info("Shutting down server...")
}

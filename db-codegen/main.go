package main

import (
	"log/slog"

	"db-codegen/generator"
)

func main() {
	gen := &generator.CodeGenerator{
		ConnString: "host=localhost user=postgres password=password dbname=postgres port=5432 sslmode=disable",
		TempDB:     "gopher_patterns_gen",
	}

	if err := gen.Run(); err != nil {
		slog.Error("Code generation failed", "error", err)
		return
	}
}

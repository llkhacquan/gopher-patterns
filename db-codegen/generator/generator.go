package generator

import (
	"fmt"
	"log/slog"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type CodeGenerator struct {
	ConnString string
	TempDB     string
}

func (c *CodeGenerator) Run() error {
	slog.Info("Starting database code generation")

	// Connect to admin database
	gormDB, err := gorm.Open(postgres.Open(c.ConnString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return fmt.Errorf("could not connect to db: %v", err)
	}

	// Drop and create temporary database
	if err := gormDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", c.TempDB)).Error; err != nil {
		slog.Warn("drop database error", "error", err)
	}
	if err := gormDB.Exec(fmt.Sprintf("CREATE DATABASE %s", c.TempDB)).Error; err != nil {
		return fmt.Errorf("create database error: %v", err)
	}
	defer gormDB.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", c.TempDB))

	// Connect to temporary database
	tempConnString := fmt.Sprintf("host=localhost user=postgres password=password dbname=%s port=5432 sslmode=disable", c.TempDB)
	tempDB, err := gorm.Open(postgres.Open(tempConnString), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return fmt.Errorf("could not open temp gorm db: %v", err)
	}

	// Create database schema
	if err := c.createSchema(tempDB); err != nil {
		return err
	}

	// Generate code
	if err := c.generateCode(tempDB); err != nil {
		return err
	}

	slog.Info("Code generation completed")

	// Close database connection before cleanup
	if sqlDB, err := tempDB.DB(); err == nil {
		sqlDB.Close()
	}

	return nil
}

// createSchema creates dummy tables for code generation only. In real projects, you should use your actual database schema.
func (c *CodeGenerator) createSchema(db *gorm.DB) error {
	if err := db.Exec(`
		CREATE TABLE users (
			id BIGSERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			email VARCHAR(100) UNIQUE NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create users table: %v", err)
	}

	if err := db.Exec(`
		CREATE TABLE orders (
			id BIGSERIAL PRIMARY KEY,
			user_id BIGINT NOT NULL,
			product VARCHAR(100) NOT NULL,
			quantity INTEGER NOT NULL DEFAULT 1,
			price DECIMAL(10,2) NOT NULL,
			status VARCHAR(20) DEFAULT 'pending',
			created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
		)
	`).Error; err != nil {
		return fmt.Errorf("failed to create orders table: %v", err)
	}

	return nil
}

func (c *CodeGenerator) generateCode(db *gorm.DB) error {
	var genConfig = gen.Config{
		OutPath:           "query",
		OutFile:           "gen.go",
		FieldSignable:     false,
		FieldWithIndexTag: false,
		FieldWithTypeTag:  true,
		Mode:              gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	}

	g := gen.NewGenerator(genConfig)
	g.UseDB(db)

	user := g.GenerateModel("users")
	order := g.GenerateModel("orders")

	g.ApplyBasic(user, order)
	g.Execute()

	return nil
}

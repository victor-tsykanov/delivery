package testutils

import (
	"context"
	"fmt"
	"log"

	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/victor-tsykanov/delivery/migrations"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBTestSuite struct {
	suite.Suite
	postgresContainer *postgres.PostgresContainer
	db                *gorm.DB
	ctx               context.Context
}

func (s *DBTestSuite) SetupSuite() {
	s.ctx = context.Background()

	container, err := s.startPostgresContainer()
	if err != nil {
		log.Fatalf("failed to start Postgres conteiner: %v", err)
	}

	s.postgresContainer = container

	db, err := s.connect()
	if err != nil {
		s.TearDownSuite()
		log.Fatal(err)
	}

	s.db = db
}

func (s *DBTestSuite) TearDownSuite() {
	if s.postgresContainer != nil {
		err := s.postgresContainer.Terminate(s.ctx)
		if err != nil {
			fmt.Printf("failed to terminate Postgres conteiner: %v", err)
		}
	}
}

func (s *DBTestSuite) SetupTest() {
	err := s.applyMigrations()
	if err != nil {
		s.Require().NoError(err)
	}
}

func (s *DBTestSuite) TearDownTest() {
	err := s.resetDB()
	if err != nil {
		s.Require().NoError(err)
	}
}

func (s *DBTestSuite) DB() *gorm.DB {
	return s.db
}

func (s *DBTestSuite) connect() (*gorm.DB, error) {
	connectionString, err := s.postgresContainer.ConnectionString(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	db, err := gorm.Open(gormPostgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return db, nil
}

func (s *DBTestSuite) startPostgresContainer() (*postgres.PostgresContainer, error) {
	container, err := postgres.Run(
		s.ctx,
		"postgres:17.4-alpine",
		postgres.WithDatabase("delivery"),
		postgres.WithUsername("delivery"),
		postgres.WithPassword("secret"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, err
	}

	return container, nil
}

func (s *DBTestSuite) applyMigrations() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	goose.SetBaseFS(migrations.EmbeddedMigrations)

	err = goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	err = goose.Up(sqlDB, ".")
	if err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	return nil
}

func (s *DBTestSuite) resetDB() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return fmt.Errorf("failed to reset database: %w", err)
	}

	goose.SetBaseFS(migrations.EmbeddedMigrations)

	err = goose.SetDialect("postgres")
	if err != nil {
		return fmt.Errorf("failed to reset database: %w", err)
	}

	err = goose.Reset(sqlDB, ".")
	if err != nil {
		return fmt.Errorf("failed to reset database: %w", err)
	}

	return nil
}

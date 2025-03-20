package testing

import (
	"context"
	"fmt"
	"log"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	gormPostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type RepositoryTestSuite struct {
	suite.Suite
	postgresContainer *postgres.PostgresContainer
	db                *gorm.DB
	ctx               context.Context
}

func (s *RepositoryTestSuite) SetupSuite() {
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

func (s *RepositoryTestSuite) TearDownSuite() {
	if s.postgresContainer != nil {
		err := s.postgresContainer.Terminate(s.ctx)
		if err != nil {
			fmt.Printf("failed to terminate Postgres conteiner: %v", err)
		}
	}
}

func (s *RepositoryTestSuite) DB() *gorm.DB {
	return s.db
}

func (s *RepositoryTestSuite) connect() (*gorm.DB, error) {
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

func (s *RepositoryTestSuite) startPostgresContainer() (*postgres.PostgresContainer, error) {
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

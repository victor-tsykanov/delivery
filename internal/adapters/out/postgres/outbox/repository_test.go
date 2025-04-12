package outbox_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/victor-tsykanov/delivery/internal/adapters/out/postgres/outbox"
	"github.com/victor-tsykanov/delivery/internal/common/errors"
	"github.com/victor-tsykanov/delivery/internal/common/testutils"
)

type OutboxRepositoryTestSuite struct {
	testutils.DBTestSuite
	repository *outbox.Repository
}

func (s *OutboxRepositoryTestSuite) SetupTest() {
	s.DBTestSuite.SetupTest()

	repository, err := outbox.NewRepository(s.DB())
	s.Require().NoError(err)

	s.repository = repository
}

func TestCourierRepositoryTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(OutboxRepositoryTestSuite))
}

func (s *OutboxRepositoryTestSuite) TestCreate() {
	// Arrange
	id := uuid.New()
	messageType := "message-type"
	payload := []byte("payload")
	message := &outbox.Message{
		ID:      id,
		Type:    messageType,
		Payload: payload,
	}

	// Act
	err := s.repository.Create(context.Background(), message)

	// Assert
	s.Require().NoError(err)

	messageFromDB := &outbox.Message{}
	err = s.DB().
		First(messageFromDB, "id = ?", id).
		Error
	s.Require().NoError(err)

	s.Assert().Equal(id, messageFromDB.ID)
	s.Assert().Equal(messageType, messageFromDB.Type)
	s.Assert().Equal(payload, messageFromDB.Payload)
	s.Assert().Equal(sql.NullTime{}, messageFromDB.ProcessedAt)
}

func (s *OutboxRepositoryTestSuite) TestMarkProcessedExistingMessage() {
	// Arrange
	ctx := context.Background()
	message := &outbox.Message{ID: uuid.New()}
	s.DB().Create(message)

	// Act
	err := s.repository.MarkProcessed(ctx, message.ID)

	// Assert
	s.Require().NoError(err)

	var messageFromDB = &outbox.Message{}
	s.DB().Find(messageFromDB, message.ID)
	s.Assert().True(messageFromDB.ProcessedAt.Valid)
	s.Assert().WithinDuration(time.Now(), messageFromDB.ProcessedAt.Time, time.Second)
}

func (s *OutboxRepositoryTestSuite) TestMarkProcessedNonexistentMessage() {
	// Arrange
	ctx := context.Background()
	id := uuid.New()

	// Act
	err := s.repository.MarkProcessed(ctx, id)

	// Assert
	s.Require().Equal(errors.NewEntityNotfoundError("message", id), err)
}

func (s *OutboxRepositoryTestSuite) TestFindUnprocessed() {
	// Arrange
	messages := []*outbox.Message{
		{ID: uuid.New()},
		{ID: uuid.New(), ProcessedAt: sql.NullTime{Valid: true, Time: time.Now()}},
		{ID: uuid.New()},
	}
	s.DB().Create(messages)

	// Act
	unprocessedMessages, hasMore, err := s.repository.FindUnprocessed(10)

	// Assert
	s.Require().NoError(err)
	s.Assert().False(hasMore)
	s.Assert().Equal(messages[0].ID, unprocessedMessages[0].ID)
	s.Assert().Equal(messages[2].ID, unprocessedMessages[1].ID)
}

func (s *OutboxRepositoryTestSuite) TestFindUnprocessedWithLimitOverflow() {
	// Arrange
	messages := []*outbox.Message{
		{ID: uuid.New()},
		{ID: uuid.New(), ProcessedAt: sql.NullTime{Valid: true, Time: time.Now()}},
		{ID: uuid.New()},
	}
	s.DB().Create(messages)

	// Act
	unprocessedMessages, hasMore, err := s.repository.FindUnprocessed(1)

	// Assert
	s.Require().NoError(err)
	s.Assert().True(hasMore)
	s.Assert().Equal(messages[0].ID, unprocessedMessages[0].ID)
}

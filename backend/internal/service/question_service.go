package service

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/gbplantwiki/gbplantwiki/internal/constants"
	"github.com/gbplantwiki/gbplantwiki/internal/model"
	"github.com/gbplantwiki/gbplantwiki/internal/repository"
	"github.com/gbplantwiki/gbplantwiki/internal/util"
)

// QuestionService implements Q&A question logic.
type QuestionService struct {
	repo        *repository.QuestionRepository
	answerRepo  *repository.AnswerRepository
	userService *UserService
	logger      *slog.Logger
}

// NewQuestionService creates a QuestionService.
func NewQuestionService(repo *repository.QuestionRepository, answerRepo *repository.AnswerRepository, userService *UserService, logger *slog.Logger) *QuestionService {
	return &QuestionService{repo: repo, answerRepo: answerRepo, userService: userService, logger: logger}
}

// Create publishes a question for the current user.
func (s *QuestionService) Create(userID uint, q *model.Question) (*model.Question, error) {
	q.UserID = userID
	if q.Images == "" {
		q.Images = "[]"
	}
	if q.Status == "" {
		q.Status = "open"
	}
	if err := s.repo.Create(q); err != nil {
		s.logger.Error(fmt.Sprintf(constants.LogQuestionCreateFailed, q.Title), "error", err)
		return nil, fmt.Errorf("question create: %w", err)
	}
	s.logger.Info(fmt.Sprintf(constants.LogQuestionCreateSuccess, q.Title), "id", q.ID)
	return q, nil
}

// Get returns a question by id.
func (s *QuestionService) Get(id uint) (*model.Question, error) {
	q, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return nil, util.NewAppError(404, constants.CodeNotFound, fmt.Sprintf("Question[id=%d] not found", id))
		}
		return nil, fmt.Errorf("question get: %w", err)
	}
	return q, nil
}

// List returns paginated questions.
func (s *QuestionService) List(page, pageSize int) ([]model.Question, int64, error) {
	items, total, err := s.repo.List(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("question list: %w", err)
	}
	return items, total, nil
}

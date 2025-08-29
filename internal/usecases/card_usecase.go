package usecases

import (
	"errors"
	"time"

	"interview-tracker/internal/adapters/repositories"
	"interview-tracker/internal/entities"
	"interview-tracker/internal/pkg/errs"

	"github.com/google/uuid"
)

type CardUsecase struct {
	repo repositories.CardRepository
}

func NewCardUsecase(r repositories.CardRepository) *CardUsecase { return &CardUsecase{repo: r} }

func (uc *CardUsecase) Create(card *entities.Card) error {
	var txnDtm = time.Now()
	card.StatusCode = "todo"
	card.CreatedAt = txnDtm
	card.UpdatedAt = txnDtm
	if err := uc.repo.Create(card); err != nil {
		return err
	}
	if err := uc.addHistory(card.CreatedBy, card.ID, card.StatusCode, card.Description); err != nil {
		return err
	}
	return nil
}

func (uc *CardUsecase) UpdatePartial(id uuid.UUID, patch map[string]any) (*entities.Card, error) {
	card, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if v, ok := patch["title"].(string); ok {
		card.Title = v
	}
	if v, ok := patch["description"].(string); ok {
		card.Description = v
	}
	if v, ok := patch["scheduled_at"].(time.Time); ok {
		card.ScheduledAt = v
	}
	card.UpdatedAt = time.Now()
	if err := uc.repo.Update(card); err != nil {
		return nil, err
	}

	if err := uc.addHistory(card.CreatedBy, card.ID, card.StatusCode, card.Description); err != nil {
		return nil, err
	}
	return card, nil
}

func (uc *CardUsecase) UpdateStatus(id uuid.UUID, status string, actor uuid.UUID) (*entities.Card, error) {
	if status != "todo" && status != "in_progress" && status != "done" {
		return nil, errors.New("invalid status")
	}
	card, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	card.StatusCode = status
	card.UpdatedAt = time.Now()
	if err := uc.repo.Update(card); err != nil {
		return nil, err
	}

	// add history
	if err := uc.addHistory(actor, id, card.StatusCode, card.Description); err != nil {
		return nil, err
	}
	return card, nil
}

func (uc *CardUsecase) GetByID(id uuid.UUID) (*entities.Card, error) { return uc.repo.GetByID(id) }

func (uc *CardUsecase) List(status string, page, size int) ([]*entities.Card, int64, error) {
	return uc.repo.List(status, page, size)
}

func (uc *CardUsecase) AddComment(authorID uuid.UUID, cardID, content string) error {
	cmt := &entities.CardComment{
		CardID:    cardID,
		AuthorID:  authorID,
		Content:   content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		CreatedBy: authorID,
		UpdatedBy: authorID,
	}
	return uc.repo.AddComment(cmt)
}

func (uc *CardUsecase) UpdateComment(authorID, commentId uuid.UUID, content string) error {
	comment, err := uc.repo.GetCommentByID(commentId)
	if err != nil {
		return err
	}
	if comment.AuthorID != authorID {
		return errs.Unauthorized("unauthorized")
	}
	cmt := &entities.CardComment{
		CardID:    comment.CardID,
		AuthorID:  authorID,
		Content:   content,
		UpdatedAt: time.Now(),
		UpdatedBy: authorID,
	}
	return uc.repo.UpdateComment(cmt)
}

func (uc *CardUsecase) ListComments(cardID string, page, size int) ([]*entities.CardComment, int64, error) {
	return uc.repo.ListComments(cardID, page, size)
}

func (uc *CardUsecase) addHistory(actorID, cardID uuid.UUID, statusCode, description string) error {
	p := &entities.CardHistoryLogs{
		CardID:      cardID,
		ActorID:     actorID,
		StatusCode:  statusCode,
		Description: description,
		CreatedBy:   actorID,
		UpdatedBy:   actorID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return uc.repo.AddHistory(p)
}

func (uc *CardUsecase) ListHistory(cardID string, page, size int) ([]*entities.CardHistoryLogs, int64, error) {
	return uc.repo.ListHistory(cardID, page, size)
}

func (uc *CardUsecase) Keep(cardID string) error {
	return uc.repo.DeleteCard(cardID)
}

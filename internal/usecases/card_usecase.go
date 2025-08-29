package usecases

import (
	"errors"
	"time"

	"interview-tracker/internal/adapters/repositories"
	"interview-tracker/internal/entities"

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
	return uc.repo.Create(card)
}

func (uc *CardUsecase) UpdatePartial(id string, patch map[string]any) (*entities.Card, error) {
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
	return card, nil
}

func (uc *CardUsecase) UpdateStatus(id, status string, actor uuid.UUID) (*entities.Card, error) {
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

	// log progress
	_ = uc.repo.AddProgress(&entities.CardProgressLogs{
		CardID:    id,
		ActorID:   actor,
		Message:   "status changed to " + status,
		CreatedAt: time.Now(),
	})
	return card, nil
}

func (uc *CardUsecase) GetByID(id string) (*entities.Card, error) { return uc.repo.GetByID(id) }

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
	}
	return uc.repo.AddComment(cmt)
}

func (uc *CardUsecase) ListComments(cardID string, page, size int) ([]*entities.CardComment, int64, error) {
	return uc.repo.ListComments(cardID, page, size)
}

func (uc *CardUsecase) AddProgress(actorID uuid.UUID, cardID, msg string) error {
	p := &entities.CardProgressLogs{
		CardID:    cardID,
		ActorID:   actorID,
		Message:   msg,
		CreatedAt: time.Now(),
	}
	return uc.repo.AddProgress(p)
}

func (uc *CardUsecase) ListProgress(cardID string, page, size int) ([]*entities.CardProgressLogs, int64, error) {
	return uc.repo.ListProgress(cardID, page, size)
}

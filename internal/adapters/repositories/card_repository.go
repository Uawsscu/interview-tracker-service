package repositories

import (
	"interview-tracker/internal/entities"

	"gorm.io/gorm"
)

type CardRepository interface {
	Create(card *entities.Card) error
	Update(card *entities.Card) error
	GetByID(id string) (*entities.Card, error)
	List(status string, page, size int) ([]*entities.Card, int64, error)

	AddComment(c *entities.CardComment) error
	ListComments(cardID string, page, size int) ([]*entities.CardComment, int64, error)

	AddProgress(p *entities.CardProgressLogs) error
	ListProgress(cardID string, page, size int) ([]*entities.CardProgressLogs, int64, error)
}

type cardRepo struct{ db *gorm.DB }

func NewCardRepo(db *gorm.DB) CardRepository { return &cardRepo{db} }

func (r *cardRepo) Create(c *entities.Card) error { return r.db.Create(c).Error }

func (r *cardRepo) Update(c *entities.Card) error { return r.db.Save(c).Error }

func (r *cardRepo) GetByID(id string) (*entities.Card, error) {
	var card entities.Card
	if err := r.db.First(&card, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *cardRepo) List(status string, page, size int) ([]*entities.Card, int64, error) {
	var list []*entities.Card
	var total int64
	qb := r.db.Model(&entities.Card{})
	if status != "" {
		qb = qb.Where("status_code = ?", status)
	}
	if err := qb.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := qb.Order("scheduled_at asc").
		Limit(size).
		Offset((page - 1) * size).
		Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *cardRepo) AddComment(cmt *entities.CardComment) error {
	return r.db.Create(cmt).Error
}

func (r *cardRepo) ListComments(cardID string, page, size int) ([]*entities.CardComment, int64, error) {
	var list []*entities.CardComment
	var total int64
	qb := r.db.Model(&entities.CardComment{}).Where("card_id = ?", cardID)
	if err := qb.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := qb.Order("created_at desc").Limit(size).Offset((page - 1) * size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *cardRepo) AddProgress(p *entities.CardProgressLogs) error {
	return r.db.Create(p).Error
}

func (r *cardRepo) ListProgress(cardID string, page, size int) ([]*entities.CardProgressLogs, int64, error) {
	var list []*entities.CardProgressLogs
	var total int64
	qb := r.db.Model(&entities.CardProgressLogs{}).Where("card_id = ?", cardID)
	if err := qb.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := qb.Order("created_at desc").Limit(size).Offset((page - 1) * size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

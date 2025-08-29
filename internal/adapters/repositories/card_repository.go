package repositories

import (
	"interview-tracker/internal/entities"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CardRepository interface {
	Create(card *entities.Card) error
	Update(card *entities.Card) error
	GetByID(id uuid.UUID) (*entities.Card, error)
	List(status string, page, size int) ([]*entities.Card, int64, error)
	DeleteCard(id string) error

	AddComment(c *entities.CardComment) error
	UpdateComment(c *entities.CardComment) error
	GetCommentByID(commentId uuid.UUID) (*entities.CardComment, error)
	ListComments(cardID string, page, size int) ([]*entities.CardComment, int64, error)

	AddHistory(p *entities.CardHistoryLogs) error
	ListHistory(cardID string, page, size int) ([]*entities.CardHistoryLogs, int64, error)
}

type cardRepo struct{ db *gorm.DB }

func NewCardRepo(db *gorm.DB) CardRepository { return &cardRepo{db} }

func (r *cardRepo) Create(c *entities.Card) error { return r.db.Create(c).Error }

func (r *cardRepo) Update(c *entities.Card) error { return r.db.Save(c).Error }

func (r *cardRepo) GetByID(id uuid.UUID) (*entities.Card, error) {
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

func (r *cardRepo) AddHistory(p *entities.CardHistoryLogs) error {
	return r.db.Create(p).Error
}

func (r *cardRepo) ListHistory(cardID string, page, size int) ([]*entities.CardHistoryLogs, int64, error) {
	var list []*entities.CardHistoryLogs
	var total int64
	qb := r.db.Model(&entities.CardHistoryLogs{}).Where("card_id = ?", cardID)
	if err := qb.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := qb.Order("created_at desc").Limit(size).Offset((page - 1) * size).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *cardRepo) DeleteCard(id string) error {
	return r.db.Delete(&entities.Card{}, "id = ?", id).Error
}

func (r *cardRepo) UpdateComment(c *entities.CardComment) error {
	return r.db.Model(&entities.CardComment{}).
		Where("id = ?", c.ID).
		Updates(map[string]interface{}{
			"content":    c.Content,
			"author_id":  c.AuthorID,
			"updated_at": c.UpdatedAt,
			"updated_by": c.UpdatedBy,
		}).Error
}

func (r *cardRepo) GetCommentByID(commentId uuid.UUID) (*entities.CardComment, error) {
	var comment entities.CardComment
	if err := r.db.First(&comment, "id = ?", commentId).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}

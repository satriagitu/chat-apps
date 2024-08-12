package repository

import (
	"chat-apps/internal/domain"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	GetArticleList() ([]domain.ArticleList, error)
}

type articleRepository struct {
	db *gorm.DB
}

func NewArtikelRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) GetArticleList() ([]domain.ArticleList, error) {
	articles := []domain.ArticleList{}
	query := `SELECT a.id, m.name AS menu, sm.name AS sub_menu, a.title, a.image, a.likes, c.comment_count, 
				a.created_at
				FROM articles AS a
				LEFT JOIN menus m ON a.menu_id = m.id
				LEFT JOIN sub_menus sm ON a.sub_menu_id = sm.id
				LEFT JOIN (
					SELECT article_id, COUNT(id) AS comment_count
					FROM comments
					GROUP BY article_id
				) c ON a.id = c.article_id
				ORDER BY a.id desc`

	if err := r.db.Raw(query).Scan(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

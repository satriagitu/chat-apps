package repository

import (
	"chat-apps/internal/domain"
	"fmt"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	GetArticleList(search string) ([]domain.ArticleList, error)
}

type articleRepository struct {
	db *gorm.DB
}

func NewArtikelRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) GetArticleList(search string) ([]domain.ArticleList, error) {
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
				`

	// Tambahkan kondisi pencarian jika parameter search tidak kosong
	if search != "" {
		search = "%" + search + "%" // Menambahkan wildcard untuk LIKE
		query += fmt.Sprintf(" WHERE a.title ILIKE '%s' OR m.name ILIKE '%s' OR sm.name ILIKE '%s'", search, search, search)
	}

	// Menambahkan ORDER BY
	query += " ORDER BY a.id DESC"

	if err := r.db.Raw(query).Scan(&articles).Error; err != nil {
		return nil, err
	}

	return articles, nil
}

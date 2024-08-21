package service

import (
	"chat-apps/internal/domain"
	"chat-apps/internal/repository"
	"chat-apps/internal/util"
)

type ArticleService interface {
	GetArticleList(search string) ([]domain.ArticleList, error)
}

type artikelService struct {
	articleRepo repository.ArticleRepository
}

func NewArtikelService(ar repository.ArticleRepository) ArticleService {
	return &artikelService{articleRepo: ar}
}

func (s *artikelService) GetArticleList(search string) ([]domain.ArticleList, error) {
	articles, err := s.articleRepo.GetArticleList(search)
	if err != nil {
		return nil, err
	}

	articleList := make([]domain.ArticleList, len(articles))

	for i, article := range articles {
		articleList[i] = domain.ArticleList{
			ID:           article.ID,
			Menu:         article.Menu,
			SubMenu:      article.SubMenu,
			Title:        article.Title,
			Image:        article.Image,
			Likes:        article.Likes,
			CommentCount: article.CommentCount,
			TimeAgo:      util.TimeAgo(article.CreatedAt),
		}
	}

	return articleList, nil

}

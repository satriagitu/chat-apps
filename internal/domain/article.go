package domain

import "time"

type ArticleList struct {
	ID           int       `json:"id"`
	Menu         string    `json:"menu"`
	SubMenu      string    `json:"sub_menu"`
	Title        string    `json:"title"`
	Image        string    `json:"image"`
	TimeAgo      string    `json:"time_ago"`
	CreatedAt    time.Time `json:"-"`
	Likes        int       `json:"likes"`
	CommentCount int       `json:"comment_count"`
}

type Article struct {
	ID        int       `json:"id"`
	MenuID    int       `json:"menu_id"`
	SubMenuID int       `json:"sub_menu_id"`
	Title     string    `json:"title"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
	Likes     int       `json:"likes"`
}

type Menu struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type SubMenu struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Comment struct {
	ID        int    `json:"id"`
	ArticleID int    `json:"article_id"`
	Content   string `json:"content"`
}

package infra

import (
	"context"
	"fmt"

	"github.com/go-gorp/gorp"
	"github.com/openhacku-saboten/OmnisCode-backend/domain/entity"
	"github.com/openhacku-saboten/OmnisCode-backend/repository"
)

var _ repository.Post = (*PostRepository)(nil)

// PostRepository はrepository.UserRepositoryを満たすstructです
type PostRepository struct {
	dbMap *gorp.DbMap
}

// NewPostRespository はPostRepositoryに対するポインタを生成します
func NewPostRepository(dbMap *gorp.DbMap) *PostRepository {
	return &PostRepository{dbMap: dbMap}
}

// GetAll はMySQLサーバに接続して、全てのPostを取得して返すメソッドです
func (p *PostRepository) GetAll(context.Context) ([]*entity.Post, error) {
	var resDTO []*entity.Post

	const query = `SELECT * FROM posts`
	rows, err := p.dbMap.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed dbMap.Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var postDTO PostDTO
		if err := rows.Scan(&postDTO); err != nil {
			return nil, fmt.Errorf("failed dbMap.Scan: %w", err)
		}
		resDTO = append(resDTO, &entity.Post{
			ID:       postDTO.ID,
			UserID:   postDTO.UserID,
			Title:    postDTO.Title,
			Code:     postDTO.Code,
			Language: postDTO.Language,
			Content:  postDTO.Content,
			Source:   postDTO.Source,
		})
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to Err: %w", err)
	}
	return resDTO, nil
}

// PostDTO はDBとやり取りするためのDataTransferObjectです。
type PostDTO struct {
	ID       string `db:"id"`
	UserID   string `db:"user_id"`
	Title    string `db:"title"`
	Code     string `db:"code"`
	Language string `db:"language"`
	Content  string `db:"content"`
	Source   string `db:"source"`
}

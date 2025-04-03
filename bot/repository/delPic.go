package repository

import (
	"database/sql"
	"errors"
)

func (p *PicRepository) DelPic(UserID int64, tag string) (string, error) {
	var picKey string

	q := `DELETE FROM user_pictures WHERE user_id = $1 AND tag = $2 RETURNING pic_key`

	err := p.db.QueryRow(q, UserID, tag).Scan(&picKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", PicNotFound
		}
		return "", err
	}

	return picKey, nil
}

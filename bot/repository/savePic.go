package repository

import "github.com/lib/pq"

func (p *PicRepository) SavePic(userID int64, picKey string, tag string) error {
	q := `INSERT INTO user_pictures (user_id, pic_key, tag) VALUES ($1, $2, $3)`

	_, err := p.db.Exec(q, userID, picKey, tag)
	if err != nil {
		if isDuplicateEntryError(err) {
			return PicExists
		}
		return err
	}

	return nil
}

func isDuplicateEntryError(err error) bool {
	if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
		return true
	}

	return false
}

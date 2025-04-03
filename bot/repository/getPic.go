package repository

import (
	"MemoryPicBot/bot"
	"database/sql"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
)

func (p *PicRepository) GetPic(UserID int64, tag string) (string, error) {
	var picKey string

	q := `SELECT pic_key FROM user_pictures WHERE user_id = $1 AND tag = $2`

	err := p.db.QueryRow(q, UserID, tag).Scan(&picKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", PicNotFound
		}
		return "", err
	}

	return picKey, nil
}

func (p *PicRepository) GetPicList() ([]bot.UserPictures, error) {
	var picList []bot.UserPictures

	users, err := p.getAllUsers()
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	for _, userID := range users {
		userPic, err := p.getRandomUserPic(userID)
		if err != nil {
			logrus.Warnf("failed to get random picture for user %d: %v", userID, err)

			continue
		}

		picList = append(picList, userPic)
	}

	return picList, nil
}

func (p *PicRepository) getAllUsers() ([]int64, error) {
	var users []int64
	q := `SELECT DISTINCT user_id FROM user_pictures`

	rows, err := p.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID int64
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		users = append(users, userID)
	}

	return users, nil
}

func (p *PicRepository) getRandomUserPic(userId int64) (bot.UserPictures, error) {
	var userPic bot.UserPictures
	q := `SELECT user_id, pic_key, tag FROM user_pictures WHERE user_id = $1 ORDER BY RANDOM() LIMIT 1`

	if err := p.db.QueryRow(q, userId).Scan(&userPic.UserId, &userPic.PicKey, &userPic.Tag); err != nil {
		return bot.UserPictures{}, err
	}

	return userPic, nil
}

func (p *PicRepository) GetAllTags(userID int64) ([]string, error) {
	var tags []string

	logrus.Infof("USERID %d get tags", userID)
	if !p.userExists(userID) {
		return nil, UserNotFound
	}

	q := `SELECT tag FROM user_pictures WHERE user_id = $1`

	rows, err := p.db.Query(q, userID)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (p *PicRepository) userExists(userID int64) bool {
	var exists bool

	err := p.db.QueryRow("SELECT EXISTS(SELECT 1 FROM user_pictures WHERE user_id = $1)", userID).Scan(&exists)
	if err != nil {
		logrus.Warnf("failed to check user existence: %v", err)

		return false
	}

	return exists
}

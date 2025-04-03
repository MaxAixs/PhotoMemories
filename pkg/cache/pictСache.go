package cache

import (
	"context"
	"strconv"
	"time"
)

func (c *Client) SavePicture(userID int64, picKey string) error {
	return c.rdb.Set(context.Background(), strconv.FormatInt(userID, 10), picKey, 5*time.Minute).Err()
}

func (c *Client) GetPicture(userID int64) (string, error) {
	return c.rdb.Get(context.Background(), strconv.FormatInt(userID, 10)).Result()
}

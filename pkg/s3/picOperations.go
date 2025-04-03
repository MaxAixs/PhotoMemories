package s3

import (
	"bytes"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"golang.org/x/net/context"
	"io"
)

func (c *Client) UploadPic(picKey string, fileDate []byte) error {
	ctx := context.TODO()

	_, err := c.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &c.bucketName,
		Key:    &picKey,
		Body:   bytes.NewReader(fileDate),
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeletePic(picKey string) error {
	ctx := context.TODO()

	_, err := c.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &c.bucketName,
		Key:    &picKey,
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DownLoadPic(picKey string) ([]byte, error) {
	ctx := context.TODO()

	resp, err := c.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &c.bucketName,
		Key:    &picKey,
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

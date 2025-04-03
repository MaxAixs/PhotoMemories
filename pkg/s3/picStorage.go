package s3

type PhotoStorage interface {
	UploadPic(picKey string, fileDate []byte) error
	DeletePic(picKey string) error
	DownLoadPic(picKey string) ([]byte, error)
}

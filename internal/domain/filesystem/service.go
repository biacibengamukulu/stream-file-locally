package filesystem

import "github.com/biacibengamukulu/stream-file-locally/sharded/config"

type FileService interface {
	Upload(file FileDTO) (File, error)
	Get(id string) (File, error)
}

type FileServiceImpl struct {
	cfg config.Config
}

func NewFileService(cfg config.Config) *FileServiceImpl {
	return &FileServiceImpl{cfg: cfg}
}

func (f FileServiceImpl) Upload(file FileDTO) (File, error) {
	//TODO implement me
	panic("implement me")
}

func (f FileServiceImpl) Get(id string) (File, error) {
	//TODO implement me
	panic("implement me")
}

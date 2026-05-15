package filesystem

type Storage interface {
	Save(file StoredFile) (File, error)
	Get(id string) (File, error)
	Read(id string) (StoredFile, error)
	Delete(id string) error
}

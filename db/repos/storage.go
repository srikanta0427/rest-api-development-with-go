package db

type Storage struct {
	UserRepository UserRepository
	OrganizerRepo  OrganizerRepo
}

func NewStorage() *Storage {
	return &Storage{
		UserRepository: &UserRepositoryImpl{},
		OrganizerRepo:  &OrganizerImpl{},
	}
}

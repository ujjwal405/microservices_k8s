package user

type queue interface {
	AddToQueue(mail Mail) error
}
type Repo interface {
	CountDocuments(email string) (bool, error)
	InsertUser(user User) error
	FetchUser(email string) (User, error)
}
type Cache interface {
	AddClient(code string, detail *Details)
	Loop()
	Check(code string) (Details, bool)
	Awake()
}

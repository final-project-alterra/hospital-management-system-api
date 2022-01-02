package auth

type IBusiness interface {
	Login(email string, password string) (string, error)
}

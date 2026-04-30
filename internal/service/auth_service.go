package service

type IAuthService interface {
	RegisterNewCustomer(string, string, string, string, string) error
	Login(string, string) (int, string, string, error)
	RefreshToken(string) (string, error)
	Logout(string) error
}

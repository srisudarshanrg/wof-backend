package app

import "golang.org/x/crypto/bcrypt"

func (app *Application) HashPassword(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", nil
	}
	return string(hashPassword), nil
}

func (app *Application) ComparePasswordWithHash(password string, hashPassword string) (error) {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}

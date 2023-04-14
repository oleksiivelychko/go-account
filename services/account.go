package services

import (
	"errors"
	"github.com/badoux/checkmail"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type Account struct {
	repo *repositories.Account
}

func NewAccount(accountRepo *repositories.Account) *Account {
	return &Account{accountRepo}
}

/*
VerifyPassword
modelAccount *models.Account has hashed password
password string has non-hashed password
*/
func (service *Account) VerifyPassword(account *models.Account, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
}

func (service *Account) Serialize(account *models.Account) *models.AccountSerialized {
	accountSerialized := &models.AccountSerialized{
		ID:    account.ID,
		Email: account.Email,
	}

	for _, role := range account.Roles {
		accountSerialized.Roles = append(accountSerialized.Roles, models.RoleSerialized{Name: role.Name})
	}

	return accountSerialized
}

func (service *Account) Validate(account *models.Account) error {
	if account.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		account.Password = string(hashedPassword)
	}

	if account.Email != "" {
		if err := checkmail.ValidateFormat(account.Email); err != nil {
			return err
		}

		accountFound, err := service.repo.FindOneByEmail(account.Email, false)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if accountFound.Email != "" {
			return errors.New("account already exists")
		}
	}

	return nil
}

func (service *Account) Create(account *models.Account) (*models.Account, error) {
	err := service.Validate(account)
	if err != nil {
		return &models.Account{}, err
	}

	err = service.repo.Create(account)
	if err != nil {
		return &models.Account{}, err
	}

	return account, nil
}

func (service *Account) Update(account *models.Account) (*models.Account, error) {
	err := service.Validate(account)
	if err != nil {
		return &models.Account{}, err
	}

	err = service.repo.Update(account, map[string]interface{}{
		"password":   account.Password,
		"email":      account.Email,
		"updated_at": time.Now(),
	})

	if err != nil {
		return &models.Account{}, err
	}

	return account, nil
}

func (service *Account) Delete(account *models.Account) (int64, error) {
	return service.repo.Delete(account)
}

func (service *Account) Auth(email, password string) (*models.AccountSerialized, error) {
	account, err := service.repo.FindOneByEmail(email, true)

	if err == nil {
		err = service.VerifyPassword(account, password)
		if err != nil {
			return nil, err
		}
	}

	return service.Serialize(account), err
}

func (service *Account) HasRoles(account *models.Account, roles []string) bool {
	return service.repo.HasRoles(account, roles)
}

func (service *Account) GetRepository() *repositories.Account {
	return service.repo
}

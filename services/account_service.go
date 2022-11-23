package services

import (
	"errors"
	"fmt"
	"github.com/badoux/checkmail"
	"github.com/oleksiivelychko/go-account/models"
	"github.com/oleksiivelychko/go-account/repositories"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type AccountService struct {
	repository *repositories.AccountRepository
}

func NewAccountService(ar *repositories.AccountRepository) *AccountService {
	return &AccountService{ar}
}

func (service *AccountService) VerifyPassword(model *models.Account, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(model.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}

func (service *AccountService) Serialize(model *models.Account) *models.AccountSerialized {
	accountSerialized := &models.AccountSerialized{
		ID:    model.ID,
		Email: model.Email,
	}

	for _, role := range model.Roles {
		accountSerialized.Roles = append(accountSerialized.Roles, models.RoleSerialized{Name: role.Name})
	}

	return accountSerialized
}

func (service *AccountService) Validate(account *models.Account) error {
	if account.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		account.Password = string(hashedPassword)
	}

	if account.Email != "" {
		if err := checkmail.ValidateFormat(account.Email); err != nil {
			return errors.New("invalid:email")
		}

		existsAccount, err := service.repository.FindOneByEmail(account.Email, false)
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}

		if existsAccount.Email != "" {
			return errors.New("email address already exists")
		}
	}

	return nil
}

func (service *AccountService) Create(account *models.Account) (*models.Account, error) {
	err := service.Validate(account)
	if err != nil {
		return &models.Account{}, err
	}

	return service.repository.Create(account)
}

func (service *AccountService) Update(account *models.Account) (*models.Account, error) {
	err := service.Validate(account)
	if err != nil {
		return &models.Account{}, err
	}

	data := map[string]interface{}{
		"password":   account.Password,
		"email":      account.Email,
		"updated_at": time.Now(),
	}

	return service.repository.Update(account, data)
}

func (service *AccountService) Delete(account *models.Account) (int64, error) {
	return service.repository.Delete(account)
}

func (service *AccountService) Auth(email, password string) (accountSerialized *models.AccountSerialized, err error) {
	account, err := service.repository.FindOneByEmail(email, true)

	if err == nil {
		err = service.VerifyPassword(account, password)
		if err != nil {
			return nil, fmt.Errorf("invalid password: %s", err)
		}
	}

	return service.Serialize(account), err
}

func (service *AccountService) GetRepository() *repositories.AccountRepository {
	return service.repository
}

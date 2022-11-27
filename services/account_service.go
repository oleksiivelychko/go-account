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

/*
VerifyPassword
modelAccount *models.Account contains hashed password
password string - contains non-hashed password
*/
func (service *AccountService) VerifyPassword(modelAccount *models.Account, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(modelAccount.Password), []byte(password))
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

func (service *AccountService) Create(modelAccount *models.Account) (*models.Account, error) {
	err := service.Validate(modelAccount)
	if err != nil {
		return &models.Account{}, err
	}

	err = service.repository.Create(modelAccount)
	if err != nil {
		return &models.Account{}, err
	}

	return modelAccount, nil
}

func (service *AccountService) Update(modelAccount *models.Account) (*models.Account, error) {
	err := service.Validate(modelAccount)
	if err != nil {
		return &models.Account{}, err
	}

	data := map[string]interface{}{
		"password":   modelAccount.Password,
		"email":      modelAccount.Email,
		"updated_at": time.Now(),
	}

	err = service.repository.Update(modelAccount, data)
	if err != nil {
		return &models.Account{}, err
	}

	return modelAccount, nil
}

func (service *AccountService) Delete(modelAccount *models.Account) (int64, error) {
	return service.repository.Delete(modelAccount)
}

func (service *AccountService) Auth(email, password string) (accountSerialized *models.AccountSerialized, err error) {
	modelAccount, err := service.repository.FindOneByEmail(email, true)

	if err == nil {
		err = service.VerifyPassword(modelAccount, password)
		if err != nil {
			return nil, fmt.Errorf("invalid password: %s", err)
		}
	}

	return service.Serialize(modelAccount), err
}

func (service *AccountService) GetRepository() *repositories.AccountRepository {
	return service.repository
}

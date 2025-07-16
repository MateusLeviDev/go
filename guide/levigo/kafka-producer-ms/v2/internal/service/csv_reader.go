package service

import (
	"github.com/MateusLeviDev/internal/models"
	"github.com/MateusLeviDev/pkg/utils"
)

func GetUSers() ([]models.User, error) {
	usersChan := make(chan []models.User, 1)
	errorsChan := make(chan error, 1)

	go func() {
		users, err := utils.ReadCSV()
		if err != nil {
			errorsChan <- err
			close(usersChan)
			return
		}
		usersChan <- users
	}()

	select {
	case err := <-errorsChan:
		return nil, err
	case users := <-usersChan:
		return users, nil
	}
}

package utils

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/MateusLeviDev/internal/models"
	"github.com/MateusLeviDev/pkg/constants"
	"github.com/MateusLeviDev/pkg/logger"
	"github.com/pquerna/ffjson/ffjson"
)

func ReadCSV() ([]models.User, error) {
	file, err := os.Open(constants.UsersFile)
	if err != nil {
		return nil, fmt.Errorf(constants.FileOpenErrMessage, err)
	}
	defer safelyClose(file)
	var users []models.User
	reader := csv.NewReader(file)
	reader.Comma = constants.Separator
	records, err := reader.ReadAll()

	if err != nil {
		return nil, fmt.Errorf(constants.RecordsReadErrMessage, err)
	}
	for _, record := range records[1:] {
		user, err := createUserFromRecord(record)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil

}

func safelyClose(file *os.File) {
	err := file.Close()
	if err != nil {
		logger.ErrorAsync(fmt.Sprintf("Warning - Error while closing the file: %v", err))
		return
	}
}

func createUserFromRecord(record []string) (models.User, error) {
	id, err := strconv.Atoi(record[0])
	if err != nil {
		return models.User{}, fmt.Errorf("error while converting the id: %v", err)
	}

	user := models.User{
		ID:       id,
		Username: record[1],
		Email:    record[2],
	}

	return user, nil
}

func DisplayUsersAsJSON(users []models.User) {
	jsonData, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		logger.ErrorAsync(fmt.Sprintf("Error during JSON conversion: %v", err))
	}
	fmt.Println(string(jsonData))
}

// WriteUsersToJSONFile writes a list of users to a JSON file using ffjson
func WriteUsersToJSONFile(users []models.User, filename string) {
	// Serialize using ffjson
	jsonData, err := ffjson.Marshal(users)
	if err != nil {
		logger.ErrorAsync(fmt.Sprintf("Error during JSON conversion: %v", err))
		return
	}

	// Indent the JSON using encoding/json
	var indentedData bytes.Buffer
	err = json.Indent(&indentedData, jsonData, "", "  ")
	if err != nil {
		logger.ErrorAsync(fmt.Sprintf("Error during JSON indentation: %v", err))
		return
	}

	file, err := os.Create(filename)
	if err != nil {
		logger.ErrorAsync(fmt.Sprintf("Error while creating the file: %v", err))
		return
	}
	defer safelyClose(file)

	_, err = file.Write(indentedData.Bytes())
	if err != nil {
		logger.ErrorAsync(fmt.Sprintf("Error while writing to the file: %v", err))
		return
	}

	logger.InfoAsync(fmt.Sprintf("Data successfully written to file %s\n", filename))
}

// BatchUsers splits a slice of users into batches of a given size
func BatchUsers(users []models.User, batchSize int) [][]models.User {
	var batches [][]models.User
	for batchSize < len(users) {
		users, batches = users[batchSize:], append(batches, users[0:batchSize:batchSize])
	}
	return append(batches, users)
}

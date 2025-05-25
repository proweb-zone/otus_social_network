package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"otus_social_network/app/internal/app/entity"
	"otus_social_network/app/internal/app/repository"
	"otus_social_network/app/internal/config"
	"otus_social_network/app/internal/db/postgres"
	"otus_social_network/app/internal/utils"
	"strconv"
	"strings"
	"time"
)

func main() {

	currentDir := utils.GetProjectPath()
	configPath := config.PathDefault(currentDir, nil)
	config := config.MustInit(configPath)

	filePath := currentDir + "/app/cmd/build_data/people.v2.csv"

	// read `csv file`
	reader, err := readCsvFile(filePath)
	if err != nil {
		fmt.Println("Error Read Csv File")
		return
	}

	sqlDb := postgres.Connect(config)

	userRepository := repository.InitPostgresRepository(sqlDb)

	numerator := 999931
	batch := 999931

	// остаток от деления batch
	remainder := math.Mod(float64(numerator), float64(batch))

	// count batch
	countBatch := numerator / batch

	// countBatch = 1 // переопределение countBatch для ограничения кол-ва записи в БД

	// create users slice
	users := make([]*entity.Users, 0)

	index := 1
	batchNumber := 1
	start := time.Now()
	for {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error read IO:", err)
			return
		}

		person, err := buildPersonData(record, index)
		if err != nil {
			fmt.Println("error build personal")
			return
		}

		users = append(users, person)

		if index == batch && batchNumber <= countBatch {
			fmt.Println("Writing batch", batchNumber, len(users))
			userRepository.BatchInsertUsers(users)
			users = make([]*entity.Users, 0)
			// time.Sleep(5 * time.Second)
			batchNumber++
			index = 1
		} else if index == int(remainder) && batchNumber > countBatch {
			fmt.Println("Записываем остатки штук", index)
			userRepository.BatchInsertUsers(users)
			users = make([]*entity.Users, 0)
			// time.Sleep(5 * time.Second)
			index = 1
		} else {
			index++
		}

	}

	elapsed := time.Since(start)
	fmt.Printf("Загрузка завершена за %s\n", elapsed)
}

func buildPersonData(record []string, index int) (*entity.Users, error) {
	var person entity.Users

	hashPass, err := utils.HashPassword("123645Xc")
	if err != nil {
		return nil, fmt.Errorf("Error: hash password", err)
	}

	firstAndLastName := strings.Split(record[0], " ")
	person.First_name = firstAndLastName[1]
	person.Last_name = firstAndLastName[0]
	person.City = record[2]
	person.Email = strconv.Itoa(index) + "@gmail.com"
	person.Password = hashPass

	parsedTime, err := time.Parse("2006-01-02", strings.TrimSpace(record[1]))
	if err != nil {
		return nil, fmt.Errorf("Error: Incorect date in field birth_date")
	}

	person.Birth_date = parsedTime

	return &person, nil
}

func readCsvFile(filePath string) (*csv.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	//defer file.Close()
	reader := csv.NewReader(file)
	return reader, nil
}

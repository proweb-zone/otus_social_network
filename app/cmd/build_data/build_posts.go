package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"otus_social_network/app/internal/db/redis"
	"otus_social_network/app/internal/utils"
	"strconv"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run your_program.go <user_id_1> <user_id_2> ...")
		return
	}

	userIDs := buildUserIds() // преобразовываем аргументы массив идентификаторов

	currentDir := utils.GetProjectPath()
	filePath := currentDir + "/app/cmd/build_data/posts.txt"

	buildPostsData(filePath, userIDs) // парсим файл с постами и наполняем redis моковыми данными

}

func buildUserIds() []int {
	userIDs := os.Args[1:]
	var intIDs []int
	for _, idStr := range userIDs {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			log.Fatalf("Invalid user ID: %s.  Must be an integer: %v", idStr, err)
		}
		intIDs = append(intIDs, id)
	}

	return intIDs
}

func buildPostsData(filePath string, userIDs []int) {

	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}
	defer file.Close()

	redis, errConnRedis := redis.InitRedisDb()
	if errConnRedis != nil {
		fmt.Println(errConnRedis)
		return
	}

	// errAddMsg := redis.AddMsg(request.User_id, request.Text)

	// if errAddMsg != nil {
	// 	return errAddMsg
	// }

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		sentences := strings.Split(line, ". ")
		for _, post := range sentences {

			//post := strings.Fields(sentence)
			if len(post) > 0 {
				randomUserId, err := randomUserIDs(userIDs)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}
				// TODO записываем посты в Redis
				redis.AddMsg(randomUserId, post)
			}
		}
		//parts := strings.Fields(line)
	}
}

func randomUserIDs(userIDs []int) (int, error) {

	if len(userIDs) == 0 {
		return 0, fmt.Errorf("cannot choose from an empty slice")
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(userIDs))

	return userIDs[randomIndex], nil
}

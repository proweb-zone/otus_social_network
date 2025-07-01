package service

import (
	"context"
	"fmt"
	"log"
	"otus_social_network/app/internal/app/dto"
	"otus_social_network/app/internal/app/entity"
	"otus_social_network/app/internal/app/repository"
	"otus_social_network/app/internal/db/rabbitmq"
	"otus_social_network/app/internal/db/redis"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type PostService struct {
	repo *repository.PostRepository
}

func InitPostService(repo *repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (p *PostService) CreatePost(ctx context.Context, request *dto.PostRequestDto) error {
	_, errorCreatePost := p.repo.CreatePost(
		ctx,
		&entity.Posts{
			User_id: request.User_id,
			Text:    request.Text,
		},
	)

	if errorCreatePost != nil {
		return errorCreatePost
	}

	// redis, errConnRedis := redis.InitRedisDb()
	// if errConnRedis != nil {
	// 	return errConnRedis
	// }

	// errAddMsg := redis.AddMsg(request.User_id, request.Text)

	// if errAddMsg != nil {
	// 	return errAddMsg
	// }

	// работаем через очереди RabbitMq
	rabbitmq, errConnRabbitmq := rabbitmq.InitRabbitMqDb()
	if errConnRabbitmq != nil {
		return errConnRabbitmq
	}

	errAddMsgRabbitmq := rabbitmq.AddMsg(request.User_id, request.Text)
	if errAddMsgRabbitmq != nil {
		return errAddMsgRabbitmq
	}

	return nil
}

func (p *PostService) UpdatePost(ctx context.Context, request *dto.PostRequestDto) error {
	var updatedAt time.Time
	if len(request.Updated_at) > 0 {
		parsedTime, err := time.Parse("2006-01-02", strings.TrimSpace(request.Updated_at))
		if err != nil {
			return fmt.Errorf("Error: Incorect date in field updated_at")
		}
		updatedAt = parsedTime
	}

	return p.repo.UpdatePost(
		ctx,
		&entity.Posts{
			ID:        request.ID,
			User_id:   request.User_id,
			Text:      request.Text,
			UpdatedAt: updatedAt,
		},
	)
}

func (p *PostService) DeletePost(userId int, postId int) error {
	return p.repo.DeletePost(userId, postId)
}

func (p *PostService) GetPost(userId int, postId int) (*entity.Posts, error) {
	return p.repo.GetPostById(userId, postId)
}

func (p *PostService) FeedPost(ids []int) ([]*entity.Posts, error) {
	p.repo.FeedPost(ids)

	redis, errConnRedis := redis.InitRedisDb()
	if errConnRedis != nil {
		return nil, errConnRedis
	}

	posts, errGetRedisMsg := redis.GetMessages(ids)
	if errGetRedisMsg != nil {
		return nil, errGetRedisMsg
	}

	return posts, nil
}

func (p *PostService) FeedPostWebSocket(ids []int) {
	connStr := "amqp://guest:guest@172.33.0.6:5672/"

	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Fatalf("Не удалось подключиться к RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Не удалось получить канал: %v", err)
	}

	var wg sync.WaitGroup
	for _, id := range ids {
		wg.Add(1)
		idStr := strconv.FormatInt(int64(id), 10)
		go consumeMessages(ch, "userid_"+idStr, &wg)
	}

	wg.Wait()
}

func consumeMessages(ch *amqp.Channel, queueName string, wg *sync.WaitGroup) {
	defer wg.Done()

	q, err := ch.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	forever := make(chan bool)

	go func() {
		for msg := range msgs {
			fmt.Printf("Получена новость: %s от пользователя %s\n", msg.Body, queueName)
			// time.Sleep(2 * time.Second)
			// msg.Ack(false)
		}
	}()

	fmt.Println("Подписан на очередь", queueName)
	fmt.Println("Ждем сообщений...")
	<-forever
}

package infrastructure

import (
	"chat-apps/internal/controller"
	"chat-apps/internal/rabbitmq"
	"chat-apps/internal/repository"
	"chat-apps/internal/service"
	"chat-apps/internal/util"
	"log"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	db, err := NewDB()
	if err != nil {
		panic(err)
	}

	_, ch, q, err := rabbitmq.InitRabbitMQ()
	if err != nil {
		log.Fatal(err)
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	messageRepo := repository.NewMessageRepository(db)
	messageService := service.NewMessageService(messageRepo)
	messageController := controller.NewMessageController(messageService)

	conversationRepo := repository.NewConversationRepository(db)
	conversationService := service.NewConversationService(conversationRepo)
	conversationController := controller.NewConversationController(conversationService)

	fileRepo := repository.NewFileRepository(db)
	fileService := service.NewFileService(fileRepo, userRepo)
	fileController := controller.NewFileController(fileService)

	jobRepo := repository.NewJobRepository(db)
	notificationRepo := repository.NewNotificationRepository(db)

	jobService := service.NewJobService(jobRepo, notificationRepo, userRepo, ch)
	notificationService := service.NewNotificationService(notificationRepo, userRepo)

	jobController := controller.NewJobController(jobService)
	notificationController := controller.NewNotificationController(notificationService)

	r := gin.Default()

	r.POST("/users", userController.CreateUser)
	r.GET("/users/:id", userController.GetUserByID)

	r.POST("/conversations/:conversationId/messages", messageController.CreateMessage)
	r.GET("/conversations/:conversationId/messages", messageController.GetMessagesByConversationID)

	r.POST("/conversations", conversationController.CreateConversation)
	r.GET("/conversations/:conversationId", conversationController.GetConversationByID)

	r.POST("/files/upload", fileController.UploadFile)
	r.GET("/files/:id", fileController.GetFileByID)

	r.POST("/notifications", notificationController.SendNotification)
	r.GET("/notifications/:userId", notificationController.GetNotificationsByUserID)

	r.POST("/notifications/broadcast", jobController.QueueBroadcastNotification)
	r.GET("/jobs/:id", jobController.GetJobStatus)

	notificationWorker := util.NewNotificationWorker(notificationRepo, userRepo, jobRepo)
	if err := rabbitmq.StartConsumer(ch, q.Name, notificationWorker); err != nil {
		log.Fatal(err)
	}

	return r
}

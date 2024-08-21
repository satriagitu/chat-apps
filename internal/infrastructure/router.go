package infrastructure

import (
	"chat-apps/internal/controller"
	"chat-apps/internal/repository"
	"chat-apps/internal/service"
	"chat-apps/internal/third_party"
	"chat-apps/internal/util"
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	db, err := NewDB()
	if err != nil {
		panic(err)
	}

	rabbit, err := third_party.NewRabbitMQ()
	if err != nil {
		log.Fatal(err)
	}

	rdb := third_party.InitRedis()

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

	jobService := service.NewJobService(jobRepo, notificationRepo, userRepo, rabbit.Channel)
	notificationService := service.NewNotificationService(notificationRepo, userRepo)

	jobController := controller.NewJobController(jobService)
	notificationController := controller.NewNotificationController(notificationService)

	artikeRepo := repository.NewArtikelRepository(db)
	artikelService := service.NewArtikelService(artikeRepo)
	artikelController := controller.NewArticleController(artikelService)

	externalPostController := controller.NewExternalPostController(os.Getenv("ExternalAPI"))

	ctx := context.Background()
	cacheRepo := repository.NewRedisRepository(rdb, ctx)
	cacheService := service.NewCacheService(cacheRepo)
	cacheController := controller.NewCacheController(cacheService)

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

	r.GET("/article-list", artikelController.GetArticleList)

	r.GET("/external-post", externalPostController.GetExternalPosts)

	r.POST("/cache-redis", cacheController.SetCache)
	r.GET("/cache-redis/:key", cacheController.GetCache)

	notificationWorker := util.NewNotificationWorker(notificationRepo, userRepo, jobRepo)
	if err := rabbit.StartConsumer(rabbit.Queue.Name, notificationWorker); err != nil {
		log.Fatal(err)
	}

	return r
}

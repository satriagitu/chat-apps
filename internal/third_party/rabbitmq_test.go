package third_party_test

import (
	"chat-apps/internal/domain"
	"chat-apps/internal/third_party"
	"chat-apps/internal/util"
	"encoding/json"
	"errors"
	"testing"

	"chat-apps/internal/repository/mocks"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStartConsumer(t *testing.T) {
	// Create a mock RabbitMQ channel
	mockChannel := new(mocks.RabbitMQChannel)

	// Create a mock NotificationWorker
	notificationRepo := new(mocks.NotificationRepository)
	userRepo := new(mocks.UserRepository)
	jobRepo := new(mocks.JobRepository)
	worker := util.NewNotificationWorker(notificationRepo, userRepo, jobRepo)

	// Prepare the mock for Consume
	msgs := make(chan amqp.Delivery)
	mockChannel.On("Consume", "test_queue", "", false, false, false, false, mock.Anything).Return((<-chan amqp.Delivery)(msgs), nil)

	// Mock the Close method (expecting a single call and returning nil to simulate no error)
	mockChannel.On("Close").Return(nil)

	// Create a RabbitMQ struct using the mock channel
	rabbitmq := &third_party.RabbitMQ{
		Connection: nil, // No need for a real connection in the test
		Channel:    mockChannel,
		Queue:      amqp.Queue{Name: "test_queue"},
	}

	// Start the consumer
	go func() {
		err := rabbitmq.StartConsumer("test_queue", worker)
		assert.NoError(t, err)
	}()

	// Prepare the mock for GetAllUsers
	users := []domain.User{
		{ID: 1, Username: "User1"},
		{ID: 2, Username: "User2"},
	}

	userRepo.On("GetAllUsers").Return(users, nil)

	// Prepare the mock for CreateNotification
	notificationRepo.On("CreateNotification", mock.AnythingOfType("domain.Notification")).Return(domain.Notification{}, nil)

	// Prepare the mock for UpdateJobStatus
	jobRepo.On("UpdateJobStatus", 1, "completed", mock.AnythingOfType("time.Time")).Return(nil)

	// Simulate message processing
	payload := map[string]interface{}{
		"job_id":  1,
		"message": "Test Message",
	}
	body, _ := json.Marshal(payload)
	msgs <- amqp.Delivery{Body: body}
	// Close the channel to stop the consumer
	close(msgs)

	// Close the RabbitMQ connection and channel
	err := rabbitmq.Close()
	assert.NoError(t, err) // Ensure no error during closing

	// Assert expectations
	mockChannel.AssertExpectations(t)
	userRepo.AssertExpectations(t)
	notificationRepo.AssertExpectations(t)
	jobRepo.AssertExpectations(t)
}

func TestStartConsumer_ReturnsErrorOnConsumeFailure(t *testing.T) {
	// Create a mock RabbitMQ channel
	mockChannel := new(mocks.RabbitMQChannel)

	// Create a mock NotificationWorker
	notificationRepo := new(mocks.NotificationRepository)
	userRepo := new(mocks.UserRepository)
	jobRepo := new(mocks.JobRepository)
	worker := util.NewNotificationWorker(notificationRepo, userRepo, jobRepo)

	// Simulate an error in the Consume method
	expectedErr := errors.New("failed to consume messages")
	mockChannel.On("Consume", "test_queue", "", false, false, false, false, mock.Anything).Return(nil, expectedErr)

	// Create a RabbitMQ struct using the mock channel
	rabbitmq := &third_party.RabbitMQ{
		Connection: nil, // No need for a real connection in the test
		Channel:    mockChannel,
		Queue:      amqp.Queue{Name: "test_queue"},
	}

	// Call StartConsumer and expect it to return the error
	err := rabbitmq.StartConsumer("test_queue", worker)

	// Assert that the error is returned as expected
	assert.EqualError(t, err, expectedErr.Error())

	// Assert expectations
	mockChannel.AssertExpectations(t)
}

func TestNotificationWorker_ProcessTask(t *testing.T) {
	// Initialize mocks
	notificationRepo := new(mocks.NotificationRepository)
	userRepo := new(mocks.UserRepository)
	jobRepo := new(mocks.JobRepository)

	// Create a NotificationWorker with the mocked repositories
	worker := util.NewNotificationWorker(notificationRepo, userRepo, jobRepo)

	t.Run("successful processing", func(t *testing.T) {
		// Prepare the delivery payload
		jobID := 1
		message := "Test Message"
		payload := map[string]interface{}{
			"job_id":  jobID,
			"message": message,
		}
		body, _ := json.Marshal(payload)

		delivery := amqp.Delivery{
			Body: body,
		}

		// Mock behavior for GetAllUsers
		users := []domain.User{
			{ID: 1, Username: "User1"},
			{ID: 2, Username: "User2"},
		}
		userRepo.On("GetAllUsers").Return(users, nil)

		// Mock behavior for CreateNotification
		notificationRepo.On("CreateNotification", mock.AnythingOfType("domain.Notification")).Return(domain.Notification{}, nil)

		// Mock behavior for UpdateJobStatus
		jobRepo.On("UpdateJobStatus", jobID, "completed", mock.AnythingOfType("time.Time")).Return(nil)

		// Call the method under test
		worker.ProcessTask(delivery)

		// Assert expectations
		userRepo.AssertExpectations(t)
		notificationRepo.AssertExpectations(t)
		jobRepo.AssertExpectations(t)
	})

	t.Run("error handling GetAllUsers failure", func(t *testing.T) {
		// Prepare the delivery payload
		jobID := 1
		message := "Test Message"
		payload := map[string]interface{}{
			"job_id":  jobID,
			"message": message,
		}
		body, _ := json.Marshal(payload)

		delivery := amqp.Delivery{
			Body: body,
		}

		// Mock behavior for GetAllUsers failure
		userRepo.On("GetAllUsers").Return(nil, errors.New("failed to get users"))

		// Call the method under test
		worker.ProcessTask(delivery)

		// Assert that GetAllUsers was called
		userRepo.AssertCalled(t, "GetAllUsers")

		// Assert that CreateNotification was not called
		notificationRepo.AssertNotCalled(t, "CreateNotification")

		// Assert that UpdateJobStatus was not called
		jobRepo.AssertNotCalled(t, "UpdateJobStatus", jobID, "completed")
	})

	t.Run("error handling CreateNotification failure", func(t *testing.T) {
		// Prepare the delivery payload
		jobID := 1
		message := "Test Message"
		payload := map[string]interface{}{
			"job_id":  jobID,
			"message": message,
		}
		body, _ := json.Marshal(payload)

		delivery := amqp.Delivery{
			Body: body,
		}

		// Mock behavior for GetAllUsers
		users := []domain.User{
			{ID: 1, Username: "User1"},
			{ID: 2, Username: "User2"},
		}
		userRepo.On("GetAllUsers").Return(users, nil)

		// Mock behavior for CreateNotification failure
		notificationRepo.On("CreateNotification", mock.AnythingOfType("domain.Notification")).Return(domain.Notification{}, errors.New("failed to create notification"))

		// Call the method under test
		worker.ProcessTask(delivery)

		// Assert that CreateNotification was called, even if it failed
		notificationRepo.AssertExpectations(t)
		// Assert that job status was still updated even though notifications failed
		jobRepo.AssertCalled(t, "UpdateJobStatus", jobID, "completed", mock.AnythingOfType("time.Time"))
	})

	t.Run("error handling UpdateJobStatus failure", func(t *testing.T) {
		// Prepare the delivery payload
		jobID := 1
		message := "Test Message"
		payload := map[string]interface{}{
			"job_id":  jobID,
			"message": message,
		}
		body, _ := json.Marshal(payload)

		delivery := amqp.Delivery{
			Body: body,
		}

		// Mock behavior for GetAllUsers
		users := []domain.User{
			{ID: 1, Username: "User1"},
			{ID: 2, Username: "User2"},
		}
		userRepo.On("GetAllUsers").Return(users, nil)

		// Mock behavior for CreateNotification
		notificationRepo.On("CreateNotification", mock.AnythingOfType("domain.Notification")).Return(domain.Notification{}, nil)

		// Mock behavior for UpdateJobStatus failure
		jobRepo.On("UpdateJobStatus", jobID, "completed", mock.AnythingOfType("time.Time")).Return(errors.New("failed to update job status"))

		// Call the method under test
		worker.ProcessTask(delivery)

		// Assert that UpdateJobStatus was called and failed
		jobRepo.AssertExpectations(t)
	})

	t.Run("invalid payload handling", func(t *testing.T) {
		// Prepare invalid payload (missing job_id)
		invalidPayload := map[string]interface{}{
			"message": "Test Message",
		}
		invalidBody, _ := json.Marshal(invalidPayload)

		invalidDelivery := amqp.Delivery{
			Body: invalidBody,
		}

		// Call the method under test with invalid payload
		worker.ProcessTask(invalidDelivery)

		// Assert that GetAllUsers was not called due to invalid payload
		userRepo.AssertNotCalled(t, "GetAllUsers", mock.AnythingOfType("domain.Notification"))
		notificationRepo.AssertNotCalled(t, "CreateNotification")
		jobRepo.AssertNotCalled(t, "UpdateJobStatus", mock.Anything)
	})

	t.Run("empty users list", func(t *testing.T) {
		// Prepare the delivery payload
		jobID := 1
		message := "Test Message"
		payload := map[string]interface{}{
			"job_id":  jobID,
			"message": message,
		}
		body, _ := json.Marshal(payload)

		delivery := amqp.Delivery{
			Body: body,
		}

		// Mock behavior for GetAllUsers returning empty list
		userRepo.On("GetAllUsers").Return([]domain.User{}, nil)

		// Call the method under test
		worker.ProcessTask(delivery)

		// Assert that GetAllUsers was called and that no notifications were created
		userRepo.AssertExpectations(t)
		notificationRepo.AssertNotCalled(t, "CreateNotification")
	})
}

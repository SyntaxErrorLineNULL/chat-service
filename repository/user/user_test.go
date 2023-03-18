package user_test

import (
	"context"
	"fmt"
	"github.com/SyntaxErrorLineNULL/chat-service/domain"
	"github.com/SyntaxErrorLineNULL/chat-service/repository/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"testing"
)

// mongodbContainer represents the mongodb container type used in the module
type mongodbContainer struct {
	testcontainers.Container
}

// startContainer creates an instance of the mongodb container type
func startContainer(ctx context.Context) (*mongodbContainer, error) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:latest",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor: wait.ForAll(
			wait.ForLog("Waiting for connections"),
			wait.ForListeningPort("27017/tcp"),
		),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	return &mongodbContainer{Container: container}, nil
}

func TestUserRepository_Create(t *testing.T) {
	ctx := context.Background()

	container, err := startContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	endpoint, err := container.Endpoint(ctx, "mongodb")
	if err != nil {
		t.Error(fmt.Errorf("failed to get endpoint: %w", err))
	}

	// Create a MongoDB client
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(endpoint))
	if err != nil {
		panic(err)
	}

	r := user.NewDefaultUserRepository(client, zap.NewNop())
	col := client.Database("chat-service").Collection("user")
	t.Run("user repository", func(t *testing.T) {
		t.Run("should create a new user", func(t *testing.T) {
			// Arrange
			id := uuid.New().String()
			newUser := &domain.User{
				ID:        id,
				Email:     "test@example.com",
				UserName:  "test_user",
				FirstName: "Test",
				LastName:  "User",
			}
			err := r.Create(ctx, newUser)
			assert.NoError(t, err)

			// check exist record in database
			var res domain.User
			errFind := col.FindOne(ctx, bson.M{"id": id}).Decode(&res)
			assert.NoError(t, errFind)
			assert.Equal(t, newUser, &res)
		})

		t.Run("should return error when user is nil", func(t *testing.T) {
			var nilUser *domain.User
			err := r.Create(ctx, nilUser)
			assert.EqualError(t, err, user.ErrInvalidArgument.Error())
		})

		expect := findTestUser(ctx, col)
		t.Run("should find an existing user by email", func(t *testing.T) {
			res, err := r.Find(ctx, &domain.User{Email: "cyb_orange190@gmail.com"})
			assert.NoError(t, err)
			assert.Equal(t, expect, res)
		})
		t.Run("should find an existing user by id", func(t *testing.T) {
			res, err := r.Find(ctx, &domain.User{ID: "7d05392e-1675-43df-a3e6-bd3a834dd729"})
			assert.NoError(t, err)
			assert.Equal(t, expect, res)
		})
	})

}

func findTestUser(ctx context.Context, col *mongo.Collection) *domain.User {
	expect := &domain.User{
		ID:        "7d05392e-1675-43df-a3e6-bd3a834dd729",
		Email:     "cyb_orange190@gmail.com",
		UserName:  "test_user",
		FirstName: "Test",
		LastName:  "User",
	}

	_, errInsert := col.InsertOne(ctx, expect)
	if errInsert != nil {
		panic(1)
	}

	return expect
}

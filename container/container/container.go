package container

import (
	"fmt"
	"github.com/SyntaxErrorLineNULL/chat-service/config"
	"go.uber.org/dig"
	"go.uber.org/zap"
)

func BuildContainer(config *config.Config) {
	container := dig.New()
	cfg := zap.NewProductionConfig()
	fmt.Print(cfg, container)
}

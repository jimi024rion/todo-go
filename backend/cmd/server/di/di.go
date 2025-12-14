//go:build wireinject
// +build wireinject

//go:generate kessoku di.go

package di

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// Build initializes the Gin engine with all its dependencies.
func Build() (*gin.Engine, func(), error) {
	wire.Build(
		gin.New,
	)

	// ここはkessoku generateによって置き換えられる
	return nil, nil, nil
}

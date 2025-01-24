package utils

import (
	"context"
	"fp_kata/pkg/log"
	"math/rand"
)

func GenerateNewId() int {
	return rand.Intn(100) + 1
}

func LogAction(ctx context.Context, component, function string) {
	log.GetLogger(ctx).Debug().Str(log.Comp, component).Str(log.Func, function).Send()
}

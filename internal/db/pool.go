package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Aiszhio/Task/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	Client *pgxpool.Pool
}

func NewPool(name string, ctx context.Context) (*Pool, error) {
	var pingDur = 5 * time.Second

	dsn, err := utils.GetEnv(name)
	if err != nil {
		log.Println(err)
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Println(err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Println(err)
	}

	ctxPing, cancel := context.WithTimeout(ctx, pingDur)
	defer cancel()

	if err = pool.Ping(ctxPing); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database %w: ", err)
	}

	return &Pool{
		Client: pool,
	}, nil
}

package main

import (
	"github.com/tarantool/go-tarantool/v2"

	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"
)

func InitTarantool() (*tarantool.Connection, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	dialer := tarantool.NetDialer{
		Address:  os.Getenv("TARANTOOL_ADDRESS"),
		User:     os.Getenv("TARANTOOL_USER_NAME"),
		Password: os.Getenv("TARANTOOL_USER_PASSWORD"),
	}

	opts := tarantool.Opts{
		Timeout:       5 * time.Second,
		Reconnect:     1 * time.Second,
		MaxReconnects: 5,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %v", err)
	}

	resp, err := conn.Do(tarantool.NewEvalRequest("return box.space.polls ~= nil")).Get()
	if err != nil {
		return nil, fmt.Errorf("space check failed: %v", err)
	}

	if len(resp) == 0 || !resp[0].(bool) {
		return nil, errors.New("space 'polls' does not exist")
	}

	log.Println("Successfully connected to Tarantool and verified 'polls' space exists")
	return conn, nil
}

package main

import (
	"context"
	"crud_service/app"
	"crud_service/app/repos"
	"crud_service/cmd/crud/config"
	"crud_service/server"
	"crud_service/storage/db"
	fs "crud_service/storage/file"
	"crud_service/storage/mem"
	"fmt"
	"os"
	"os/signal"
	"sync"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	var st repos.Storage
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %v runtime|database|file CONFIG_PATH\n", os.Args[0])
		return
	}

	cfg, err := config.ParseConfig(os.Args[2])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}

	switch os.Args[1] {
	case "runtime":
		st = mem.New()
	case "database":
		connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
			cfg.Database.User,
			cfg.Database.Pass,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Db,
		)
		conn, err := pgx.Connect(ctx, connStr)
		defer conn.Close(ctx)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}
		st = db.New(conn)

	case "file":
		st, err = fs.New(cfg.Folder.Path)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return
		}

	default:
		fmt.Printf("Usage: %v runtime|database|file\n", os.Args[0])
		return
	}
	conStr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := server.New(conStr)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	app := app.New(st)
	go app.Serve(ctx, wg, srv)

	<-ctx.Done()
	wg.Wait()
}

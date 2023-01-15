package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"github.com/hizani/crud_service/cmd/storage_service/config"
	"github.com/hizani/crud_service/storage_service"
	"github.com/hizani/crud_service/storage_service/db"
	fs "github.com/hizani/crud_service/storage_service/file"
	"github.com/hizani/crud_service/storage_service/mem"
	"github.com/hizani/crud_service/storage_service/model"
	"github.com/jackc/pgx/v5"
)

func main() {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	var st model.Storage
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

	ss := storage_service.New(st)
	ss.Start(conStr)

	<-ctx.Done()
	ss.Wg.Wait()
}

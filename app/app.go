package app

import (
	"context"
	"crud_service/app/repos"
	"sync"
)

type App struct {
	cs *repos.Customers
	ss *repos.Shops
}

func New(ust repos.Storage) *App {
	a := &App{
		cs: repos.NewCustomers(ust),
		ss: repos.NewShops(ust),
	}
	return a
}

type HTTPServer interface {
	Start(us *repos.Customers, cs *repos.Shops)
	Stop()
}

func (a *App) Serve(ctx context.Context, wg *sync.WaitGroup, hs HTTPServer) {
	defer wg.Done()
	hs.Start(a.cs, a.ss)
	<-ctx.Done()
	hs.Stop()
}

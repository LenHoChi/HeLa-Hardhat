package main

import (
	"context"
	"fmt"
	"hela-bank-sc/internal/blockchain"
	clientpkg "hela-bank-sc/internal/blockchain/client"
	blockchainevent "hela-bank-sc/internal/blockchain/event"
	txpkg "hela-bank-sc/internal/blockchain/transaction"
	"hela-bank-sc/internal/config"
	"hela-bank-sc/internal/database"
	"hela-bank-sc/internal/httpserver"
	repositorybank "hela-bank-sc/internal/repository/bank"
	"hela-bank-sc/internal/router"
	banksvc "hela-bank-sc/internal/service/bank"
	"log"
	"net/http"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	clientpkg.Init(cfg)
	txpkg.InitWallet()
	fmt.Println("✅ Setup done")

	dbConn, err := database.New(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.DB.Close()

	// Start listening events in background
	ctx := context.Background()
	go blockchainevent.ListenExplorer(ctx)

	txRepo := repositorybank.New(dbConn.DB)
	chainGateway := blockchain.New()

	bankSvc := banksvc.New(txRepo, chainGateway)
	rtr, err := initRouter(ctx, bankSvc)
	if err != nil {
		log.Fatal(err)
	}
	r := rtr.Routes()

	srv := httpserver.New(httpserver.Config{
		Addr:              ":" + cfg.AppPort,
		ReadHeaderTimeout: cfg.HTTPReadHeaderTimeout,
		ReadTimeout:       cfg.HTTPReadTimeout,
		WriteTimeout:      cfg.HTTPWriteTimeout,
		IdleTimeout:       cfg.HTTPIdleTimeout,
	}, r)

	if err := srv.Start(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func initRouter(ctx context.Context, bankSvc banksvc.Service) (router.Router, error) {
	return router.Router{
		Ctx:     ctx,
		BankSvc: bankSvc,
	}, nil
}

package main

import (
	"context"
	"fmt"
	"hela-bank-sc/internal/blockchain"
	clientpkg "hela-bank-sc/internal/blockchain/client"
	blockchainevent "hela-bank-sc/internal/blockchain/event"
	txpkg "hela-bank-sc/internal/blockchain/transaction"
	"hela-bank-sc/internal/database"
	"hela-bank-sc/internal/httpserver"
	"hela-bank-sc/internal/repository/transaction"
	"hela-bank-sc/internal/router"
	banksvc "hela-bank-sc/internal/service/bank"
	"log"
	"net/http"
)

func main() {
	clientpkg.Init()
	txpkg.InitWallet()
	fmt.Println("✅ Setup done")

	dbConn, err := database.New()
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.DB.Close()

	// Start listening events in background
	ctx := context.Background()
	go blockchainevent.ListenExplorer(ctx)

	txRepo := transaction.New(dbConn.DB)
	chainGateway := blockchain.New()

	bankSvc := banksvc.New(txRepo, chainGateway)
	rtr, err := initRouter(ctx, bankSvc)
	if err != nil {
		log.Fatal(err)
	}
	r := rtr.Routes()

	srv := httpserver.New(":8080", r)

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

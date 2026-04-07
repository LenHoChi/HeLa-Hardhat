package main

import (
	"context"
	"fmt"
	bank "hela-bank-sc/internal/blockchain"
	"hela-bank-sc/internal/router"
	banksvc "hela-bank-sc/internal/service/bank"
	"net/http"
)

func main() {
	bank.Init()
	bank.InitWallet()
	fmt.Println("✅ Setup done")

	// Start listening events in background
	ctx := context.Background()
	go bank.ListenEventsExplorer(ctx)

	bankSvc := banksvc.New()
	r := router.NewRouter(bankSvc)

	http.ListenAndServe(":8080", r)
}

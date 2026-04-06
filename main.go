package main

import (
	"fmt"
	"hela-bank-sc/internal/bank"
	"hela-bank-sc/internal/router"
	"net/http"
)

func main() {
	bank.Init()
	bank.InitWallet()
	fmt.Println("✅ Setup done")

	// flow:
	// main:    -> call DIRECT router func router
	// router:  -> call DIRECT handler func
	// handler: -> call DIRECT bank func
	// r := handler.NewRouter()

	// flow:
	// main:    -> create svc, call DIRECT router func (with svc param)
	// router:  -> create router struct with handler (from svc param), use this call handler method
	// handler: -> call bank func via h.bankSvc.Deposit (interface)
	bankSvc := bank.NewService()
	r := router.NewRouter(bankSvc)

	http.ListenAndServe(":8080", r)
}

// func main() {
// 	bank.Init()
// 	bank.InitWallet()
// 	fmt.Println("✅ Setup done")

// 	userAddr := common.HexToAddress("0x32A413fc36E202849B4eDffdB111802804fC7AEe")

// 	// Balance trước
// 	bank.PrintBalance(userAddr)

// 	// Deposit 0.1 ETH
// 	txHash, err := bank.Deposit(0.1)
// 	if err != nil {
// 		fmt.Println("Deposit error:", err)
// 		return
// 	}
// 	bank.WaitForTx(txHash)
// 	bank.PrintBalance(userAddr)

// 	// Withdraw 0.05 ETH
// 	txHash, err = bank.Withdraw(0.05)
// 	if err != nil {
// 		fmt.Println("Withdraw error:", err)
// 		return
// 	}
// 	bank.WaitForTx(txHash)
// 	bank.PrintBalance(userAddr)
// }

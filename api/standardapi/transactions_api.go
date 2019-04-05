// Package standardapi defines the summercash-wallet-server API.
package standardapi

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/SummerCash/summercash-wallet-server/common"
	"github.com/SummerCash/summercash-wallet-server/transactions"

	"github.com/valyala/fasthttp"

	summercashCommon "github.com/SummerCash/go-summercash/common"
)

/* BEGIN EXPORTED METHODS */

// SetupTransactionsRoutes sets up all the transactions api-related routes.
func (api *JSONHTTPAPI) SetupTransactionsRoutes() error {
	transactionsAPIRoot := "/api/transactions" // Get transactions API root path

	api.Router.POST(fmt.Sprintf("%s/NewTransaction", transactionsAPIRoot), api.NewTransaction) // Set NewTransaction post

	return nil // No error occurred, return nil
}

// NewTransaction handles a NewTransaction request.
func (api *JSONHTTPAPI) NewTransaction(ctx *fasthttp.RequestCtx) {
	var recipient summercashCommon.Address // Init recipient buffer
	var err error                          // Init error buffer

	if !strings.Contains(string(common.GetCtxValue(ctx, "recipient")), "0x") { // Check is sending to username
		recipientAccount, err := api.AccountsDatabase.QueryAccountByUsername(string(common.GetCtxValue(ctx, "recipient"))) // Query account

		if err != nil { // Check for errors
			logger.Errorf("errored while handling NewTransaction request with username %s: %s", string(common.GetCtxValue(ctx, "username")), err.Error()) // Log error

			panic(err) // Panic
		}

		recipient = recipientAccount.Address // Set address
	}

	recipient, err = summercashCommon.StringToAddress(string(common.GetCtxValue(ctx, "recipient"))) // Parse recipient

	if err != nil { // Check for errors
		logger.Errorf("errored while handling NewTransaction request with username %s: %s", string(common.GetCtxValue(ctx, "username")), err.Error()) // Log error

		panic(err) // Panic
	}

	amount, err := strconv.ParseFloat(string(common.GetCtxValue(ctx, "amount")), 64) // Parse amount

	if err != nil { // Check for errors
		logger.Errorf("errored while handling NewTransaction request with username %s: %s", string(common.GetCtxValue(ctx, "username")), err.Error()) // Log error

		panic(err) // Panic
	}

	transaction, err := transactions.NewTransaction(api.AccountsDatabase, string(common.GetCtxValue(ctx, "username")), string(common.GetCtxValue(ctx, "password")), &recipient, amount, common.GetCtxValue(ctx, "payload")) // Initialize transaction

	if err != nil { // Check for errors
		logger.Errorf("errored while handling NewTransaction request with username %s: %s", string(common.GetCtxValue(ctx, "username")), err.Error()) // Log error

		panic(err) // Panic
	}

	fmt.Fprintf(ctx, transaction.String()) // Write tx string value
}

/* END EXPORTED METHODS */

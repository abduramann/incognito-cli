package main

import (
	"github.com/incognitochain/go-incognito-sdk-v2/common"
	"log"
	"os"
	"sort"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "incognito-cli",
		Usage:   "A simple CLI application for the Incognito network",
		Version: "v0.0.2",
		Description: "A simple CLI application for the Incognito network. With this tool, you can run some basic functions" +
			" on your computer to interact with the Incognito network such as checking balances, transferring PRV or tokens," +
			" consolidating and converting your UTXOs, transferring tokens, manipulating with the pDEX, etc.",
		Authors: []*cli.Author{
			{
				Name: "Incognito Devs Team",
			},
		},
		Copyright: "This tool is developed and maintained by the Incognito Devs Team. It is free for anyone. However, any " +
			"commercial usages should be acknowledged by the Incognito Devs Team.",
	}

	// set app defaultFlags
	app.Flags = []cli.Flag{
		defaultFlags[networkFlag],
		defaultFlags[hostFlag],
		defaultFlags[clientVersionFlag],
		defaultFlags[debugFlag],
	}

	// all account-related commands
	accountCommands := []*cli.Command{
		{
			Name:     "keyinfo",
			Usage:    "Print all related-keys of a private key.",
			Category: accountCat,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     privateKeyFlag,
					Aliases:  aliases[privateKeyFlag],
					Usage:    "a base58-encoded private key",
					Required: true,
				},
			},
			Action: keyInfo,
		},
		{
			Name:     "balance",
			Usage:    "Check the balance of an account.",
			Category: accountCat,
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     privateKeyFlag,
					Aliases:  aliases[privateKeyFlag],
					Usage:    "a base58-encoded private key",
					Required: true,
				},
				&cli.StringFlag{
					Name:  tokenIDFlag,
					Usage: "ID of the token",
					Value: common.PRVIDStr,
				},
			},
			Action: checkBalance,
		},
		{
			Name:     "outcoin",
			Usage:    "Print the output coins of an account.",
			Category: accountCat,
			Flags: []cli.Flag{
				defaultFlags[addressFlag],
				defaultFlags[otaKeyFlag],
				defaultFlags[readonlyKeyFlag],
				defaultFlags[tokenIDFlag],
			},
			Action: getOutCoins,
		},
		{
			Name:     "utxo",
			Usage:    "Print the UTXOs of an account.",
			Category: accountCat,
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[tokenIDFlag],
			},
			Action: checkUTXOs,
		},
		{
			Name:    "consolidate",
			Aliases: []string{"csl"},
			Usage:   "Consolidate UTXOs of an account.",
			Description: "This function helps consolidate UTXOs of an account. It consolidates a version of UTXOs at a time, users need to specify which version they need to consolidate. " +
				"Please note that this process is time-consuming and requires a considerable amount of CPU.",
			Category: accountCat,
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[tokenIDFlag],
				defaultFlags[versionFlag],
				defaultFlags[numThreadsFlag],
				defaultFlags[enableLogFlag],
				defaultFlags[logFileFlag],
			},
			Action: consolidateUTXOs,
		},
		{
			Name:    "history",
			Aliases: []string{"hst"},
			Usage:   "Retrieve the history of an account.",
			Description: "This function helps retrieve the history of an account w.r.t a tokenID. " +
				"Please note that this process is time-consuming and requires a considerable amount of CPU.",
			Category: accountCat,
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				&cli.StringFlag{
					Name:  "tokenID",
					Usage: "ID of the token",
					Value: common.PRVIDStr,
				},
				defaultFlags[numThreadsFlag],
				defaultFlags[enableLogFlag],
				defaultFlags[logFileFlag],
				defaultFlags[csvFileFlag],
			},
			Action: getHistory,
		},
		{
			Name:        "generateaccount",
			Aliases:     []string{"genacc"},
			Usage:       "Generate a new Incognito account.",
			Description: "This function helps generate a new mnemonic phrase and its Incognito account.",
			Category:    accountCat,
			Flags: []cli.Flag{
				defaultFlags[numShardsFlags],
			},
			Action:      genKeySet,
		},
		{
			Name:    "submitkey",
			Aliases: []string{"sub"},
			Usage:   "Submit an ota key to the full-node.",
			Description: "This function submits an otaKey to the full-node to use the full-node's cache. If an access token " +
				"is provided, it will submit the ota key in an authorized manner. See " +
				"https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/accounts/submit_key.md " +
				"for more details.",
			Category: accountCat,
			Flags: []cli.Flag{
				defaultFlags[otaKeyFlag],
				defaultFlags[accessTokenFlag],
				defaultFlags[fromHeightFlag],
				defaultFlags[isResetFlag],
			},
			Action: submitKey,
		},
	}

	// all committee-related commands
	committeeCommands := []*cli.Command{
		{
			Name:     "checkrewards",
			Usage:    "Get all rewards of a payment address.",
			Category: committeeCat,
			Flags: []cli.Flag{
				defaultFlags[addressFlag],
			},
			Action: checkRewards,
		},
		{
			Name:     "withdrawreward",
			Usage:    "Withdraw the reward of a privateKey w.r.t to a tokenID.",
			Category: committeeCat,
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				&cli.StringFlag{
					Name:    addressFlag,
					Aliases: aliases[addressFlag],
					Usage:   "the payment address of a candidate (default: the payment address of the privateKey)",
				},
				defaultFlags[tokenIDFlag],
				defaultFlags[versionFlag],
			},
			Action: withdrawReward,
		},
	}

	// all tx-related commands
	txCommands := []*cli.Command{
		{
			Name:  "send",
			Usage: "Send an amount of PRV or token from one wallet to another wallet.",
			Description: "This function sends an amount of PRV or token from one wallet to another wallet. By default, " +
				"it used 100 nano PRVs to pay the transaction fee.",
			Category: transactionCat,
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[addressFlag],
				defaultFlags[amountFlag],
				defaultFlags[tokenIDFlag],
				defaultFlags[feeFlag],
				defaultFlags[versionFlag],
			},
			Action: send,
		},
		{
			Name:  "convert",
			Usage: "Convert UTXOs of an account w.r.t a tokenID.",
			Description: "This function helps convert UTXOs v1 of a user to UTXO v2 w.r.t a tokenID. " +
				"Please note that this process is time-consuming and requires a considerable amount of CPU.",
			Category: transactionCat,
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[tokenIDFlag],
				defaultFlags[numThreadsFlag],
				defaultFlags[enableLogFlag],
				defaultFlags[logFileFlag],
			},
			Action: convertUTXOs,
		},
		{
			Name:  "convertall",
			Usage: "Convert UTXOs of an account for all assets.",
			Description: "This function helps convert UTXOs v1 of a user to UTXO v2 for all assets. " +
				"It will automatically check for all UTXOs v1 of all tokens and convert them. " +
				"Please note that this process is time-consuming and requires a considerable amount of CPU.",
			Category: transactionCat,
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[numThreadsFlag],
				defaultFlags[logFileFlag],
			},
			Action: convertAll,
		},
		{
			Name:  "checkreceiver",
			Usage: "Check if an OTA key is a receiver of a transaction.",
			Description: "This function checks if an OTA key is a receiver of a transaction. If so, it will try to decrypt " +
				"the received outputs and return the receiving info.",
			Category: transactionCat,
			Flags: []cli.Flag{
				defaultFlags[txHashFlag],
				defaultFlags[otaKeyFlag],
				defaultFlags[readonlyKeyFlag],
			},
			Action: checkReceiver,
		},
	}

	// pDEX command
	pDEXCommands := []*cli.Command{
		{
			Name:  "pdecheckprice",
			Usage: "Check the price between two tokenIDs",
			Description: "This function checks the price of a pair of tokenIds. It must be supplied with the selling amount " +
				"since the pDEX uses the AMM algorithm.",
			Category: pDEXCat,
			Flags: []cli.Flag{
				defaultFlags[tokenIDToSellFlag],
				defaultFlags[tokenIDToBuyFlag],
				defaultFlags[sellingAmountFlag],
			},
			Action: pDEXCheckPrice,
		},
		{
			Name:        "pdetrade",
			Usage:       "Create a trade transaction",
			Description: "This function creates a trade transaction on the pDEX.",
			Category:    pDEXCat,
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[tokenIDToSellFlag],
				defaultFlags[tokenIDToBuyFlag],
				defaultFlags[sellingAmountFlag],
				defaultFlags[minAcceptableAmountFlag],
				defaultFlags[tradingFeeFlag],
			},
			Action: pDEXTrade,
		},
		{
			Name:        "pdecontribute",
			Usage:       "Create a pDEX contributing transaction",
			Description: "This function creates a pDEX contributing transaction. See more about this transaction: https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/pdex/contribute.md",
			Category:    pDEXCat,
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[pairIDFlag],
				defaultFlags[tokenIDFlag],
				defaultFlags[amountFlag],
				defaultFlags[versionFlag],
			},
			Action: pDEXContribute,
		},
		{
			Name:        "pdewithdraw",
			Usage:       "Create a pDEX withdrawal transaction",
			Description: "This function creates a transaction withdrawing an amount of `shared` from the pDEX. See more about this transaction: https://github.com/incognitochain/go-incognito-sdk-v2/blob/master/tutorials/docs/pdex/withdrawal.md",
			Category:    pDEXCat,
			Flags: []cli.Flag{
				defaultFlags[privateKeyFlag],
				defaultFlags[amountFlag],
				defaultFlags[tokenID1Flag],
				defaultFlags[tokenID2Flag],
				defaultFlags[versionFlag],
			},
			Action: pDEXWithdraw,
		},
		{
			Name:        "pdeshare",
			Usage:       "Retrieve the share amount of a pDEX pair",
			Description: "This function returns the share amount of a user within a pDEX pair.",
			Category:    pDEXCat,
			Flags: []cli.Flag{
				defaultFlags[addressFlag],
				defaultFlags[tokenID1Flag],
				defaultFlags[tokenID2Flag],
			},
			Action: pDEXGetShare,
		},
		{
			Name:  "pdetradestatus",
			Usage: "Get the status of a trade",
			Description: "This function returns the status of a trade (1: successful, 2: failed). If a `not found` error occurs, " +
				"it means that the trade has not been acknowledged by the beacon chain. Just wait and check again later.",
			Category: pDEXCat,
			Flags: []cli.Flag{
				defaultFlags[txHashFlag],
			},
			Action: pDEXTradeStatus,
		},
	}

	app.Commands = make([]*cli.Command, 0)
	app.Commands = append(app.Commands, accountCommands...)
	app.Commands = append(app.Commands, committeeCommands...)
	app.Commands = append(app.Commands, txCommands...)
	app.Commands = append(app.Commands, pDEXCommands...)

	for _, command := range app.Commands {
		buildUsageTextFromCommand(command)
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	//_ = generateDocsToFile(app, "commands.md") // un-comment this line to generate docs for the app's commands.

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

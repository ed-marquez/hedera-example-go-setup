package main

import (
	"fmt"
	"log"
	"os"

	hedera "github.com/hashgraph/hedera-sdk-go/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Retrieve operator ID from environment variables
	operatorIDStr := os.Getenv("OPERATOR_ID")
	if operatorIDStr == "" {
		log.Fatal("Environment variable OPERATOR_ID is not set")
	}
	operatorID, err := hedera.AccountIDFromString(operatorIDStr)
	if err != nil {
		log.Fatal("Invalid OPERATOR_ID:", err)
	}

	// Retrieve operator key from environment variables
	operatorKeyStr := os.Getenv("OPERATOR_KEY")
	if operatorKeyStr == "" {
		log.Fatal("Environment variable OPERATOR_KEY is not set")
	}
	operatorKey, err := hedera.PrivateKeyFromStringEd25519(operatorKeyStr)
	if err != nil {
		log.Fatal("Invalid OPERATOR_KEY:", err)
	}

	// Retrieve network from environment variables
	network := os.Getenv("NETWORK")
	if network == "" {
		log.Fatal("Environment variable NETWORK is not set")
	}

	// Create a Hedera client based on the network
	var client *hedera.Client
	switch network {
	case "testnet":
		client = hedera.ClientForTestnet()
	case "mainnet":
		client = hedera.ClientForMainnet()
	default:
		log.Fatalf("Unsupported network: %s. Use 'testnet' or 'mainnet'.", network)
	}
	client.SetOperator(operatorID, operatorKey)

	// Create Alice's account
	aliceKey, err := hedera.PrivateKeyGenerateEd25519()
	if err != nil {
		log.Fatal("Error generating Alice's private key:", err)
	}

	aliceAccount, err := hedera.NewAccountCreateTransaction().
		SetKey(aliceKey.PublicKey()).
		SetInitialBalance(hedera.NewHbar(1)).
		Execute(client)
	if err != nil {
		log.Fatal("Error creating Alice's account:", err)
	}

	aliceReceipt, err := aliceAccount.GetReceipt(client)
	if err != nil {
		log.Fatal("Error getting Alice's account receipt:", err)
	}

	if aliceReceipt.Status != hedera.StatusSuccess {
		log.Fatalf("Failed to create Alice's account: %v", aliceReceipt.Status)
	}

	aliceID := *aliceReceipt.AccountID
	fmt.Printf("Alice's account ID: %v\n", aliceID)

	// Create Bob's account
	bobKey, err := hedera.PrivateKeyGenerateEd25519()
	if err != nil {
		log.Fatal("Error generating Bob's private key:", err)
	}

	bobAccount, err := hedera.NewAccountCreateTransaction().
		SetKey(bobKey.PublicKey()).
		SetInitialBalance(hedera.NewHbar(1)).
		Execute(client)
	if err != nil {
		log.Fatal("Error creating Bob's account:", err)
	}

	bobReceipt, err := bobAccount.GetReceipt(client)
	if err != nil {
		log.Fatal("Error getting Bob's account receipt:", err)
	}

	if bobReceipt.Status != hedera.StatusSuccess {
		log.Fatalf("Failed to create Bob's account: %v", bobReceipt.Status)
	}

	bobID := *bobReceipt.AccountID
	fmt.Printf("Bob's account ID: %v\n", bobID)

	// Create a fungible HTS token
	tokenCreate, err := hedera.NewTokenCreateTransaction().
		SetTokenName("MyToken").
		SetTokenSymbol("MTK").
		SetTokenType(hedera.TokenTypeFungibleCommon).
		SetDecimals(0).
		SetInitialSupply(1000).
		SetTreasuryAccountID(operatorID).
		SetSupplyType(hedera.TokenSupplyTypeFinite).
		SetAdminKey(operatorKey.PublicKey()).
		SetFreezeKey(operatorKey.PublicKey()).
		SetWipeKey(operatorKey.PublicKey()).
		SetKycKey(operatorKey.PublicKey()).
		SetSupplyKey(operatorKey.PublicKey()).
		SetFreezeDefault(false).
		Execute(client)
	if err != nil {
		log.Fatal("Error creating token:", err)
	}

	tokenCreateReceipt, err := tokenCreate.GetReceipt(client)
	if err != nil {
		log.Fatal("Error getting token creation receipt:", err)
	}

	if tokenCreateReceipt.Status != hedera.StatusSuccess {
		log.Fatalf("Failed to create token: %v", tokenCreateReceipt.Status)
	}

	tokenID := *tokenCreateReceipt.TokenID
	fmt.Printf("Token ID: %v\n", tokenID)

	// Associate Alice's account with the token
	associateAliceTx, err := hedera.NewTokenAssociateTransaction().
		SetAccountID(aliceID).
		SetTokenIDs(tokenID).
		FreezeWith(client)
	if err != nil {
		log.Fatal("Error freezing association transaction:", err)
	}

	// Sign the transaction with Alice's key
	associateAliceTxSigned := associateAliceTx.Sign(aliceKey)

	// Submit the transaction
	associateAliceTxSubmit, err := associateAliceTxSigned.Execute(client)
	if err != nil {
		log.Fatal("Error submitting association transaction:", err)
	}

	associateAliceTxReceipt, err := associateAliceTxSubmit.GetReceipt(client)
	if err != nil {
		log.Fatal("Error getting association transaction receipt:", err)
	}

	if associateAliceTxReceipt.Status != hedera.StatusSuccess {
		log.Fatalf("Failed to associate Alice's account with the token: %v", associateAliceTxReceipt.Status)
	}

	fmt.Printf("Associate Alice transaction status: %v\n", associateAliceTxReceipt.Status)

	// Transfer the token to Alice's account
	transferTx, err := hedera.NewTransferTransaction().
		AddTokenTransfer(tokenID, operatorID, -100).
		AddTokenTransfer(tokenID, aliceID, 100).
		Execute(client)
	if err != nil {
		log.Fatal("Error transferring tokens:", err)
	}

	transferReceipt, err := transferTx.GetReceipt(client)
	if err != nil {
		log.Fatal("Error getting transfer transaction receipt:", err)
	}

	if transferReceipt.Status != hedera.StatusSuccess {
		log.Fatalf("Failed to transfer tokens: %v", transferReceipt.Status)
	}

	fmt.Printf("Transfer transaction status: %v\n", transferReceipt.Status)
}

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

        // Retrieve operator ID and key from environment variables
        operatorID, err := hedera.AccountIDFromString(os.Getenv("OPERATOR_ID"))
        if err != nil {
                log.Fatal("Invalid OPERATOR_ID:", err)
        }

        operatorKey, err := hedera.PrivateKeyFromStringEd25519(os.Getenv("OPERATOR_KEY"))
        if err != nil {
                log.Fatal("Invalid OPERATOR_KEY:", err)
        }

		network, err = os.Getenv("NETWORK")

        // Create a Hedera client
        client := hedera.ClientForNetwork(network)
        client.SetOperator(operatorID, operatorKey)

        // Create Alice's account
        aliceKey, err := hedera.PrivateKeyGenerateEd25519()
        if err != nil {
                log.Fatal("Error generating private key:", err)
        }

        aliceAccount, err := hedera.NewAccountCreateTransaction().
                SetKey(aliceKey.PublicKey()).
                SetInitialBalance(hedera.NewHbar(10)).
                Execute(client)
        if err != nil {
                log.Fatal("Error creating account:", err)
        }

        aliceReceipt, err := aliceAccount.GetReceipt(client)
        if err != nil {
                log.Fatal("Error getting receipt:", err)
        }

        aliceID := *aliceReceipt.AccountID
        fmt.Printf("Alice's account ID: %v\n", aliceID)

        // Create Bob's account
        bobKey, err := hedera.PrivateKeyGenerateEd25519()
        if err != nil {
                log.Fatal("Error generating private key:", err)
        }

        bobAccount, err := hedera.NewAccountCreateTransaction().
                SetKey(bobKey.PublicKey()).
                SetInitialBalance(hedera.NewHbar(10)).
                Execute(client)
        if err != nil {
                log.Fatal("Error creating account:", err)
        }

        bobReceipt, err := bobAccount.GetReceipt(client)
        if err != nil {
                log.Fatal("Error getting receipt:", err)
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
                log.Fatal("Error getting receipt:", err)
        }

        tokenID := *tokenCreateReceipt.TokenID
        fmt.Printf("Token ID: %v\n", tokenID)

        // Associate Alice's account with the token
        associateAliceTx, err := hedera.NewTokenAssociateTransaction().
                SetAccountID(aliceID).
                SetTokenIDs(tokenID).
                FreezeWith(client)
        if err != nil {
                log.Fatal("Error freezing transaction:", err)
        }

        associateAliceTxSign, err := associateAliceTx.Sign(aliceKey)
        if err != nil {
                log.Fatal("Error signing transaction:", err)
        }

        associateAliceTxSubmit, err := associateAliceTxSign.Execute(client)
        if err != nil {
                log.Fatal("Error submitting transaction:", err)
        }

        associateAliceTxReceipt, err := associateAliceTxSubmit.GetReceipt(client)
        if err != nil {
                log.Fatal("Error getting receipt:", err)
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
                log.Fatal("Error getting receipt:", err)
        }

        fmt.Printf("Transfer transaction status: %v\n", transferReceipt.Status)
}
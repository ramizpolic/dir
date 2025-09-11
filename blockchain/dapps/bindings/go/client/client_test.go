// Copyright AGNTCY Contributors (https://github.com/agntcy)
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"math/big"
	"testing"

	"github.com/agntcy/dir/blockchain/dapps/bindings/go/agentstore"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
)

func setupTestClient(t *testing.T) *Client {
	// These values should be set to your testnet/devnet or local chain
	chainURL := "http://localhost:8545"
	agentStoreAddress := "0x5fbdb2315678afecb367f032d93f642f64180aa3"                // Replace with deployed contract address
	privKeyHex := "8166f546bab6da521a8369cab06c5d2b9e46670292d85c875ee9ec20e84ffb61" // no 0x prefix
	privKey, err := crypto.HexToECDSA(privKeyHex)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}
	txOpts, err := bind.NewKeyedTransactorWithChainID(privKey, big.NewInt(123321))
	if err != nil {
		t.Fatalf("failed to create tx opts: %v", err)
	}
	client, err := New(
		WithChainURL(chainURL),
		WithAgentStore(agentStoreAddress),
		func(o *Options) error { o.txOpts = txOpts; return nil },
	)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	return client
}

func TestAgentStore_Get(t *testing.T) {
	client := setupTestClient(t)
	callOpts := &bind.CallOpts{Context: context.Background()}
	_, err := client.AgentStore().Get(callOpts, "test-agent-id")
	if err != nil {
		t.Logf("expected error or empty result: %v", err)
	}
}

func TestAgentStore_Total(t *testing.T) {
	client := setupTestClient(t)
	callOpts := &bind.CallOpts{Context: context.Background()}
	total, err := client.AgentStore().Total(callOpts)
	if err != nil {
		t.Fatalf("failed to call Total: %v", err)
	}
	t.Logf("Total agents: %v", total)
}

func TestAgentStore_Add(t *testing.T) {
	client := setupTestClient(t)
	agent := agentstore.Agent{
		Id:        "test-agent-id",
		Signature: "test-signature",
		Owner:     client.TxOpts().From,
	}
	_, err := client.AgentStore().Add(client.TxOpts(), agent)
	if err != nil {
		t.Logf("expected error if contract is not deployed or tx fails: %v", err)
	}
}

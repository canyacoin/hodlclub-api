package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func getCanBalanceAtWalletAddress(address common.Address) (*big.Float, int64, error) {
	var err error
	b, err := nodeConnection.CodeAt(context.Background(), address, nil)
	if err != nil {
		return nil, 0, err
	}

	if len(b) > 0 {
		return nil, 0, err
	}

	header, err := nodeConnection.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	blockHeight, _ := strconv.Atoi(header.Number.String())

	balance, err := contractInstance.BalanceOf(&bind.CallOpts{}, address)
	if err != nil {
		return nil, 0, err
	}

	return canAmountToFloat(balance), int64(blockHeight), err
}

func canAmountToFloat(amount *big.Int) *big.Float {
	famount := new(big.Float)
	famount.SetString(amount.String())
	return new(big.Float).Quo(famount, big.NewFloat(math.Pow10(6)))
}

/*
 * TODO: We have a number of these in our projects, move to a package
 */
func mustGetenv(k string) string {
	v := os.Getenv(k)
	if v == "" {
		panic(fmt.Sprintf("%s environment variable not set.", k))
	}
	return v
}

func bf2f64(a *big.Float) float64 {
	x := fmt.Sprintf("%v", a)
	if strings.Contains(x, ".") {
		y := strings.Split(x, ".")
		x = y[0]
	}
	i, _ := strconv.ParseFloat(x, 64)
	return i
}

func updateLastBlockProcessed(i int64) error {
	lb := LastBlock{BlockHeight: i, Key: "block"}
	// err := db.Table("last_blocks").Where(LastBlock{Key: "block"}).Assign(&lb).FirstOrCreate(&lb).Error
	err := db.Table("last_blocks").Where(LastBlock{Key: "block"}).FirstOrCreate(&lb).Error
	if err != nil {
		logger.Errorf("unable to set last_blocks default value: %s", err.Error())
	}

	err = db.Exec("UPDATE last_blocks SET block_height = ?", i).Error
	if err != nil {
		logger.Errorf("unable to save last_block record: %s", err.Error())
	}

	logger.Debugf("updated last block height to: %d", i)
	return err
}

func upsertClubMember(br BalanceRecord) error {
	logger.Debugf("upserting balance record for hash: %s with balance: %d at block height: %d", br.Hash, br.Balance, br.BlockHeight)

	// err := db.Table("balances").Where(BalanceRecord{Hash: br.Hash}).Assign(BalanceRecord{Balance: br.Balance, BlockHeight: br.BlockHeight}).FirstOrCreate(&br).Error

	err := db.Exec("UPDATE balances SET balance = ?, block_height = ? WHERE hash = ?", br.Balance, br.BlockHeight, br.Hash).Error

	return err
}

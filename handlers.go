package main

import (
	"fmt"
	"math"
	"math/big"
	"net/http"

	humanize "github.com/dustin/go-humanize"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
)

func rootHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index", gin.H{
		"title": "CanYa HODL Club",
	})
}

func hodlersHandler(c *gin.Context) {
	var err error
	var balances []BalanceRecord
	q := c.Request.URL.Query()

	if q["tier"] != nil && q["tier"][0] == "OG" {
		logger.Debugf("get the og")
		err = db.Where("balance >= 10000").Order("balance desc").Find(&balances).Error
	} else {
		logger.Debugf("get t2")
		err = db.Where("balance >= 2500 AND balance < 10000").Order("balance desc").Find(&balances).Error
	}

	if err != nil {
		logger.Warningf("unable to query balances, error was: %s", err.Error())
	}

	c.JSON(http.StatusOK, balances)

}
func processBlockchainHandler(c *gin.Context) {
	var err error
	verifiedHashes := make(map[string]BalanceRecord)
	lowestBlockHeight := int64(0)
	blockHeight := int64(0)

	lb := LastBlock{}
	err = db.Table("last_blocks").First(&lb).Error
	if err != nil {
		logger.Warningf("unable to count last_blocks, error was: %s using starting block: %s", err.Error(), startingBlock)
	}

	nodeConnection, err = ethclient.Dial(ethNodeURL)

	if err != nil {
		m := fmt.Sprintf("error connecting to infura at: %s error was: %s", ethNodeURL, err.Error())
		logger.Errorf(m)
		c.JSON(http.StatusInternalServerError, gin.H{"message": m})
		return
	}

	contractInstance, err = NewCanYaCoin(common.HexToAddress(canyCoinContract), nodeConnection)
	if err != nil {
		m := fmt.Sprintf("error loading contract: %s error was: %s", canyCoinContract, err.Error())
		logger.Errorf(m)
		c.JSON(http.StatusInternalServerError, gin.H{"message": m})
		return
	}

	dc, err := contractInstance.Decimals(&bind.CallOpts{})
	if err != nil {
		fmt.Printf("ERR: %+v\n", err)
	}
	fmt.Printf("DUMP: %+v\n", dc)

	var blacklisted []BlacklistRecord
	err = db.Find(&blacklisted).Error
	if err != nil {
		logger.Warningf("unable to reteive blacklisted records, error was: %s", err.Error())
	}
	blacklistedHashes := make(map[string]bool)
	for _, blh := range blacklisted {
		blacklistedHashes[blh.Hash] = true
	}

	logger.Infof("analysing blockchain from block: %d", uint64(lb.BlockHeight))

	itr, err := contractInstance.FilterTransfer(&bind.FilterOpts{Start: uint64(lb.BlockHeight)}, []common.Address{}, []common.Address{})
	// itr, err := contractInstance.FilterTransfer(&bind.FilterOpts{Start: uint64(6705124)}, []common.Address{}, []common.Address{})
	if err != nil {
		logger.Fatalf("error getting FilterTransfer, error was:  %s", err.Error())
	}

	for itr.Next() {
		canValue := new(big.Float)
		var err error
		addresses := [2]common.Address{itr.Event.To, itr.Event.From}

		for _, a := range addresses {

			if blacklistedHashes[a.String()] == true {
				logger.Debugf("skipping address: %s blacklisted", a.String())
				continue
			}

			if verifiedHashes[a.String()].Hash != "" {
				logger.Debugf("skipping address: %s already analysed", a.String())
				continue
			}

			canValue, blockHeight, err = getCanBalanceAtWalletAddress(a)
			if lowestBlockHeight == 0 {
				lowestBlockHeight = blockHeight
			}

			if err != nil {
				logger.Warningf("unable to get balance at address: %s error was: %s", a.String(), err.Error())
			}
			logger.Debugf("address: %s has balance of: %v", a.String(), canValue)

			canf64 := bf2f64(canValue)
			verifiedHashes[a.String()] = BalanceRecord{
				Hash:        a.String(),
				Balance:     int64(math.Floor(canf64)),
				BlockHeight: blockHeight,
			}
		}
	}
	itr.Close()

	minCanValue := int64(2500)
	insertedCount := 0
	for key := range verifiedHashes {
		if verifiedHashes[key].Balance >= minCanValue {
			err := upsertClubMember(verifiedHashes[key])
			if err != nil {
				logger.Errorf("unable to find transaction record by hash, error: %s", err.Error())
			}
			insertedCount++
		}
	}
	logger.Infof("inserted: %d hodl club member hashes", insertedCount)

	err = updateLastBlockProcessed(lowestBlockHeight)
	if err != nil {
		logger.Errorf("error updating last processed block height: %s", err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"status":            "OK",
		"memberCount":       insertedCount,
		"lowestBlockHeight": lowestBlockHeight,
	})
}

func statusHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":         "OK",
		"serviceID":      serviceID,
		"serviceStarted": humanize.Time(startedAt),
	})
}

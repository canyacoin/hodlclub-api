package main

// BalanceRecord describes the data structure for an Ethereum wallet address
// Represented as a database table, and json response from etherescan
type BalanceRecord struct {
	Hash        string `gorm:"column:hash;type:varchar(255);unique_index" json:"hash,omitempty"`
	Balance     int64  `gorm:"column:balance" json:"balance,omitempty"`
	BlockHeight int64  `gorm:"column:block_height" json:"blockNumber,omitempty"`
	IsOG        int64  `gorm:"column:is_og" json:"isOg,omitempty"`
}

// TableName sets the default table name
func (BalanceRecord) TableName() string {
	return "balances"
}

// LastBlock describes the data structure for last procedd block
// Represented as a database table, and json response from etherescan
// It is intended that this table will hold a single record
type LastBlock struct {
	Key         string `gorm:"column:key;type:varchar(255);unique_index" json:"key,omitempty"`
	BlockHeight int64  `gorm:"column:block_height" json:"blockNumber,omitempty"`
}

// BlacklistRecord describes the data structure for Ethereum wallet addresses
// that are blaclisted from the hodl club
type BlacklistRecord struct {
	Hash string `gorm:"column:hash;type:varchar(255);unique_index" json:"hash,omitempty"`
}

// IcoOgRecord describes the data structure for Ethereum wallet addresses
// that are promised hodl club OG membership during the ICO
type IcoOgRecord struct {
	Hash string `gorm:"column:hash;type:varchar(255)" json:"hash,omitempty"`
}

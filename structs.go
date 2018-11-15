package main

// BalanceRecord describes the data structure for an Ethereum wallet address
// Represented as a database table, and json response from etherescan
type BalanceRecord struct {
	Hash        string `gorm:"column:hash;type:varchar(255);unique_index" json:"hash,omitempty"`
	Balance     int64  `gorm:"column:balance" json:"balance,omitempty"`
	BlockHeight int64  `gorm:"column:block_height" json:"blockNumber,omitempty"`
}

// TableName sets the default table name
func (BalanceRecord) TableName() string {
	return "balances"
}

// LastBlock describes the data structure for last procedd block
// Represented as a database table, and json response from etherescan
// It is intended that this table will hold a single record
type LastBlock struct {
	BlockHeight int64 `gorm:"column:block_height" json:"blockNumber,omitempty"`
}

// BlacklistRecord describes the data structure for Ethereum wallet addresses
// that are blaclisted from the hodl club
type BlacklistRecord struct {
	Hash string `gorm:"column:hash;type:varchar(255);unique_index" json:"hash,omitempty"`
}

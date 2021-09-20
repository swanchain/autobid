package models

import (
	"go-swan/database"
	"go-swan/logs"
)

type AutoBidHistory struct {
	Id       int `json:"id"`
	TaskId   int `json:"task_id"`
	MinerId  int `json:"miner_id"`
	CreateAt int `json:"create_at"`
}

func AddAutoBidHistory(autoBidHistory *AutoBidHistory) error {
	err := database.GetDB().Create(autoBidHistory).Error

	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return err
}

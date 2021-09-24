package models

import (
	"go-swan/database"
	"go-swan/logs"
)

type Miner struct {
	Id                     int      `json:"id"`
	Score                  int      `json:"score"`
	Price                  *float64 `json:"price"`
	VerifiedPrice          *float64 `json:"verified_price"`
	MinerFid               string   `json:"miner_fid"`
	UpdateTimeStr          string   `json:"update_time_str"`
	Status                 string   `json:"status"`
	MinPieceSize           *string  `json:"min_piece_size"`
	MaxPieceSize           *string  `json:"max_piece_size"`
	Location               string   `json:"location"`
	SwanMinerId            int      `json:"swan_miner_id"`
	OfflineDealAvailable   int      `json:"offline_deal_available"`
	DailySealingCapability int      `json:"daily_sealing_capabilty"`
	AdjustedPower          string   `json:"adjusted_power"`
	ReachableCount         int      `json:"reachable_count"`
	UnreachableCount       int      `json:"unreachable_count"`
	DailyReward            float64  `json:"daily_reward"`
	SectorLiveCount        int      `json:"sector_live_count"`
	SectorFaultyCount      int      `json:"sector_faulty_count"`
	SectorActiveCount      int      `json:"sector_active_count"`
	BidMode                int      `json:"bid_mode"`
	StartEpoch             *int     `json:"start_epoch"`
	AddressBalance         float64  `json:"address_balance"`
	AutoBidTaskPerDay      int      `json:"auto_bid_task_per_day"`
	AutoBidTaskCnt         int      `json:"auto_bid_task_cnt"`
	LastAutoBidAt          int64    `json:"last_auto_bid_at"` //millisecond of last auto-bid task for this miner
	ExpectedSealingTime    *int     `json:"expected_sealing_time"`
	StartEpochAbs          *int
	MinPieceSizeByte       *float64
	MaxPieceSizeByte       *float64
	ScoreSumBefore         int
	IsScanned              bool
}

func GetMiners(pageNum int, pageSize int, status string) ([]*Miner, error) {
	var miners []*Miner
	err := database.GetDB().Where("Status=?", status).Offset(pageNum).Limit(pageSize).Find(&miners).Error

	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return miners, nil
}

func GetAutoBidMinersOrderByScore(status string) ([]*Miner, error) {
	var miners []*Miner
	var notNulCols []string
	notNulCols = append(notNulCols, "price")
	notNulCols = append(notNulCols, "verified_price")
	notNulCols = append(notNulCols, "expected_sealing_time")
	notNulCols = append(notNulCols, "min_piece_size")
	notNulCols = append(notNulCols, "max_piece_size")
	notNulCols = append(notNulCols, "start_epoch")
	filter := "bid_mode=1 and status=? and offline_deal_available=1"
	for i := range notNulCols {
		filter = filter + " and " + notNulCols[i] + " is not null"
	}
	err := database.GetDB().Where(filter, status).Find(&miners).Error

	if err != nil {
		logs.GetLogger().Error(err)
		return nil, err
	}

	return miners, nil
}

func MinerUpdateLastAutoBidInfo(minerId, autoBidTaskCnt int, lastAutoBidAt int64) error {
	lastAutoBidInfo := make(map[string]interface{})
	lastAutoBidInfo["auto_bid_task_cnt"] = autoBidTaskCnt
	lastAutoBidInfo["last_auto_bid_at"] = lastAutoBidAt
	err := database.GetDB().Model(Miner{}).Where("id=?", minerId).Update(lastAutoBidInfo).Error

	if err != nil {
		logs.GetLogger().Error(err)
		return err
	}

	return err
}

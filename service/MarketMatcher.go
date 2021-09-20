package service

import (
	"go-swan/common/constants"
	"go-swan/common/utils"
	"go-swan/logs"
	"go-swan/models"
	"math/rand"
	"time"
)

var miners []*models.Miner = nil

func FindMiners() {
	go FindMiner4AllTasks()
}

func FindMiner4AllTasks() {
	for {
		taskCntFilled := 0
		for {
			taskCntFilled1Time := FindMiner4Tasks()

			taskCntFilled = taskCntFilled + taskCntFilled1Time

			if taskCntFilled1Time > 0 {
				logs.GetLogger().Info(taskCntFilled1Time, " auto bidding tasks are matched.")
				continue
			}

			if taskCntFilled1Time == 0 {
				logs.GetLogger().Info("All auto bidding tasks are matched.")
				break
			}
		}

		logs.GetLogger().Info(taskCntFilled, " auto bidding tasks are matched in total this time.")

		for _, miner := range miners {
			models.MinerUpdateLastAutoBidInfo(miner.Id, miner.AutoBidTaskCnt, miner.LastAutoBidAt)
		}

		miners = nil

		time.Sleep(2 * time.Minute)
	}
}

func FindMiner4Tasks() int {
	tasks, err := models.GetAutoBidTasks(0, 100, constants.TASK_STATUS_CREATED)
	if err != nil {
		logs.GetLogger().Info(err)
		return 0
	}

	if tasks == nil || len(tasks) == 0 {
		return 0
	}

	if miners == nil || len(miners) == 0 {
		GetMiners()

		if miners == nil || len(miners) == 0 {
			logs.GetLogger().Info("No miners found")
			return 0
		}
	}

	for _, task := range tasks {
		miner := FindMiner4OneTask(task)

		if miner == nil {
			logs.GetLogger().Error("Did not find miner for task: ", task.TaskName, " id=", task.ID)
			continue
		}
		err = models.TaskAssignMiner(task.ID, miner.Id)
		if err != nil {
			logs.GetLogger().Error(err)
		}

		currentNanoSec := time.Now().UnixNano()
		if utils.IsSameDay(miner.LastAutoBidAt, currentNanoSec) {
			miner.AutoBidTaskCnt = miner.AutoBidTaskCnt + 1
		} else {
			miner.AutoBidTaskCnt = 1
		}
		miner.LastAutoBidAt = currentNanoSec
	}

	return len(tasks)
}

func FindMiner4OneTask(task *models.Task) *models.Miner {
	if task.MaxPrice == nil {
		err := models.TaskUpdateStatus(task.ID, constants.TASK_STATUS_ACTION_REQUIRED)
		if err != nil {
			logs.GetLogger().Error(err)
		}
		return nil
	}

	var matchedMiners []*models.Miner

	offlineDeals, err := models.GetOfflineDealByTaskId(task.ID)
	if err != nil {
		logs.GetLogger().Error(err)
		err = models.TaskUpdateStatus(task.ID, constants.TASK_STATUS_AUTO_BID_FAILED)
		if err != nil {
			logs.GetLogger().Error(err)
		}
		return nil
	}

	if len(offlineDeals) == 0 {
		err = models.TaskUpdateStatus(task.ID, constants.TASK_STATUS_AUTO_BID_FAILED)
		if err != nil {
			logs.GetLogger().Error(err)
		}
		return nil
	}

	for _, miner := range miners {
		if miner.AutoBidTaskCnt >= miner.AutoBidTaskPerDay && utils.IsSameDay(miner.LastAutoBidAt, time.Now().UnixNano()) {
			continue
		}

		minerMatch := true
		for _, offlineDeal := range offlineDeals {
			var fileSize float64
			fileSize, err = utils.GetFloat64FromStr(offlineDeal.FileSize)
			if err != nil {
				logs.GetLogger().Error(err)
				err = models.TaskUpdateStatus(task.ID, constants.TASK_STATUS_ACTION_REQUIRED)
				if err != nil {
					logs.GetLogger().Error(err)
				}
				return nil
			}

			if fileSize < 0 {
				logs.GetLogger().Error("DealCid:", offlineDeal.DealCid, " no valid file size")
				err = models.TaskUpdateStatus(task.ID, constants.TASK_STATUS_ACTION_REQUIRED)
				if err != nil {
					logs.GetLogger().Error(err)
				}
				return nil
			}

			if fileSize < miner.ByteSizeMin || fileSize > miner.ByteSizeMax {
				minerMatch = false
				break
			}

			if offlineDeal.StartEpoch != nil && *offlineDeal.StartEpoch <= miner.StartEpoch {
				minerMatch = false
				break
			}
		}

		if !minerMatch {
			continue
		}

		if !(miner.Price < *task.MaxPrice || miner.VerifiedPrice < *task.MaxPrice) {
			continue
		}

		matchedMiners = append(matchedMiners, miner)
	}

	if len(matchedMiners) == 0 {
		err = models.TaskUpdateStatus(task.ID, constants.TASK_STATUS_AUTO_BID_FAILED)
		if err != nil {
			logs.GetLogger().Error(err)
		}
		return nil
	}

	ratio := rand.ExpFloat64()
	for _, matchedMiner := range matchedMiners {
		if ratio < matchedMiner.ScorePercent {
			logs.GetLogger().Info(matchedMiner.Id, " ScorePercent=", matchedMiner.ScorePercent, " MinerFid=", matchedMiner.MinerFid, " ratio=", ratio, " is selected")
			return matchedMiner
		}
		//logs.GetLogger().Info(ratio)
	}

	return matchedMiners[len(matchedMiners)-1]
}

func GetMiners() int {
	var err error
	miners, err = models.GetAllMinersOrderByScore(constants.MINER_STATUS_ACTIVE)
	if err != nil {
		logs.GetLogger().Error(err)
		return 0
	}

	totalScore := 0

	for _, miner := range miners {
		totalScore = totalScore + miner.Score

		miner.ByteSizeMin = utils.GetByteSizeFromStr(miner.MinPieceSize)
		miner.ByteSizeMax = utils.GetByteSizeFromStr(miner.MaxPieceSize)
	}

	var minScorePercent float64 = 1
	var maxScorePercent float64 = 0

	for _, miner := range miners {
		miner.ScorePercent = float64(miner.Score) / float64(totalScore)
		if miner.ScorePercent < minScorePercent {
			minScorePercent = miner.ScorePercent
		}

		if miner.ScorePercent > maxScorePercent {
			maxScorePercent = miner.ScorePercent
		}
	}

	multiple := 1.0 / maxScorePercent

	for _, miner := range miners {
		//oldScorePercent := miner.ScorePercent
		miner.ScorePercent = miner.ScorePercent * multiple
		//fmt.Println(oldScorePercent,miner.Score, miner.ScorePercent)
	}

	return len(miners)
}

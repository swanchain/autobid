package service

import (
	"errors"
	"go-swan/common/constants"
	"go-swan/common/utils"
	"go-swan/config"
	"go-swan/logs"
	"go-swan/models"
	"time"
)

var miners []*models.Miner = nil

type MinerStat struct {
	miners []*models.Miner
}

var minScore = 100
var maxScore = 0
var minerStat [101]*MinerStat

const DEFAULT_TASK_PAGE_SIZE = 100

func FindMiners() {
	for {
		FindMiner4AllTasks()
		time.Sleep(config.GetConfig().AutoBidIntervalSec * time.Second)
	}
}

func FindMiner4AllTasks() {
	miners = nil
	minScore = 100
	maxScore = 0
	for i := 0; i < len(minerStat); i++ {
		minerStat[i] = nil
	}
	GetMiners()

	taskCntDealt := 0

	for taskCntDealt1Time := DEFAULT_TASK_PAGE_SIZE; taskCntDealt1Time >= DEFAULT_TASK_PAGE_SIZE; {
		taskCntDealt1Time = FindMiner4Tasks()
		taskCntDealt = taskCntDealt + taskCntDealt1Time
		logs.GetLogger().Info(taskCntDealt1Time, " auto bidding tasks are dealt.")
	}

	/*	for _, miner := range miners {
		models.MinerUpdateLastAutoBidInfo(miner.Id, miner.AutoBidTaskCnt, miner.LastAutoBidAt)
	}*/

	logs.GetLogger().Info(taskCntDealt, " auto bidding tasks are dealt in total this time.")
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
		logs.GetLogger().Info("No miners")
		return 0
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
			continue
		}

		currentNanoSec := time.Now().UnixNano()
		if utils.IsSameDay(miner.LastAutoBidAt, currentNanoSec) {
			miner.AutoBidTaskCnt = miner.AutoBidTaskCnt + 1
		} else {
			miner.AutoBidTaskCnt = 1
		}
		miner.LastAutoBidAt = currentNanoSec
		models.MinerUpdateLastAutoBidInfo(miner.Id, miner.AutoBidTaskCnt, miner.LastAutoBidAt)
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

	for _, offlineDeal := range offlineDeals {
		offlineDeal.FileSizeNum, err = utils.GetFloat64FromStr(offlineDeal.FileSize)
		if err != nil {
			logs.GetLogger().Error(err)
			err = models.TaskUpdateStatus(task.ID, constants.TASK_STATUS_ACTION_REQUIRED)
			if err != nil {
				logs.GetLogger().Error(err)
			}
			err = errors.New("file size invalid, deal cid: " + offlineDeal.DealCid)
			return nil
		}

		if offlineDeal.FileSizeNum < 0 {
			logs.GetLogger().Error("DealCid:", offlineDeal.DealCid, " no valid file size")
			err = models.TaskUpdateStatus(task.ID, constants.TASK_STATUS_ACTION_REQUIRED)
			if err != nil {
				logs.GetLogger().Error(err)
			}
			err = errors.New("file size invalid, deal cid: " + offlineDeal.DealCid)
			return nil
		}
	}

	miner := SelectMiner(task, offlineDeals)

	return miner
}

func SelectMiner(task *models.Task, offlineDeals []*models.OfflineDeals) *models.Miner {
	score := utils.GetRandInRange(minScore, maxScore)
	scoreUp := score
	scoreDown := score - 1
	for scoreUp <= 100 || scoreDown >= 0 {
		if scoreUp <= 100 {
			miner := GetMatchedMiner(scoreUp, task, offlineDeals)
			if miner != nil {
				return miner
			}
		}

		if scoreDown >= 0 {
			miner := GetMatchedMiner(scoreDown, task, offlineDeals)
			if miner != nil {
				return miner
			}
		}

		scoreUp++
		scoreDown--
	}

	return nil
}

func GetMatchedMiner(score int, task *models.Task, offlineDeals []*models.OfflineDeals) *models.Miner {
	if minerStat[score] == nil {
		return nil
	}

	if minerStat[score].miners == nil {
		return nil
	}

	cntScoreMiner := len(minerStat[score].miners)

	if cntScoreMiner == 0 {
		return nil
	}

	var isChecked []bool
	for i := 0; i < cntScoreMiner; i++ {
		isChecked = append(isChecked, false)
	}

	for i := 0; i < cntScoreMiner; i++ {
		minerIndex := utils.GetRandInRange(0, cntScoreMiner-1)
		if isChecked[minerIndex] {
			continue
		}

		isChecked[minerIndex] = true
		miner := minerStat[score].miners[minerIndex]
		isMinerMatch := IsMinerMatch(miner, task, offlineDeals)
		if isMinerMatch {
			return miner
		}
	}

	for i := 0; i < cntScoreMiner; i++ {
		if isChecked[i] {
			continue
		}

		miner := minerStat[score].miners[i]
		isMinerMatch := IsMinerMatch(miner, task, offlineDeals)
		if isMinerMatch {
			return miner
		}
	}

	return nil
}

func IsMinerMatch(miner *models.Miner, task *models.Task, offlineDeals []*models.OfflineDeals) bool {
	if miner.AutoBidTaskCnt >= miner.AutoBidTaskPerDay && utils.IsSameDay(miner.LastAutoBidAt, time.Now().UnixNano()) {
		return false
	}

	for _, offlineDeal := range offlineDeals {
		if offlineDeal.FileSizeNum < miner.MinPieceSizeByte || offlineDeal.FileSizeNum > miner.MaxPieceSizeByte {
			return false
		}

		if offlineDeal.StartEpoch != nil && *offlineDeal.StartEpoch <= miner.StartEpoch {
			return false
		}
	}

	if !(miner.Price < *task.MaxPrice || miner.VerifiedPrice < *task.MaxPrice) {
		return false
	}

	return true
}

func GetMiners() []*models.Miner {
	var err error
	miners, err = models.GetAllMinersOrderByScore(constants.MINER_STATUS_ACTIVE)
	if err != nil {
		logs.GetLogger().Error(err)
		return nil
	}

	if miners == nil || len(miners) == 0 {
		return nil
	}

	for _, miner := range miners {
		if minerStat[miner.Score] == nil {
			minerStat[miner.Score] = &MinerStat{}
		}

		miner.MinPieceSizeByte = utils.GetByteSizeFromStr(miner.MinPieceSize)
		miner.MaxPieceSizeByte = utils.GetByteSizeFromStr(miner.MaxPieceSize)

		minerStat[miner.Score].miners = append(minerStat[miner.Score].miners, miner)
		if miner.Score > maxScore {
			maxScore = miner.Score
		}

		if miner.Score < minScore {
			minScore = miner.Score
		}
	}

	return miners
}

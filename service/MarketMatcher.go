package service

import (
	"autobid/common/constants"
	"autobid/common/utils"
	"autobid/config"
	"autobid/logs"
	"autobid/models"
	"fmt"
	"strconv"
	"time"

	"github.com/filswan/go-swan-lib/client/lotus"
)

var miners []*models.Miner = nil

const DEFAULT_TASK_PAGE_SIZE = 10000

func FindMiners() {
	for {
		FindMiner4AllTasks()
		time.Sleep(config.GetConfig().AutoBidIntervalSec * time.Second)
	}
}

func FindMiner4AllTasks() {
	miners = nil
	GetMiners()

	taskCntDealt := 0

	for taskCntDealt1Time := DEFAULT_TASK_PAGE_SIZE; taskCntDealt1Time >= DEFAULT_TASK_PAGE_SIZE; {
		taskCntDealt1Time = FindMiner4Tasks()
		taskCntDealt = taskCntDealt + taskCntDealt1Time
		logs.GetLogger().Info(taskCntDealt1Time, " auto bidding tasks are dealt.")
	}

	logs.GetLogger().Info(taskCntDealt, " auto bidding tasks are dealt in total this time.")
}

func FindMiner4Tasks() int {
	tasks, err := models.GetAutoBidTasks(0, DEFAULT_TASK_PAGE_SIZE, constants.TASK_STATUS_CREATED)
	if err != nil {
		logs.GetLogger().Info(err)
		return 0
	}

	if len(tasks) == 0 {
		return 0
	}

	if len(miners) == 0 {
		logs.GetLogger().Info("No miners")
		return 0
	}

	for _, task := range tasks {
		miner := FindMiner4OneTask(task)

		if miner == nil {
			logs.GetLogger().Error("Did not find miner for task: ", task.TaskName, " id=", task.Id)
			continue
		}

		logs.GetLogger().Info("Allocated miner:", miner.Id, " to task:", task.Id)

		currentNanoSec := time.Now().UnixNano()
		autoBidTaskCnt := 1
		if utils.IsSameDay(miner.LastAutoBidAt, currentNanoSec) {
			autoBidTaskCnt = miner.AutoBidTaskCnt + 1
		}

		err = models.TaskAssignMiner(task.Id, miner.Id, autoBidTaskCnt, currentNanoSec)
		if err != nil {
			logs.GetLogger().Error(err)
			continue
		}

		miner.AutoBidTaskCnt = autoBidTaskCnt
		miner.LastAutoBidAt = currentNanoSec
	}

	return len(tasks)
}

func FindMiner4OneTask(task *models.Task) *models.Miner {
	if task.MaxPrice == nil {
		logs.GetLogger().Error("Task:", task.Id, " no max price")
		err := models.TaskUpdateStatus(task.Id, constants.TASK_STATUS_ACTION_REQUIRED)
		if err != nil {
			logs.GetLogger().Error(err)
		}
		return nil
	}

	if task.FastRetrieval == nil {
		logs.GetLogger().Error("Task:", task.Id, " no fast retrieval")
		err := models.TaskUpdateStatus(task.Id, constants.TASK_STATUS_ACTION_REQUIRED)
		if err != nil {
			logs.GetLogger().Error(err)
		}
		return nil
	}

	if task.Type == nil || (*task.Type != constants.TASK_TYPE_VERIFIED && *task.Type != constants.TASK_TYPE_REGULAR) {
		logs.GetLogger().Error("Task:", task.Id, " invalid task type")
		err := models.TaskUpdateStatus(task.Id, constants.TASK_STATUS_ACTION_REQUIRED)
		if err != nil {
			logs.GetLogger().Error(err)
		}
		return nil
	}

	offlineDeals, err := models.GetOfflineDealByTaskId(task.Id)
	if err != nil {
		logs.GetLogger().Error(err)
		err = models.TaskUpdateStatus(task.Id, constants.TASK_STATUS_ACTION_REQUIRED)
		if err != nil {
			logs.GetLogger().Error(err)
		}
		return nil
	}

	if len(offlineDeals) == 0 {
		err = models.TaskUpdateStatus(task.Id, constants.TASK_STATUS_ACTION_REQUIRED)
		if err != nil {
			logs.GetLogger().Error(err)
		}
		return nil
	}

	for _, offlineDeal := range offlineDeals {
		if offlineDeal.FileSize == nil || *offlineDeal.FileSize == "" {
			err := fmt.Errorf("file size is empty, offline deal id:%d", offlineDeal.Id)
			logs.GetLogger().Error(err)
			err = models.TaskUpdateStatus(task.Id, constants.TASK_STATUS_ACTION_REQUIRED)
			if err != nil {
				logs.GetLogger().Error(err)
			}
			return nil
		}

		offlineDeal.FileSizeNumer, err = strconv.ParseInt(*offlineDeal.FileSize, 10, 64)
		if err != nil {
			logs.GetLogger().Error(err)
			err = models.TaskUpdateStatus(task.Id, constants.TASK_STATUS_ACTION_REQUIRED)
			if err != nil {
				logs.GetLogger().Error(err)
			}
			return nil
		}

		if offlineDeal.FileSizeNumer <= 0 {
			logs.GetLogger().Error("Deal:", offlineDeal.Id, " no valid file size")
			err = models.TaskUpdateStatus(task.Id, constants.TASK_STATUS_ACTION_REQUIRED)
			if err != nil {
				logs.GetLogger().Error(err)
			}
			return nil
		}

		if utils.IsStrEmpty(offlineDeal.FileSourceUrl) {
			logs.GetLogger().Error("Deal:", offlineDeal.Id, " no file source url")
			err = models.TaskUpdateStatus(task.Id, constants.TASK_STATUS_ACTION_REQUIRED)
			if err != nil {
				logs.GetLogger().Error(err)
			}
			return nil
		}

		if offlineDeal.StartEpoch == nil {
			logs.GetLogger().Error("Deal:", offlineDeal.Id, " no StartEpoch")
			err = models.TaskUpdateStatus(task.Id, constants.TASK_STATUS_ACTION_REQUIRED)
			if err != nil {
				logs.GetLogger().Error(err)
			}
			return nil
		}

		if utils.IsStrEmpty(offlineDeal.PayloadCid) {
			logs.GetLogger().Error("Deal:", offlineDeal.Id, " no payload cid")
			err = models.TaskUpdateStatus(task.Id, constants.TASK_STATUS_ACTION_REQUIRED)
			if err != nil {
				logs.GetLogger().Error(err)
			}
			return nil
		}

		if utils.IsStrEmpty(offlineDeal.PieceCid) {
			logs.GetLogger().Error("Deal:", offlineDeal.Id, " no piece cid")
			err = models.TaskUpdateStatus(task.Id, constants.TASK_STATUS_ACTION_REQUIRED)
			if err != nil {
				logs.GetLogger().Error(err)
			}
			return nil
		}
	}

	miner := SelectMiner(task, offlineDeals)

	return miner
}

func SelectMiner(task *models.Task, offlineDeals []*models.OfflineDeals) *models.Miner {
	for _, miner := range miners {
		miner.IsScanned = false
	}

	for i := 0; i < len(miners)/2; i++ {
		miner := FindMiner(task, offlineDeals)
		if miner != nil {
			return miner
		}
	}

	for _, miner := range miners {
		if miner.IsScanned {
			continue
		}

		miner.IsScanned = true
		isMinerMatch := IsMinerMatch(miner, task, offlineDeals)
		if isMinerMatch {
			return miner
		}
	}

	return nil
}

func FindMiner(task *models.Task, offlineDeals []*models.OfflineDeals) *models.Miner {
	maxRand := miners[len(miners)-1].ScoreSumBefore
	scoreSumRand := utils.GetRandInRange(1, maxRand)

	start := 0
	end := len(miners) - 1
	for start <= end {
		mid := (start + end) / 2
		miner := miners[mid]
		scoreSumEnd := miner.ScoreSumBefore
		scoreSumStart := scoreSumEnd - miner.Score
		if scoreSumRand > scoreSumEnd {
			start = mid + 1
			continue
		}

		if scoreSumRand <= scoreSumStart {
			end = mid - 1
			continue
		}

		if miner.IsScanned {
			return nil
		}

		miner.IsScanned = true
		isMinerMatch := IsMinerMatch(miner, task, offlineDeals)
		if isMinerMatch {
			return miner
		}
		return nil
	}

	return nil
}

func IsMinerMatch(miner *models.Miner, task *models.Task, offlineDeals []*models.OfflineDeals) bool {
	if miner.Price == nil || miner.VerifiedPrice == nil || miner.ExpectedSealingTime == nil ||
		miner.MinPieceSize == nil || miner.MaxPieceSize == nil || miner.StartEpoch == nil {
		return false
	}

	if miner.AutoBidTaskCnt >= miner.AutoBidTaskPerDay && utils.IsSameDay(miner.LastAutoBidAt, time.Now().UnixNano()) {
		logs.GetLogger().Info("Auto-bid tasks are enough. miner:", miner.Id, " task:", task.Id)
		return false
	}

	if *task.Type == constants.TASK_TYPE_REGULAR && (*task.MaxPrice).Cmp(*miner.Price) < 0 {
		logs.GetLogger().Info("*task.MaxPrice:", *task.MaxPrice, " < *miner.Price", *miner.Price, " miner:", miner.Id, " task:", task.Id)
		return false
	}

	if *task.Type == constants.TASK_TYPE_VERIFIED && (*task.MaxPrice).Cmp(*miner.VerifiedPrice) < 0 {
		logs.GetLogger().Info("*task.MaxPrice:", *task.MaxPrice, " < *miner.VerifiedPrice", *miner.VerifiedPrice, " miner:", miner.Id, " task:", task.Id)
		return false
	}

	if task.ExpireDays != nil {
		taskExpireEpoch := utils.GetEpochFromDay(*task.ExpireDays)
		if taskExpireEpoch < *miner.ExpectedSealingTime {
			logs.GetLogger().Info("taskExpireEpoch:", taskExpireEpoch, " < *miner.ExpectedSealingTime:", *miner.ExpectedSealingTime, " miner:", miner.Id, " task:", task.Id)
			return false
		}
	}

	for _, offlineDeal := range offlineDeals {
		if miner.MinPieceSizeByte == nil || offlineDeal.FileSizeNumer < *miner.MinPieceSizeByte {
			if miner.MinPieceSizeByte == nil {
				logs.GetLogger().Info("miner.MinPieceSizeByte == nil. miner:", miner.Id, " task:", task.Id, " offlineDeal:", offlineDeal.Id)
			} else {
				logs.GetLogger().Info("offlineDeal.FileSizeNumer:", offlineDeal.FileSizeNumer, " < *miner.MinPieceSizeByte:", *miner.MinPieceSizeByte, " miner:", miner.Id, " task:", task.Id, " offlineDeal:", offlineDeal.Id)
			}

			return false
		}

		if miner.MaxPieceSizeByte == nil || offlineDeal.FileSizeNumer > *miner.MaxPieceSizeByte {
			if miner.MaxPieceSizeByte == nil {
				logs.GetLogger().Info("miner.MaxPieceSizeByte == nil. miner:", miner.Id, " task:", task.Id, " offlineDeal:", offlineDeal.Id)
			} else {
				logs.GetLogger().Info("offlineDeal.FileSizeNumer:", offlineDeal.FileSizeNumer, " > *miner.MaxPieceSizeByte:", *miner.MaxPieceSizeByte, " miner:", miner.Id, " task:", task.Id, " offlineDeal:", offlineDeal.Id)
			}

			return false
		}

		if *offlineDeal.StartEpoch < *miner.StartEpochAbs {
			logs.GetLogger().Info("*offlineDeal.StartEpoch:", *offlineDeal.StartEpoch, " < *miner.StartEpochAbs:", *miner.StartEpochAbs, " miner:", miner.Id, " task:", task.Id, " offlineDeal:", offlineDeal.Id)
			return false
		}
	}

	return true
}

func GetMiners() []*models.Miner {
	var err error
	miners, err = models.GetAutoBidMiners()
	if err != nil {
		logs.GetLogger().Error(err)
		return nil
	}

	if len(miners) == 0 {
		return nil
	}

	lotusClient, err := lotus.LotusGetClient(config.GetConfig().Lotus.ClientApiUrl, "")
	if err != nil {
		logs.GetLogger().Error(err)
		return nil
	}

	for i, miner := range miners {
		logs.GetLogger().Info("getting price for miner:", miner.MinerFid)
		minerConfig, err := lotusClient.LotusClientQueryAsk(miner.MinerFid)
		if err != nil {
			logs.GetLogger().Error(err)
			return nil
		}

		miner.MinPieceSizeByte = &minerConfig.MinPieceSize
		miner.MaxPieceSizeByte = &minerConfig.MaxPieceSize
		miner.Price = &minerConfig.Price
		miner.VerifiedPrice = &minerConfig.VerifiedPrice

		if miner.StartEpoch != nil {
			epochAbs := utils.GetCurrentEpoch() + *miner.StartEpoch + constants.EPOCH_PER_HOUR
			miner.StartEpochAbs = &epochAbs
		}

		if i == 0 {
			miner.ScoreSumBefore = miner.Score
		} else {
			miner.ScoreSumBefore = miner.Score + miners[i-1].ScoreSumBefore
		}

		miner.IsScanned = false
	}

	return miners
}

package test

import (
	"fmt"
	"go-swan/common/constants"
	"go-swan/common/utils"
	"go-swan/logs"
	"go-swan/models"
	"go-swan/service"
	"math/rand"
)

func Test() {
	//TestTask_GetTasks()
	//TestTask_GetAutoBidTasks()
	TestMiner_GetAllMiners()
	//service.FindMiner4AllTasks()
	//models.MinerUpdateLastAutoBidInfo(726, 19, time.Now().UnixNano())
	//models.TaskAssignMiner(1389, 14)
}

func TestTask_GetTasks() {
	tasks, err := models.GetTasks(0, 10, constants.TASK_STATUS_CREATED)
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	for _, task := range tasks {
		logs.GetLogger().Info(utils.ToJson(task))
	}
}

func TestTask_GetAutoBidTasks() {
	tasks, err := models.GetAutoBidTasks(0, 10, constants.TASK_STATUS_CREATED)
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	for _, task := range tasks {
		logs.GetLogger().Info(utils.ToJson(task))
	}
}

func TestTask_GetTaskById() {
	task, err := models.GetTaskById(1)
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	logs.GetLogger().Info(utils.ToJson(task))
}

func TestMiner_GetMiners() {
	miners, err := models.GetMiners(0, 10, constants.MINER_STATUS_ACTIVE)
	if err != nil {
		logs.GetLogger().Error(err)
		return
	}

	for _, miner := range miners {
		logs.GetLogger().Info(utils.ToJson(miner))
	}
}

func TestMiner_GetAllMiners() {
	miners := service.GetMiners()

	for j := 0; j < 100; j++ {
		ratio := rand.ExpFloat64()
		for i, miner := range miners {
			//score := miner.ScorePercent * multiple
			//logs.GetLogger().Info("Score:", ratio, randNum)
			if ratio < miner.ScorePercent {
				fmt.Println(i, " ScorePercent=", miner.ScorePercent, " MinerFid=", miner.MinerFid, " ratio=", ratio, " is selected")
				break
			}
		}
	}
}

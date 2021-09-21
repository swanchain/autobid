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
	//TestMiner_GetAllMiners()
	//service.FindMiner4AllTasks()
	//models.MinerUpdateLastAutoBidInfo(726, 19, time.Now().UnixNano())
	//models.TaskAssignMiner(1389, 14)
	cnt5 := 0
	cnt6 := 0
	cnt7 := 0
	cnt8 := 0
	cnt9 := 0
	cnt10 := 0
	cntOther := 0
	for i := 0; i < 100; i++ {
		val := utils.GetRandInRange(5, 10)
		switch val {
		case 5:
			cnt5++
		case 6:
			cnt6++
		case 7:
			cnt7++
		case 8:
			cnt8++
		case 9:
			cnt9++
		case 10:
			cnt10++
		default:
			cntOther++
		}
		//fmt.Println(val)
	}
	fmt.Println(cnt5, cnt6, cnt7, cnt8, cnt9, cnt10, cntOther)
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

	maxRand := 1e10
	for j := 0; j < 100; j++ {
		ratio := float64(rand.Int63n(int64(maxRand))) / maxRand
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

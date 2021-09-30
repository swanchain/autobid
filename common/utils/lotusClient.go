package utils

import (
	"go-swan/logs"
	"go-swan/models"
	"strings"
)

func LotusGetMinerInfo(miner *models.Miner) bool {
	cmd := "lotus client query-ask " + miner.MinerFid
	logs.GetLogger().Info(cmd)

	result, err := ExecOsCmd(cmd)

	if err != nil {
		logs.GetLogger().Error(err)
		return false
	}

	if len(result) == 0 {
		logs.GetLogger().Error("Failed to get info for:", miner.MinerFid)
		return false
	}

	lines := strings.Split(result, "\n")
	logs.GetLogger().Info(lines)

	for _, line := range lines {
		if strings.Contains(line, "Verified Price per GiB:") {
			miner.VerifiedPrice = SearchFloat64FromStr(line)
			if miner.VerifiedPrice != nil {
				logs.GetLogger().Info("miner VerifiedPrice: ", *miner.VerifiedPrice)
			} else {
				logs.GetLogger().Error("Failed to get miner VerifiedPrice from lotus")
			}

			continue
		}

		if strings.Contains(line, "Price per GiB:") {
			miner.Price = SearchFloat64FromStr(line)
			if miner.Price != nil {
				logs.GetLogger().Info("miner Price: ", *miner.Price)
			} else {
				logs.GetLogger().Error("Failed to get miner Price from lotus")
			}

			continue
		}

		if strings.Contains(line, "Max Piece size:") {
			words := strings.Split(line, ":")
			if len(words) == 2 {
				maxPieceSize := strings.Trim(words[1], " ")
				miner.MaxPieceSize = &maxPieceSize
				if miner.MaxPieceSize != nil {
					logs.GetLogger().Info("miner MaxPieceSize: ", *miner.MaxPieceSize)
				} else {
					logs.GetLogger().Error("Failed to get miner MaxPieceSize from lotus")
				}
			}
			continue
		}

		if strings.Contains(line, "Min Piece size:") {
			words := strings.Split(line, ":")
			if len(words) == 2 {
				minPieceSize := strings.Trim(words[1], " ")
				miner.MinPieceSize = &minPieceSize
				if miner.MinPieceSize != nil {
					logs.GetLogger().Info("miner MinPieceSize: ", *miner.MinPieceSize)
				} else {
					logs.GetLogger().Error("Failed to get miner MinPieceSize from lotus")
				}
			}
			continue
		}
	}

	return true
}

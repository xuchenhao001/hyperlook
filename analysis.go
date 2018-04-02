package main

import (
	"log"
	"math"
	"regexp"
)

func analysisLogs(logsArray *[]HitContent) {
	// Analysis cscc, escc, lscc, qscc, vscc, and generateDockerfile keyword
	// Set their weight to 1, 2, 4, 8, 16, 32
	// So that:
	// join channel: 31
	// install chaincode: 16
	// instantiate chaincode: 42
	// invoke chaincode: 22
	// query chaincode: 6
	// upgrade chaincode: 54
	csccReg, err := regexp.Compile(".*chain=[^,].*chaincode=cscc.*")
	if err != nil {
		log.Printf("Error cannot compile regex: %s", err)
		return
	}
	esccReg, err := regexp.Compile(".*chain=[^,].*chaincode=escc.*")
	if err != nil {
		log.Printf("Error cannot compile regex: %s", err)
		return
	}
	lsccReg, err := regexp.Compile(".*chain=[^,].*chaincode=lscc.*")
	if err != nil {
		log.Printf("Error cannot compile regex: %s", err)
		return
	}
	qsccReg, err := regexp.Compile(".*chain=[^,].*chaincode=qscc.*")
	if err != nil {
		log.Printf("Error cannot compile regex: %s", err)
		return
	}
	vsccReg, err := regexp.Compile(".*chain=[^,].*chaincode=vscc.*")
	if err != nil {
		log.Printf("Error cannot compile regex: %s", err)
		return
	}
	genDocReg, err := regexp.Compile(".*generateDockerfile.*")
	if err != nil {
		log.Printf("Error cannot compile regex: %s", err)
		return
	}

	// proposal starts from "Entry" and ends with "Exit"
	startReg, err := regexp.Compile(".*ProcessProposal.*Entry.*")
	if err != nil {
		log.Printf("Error cannot compile regex: %s", err)
		return
	}
	endReg, err := regexp.Compile(".*ProcessProposal.*Exit.*")
	if err != nil {
		log.Printf("Error cannot compile regex: %s", err)
		return
	}


	// weight[6] respectively represent: cscc, escc, lscc, qscc, vscc, and generateDockerfile
	var weight[6] int
	var endTime uint64
	// analysed[6] respectively represent: join channel, install chaincode, instantiate chaincode,
	// upgrade chaincode, invoke, and query. So that the Hyperlook always exports the newest data.
	var analysed[6] bool

	for _, hit := range *logsArray {
		// because logs are in a time reverse order, so we analysis from end to start.
		matchEnd := endReg.FindString(hit.Source.Log)
		if matchEnd != "" {
			endTime = hit.Sort[0]
		}

		matchStart := startReg.FindString(hit.Source.Log)
		if matchStart != "" {
			// if there isn't a end time, go ahead
			if endTime == 0 {
				// clean weight array
				for index := range weight {
					weight[index] = 0
				}
				continue
			}
			var weightSum float64 = 0
			for index := range weight {
				weightSum += float64(weight[index]) * math.Pow(2, float64(index))
			}
			// join channel
			if math.Abs(weightSum - 31) <= 0.1 {
				if analysed[0] == true {
					continue
				}
				duration := float64(endTime - hit.Sort[0])/1000
				log.Printf("Info join channel from %d to %d, duration: %fs", hit.Sort[0], endTime, duration)
				peerJoinChannel.Set(duration)
				analysed[0] = true
			} else
			// install chaincode
			if math.Abs(weightSum - 16) <= 0.1 {
				if analysed[1] == true {
					continue
				}
				duration := float64(endTime - hit.Sort[0])/1000
				log.Printf("Info install chaincode from %d to %d, duration: %fs", hit.Sort[0], endTime, duration)
				peerInstallChaincode.Set(duration)
				analysed[1] = true
			} else
			// instantiate chaincode
			if math.Abs(weightSum - 42) <= 0.1 {
				if analysed[2] == true {
					continue
				}
				duration := float64(endTime - hit.Sort[0])/1000
				log.Printf("Info instantiate chaincode from %d to %d, duration: %fs", hit.Sort[0], endTime, duration)
				peerInstantiateChaincode.Set(duration)
				analysed[2] = true
			} else
			// upgrade chaincode
			if math.Abs(weightSum - 54) <= 0.1 {
				if analysed[3] == true {
					continue
				}
				duration := float64(endTime - hit.Sort[0])/1000
				log.Printf("Info upgrade chaincode from %d to %d, duration: %fs", hit.Sort[0], endTime, duration)
				peerUpgradeChaincode.Set(duration)
				analysed[3] = true
			} else
			// invoke
			if math.Abs(weightSum - 22) <= 0.1 {
				if analysed[4] == true {
					continue
				}
				duration := float64(endTime - hit.Sort[0])/1000
				log.Printf("Info invoke from %d to %d, duration: %fs", hit.Sort[0], endTime, duration)
				peerInvokeChaincode.Set(duration)
				analysed[4] = true
			} else
			// query
			if math.Abs(weightSum - 6) <= 0.1 {
				if analysed[5] == true {
					continue
				}
				duration := float64(endTime - hit.Sort[0])/1000
				log.Printf("Info query from %d to %d, duration: %fs", hit.Sort[0], endTime, duration)
				peerQueryChaincode.Set(duration)
				analysed[5] = true
			}

			// clean weight array
			for index := range weight {
				weight[index] = 0
			}
			endTime = 0
			continue
		}

		matchCscc := csccReg.FindString(hit.Source.Log)
		if matchCscc != "" {
			weight[0] = 1
		}
		matchEscc := esccReg.FindString(hit.Source.Log)
		if matchEscc != "" {
			weight[1] = 1
		}
		matchLscc := lsccReg.FindString(hit.Source.Log)
		if matchLscc != "" {
			weight[2] += 1
		}
		matchQscc := qsccReg.FindString(hit.Source.Log)
		if matchQscc != "" {
			weight[3] = 1
		}
		matchVscc := vsccReg.FindString(hit.Source.Log)
		if matchVscc != "" {
			weight[4] = 1
		}
		matchGenDoc := genDocReg.FindString(hit.Source.Log)
		if matchGenDoc != "" {
			weight[5] = 1
		}
	}
}

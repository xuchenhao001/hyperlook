package main

import (
	"log"
	"math"
	"regexp"
)

// refer: https://stackoverflow.com/questions/32751537/w
// hy-do-i-get-a-cannot-assign-error-when-setting-value-to-a-struct-as-a-value-i
// golang hash value weight[containerName] is not addressable, so that cannot be
// dynamically assigned value
type Weight struct {
	w [6]int
}
type Weights map[string]*Weight

type Analysis struct {
	done [6]bool
}
type Analysises map[string]*Analysis

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


	// Weight.w[6] respectively indicates: cscc, escc, lscc, qscc, vscc, and generateDockerfile
	// weights indicates all container's weights
	weights := make(Weights)
	// endTime indicates all container's end time for each fabric event
	endTime := make(map[string]uint64)
	// Analysis.done[6] respectively represent: join channel, install chaincode, instantiate chaincode,
	// upgrade chaincode, invoke, and query. So that the Hyperlook always exports the newest data.
	// analyses indicated all container's analyses
	analyses := make(Analysises)

	for _, hit := range *logsArray {
		containerName := hit.Source.Kubernetes.ContainerName
		// Since weight[containerName] points to a stable address, we must allocate for it
		// at the first time.
		if _, ok := weights[containerName]; !ok {
			weights[containerName] = new(Weight)
		}
		if _, ok := analyses[containerName]; !ok {
			analyses[containerName] = new(Analysis)
		}

		// because logs are in a time reverse order, so we analysis from end to start.
		matchEnd := endReg.FindString(hit.Source.Log)
		if matchEnd != "" {
			endTime[containerName] = hit.Sort[0]
		}

		matchStart := startReg.FindString(hit.Source.Log)
		if matchStart != "" {
			// if there isn't a end time, go ahead
			if endTime[containerName] == 0 {
				// clean weight array
				for index := range weights[containerName].w {
					weights[containerName].w[index] = 0
				}
				continue
			}
			var weightSum float64 = 0
			for index := range weights[containerName].w {
				weightSum += float64(weights[containerName].w[index]) * math.Pow(2, float64(index))
			}
			// join channel
			if math.Abs(weightSum - 31) <= 0.1 {
				if analyses[containerName].done[0] == true {
					continue
				}
				duration := float64(endTime[containerName] - hit.Sort[0])/1000
				log.Printf("Info %s join channel from %d to %d, duration: %fs",
					containerName, hit.Sort[0], endTime[containerName], duration)
				peerJoinChannel.WithLabelValues(containerName).Set(duration)
				analyses[containerName].done[0] = true
			} else
			// install chaincode
			if math.Abs(weightSum - 16) <= 0.1 {
				if analyses[containerName].done[1] == true {
					continue
				}
				duration := float64(endTime[containerName] - hit.Sort[0])/1000
				log.Printf("Info %s install chaincode from %d to %d, duration: %fs",
					containerName, hit.Sort[0], endTime[containerName], duration)
				peerInstallChaincode.WithLabelValues(containerName).Set(duration)
				analyses[containerName].done[1] = true
			} else
			// instantiate chaincode
			if math.Abs(weightSum - 42) <= 0.1 {
				if analyses[containerName].done[2] == true {
					continue
				}
				duration := float64(endTime[containerName] - hit.Sort[0])/1000
				log.Printf("Info %s instantiate chaincode from %d to %d, duration: %fs",
					containerName, hit.Sort[0], endTime[containerName], duration)
				peerInstantiateChaincode.WithLabelValues(containerName).Set(duration)
				analyses[containerName].done[2] = true
			} else
			// upgrade chaincode
			if math.Abs(weightSum - 54) <= 0.1 {
				if analyses[containerName].done[3] == true {
					continue
				}
				duration := float64(endTime[containerName] - hit.Sort[0])/1000
				log.Printf("Info %s upgrade chaincode from %d to %d, duration: %fs",
					containerName, hit.Sort[0], endTime[containerName], duration)
				peerUpgradeChaincode.WithLabelValues(containerName).Set(duration)
				analyses[containerName].done[3] = true
			} else
			// invoke
			if math.Abs(weightSum - 22) <= 0.1 {
				if analyses[containerName].done[4] == true {
					continue
				}
				duration := float64(endTime[containerName] - hit.Sort[0])/1000
				log.Printf("Info %s invoke from %d to %d, duration: %fs",
					containerName, hit.Sort[0], endTime[containerName], duration)
				peerInvokeChaincode.WithLabelValues(containerName).Set(duration)
				analyses[containerName].done[4] = true
			} else
			// query
			if math.Abs(weightSum - 6) <= 0.1 {
				if analyses[containerName].done[5] == true {
					continue
				}
				duration := float64(endTime[containerName] - hit.Sort[0])/1000
				log.Printf("Info %s query from %d to %d, duration: %fs",
					containerName, hit.Sort[0], endTime[containerName], duration)
				peerQueryChaincode.WithLabelValues(containerName).Set(duration)
				analyses[containerName].done[5] = true
			}

			// clean weight array
			for index := range weights[containerName].w {
				weights[containerName].w[index] = 0
			}
			endTime[containerName] = 0
			continue
		}

		matchCscc := csccReg.FindString(hit.Source.Log)
		if matchCscc != "" {
			weights[containerName].w[0] = 1
		}
		matchEscc := esccReg.FindString(hit.Source.Log)
		if matchEscc != "" {
			weights[containerName].w[1] = 1
		}
		matchLscc := lsccReg.FindString(hit.Source.Log)
		if matchLscc != "" {
			weights[containerName].w[2] += 1
		}
		matchQscc := qsccReg.FindString(hit.Source.Log)
		if matchQscc != "" {
			weights[containerName].w[3] = 1
		}
		matchVscc := vsccReg.FindString(hit.Source.Log)
		if matchVscc != "" {
			weights[containerName].w[4] = 1
		}
		matchGenDoc := genDocReg.FindString(hit.Source.Log)
		if matchGenDoc != "" {
			weights[containerName].w[5] = 1
		}
	}
}

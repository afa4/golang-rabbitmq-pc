package main

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func main() {
	lenOfByteArray := []int{256, 512, 1024, 2048, 5096}

	for _, lenOfByte := range lenOfByteArray {
		var wg sync.WaitGroup
		for numOfPublishers := 0; numOfPublishers <= 40; numOfPublishers++{
			invokePublisherCommand := fmt.Sprintf("go run publisher/publisher.go %d", lenOfByte)
			wg.Add(1)
			go callPublisher(invokePublisherCommand)
		}
		wg.Wait()
		// finalExecMean := 0
		// for execMean := range execMeans{
		// 	finalExecMean += execMean/40
		// }

		// fmt.Printf("Final Execution Mean (%d Bytes): %d", lenOfByte, finalExecMean)
	}
}

func callPublisher(invokePublisherCommand string){
	bash := exec.Command("bash", "-c", invokePublisherCommand)
	// output, err := bash.Output()
	// outputStr := string(output)
	// repliesDurations := strings.Split(outputStr, "\n")
	// // Revisar o len -1
	// numberOfDurations := len(repliesDurations) - 1
	// totalDuration := 0
	// for i := 0; i < numberOfDurations; i++ {
	// 	messageDuration, err := strconv.Atoi(repliesDurations[i])
	// 	if err != nil {
	// 		panic(err.Error())
	// 	}
	// 	totalDuration += replyDuration
	// }
	// execMeans <- (totalDuration/numberOfDurations)
}

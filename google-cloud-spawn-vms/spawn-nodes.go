package main

import (
	"fmt"
	"os/exec"
	"sync"
)

var seedsCount = 3
var seedImageProject = "debian-cloud"
var seedImageFamily = "debian-9"
var seedMachineType = "f1-micro"
var seedTags = "disgo-node"
var seedStartupScript = "vm-debian9-configure.sh"

// var delegatesCount = 21
// var nodesCount = 50
var vmPrefix = "test-net-1-1"

var nodeScriptConfigURL = "https://raw.githubusercontent.com/dispatchlabs/samples/master/google-cloud-spawn-vms"
var nodeScriptConfigFile1 = "vm-debian9-configure-part1.sh"
var nodeScriptConfigFile2 = "vm-debian9-configure-part2.sh"

func main() {
	var wg sync.WaitGroup

	for i := 0; i < seedsCount; i++ {

		// Command to CREATE new VM Instance
		var createVM = fmt.Sprintf(
			"gcloud compute instances create %s-seed-%d --image-project %s --image-family %s --machine-type %s --tags %s",
			vmPrefix,
			i,
			seedImageProject,
			seedImageFamily,
			seedMachineType,
			seedTags,
		)

		// Command to DOWNLOAD BASH scripts to the newly created VM
		var downloadScriptFiles = fmt.Sprintf(
			"gcloud compute ssh %s-seed-%d --command 'curl %s/%s -o %s && curl %s/%s -o %s'",
			vmPrefix,
			i,

			nodeScriptConfigURL,
			nodeScriptConfigFile1,
			nodeScriptConfigFile1,

			nodeScriptConfigURL,
			nodeScriptConfigFile2,
			nodeScriptConfigFile2,
		)

		// Commands to RUN scripts
		var execScript1 = fmt.Sprintf(
			"gcloud compute ssh %s-seed-%d --command 'bash %s'",
			vmPrefix,
			i,
			nodeScriptConfigFile1,
		)
		var execScript2 = fmt.Sprintf(
			"gcloud compute ssh %s-seed-%d --command 'bash %s'",
			vmPrefix,
			i,
			nodeScriptConfigFile2,
		)

		// Run all the commands in sequential order inside the new VM
		// Each VM is created in PARALLEL
		wg.Add(1)
		go func(cmds ...string) {
			for _, cmd := range cmds {
				exec.Command("sh", "-c", cmd).Run()
			}
			wg.Done()
		}(createVM, downloadScriptFiles, execScript1, execScript2)
	}

	wg.Wait()
}

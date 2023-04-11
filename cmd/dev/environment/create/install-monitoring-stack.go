package create

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"time"

	"github.com/power-slide/cli/cmd/util"
)

var (
	//go:embed static/002-monitoring.yaml
	monitoringHelmChart string
)

func installMonitoringStack() {
	if skipMonitoring {
		return
	}
	fmt.Print("Installing Prometheus monitoring stack (kube-prometheus-stack)... ")
	util.CreateNamespace("pwrsl-monitoring")
	util.Kubectl([]string{"apply", "-f", "-"}, monitoringHelmChart)

	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()
	for {
		if util.KubectlHasCRD(ctx, "servicemonitors.monitoring.coreos.com") {
			break
		} else if ctx.Err() != nil {
			fmt.Println()
			log.Fatalln("Unable to install monitoring stack within", cmdTimeout)
		}
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Done!")
}

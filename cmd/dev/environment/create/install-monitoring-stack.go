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
	//go:embed static/monitoring/000-prometheus-operator.yaml
	prometheusOperatorHelmChart string
)

func installMonitoringStack() {
	if skipMonitoring {
		return
	}
	fmt.Print("Installing Prometheus monitoring stack (kube-prometheus-stack)... ")
	util.CreateNamespace("pwrsl-monitoring")
	util.Kubectl([]string{"apply", "-f", "-"}, prometheusOperatorHelmChart)
	
	requiredPrometheusCRDs := []string{
		"alertmanagerconfigs.monitoring.coreos.com",
		"alertmanagers.monitoring.coreos.com",
		"podmonitors.monitoring.coreos.com",
		"probes.monitoring.coreos.com",
		"prometheuses.monitoring.coreos.com",
		"prometheusrules.monitoring.coreos.com",
		"servicemonitors.monitoring.coreos.com",
		"thanosrulers.monitoring.coreos.com",
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()
	for {
		if util.KubectlHasAllCRDs(ctx, requiredPrometheusCRDs) {
			break
		} else if ctx.Err() != nil {
			fmt.Println()
			log.Fatalln("Unable to install monitoring stack within", cmdTimeout)
		}
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Done!")
}

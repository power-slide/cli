package create

import (
	"context"
	_ "embed"
	"fmt"
	"time"

	"github.com/power-slide/cli/cmd/util"
	"github.com/power-slide/cli/pkg/logger"
)

const (
	monitoringNamespace = "pwrsl-monitoring"
)

var (
	//go:embed static/monitoring/000-prometheus-operator.yaml
	prometheusOperatorHelmChart string
	//go:embed static/monitoring/001-traefik.yaml
	traefikMonitoringManifest string
)

func installMonitoringStack() {
	if skipMonitoring {
		return
	}

	installPrometheusOperator()
	addTraefikMonitoring()
}

func installPrometheusOperator() {
	fmt.Print("Installing Prometheus operator (kube-prometheus-stack)... ")
	util.CreateNamespace(monitoringNamespace)
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
		if ctx.Err() != nil {
			fmt.Println("Error!")
			logger.CheckErr(ctx.Err())
		}

		if util.KubectlHasAllCRDs(ctx, requiredPrometheusCRDs) {
			break
		}

		time.Sleep(1 * time.Second)
	}

	ctx, cancel = context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()
	for {
		if ctx.Err() != nil {
			fmt.Println("Error!")
			logger.CheckErr(ctx.Err())
		}

		result := util.KubectlJSON(
			ctx,
			[]string{
				"get", "prometheuses",
				"-n", monitoringNamespace,
			},
		)

		items := result["items"].([]any)
		if len(items) > 0 {
			break
		}

		time.Sleep(1 * time.Second)
	}

	util.Kubectl(
		[]string{
			"wait", "prometheus",
			"-n", monitoringNamespace,
			"pwrsl-monitoring-kube-prom-prometheus",
			"--for", "condition=Available",
		},
		"",
	)

	fmt.Println("Done!")
}

func addTraefikMonitoring() {
	fmt.Print("Adding service monitors for Traefik... ")
	util.Kubectl([]string{"apply", "-f", "-"}, traefikMonitoringManifest)
	fmt.Println("Done!")
}

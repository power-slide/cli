package create

import (
	"context"
	_ "embed"
	"encoding/base64"
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
	//go:embed static/monitoring/002-grafana.yaml
	grafanaStackHelmChart string
)

func installMonitoringStack() {
	if skipMonitoring {
		return
	}

	installPrometheusOperator()
	addTraefikMonitoring()
	addGrafanaStack()
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
			"pwrsl-prometheus-kube-prom-prometheus",
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

func addGrafanaStack() {
	fmt.Print("Installing Grafana observability stack... ")
	util.Kubectl([]string{"apply", "-f", "-"}, grafanaStackHelmChart)

	var data any
	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()
	for {
		if ctx.Err() != nil {
			fmt.Println("Error!")
			logger.CheckErr(ctx.Err())
		}

		data := util.KubectlJSON(
			ctx,
			[]string{
				"get", "secret", "pwrsl-grafana",
				"-n", monitoringNamespace,
			},
		)["data"]

		if data != nil {
			break
		}

		time.Sleep(1 * time.Second)
	}

	encPass := data.(map[string]any)["admin-password"].(string)
	passBytes, err := base64.URLEncoding.DecodeString(encPass)
	logger.CheckErr(err)
	fmt.Println("Done!")
	fmt.Printf("Grafana 'admin' password is: %s\n", string(passBytes))
}

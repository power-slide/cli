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
	//go:embed static/003-argocd.yaml
	argocdHelmChart string
)

func installArgoCD() {
	if skipArgoCD {
		return
	}
	fmt.Print("Installing ArgoCD... ")
	util.CreateNamespace("pwrsl-argocd")
	util.Kubectl([]string{"apply", "-f", "-"}, argocdHelmChart)

	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()
	for {
		result := util.KubectlHasCRD(ctx, "applications.argoproj.io")
		if result {
			break
		} else if ctx.Err() != nil {
			fmt.Println()
			log.Fatalln("Unable to install ArgoCD within", cmdTimeout)
		}
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Done!")
}

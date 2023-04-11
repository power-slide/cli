package create

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/power-slide/cli/cmd/util"
)

var (
	//go:embed static/000-gateway-api.yaml
	gatewayAPI string
	//go:embed static/001-traefik.yaml
	traefikHelmConfig string
)

func createCluster(clusterName string) {
	if skipCluster {
		return
	}
	fmt.Printf("Creating local cluster... ")

	clusterCommandArgs := []string{
		"cluster",
		"create",
		clusterName,
		"-p", fmt.Sprintf("%d:80@loadbalancer", clusterPort),
	}

	cmd := exec.Command("k3d", clusterCommandArgs...)
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Println("k3d output:\n", string(out))
		log.Fatalln(err)
	}
	fmt.Println("Done!")

	fmt.Print("Modifying cluster to support PowerSlide... ")

	util.Kubectl([]string{"apply", "-f", "-"}, gatewayAPI)
	util.Kubectl([]string{"apply", "-f", "-"}, traefikHelmConfig)

	ctx, cancel := context.WithTimeout(context.Background(), cmdTimeout)
	defer cancel()
	for {
		result := util.KubectlJSON(ctx, []string{"get", "gatewayclass"})
		items := result["items"].([]any)
		if len(items) > 0 {
			break
		} else if ctx.Err() != nil {
			fmt.Println()
			log.Fatalln("Unable to configure cluster within", cmdTimeout)
		}
		time.Sleep(1 * time.Second)
	}

	fmt.Println("Done!")
}

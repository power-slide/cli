package util

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/power-slide/cli/pkg/logger"
)

func Kubectl(args []string, input string) string {
	cmd := exec.Command("kubectl", args...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	if input != "" {
		stdin, err := cmd.StdinPipe()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			defer stdin.Close()
			io.WriteString(stdin, input)
		}()
	}

	if err := cmd.Run(); err != nil {
		fmt.Println("kubectl output:")
		fmt.Println(arrayToCleanString(outb.Bytes()))
		logger.CheckErr(arrayToCleanString(errb.Bytes()))
	}

	return arrayToCleanString(outb.Bytes())
}

func KubectlJSON(ctx context.Context, args []string) map[string]any {
	cmd := exec.CommandContext(ctx, "kubectl", append(args, "-o", "json")...)
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb

	if err := cmd.Run(); err != nil {
		// var result map[string]any
		result := make(map[string]any)
		result["error"] = true
		result["items"] = []any{}
		return result
	}

	var result map[string]any
	json.Unmarshal(outb.Bytes(), &result)
	return result
}

func KubectlGetCRDS(ctx context.Context) []string {
	result := KubectlJSON(ctx, []string{"get", "crds"})
	items := result["items"].([]any)
	crds := []string{}

	for _, v := range items {
		name := v.(map[string]any)["metadata"].(map[string]any)["name"].(string)
		if name != "" {
			crds = append(crds, name)
		}
	}

	return crds
}

func KubectlHasCRD(ctx context.Context, targetCRD string) bool {
	crds := KubectlGetCRDS(ctx)
	for _, crd := range crds {
		if crd == targetCRD {
			return true
		}
	}

	return false
}

func KubectlHasAllCRDs(ctx context.Context, targetCRDs []string) bool {
	installedCRDs := KubectlGetCRDS(ctx)
	installedHash := make(map[string]bool)

	for _, crd := range installedCRDs {
		installedHash[crd] = true
	}

	for _, crd := range targetCRDs {
		if !installedHash[crd] {
			return false
		}
	}

	return true
}

func CreateNamespace(namespace string) {
	Kubectl([]string{"create", "namespace", namespace}, "")
	Kubectl([]string{"label", "namespace", namespace, "powerslide.cloud/infrastructure=true"}, "")
}

func arrayToCleanString(input []byte) string {
	return strings.TrimSpace(string(input))
}

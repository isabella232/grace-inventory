package testing

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestNow(t *testing.T) {

	err := os.Setenv("AWS_ACCESS_KEY_ID", "foobar")
	if err != nil {
		t.Fatalf("failed to set AWS_ACCESS_KEY_ID: %v", err)
	}

	err = os.Setenv("AWS_SECRET_ACCESS_KEY", "foobar")
	if err != nil {
		t.Fatalf("failed to set AWS_SECRET_ACCESS_KEY: %v", err)
	}

	motoServer := `/usr/bin/moto_server`
	cmd := exec.Command(motoServer)
	go func(t *testing.T) {
		err := cmd.Run()
		if err != nil {
			t.Fatalf("failed to execute localstack: %v", err)
		}
	}(t)
	opts := &terraform.Options{

		// The path to where our Terraform code is located
		TerraformDir: "scenarios/one",

		// Variables to pass to our Terraform code using -var options
		Vars: map[string]interface{}{
			"appenv":            "development",
			"tenant_role_name":  "0",
			"master_role_name":  "1",
			"master_account_id": "2",
		},

		// Disable colors in Terraform commands so its easier to parse stdout/stderr
		NoColor: true,
	}
	defer terraform.Destroy(t, opts)
	t.Logf("output: %s\n", terraform.InitAndApply(t, opts))

	pattern := `/tmp/localstack/data/*.json`
	matches, err := filepath.Glob(pattern)
	if err != nil {
		t.Fatalf("failed to glob files %s: %v", pattern, err)
	}
	for _, m := range matches {
		t.Logf("found file: %s\n", m)
	}
}
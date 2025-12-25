package laravel

import (
	"fmt"
	"os/exec"
)

type ReportingSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewReportingSetup(projectPath string, dryRun bool) *ReportingSetup {
	return &ReportingSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (r *ReportingSetup) Setup() error {
	if r.DryRun {
		fmt.Printf("[Dry Run] Would install Excel and PDF support\n")
		return nil
	}

	fmt.Println("📊 Installing Maatwebsite Excel...")
	if err := r.runComposerRequire("maatwebsite/excel"); err != nil {
		return err
	}

	fmt.Println("📄 Installing dompdf/dompdf...")
	if err := r.runComposerRequire("dompdf/dompdf"); err != nil {
		return err
	}

	return nil
}

func (r *ReportingSetup) runComposerRequire(pkg string) error {
	cmd := exec.Command("composer", "require", pkg, "--with-all-dependencies")
	cmd.Dir = r.ProjectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install %s: %v\nOutput: %s", pkg, err, string(output))
	}
	return nil
}

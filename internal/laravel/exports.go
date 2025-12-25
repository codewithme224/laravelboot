package laravel

import (
	"fmt"
	"os"
	"path/filepath"
)

type ExportsSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewExportsSetup(projectPath string, dryRun bool) *ExportsSetup {
	return &ExportsSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (e *ExportsSetup) Setup() error {
	if e.DryRun {
		fmt.Printf("[Dry Run] Would create app/Exports and app/Imports directories with base classes\n")
		return nil
	}

	fmt.Println("📤 Setting up Exports/Imports structure...")

	// Create directories
	exportsDir := filepath.Join(e.ProjectPath, "app/Exports")
	importsDir := filepath.Join(e.ProjectPath, "app/Imports")
	os.MkdirAll(exportsDir, 0755)
	os.MkdirAll(importsDir, 0755)

	// Create base export class
	if err := e.createBaseExport(); err != nil {
		return err
	}

	// Create base import class
	if err := e.createBaseImport(); err != nil {
		return err
	}

	return nil
}

func (e *ExportsSetup) createBaseExport() error {
	content := `<?php

namespace App\Exports;

use Maatwebsite\Excel\Concerns\FromCollection;
use Maatwebsite\Excel\Concerns\WithHeadings;
use Maatwebsite\Excel\Concerns\WithMapping;
use Maatwebsite\Excel\Concerns\ShouldAutoSize;

abstract class BaseExport implements FromCollection, WithHeadings, WithMapping, ShouldAutoSize
{
    /**
     * @return \Illuminate\Support\Collection
     */
    abstract public function collection();

    /**
     * @return array
     */
    abstract public function headings(): array;

    /**
     * @param mixed $row
     * @return array
     */
    abstract public function map($row): array;
}
`
	path := filepath.Join(e.ProjectPath, "app/Exports/BaseExport.php")
	return os.WriteFile(path, []byte(content), 0644)
}

func (e *ExportsSetup) createBaseImport() error {
	content := `<?php

namespace App\Imports;

use Maatwebsite\Excel\Concerns\ToModel;
use Maatwebsite\Excel\Concerns\WithHeadingRow;
use Maatwebsite\Excel\Concerns\WithValidation;
use Maatwebsite\Excel\Concerns\SkipsOnError;
use Maatwebsite\Excel\Concerns\SkipsErrors;

abstract class BaseImport implements ToModel, WithHeadingRow, WithValidation, SkipsOnError
{
    use SkipsErrors;

    /**
     * @param array $row
     * @return \Illuminate\Database\Eloquent\Model|null
     */
    abstract public function model(array $row);

    /**
     * @return array
     */
    abstract public function rules(): array;
}
`
	path := filepath.Join(e.ProjectPath, "app/Imports/BaseImport.php")
	return os.WriteFile(path, []byte(content), 0644)
}

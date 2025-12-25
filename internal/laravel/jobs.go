package laravel

import (
	"fmt"
	"os"
	"path/filepath"
)

type JobsSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewJobsSetup(projectPath string, dryRun bool) *JobsSetup {
	return &JobsSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (j *JobsSetup) Setup() error {
	if j.DryRun {
		fmt.Printf("[Dry Run] Would create app/Jobs directory with base job class\n")
		return nil
	}

	fmt.Println("⚙️ Setting up Jobs structure...")

	jobsDir := filepath.Join(j.ProjectPath, "app/Jobs")
	os.MkdirAll(jobsDir, 0755)

	if err := j.createBaseJob(); err != nil {
		return err
	}

	return nil
}

func (j *JobsSetup) createBaseJob() error {
	content := `<?php

namespace App\Jobs;

use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Foundation\Bus\Dispatchable;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Queue\SerializesModels;
use Illuminate\Support\Facades\Log;

abstract class BaseJob implements ShouldQueue
{
    use Dispatchable, InteractsWithQueue, Queueable, SerializesModels;

    /**
     * The number of times the job may be attempted.
     */
    public int $tries = 3;

    /**
     * The number of seconds to wait before retrying the job.
     */
    public int $backoff = 60;

    /**
     * Handle a job failure.
     */
    public function failed(\Throwable $exception): void
    {
        Log::error('Job failed: ' . static::class, [
            'exception' => $exception->getMessage(),
            'trace' => $exception->getTraceAsString(),
        ]);
    }

    /**
     * Execute the job.
     */
    abstract public function handle(): void;
}
`
	path := filepath.Join(j.ProjectPath, "app/Jobs/BaseJob.php")
	return os.WriteFile(path, []byte(content), 0644)
}

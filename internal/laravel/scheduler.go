package laravel

import (
	"fmt"
	"os"
	"path/filepath"
)

type SchedulerSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewSchedulerSetup(projectPath string, dryRun bool) *SchedulerSetup {
	return &SchedulerSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (s *SchedulerSetup) Setup() error {
	if s.DryRun {
		fmt.Printf("[Dry Run] Would setup scheduler and console commands\n")
		return nil
	}

	fmt.Println("⏰ Setting up scheduler and console commands...")

	commandsDir := filepath.Join(s.ProjectPath, "app/Console/Commands")
	os.MkdirAll(commandsDir, 0755)

	if err := s.createBaseCommand(); err != nil {
		return err
	}
	if err := s.createCleanupCommand(); err != nil {
		return err
	}
	if err := s.createHealthCheckCommand(); err != nil {
		return err
	}

	return nil
}

func (s *SchedulerSetup) createBaseCommand() error {
	content := `<?php

namespace App\Console\Commands;

use Illuminate\Console\Command;
use Illuminate\Support\Facades\Log;

abstract class BaseCommand extends Command
{
    /**
     * Log and display info message.
     */
    protected function logInfo(string $message): void
    {
        $this->info($message);
        Log::info("[{$this->signature}] {$message}");
    }

    /**
     * Log and display error message.
     */
    protected function logError(string $message): void
    {
        $this->error($message);
        Log::error("[{$this->signature}] {$message}");
    }

    /**
     * Log and display warning message.
     */
    protected function logWarning(string $message): void
    {
        $this->warn($message);
        Log::warning("[{$this->signature}] {$message}");
    }

    /**
     * Execute with timing and logging.
     */
    protected function executeWithTiming(callable $callback): int
    {
        $startTime = microtime(true);
        $this->logInfo('Starting execution...');

        try {
            $result = $callback();
            $duration = round(microtime(true) - $startTime, 2);
            $this->logInfo("Completed in {$duration}s");
            return $result ?? self::SUCCESS;
        } catch (\Exception $e) {
            $this->logError("Failed: {$e->getMessage()}");
            return self::FAILURE;
        }
    }
}
`
	path := filepath.Join(s.ProjectPath, "app/Console/Commands/BaseCommand.php")
	return os.WriteFile(path, []byte(content), 0644)
}

func (s *SchedulerSetup) createCleanupCommand() error {
	content := `<?php

namespace App\Console\Commands;

use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Storage;

class CleanupCommand extends BaseCommand
{
    protected $signature = 'app:cleanup 
                            {--days=30 : Days to keep data}
                            {--dry-run : Run without deleting}';

    protected $description = 'Clean up old data and temporary files';

    public function handle(): int
    {
        return $this->executeWithTiming(function () {
            $days = (int) $this->option('days');
            $dryRun = $this->option('dry-run');

            if ($dryRun) {
                $this->logWarning('Running in dry-run mode - no data will be deleted');
            }

            // Clean old notifications
            $this->cleanNotifications($days, $dryRun);

            // Clean old activity logs (if using spatie/activitylog)
            $this->cleanActivityLogs($days, $dryRun);

            // Clean temporary files
            $this->cleanTempFiles($dryRun);

            return self::SUCCESS;
        });
    }

    protected function cleanNotifications(int $days, bool $dryRun): void
    {
        $count = DB::table('notifications')
            ->where('created_at', '<', now()->subDays($days))
            ->count();

        if (!$dryRun && $count > 0) {
            DB::table('notifications')
                ->where('created_at', '<', now()->subDays($days))
                ->delete();
        }

        $this->logInfo("Notifications: {$count} old records " . ($dryRun ? 'would be' : '') . " deleted");
    }

    protected function cleanActivityLogs(int $days, bool $dryRun): void
    {
        if (!class_exists(\Spatie\Activitylog\Models\Activity::class)) {
            return;
        }

        $count = DB::table('activity_log')
            ->where('created_at', '<', now()->subDays($days))
            ->count();

        if (!$dryRun && $count > 0) {
            DB::table('activity_log')
                ->where('created_at', '<', now()->subDays($days))
                ->delete();
        }

        $this->logInfo("Activity logs: {$count} old records " . ($dryRun ? 'would be' : '') . " deleted");
    }

    protected function cleanTempFiles(bool $dryRun): void
    {
        $files = Storage::disk('local')->files('temp');
        $count = count($files);

        if (!$dryRun && $count > 0) {
            foreach ($files as $file) {
                Storage::disk('local')->delete($file);
            }
        }

        $this->logInfo("Temp files: {$count} files " . ($dryRun ? 'would be' : '') . " deleted");
    }
}
`
	path := filepath.Join(s.ProjectPath, "app/Console/Commands/CleanupCommand.php")
	return os.WriteFile(path, []byte(content), 0644)
}

func (s *SchedulerSetup) createHealthCheckCommand() error {
	content := `<?php

namespace App\Console\Commands;

use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Cache;
use Illuminate\Support\Facades\Http;

class HealthCheckCommand extends BaseCommand
{
    protected $signature = 'app:health-check {--notify : Send notification on failure}';
    protected $description = 'Run health checks on the application';

    public function handle(): int
    {
        $this->logInfo('Running health checks...');
        $failures = [];

        // Database check
        if (!$this->checkDatabase()) {
            $failures[] = 'Database connection failed';
        }

        // Cache check
        if (!$this->checkCache()) {
            $failures[] = 'Cache connection failed';
        }

        // Storage check
        if (!$this->checkStorage()) {
            $failures[] = 'Storage write failed';
        }

        if (count($failures) > 0) {
            foreach ($failures as $failure) {
                $this->logError($failure);
            }

            if ($this->option('notify')) {
                $this->sendFailureNotification($failures);
            }

            return self::FAILURE;
        }

        $this->logInfo('All health checks passed!');
        return self::SUCCESS;
    }

    protected function checkDatabase(): bool
    {
        try {
            DB::connection()->getPdo();
            $this->info('✓ Database: OK');
            return true;
        } catch (\Exception $e) {
            $this->error('✗ Database: FAILED');
            return false;
        }
    }

    protected function checkCache(): bool
    {
        try {
            Cache::put('health_check', 'ok', 10);
            $value = Cache::get('health_check');
            Cache::forget('health_check');
            
            if ($value === 'ok') {
                $this->info('✓ Cache: OK');
                return true;
            }
            throw new \Exception('Cache value mismatch');
        } catch (\Exception $e) {
            $this->error('✗ Cache: FAILED');
            return false;
        }
    }

    protected function checkStorage(): bool
    {
        try {
            $testFile = 'health_check_' . time() . '.txt';
            \Storage::disk('local')->put($testFile, 'test');
            \Storage::disk('local')->delete($testFile);
            $this->info('✓ Storage: OK');
            return true;
        } catch (\Exception $e) {
            $this->error('✗ Storage: FAILED');
            return false;
        }
    }

    protected function sendFailureNotification(array $failures): void
    {
        // Implement your notification logic here (Slack, email, etc.)
        $this->logWarning('Failure notification would be sent: ' . implode(', ', $failures));
    }
}
`
	path := filepath.Join(s.ProjectPath, "app/Console/Commands/HealthCheckCommand.php")
	return os.WriteFile(path, []byte(content), 0644)
}

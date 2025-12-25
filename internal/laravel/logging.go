package laravel

import (
	"fmt"
	"os"
	"path/filepath"
)

type LoggingSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewLoggingSetup(projectPath string, dryRun bool) *LoggingSetup {
	return &LoggingSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (l *LoggingSetup) Setup() error {
	if l.DryRun {
		fmt.Printf("[Dry Run] Would setup logging and debugging\n")
		return nil
	}

	fmt.Println("📝 Setting up logging and debugging...")

	if err := l.createRequestLogMiddleware(); err != nil {
		return err
	}
	if err := l.createLogService(); err != nil {
		return err
	}
	if err := l.createSlackLogHandler(); err != nil {
		return err
	}

	return nil
}

func (l *LoggingSetup) createRequestLogMiddleware() error {
	content := `<?php

namespace App\Http\Middleware;

use Closure;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Log;
use Symfony\Component\HttpFoundation\Response;

class LogRequests
{
    /**
     * Paths to exclude from logging.
     */
    protected array $except = [
        'health',
        'ready',
        '_debugbar/*',
    ];

    /**
     * Handle an incoming request.
     */
    public function handle(Request $request, Closure $next): Response
    {
        $startTime = microtime(true);

        $response = $next($request);

        if ($this->shouldLog($request)) {
            $this->logRequest($request, $response, $startTime);
        }

        return $response;
    }

    /**
     * Determine if the request should be logged.
     */
    protected function shouldLog(Request $request): bool
    {
        foreach ($this->except as $pattern) {
            if ($request->is($pattern)) {
                return false;
            }
        }

        return true;
    }

    /**
     * Log the request and response.
     */
    protected function logRequest(Request $request, Response $response, float $startTime): void
    {
        $duration = round((microtime(true) - $startTime) * 1000, 2);

        $logData = [
            'method' => $request->method(),
            'url' => $request->fullUrl(),
            'status' => $response->getStatusCode(),
            'duration_ms' => $duration,
            'ip' => $request->ip(),
            'user_agent' => $request->userAgent(),
            'user_id' => $request->user()?->id,
        ];

        // Log based on status code
        if ($response->getStatusCode() >= 500) {
            Log::channel('requests')->error('API Request', $logData);
        } elseif ($response->getStatusCode() >= 400) {
            Log::channel('requests')->warning('API Request', $logData);
        } else {
            Log::channel('requests')->info('API Request', $logData);
        }
    }
}
`
	dir := filepath.Join(l.ProjectPath, "app/Http/Middleware")
	os.MkdirAll(dir, 0755)
	return os.WriteFile(filepath.Join(dir, "LogRequests.php"), []byte(content), 0644)
}

func (l *LoggingSetup) createLogService() error {
	content := `<?php

namespace App\Services;

use Illuminate\Support\Facades\Log;

class LogService
{
    /**
     * Log an API error with context.
     */
    public static function apiError(string $message, array $context = [], ?\Throwable $exception = null): void
    {
        $data = array_merge($context, [
            'url' => request()->fullUrl(),
            'method' => request()->method(),
            'user_id' => auth()->id(),
            'ip' => request()->ip(),
        ]);

        if ($exception) {
            $data['exception'] = [
                'message' => $exception->getMessage(),
                'file' => $exception->getFile(),
                'line' => $exception->getLine(),
            ];
        }

        Log::error($message, $data);
    }

    /**
     * Log an API info message.
     */
    public static function apiInfo(string $message, array $context = []): void
    {
        $data = array_merge($context, [
            'user_id' => auth()->id(),
        ]);

        Log::info($message, $data);
    }

    /**
     * Log a model action.
     */
    public static function modelAction(string $action, $model, array $context = []): void
    {
        $data = array_merge($context, [
            'action' => $action,
            'model' => get_class($model),
            'model_id' => $model->getKey(),
            'user_id' => auth()->id(),
        ]);

        Log::info("Model {$action}", $data);
    }

    /**
     * Log a performance metric.
     */
    public static function performance(string $operation, float $durationMs, array $context = []): void
    {
        $data = array_merge($context, [
            'operation' => $operation,
            'duration_ms' => $durationMs,
        ]);

        if ($durationMs > 1000) {
            Log::warning('Slow operation detected', $data);
        } else {
            Log::info('Performance metric', $data);
        }
    }

    /**
     * Log with timing.
     */
    public static function timed(string $operation, callable $callback): mixed
    {
        $start = microtime(true);
        
        try {
            $result = $callback();
            $duration = (microtime(true) - $start) * 1000;
            self::performance($operation, $duration, ['status' => 'success']);
            return $result;
        } catch (\Exception $e) {
            $duration = (microtime(true) - $start) * 1000;
            self::performance($operation, $duration, ['status' => 'failed', 'error' => $e->getMessage()]);
            throw $e;
        }
    }

    /**
     * Log a security event.
     */
    public static function security(string $event, array $context = []): void
    {
        $data = array_merge($context, [
            'event' => $event,
            'ip' => request()->ip(),
            'user_agent' => request()->userAgent(),
            'user_id' => auth()->id(),
        ]);

        Log::channel('security')->warning('Security Event', $data);
    }
}
`
	dir := filepath.Join(l.ProjectPath, "app/Services")
	os.MkdirAll(dir, 0755)
	return os.WriteFile(filepath.Join(dir, "LogService.php"), []byte(content), 0644)
}

func (l *LoggingSetup) createSlackLogHandler() error {
	content := `<?php

namespace App\Logging;

use Monolog\Handler\SlackWebhookHandler;
use Monolog\Logger;

class SlackLogHandler
{
    /**
     * Create a custom Monolog instance.
     */
    public function __invoke(array $config): Logger
    {
        $logger = new Logger('slack');

        $webhookUrl = $config['url'] ?? env('LOG_SLACK_WEBHOOK_URL');
        $channel = $config['channel'] ?? null;
        $username = $config['username'] ?? 'Laravel Logger';
        $emoji = $config['emoji'] ?? ':boom:';
        $level = $config['level'] ?? Logger::ERROR;

        if ($webhookUrl) {
            $handler = new SlackWebhookHandler(
                $webhookUrl,
                $channel,
                $username,
                true, // useAttachment
                $emoji,
                false, // useShortAttachment
                true, // includeContextAndExtra
                $level
            );

            $logger->pushHandler($handler);
        }

        return $logger;
    }
}

/*
|--------------------------------------------------------------------------
| Add to config/logging.php channels array:
|--------------------------------------------------------------------------
|
| 'slack' => [
|     'driver' => 'custom',
|     'via' => App\Logging\SlackLogHandler::class,
|     'url' => env('LOG_SLACK_WEBHOOK_URL'),
|     'username' => 'Laravel Logger',
|     'emoji' => ':boom:',
|     'level' => 'error',
| ],
|
| 'requests' => [
|     'driver' => 'daily',
|     'path' => storage_path('logs/requests.log'),
|     'level' => 'debug',
|     'days' => 14,
| ],
|
| 'security' => [
|     'driver' => 'daily',
|     'path' => storage_path('logs/security.log'),
|     'level' => 'warning',
|     'days' => 30,
| ],
|
*/
`
	dir := filepath.Join(l.ProjectPath, "app/Logging")
	os.MkdirAll(dir, 0755)
	return os.WriteFile(filepath.Join(dir, "SlackLogHandler.php"), []byte(content), 0644)
}

package laravel

import (
	"fmt"
	"os"
	"path/filepath"
)

type EventsSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewEventsSetup(projectPath string, dryRun bool) *EventsSetup {
	return &EventsSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (e *EventsSetup) Setup() error {
	if e.DryRun {
		fmt.Printf("[Dry Run] Would setup events and listeners\n")
		return nil
	}

	fmt.Println("📡 Setting up events and listeners...")

	eventsDir := filepath.Join(e.ProjectPath, "app/Events")
	listenersDir := filepath.Join(e.ProjectPath, "app/Listeners")
	os.MkdirAll(eventsDir, 0755)
	os.MkdirAll(listenersDir, 0755)

	if err := e.createBaseEvent(); err != nil {
		return err
	}
	if err := e.createBaseListener(); err != nil {
		return err
	}
	if err := e.createUserRegisteredEvent(); err != nil {
		return err
	}
	if err := e.createSendWelcomeEmailListener(); err != nil {
		return err
	}

	return nil
}

func (e *EventsSetup) createBaseEvent() error {
	content := `<?php

namespace App\Events;

use Illuminate\Broadcasting\InteractsWithSockets;
use Illuminate\Foundation\Events\Dispatchable;
use Illuminate\Queue\SerializesModels;

abstract class BaseEvent
{
    use Dispatchable, InteractsWithSockets, SerializesModels;

    /**
     * The time the event was created.
     */
    public \DateTimeInterface $createdAt;

    public function __construct()
    {
        $this->createdAt = now();
    }

    /**
     * Get the event name for logging.
     */
    public function getEventName(): string
    {
        return class_basename(static::class);
    }

    /**
     * Get the event data for logging.
     */
    public function getEventData(): array
    {
        return [];
    }
}
`
	path := filepath.Join(e.ProjectPath, "app/Events/BaseEvent.php")
	return os.WriteFile(path, []byte(content), 0644)
}

func (e *EventsSetup) createBaseListener() error {
	content := `<?php

namespace App\Listeners;

use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Queue\InteractsWithQueue;
use Illuminate\Support\Facades\Log;

abstract class BaseListener implements ShouldQueue
{
    use InteractsWithQueue;

    /**
     * The number of times the job may be attempted.
     */
    public int $tries = 3;

    /**
     * The number of seconds to wait before retrying.
     */
    public int $backoff = 60;

    /**
     * Handle a job failure.
     */
    public function failed($event, \Throwable $exception): void
    {
        Log::error('Listener failed: ' . static::class, [
            'event' => get_class($event),
            'exception' => $exception->getMessage(),
            'trace' => $exception->getTraceAsString(),
        ]);
    }

    /**
     * Log event processing.
     */
    protected function logProcessing($event): void
    {
        Log::info('Processing event', [
            'listener' => static::class,
            'event' => get_class($event),
        ]);
    }
}
`
	path := filepath.Join(e.ProjectPath, "app/Listeners/BaseListener.php")
	return os.WriteFile(path, []byte(content), 0644)
}

func (e *EventsSetup) createUserRegisteredEvent() error {
	content := `<?php

namespace App\Events;

use App\Models\User;
use Illuminate\Broadcasting\Channel;
use Illuminate\Broadcasting\PrivateChannel;
use Illuminate\Contracts\Broadcasting\ShouldBroadcast;

class UserRegistered extends BaseEvent
{
    public User $user;

    public function __construct(User $user)
    {
        parent::__construct();
        $this->user = $user;
    }

    /**
     * Get the channels the event should broadcast on.
     */
    public function broadcastOn(): array
    {
        return [
            new PrivateChannel('users.' . $this->user->id),
        ];
    }

    /**
     * Get the event data for logging.
     */
    public function getEventData(): array
    {
        return [
            'user_id' => $this->user->id,
            'email' => $this->user->email,
        ];
    }
}
`
	path := filepath.Join(e.ProjectPath, "app/Events/UserRegistered.php")
	return os.WriteFile(path, []byte(content), 0644)
}

func (e *EventsSetup) createSendWelcomeEmailListener() error {
	content := `<?php

namespace App\Listeners;

use App\Events\UserRegistered;
use App\Notifications\WelcomeNotification;
use Illuminate\Support\Facades\Log;

class SendWelcomeEmail extends BaseListener
{
    /**
     * Handle the event.
     */
    public function handle(UserRegistered $event): void
    {
        $this->logProcessing($event);

        $user = $event->user;

        try {
            $user->notify(new WelcomeNotification($user->name));
            
            Log::info('Welcome email sent', [
                'user_id' => $user->id,
                'email' => $user->email,
            ]);
        } catch (\Exception $e) {
            Log::error('Failed to send welcome email', [
                'user_id' => $user->id,
                'error' => $e->getMessage(),
            ]);
            
            throw $e; // Re-throw to trigger retry
        }
    }
}
`
	path := filepath.Join(e.ProjectPath, "app/Listeners/SendWelcomeEmail.php")
	return os.WriteFile(path, []byte(content), 0644)
}

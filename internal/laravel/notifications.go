package laravel

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type NotificationsSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewNotificationsSetup(projectPath string, dryRun bool) *NotificationsSetup {
	return &NotificationsSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (n *NotificationsSetup) Setup() error {
	if n.DryRun {
		fmt.Printf("[Dry Run] Would setup notifications system\n")
		return nil
	}

	fmt.Println("🔔 Setting up notifications system...")

	// Create notifications table
	cmd := exec.Command("php", "artisan", "notifications:table")
	cmd.Dir = n.ProjectPath
	_ = cmd.Run()

	// Create notifications directory
	notifDir := filepath.Join(n.ProjectPath, "app/Notifications")
	os.MkdirAll(notifDir, 0755)

	if err := n.createBaseNotification(); err != nil {
		return err
	}
	if err := n.createWelcomeNotification(); err != nil {
		return err
	}
	if err := n.createNotificationService(); err != nil {
		return err
	}

	return nil
}

func (n *NotificationsSetup) createBaseNotification() error {
	content := `<?php

namespace App\Notifications;

use Illuminate\Bus\Queueable;
use Illuminate\Contracts\Queue\ShouldQueue;
use Illuminate\Notifications\Messages\MailMessage;
use Illuminate\Notifications\Notification;

abstract class BaseNotification extends Notification implements ShouldQueue
{
    use Queueable;

    /**
     * Get the notification's delivery channels.
     */
    public function via(object $notifiable): array
    {
        return ['mail', 'database'];
    }

    /**
     * Get the mail representation of the notification.
     */
    abstract public function toMail(object $notifiable): MailMessage;

    /**
     * Get the array representation of the notification (for database).
     */
    abstract public function toArray(object $notifiable): array;

    /**
     * Get notification data for broadcasting.
     */
    public function toBroadcast(object $notifiable): array
    {
        return $this->toArray($notifiable);
    }
}
`
	path := filepath.Join(n.ProjectPath, "app/Notifications/BaseNotification.php")
	return os.WriteFile(path, []byte(content), 0644)
}

func (n *NotificationsSetup) createWelcomeNotification() error {
	content := `<?php

namespace App\Notifications;

use Illuminate\Notifications\Messages\MailMessage;

class WelcomeNotification extends BaseNotification
{
    protected string $userName;

    public function __construct(string $userName)
    {
        $this->userName = $userName;
    }

    public function toMail(object $notifiable): MailMessage
    {
        return (new MailMessage)
            ->subject('Welcome to ' . config('app.name'))
            ->greeting("Hello {$this->userName}!")
            ->line('Thank you for joining our platform.')
            ->line('We are excited to have you on board.')
            ->action('Get Started', url('/'))
            ->line('If you have any questions, feel free to reach out.');
    }

    public function toArray(object $notifiable): array
    {
        return [
            'type' => 'welcome',
            'title' => 'Welcome!',
            'message' => "Welcome to our platform, {$this->userName}!",
            'action_url' => url('/'),
        ];
    }
}
`
	path := filepath.Join(n.ProjectPath, "app/Notifications/WelcomeNotification.php")
	return os.WriteFile(path, []byte(content), 0644)
}

func (n *NotificationsSetup) createNotificationService() error {
	content := `<?php

namespace App\Services;

use App\Models\User;
use Illuminate\Notifications\Notification;
use Illuminate\Support\Facades\Notification as NotificationFacade;

class NotificationService
{
    /**
     * Send notification to a single user.
     */
    public static function sendToUser(User $user, Notification $notification): void
    {
        $user->notify($notification);
    }

    /**
     * Send notification to multiple users.
     */
    public static function sendToUsers($users, Notification $notification): void
    {
        NotificationFacade::send($users, $notification);
    }

    /**
     * Send notification to all users.
     */
    public static function broadcast(Notification $notification): void
    {
        $users = User::all();
        NotificationFacade::send($users, $notification);
    }

    /**
     * Mark all notifications as read for a user.
     */
    public static function markAllAsRead(User $user): void
    {
        $user->unreadNotifications->markAsRead();
    }

    /**
     * Get unread notifications for a user.
     */
    public static function getUnread(User $user, int $limit = 10)
    {
        return $user->unreadNotifications()->take($limit)->get();
    }

    /**
     * Get all notifications for a user with pagination.
     */
    public static function getPaginated(User $user, int $perPage = 15)
    {
        return $user->notifications()->paginate($perPage);
    }

    /**
     * Delete old notifications.
     */
    public static function deleteOld(int $days = 30): int
    {
        return \DB::table('notifications')
            ->where('created_at', '<', now()->subDays($days))
            ->delete();
    }
}
`
	dir := filepath.Join(n.ProjectPath, "app/Services")
	os.MkdirAll(dir, 0755)
	return os.WriteFile(filepath.Join(dir, "NotificationService.php"), []byte(content), 0644)
}

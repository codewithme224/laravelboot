package laravel

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type MediaSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewMediaSetup(projectPath string, dryRun bool) *MediaSetup {
	return &MediaSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (m *MediaSetup) Setup() error {
	if m.DryRun {
		fmt.Printf("[Dry Run] Would install spatie/laravel-medialibrary with config and service\n")
		return nil
	}

	fmt.Println("🖼️ Installing spatie/laravel-medialibrary...")
	cmd := exec.Command("composer", "require", "spatie/laravel-medialibrary", "--with-all-dependencies")
	cmd.Dir = m.ProjectPath
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to install spatie/laravel-medialibrary: %v\nOutput: %s", err, string(output))
	}

	fmt.Println("📦 Publishing migrations...")
	cmd = exec.Command("php", "artisan", "vendor:publish", "--provider=Spatie\\MediaLibrary\\MediaLibraryServiceProvider", "--tag=medialibrary-migrations")
	cmd.Dir = m.ProjectPath
	_ = cmd.Run()

	fmt.Println("⚙️ Publishing config...")
	cmd = exec.Command("php", "artisan", "vendor:publish", "--provider=Spatie\\MediaLibrary\\MediaLibraryServiceProvider", "--tag=medialibrary-config")
	cmd.Dir = m.ProjectPath
	_ = cmd.Run()

	fmt.Println("🔧 Creating SpatieMediaService...")
	if err := m.createMediaService(); err != nil {
		return err
	}

	fmt.Println("🔧 Creating HasMedia trait...")
	if err := m.createHasMediaTrait(); err != nil {
		return err
	}

	return nil
}

func (m *MediaSetup) createMediaService() error {
	content := `<?php

namespace App\Services;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Support\Facades\File;
use Illuminate\Http\UploadedFile;

class SpatieMediaService
{
    /**
     * Add multiple files from base64 encoded strings.
     */
    public static function addMultipleFilesBase64(array $base64Files, Model $model, string $collectionName, string $disk = 'public'): void
    {
        if (!is_array($base64Files) || empty($base64Files)) {
            return;
        }

        foreach ($base64Files as $data) {
            $media = $model->addMediaFromBase64($data)->toMediaCollection($collectionName, $disk);
            $extension = File::guessExtension($media->getPath());
            $media->file_name = "{$media->file_name}.{$extension}";
            $media->save();
        }
    }

    /**
     * Add a single file from upload.
     */
    public static function addFile(UploadedFile $file, Model $model, string $collectionName, string $disk = 'public'): void
    {
        $model->addMedia($file)
            ->sanitizingFileName(fn($fileName) => strtolower(str_replace(['#', '/', '\\', ' '], '-', $fileName)))
            ->toMediaCollection($collectionName, $disk);
    }

    /**
     * Add a video file.
     */
    public static function addVideo(UploadedFile $video, Model $model, string $collectionName, string $disk = 'public'): void
    {
        $media = $model->addMedia($video)
            ->sanitizingFileName(fn($fileName) => strtolower(str_replace(['#', '/', '\\', ' '], '-', $fileName)))
            ->toMediaCollection($collectionName, $disk);

        $extension = $video->getClientOriginalExtension();
        if (empty($extension)) {
            $extension = File::guessExtension($media->getPath());
        }

        if (!str_ends_with($media->file_name, '.' . $extension)) {
            $media->file_name = "{$media->file_name}.{$extension}";
        }

        $media->save();
    }

    /**
     * Remove files from a collection.
     */
    public static function removeFiles(Model $model, ?string $collectionName = null): void
    {
        if ($collectionName) {
            $model->clearMediaCollection($collectionName);
        } else {
            $model->clearMediaCollections();
        }
    }

    /**
     * Upload and replace base64 files (removes existing first).
     */
    public static function uploadAndRemoveBase64Files(array $base64Files, Model $model, string $collectionName, string $disk = 'public'): void
    {
        self::removeFiles($model, $collectionName);
        self::addMultipleFilesBase64($base64Files, $model, $collectionName, $disk);
    }

    /**
     * Upload and replace video (removes existing first).
     */
    public static function uploadAndRemoveVideo(UploadedFile $video, Model $model, string $collectionName, string $disk = 'public'): void
    {
        self::removeFiles($model, $collectionName);
        self::addVideo($video, $model, $collectionName, $disk);
    }
}
`
	dir := filepath.Join(m.ProjectPath, "app/Services")
	os.MkdirAll(dir, 0755)
	return os.WriteFile(filepath.Join(dir, "SpatieMediaService.php"), []byte(content), 0644)
}

func (m *MediaSetup) createHasMediaTrait() error {
	content := `<?php

namespace App\Traits;

use Spatie\MediaLibrary\InteractsWithMedia;
use Spatie\MediaLibrary\MediaCollections\Models\Media;

trait HasMedia
{
    use InteractsWithMedia;

    /**
     * Register media conversions.
     */
    public function registerMediaConversions(?Media $media = null): void
    {
        $this->addMediaConversion('thumb')
            ->width(150)
            ->height(150)
            ->sharpen(10);

        $this->addMediaConversion('preview')
            ->width(400)
            ->height(400);
    }

    /**
     * Get the first media URL for a collection.
     */
    public function getMediaUrl(string $collection = 'default', string $conversion = ''): ?string
    {
        $media = $this->getFirstMedia($collection);
        
        if (!$media) {
            return null;
        }

        return $conversion ? $media->getUrl($conversion) : $media->getUrl();
    }

    /**
     * Get all media URLs for a collection.
     */
    public function getMediaUrls(string $collection = 'default', string $conversion = ''): array
    {
        return $this->getMedia($collection)->map(function ($media) use ($conversion) {
            return $conversion ? $media->getUrl($conversion) : $media->getUrl();
        })->toArray();
    }
}
`
	dir := filepath.Join(m.ProjectPath, "app/Traits")
	os.MkdirAll(dir, 0755)
	return os.WriteFile(filepath.Join(dir, "HasMedia.php"), []byte(content), 0644)
}

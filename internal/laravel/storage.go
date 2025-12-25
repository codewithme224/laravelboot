package laravel

import (
	"fmt"
	"os"
	"path/filepath"
)

type StorageSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewStorageSetup(projectPath string, dryRun bool) *StorageSetup {
	return &StorageSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (s *StorageSetup) Setup() error {
	if s.DryRun {
		fmt.Printf("[Dry Run] Would setup file storage service\n")
		return nil
	}

	fmt.Println("📁 Setting up file storage service...")

	if err := s.createFileService(); err != nil {
		return err
	}
	if err := s.createFileController(); err != nil {
		return err
	}

	return nil
}

func (s *StorageSetup) createFileService() error {
	content := `<?php

namespace App\Services;

use Illuminate\Http\UploadedFile;
use Illuminate\Support\Facades\Storage;
use Illuminate\Support\Str;

class FileService
{
    protected string $disk;

    public function __construct(string $disk = 'public')
    {
        $this->disk = $disk;
    }

    /**
     * Upload a file.
     */
    public function upload(UploadedFile $file, string $directory = 'uploads', ?string $filename = null): string
    {
        $filename = $filename ?? $this->generateFilename($file);
        $path = $file->storeAs($directory, $filename, $this->disk);
        
        return $path;
    }

    /**
     * Upload multiple files.
     */
    public function uploadMultiple(array $files, string $directory = 'uploads'): array
    {
        $paths = [];
        
        foreach ($files as $file) {
            if ($file instanceof UploadedFile) {
                $paths[] = $this->upload($file, $directory);
            }
        }
        
        return $paths;
    }

    /**
     * Delete a file.
     */
    public function delete(string $path): bool
    {
        return Storage::disk($this->disk)->delete($path);
    }

    /**
     * Delete multiple files.
     */
    public function deleteMultiple(array $paths): bool
    {
        return Storage::disk($this->disk)->delete($paths);
    }

    /**
     * Check if a file exists.
     */
    public function exists(string $path): bool
    {
        return Storage::disk($this->disk)->exists($path);
    }

    /**
     * Get the full URL of a file.
     */
    public function url(string $path): string
    {
        return Storage::disk($this->disk)->url($path);
    }

    /**
     * Get the file size in bytes.
     */
    public function size(string $path): int
    {
        return Storage::disk($this->disk)->size($path);
    }

    /**
     * Get the file's last modification time.
     */
    public function lastModified(string $path): int
    {
        return Storage::disk($this->disk)->lastModified($path);
    }

    /**
     * Copy a file.
     */
    public function copy(string $from, string $to): bool
    {
        return Storage::disk($this->disk)->copy($from, $to);
    }

    /**
     * Move a file.
     */
    public function move(string $from, string $to): bool
    {
        return Storage::disk($this->disk)->move($from, $to);
    }

    /**
     * Get file contents.
     */
    public function get(string $path): ?string
    {
        return Storage::disk($this->disk)->get($path);
    }

    /**
     * Put contents into a file.
     */
    public function put(string $path, string $contents): bool
    {
        return Storage::disk($this->disk)->put($path, $contents);
    }

    /**
     * Generate a unique filename.
     */
    protected function generateFilename(UploadedFile $file): string
    {
        $extension = $file->getClientOriginalExtension();
        return Str::uuid() . '.' . $extension;
    }

    /**
     * Upload a base64 encoded file.
     */
    public function uploadBase64(string $base64, string $directory = 'uploads', ?string $extension = null): ?string
    {
        if (preg_match('/^data:(\w+\/\w+);base64,/', $base64, $matches)) {
            $mimeType = $matches[1];
            $base64 = preg_replace('/^data:\w+\/\w+;base64,/', '', $base64);
            
            if (!$extension) {
                $extension = $this->mimeToExtension($mimeType);
            }
        }

        $contents = base64_decode($base64);
        if ($contents === false) {
            return null;
        }

        $filename = Str::uuid() . '.' . ($extension ?? 'bin');
        $path = $directory . '/' . $filename;
        
        if (Storage::disk($this->disk)->put($path, $contents)) {
            return $path;
        }

        return null;
    }

    /**
     * Convert MIME type to file extension.
     */
    protected function mimeToExtension(string $mimeType): string
    {
        $map = [
            'image/jpeg' => 'jpg',
            'image/png' => 'png',
            'image/gif' => 'gif',
            'image/webp' => 'webp',
            'application/pdf' => 'pdf',
            'text/plain' => 'txt',
            'application/json' => 'json',
        ];

        return $map[$mimeType] ?? 'bin';
    }

    /**
     * Get a temporary URL (for S3/cloud storage).
     */
    public function temporaryUrl(string $path, int $minutes = 60): string
    {
        return Storage::disk($this->disk)->temporaryUrl($path, now()->addMinutes($minutes));
    }
}
`
	dir := filepath.Join(s.ProjectPath, "app/Services")
	os.MkdirAll(dir, 0755)
	return os.WriteFile(filepath.Join(dir, "FileService.php"), []byte(content), 0644)
}

func (s *StorageSetup) createFileController() error {
	content := `<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Services\FileService;
use App\Traits\ApiResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Storage;

class FileController extends Controller
{
    use ApiResponse;

    protected FileService $fileService;

    public function __construct(FileService $fileService)
    {
        $this->fileService = $fileService;
    }

    /**
     * Upload a file.
     */
    public function upload(Request $request)
    {
        $request->validate([
            'file' => 'required|file|max:10240', // 10MB max
            'directory' => 'nullable|string|max:255',
        ]);

        $directory = $request->input('directory', 'uploads');
        $path = $this->fileService->upload($request->file('file'), $directory);

        return $this->success([
            'path' => $path,
            'url' => $this->fileService->url($path),
        ], 'File uploaded successfully');
    }

    /**
     * Upload multiple files.
     */
    public function uploadMultiple(Request $request)
    {
        $request->validate([
            'files' => 'required|array',
            'files.*' => 'file|max:10240',
            'directory' => 'nullable|string|max:255',
        ]);

        $directory = $request->input('directory', 'uploads');
        $paths = $this->fileService->uploadMultiple($request->file('files'), $directory);

        $results = array_map(fn($path) => [
            'path' => $path,
            'url' => $this->fileService->url($path),
        ], $paths);

        return $this->success($results, 'Files uploaded successfully');
    }

    /**
     * Delete a file.
     */
    public function destroy(Request $request)
    {
        $request->validate([
            'path' => 'required|string',
        ]);

        if ($this->fileService->delete($request->input('path'))) {
            return $this->success(null, 'File deleted successfully');
        }

        return $this->error('Failed to delete file', 400);
    }

    /**
     * Download a file.
     */
    public function download(Request $request)
    {
        $request->validate([
            'path' => 'required|string',
        ]);

        $path = $request->input('path');

        if (!$this->fileService->exists($path)) {
            return $this->notFound('File not found');
        }

        return Storage::disk('public')->download($path);
    }
}
`
	dir := filepath.Join(s.ProjectPath, "app/Http/Controllers/Api")
	os.MkdirAll(dir, 0755)
	return os.WriteFile(filepath.Join(dir, "FileController.php"), []byte(content), 0644)
}

package laravel

import (
	"fmt"
	"os"
	"path/filepath"
)

type SoftDeletesSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewSoftDeletesSetup(projectPath string, dryRun bool) *SoftDeletesSetup {
	return &SoftDeletesSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (s *SoftDeletesSetup) Setup() error {
	if s.DryRun {
		fmt.Printf("[Dry Run] Would setup soft deletes and trash management\n")
		return nil
	}

	fmt.Println("🗑️ Setting up soft deletes and trash management...")

	if err := s.createHasSoftDeletesTrait(); err != nil {
		return err
	}
	if err := s.createTrashService(); err != nil {
		return err
	}

	return nil
}

func (s *SoftDeletesSetup) createHasSoftDeletesTrait() error {
	content := `<?php

namespace App\Traits;

use Illuminate\Database\Eloquent\SoftDeletes;

trait HasSoftDeletes
{
    use SoftDeletes;

    /**
     * Boot the trait.
     */
    public static function bootHasSoftDeletes(): void
    {
        static::deleting(function ($model) {
            if (method_exists($model, 'beforeSoftDelete')) {
                $model->beforeSoftDelete();
            }
        });

        static::restoring(function ($model) {
            if (method_exists($model, 'beforeRestore')) {
                $model->beforeRestore();
            }
        });

        static::restored(function ($model) {
            if (method_exists($model, 'afterRestore')) {
                $model->afterRestore();
            }
        });
    }

    /**
     * Scope to get only trashed records.
     */
    public function scopeTrashed($query)
    {
        return $query->onlyTrashed();
    }

    /**
     * Restore the model.
     */
    public function restoreModel(): bool
    {
        return $this->restore();
    }

    /**
     * Force delete the model permanently.
     */
    public function forceDeleteModel(): bool
    {
        return $this->forceDelete();
    }

    /**
     * Check if the model is trashed.
     */
    public function isTrashed(): bool
    {
        return $this->trashed();
    }

    /**
     * Get the deleted by user (if tracking).
     */
    public function deletedBy()
    {
        if (!$this->deleted_by) {
            return null;
        }
        return \App\Models\User::find($this->deleted_by);
    }
}
`
	dir := filepath.Join(s.ProjectPath, "app/Traits")
	os.MkdirAll(dir, 0755)
	return os.WriteFile(filepath.Join(dir, "HasSoftDeletes.php"), []byte(content), 0644)
}

func (s *SoftDeletesSetup) createTrashService() error {
	content := `<?php

namespace App\Services;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Support\Collection;

class TrashService
{
    /**
     * Get all trashed records for a model.
     */
    public static function getTrashed(string $modelClass, int $perPage = 15)
    {
        return $modelClass::onlyTrashed()->paginate($perPage);
    }

    /**
     * Restore a trashed record.
     */
    public static function restore(string $modelClass, int|string $id): bool
    {
        $model = $modelClass::onlyTrashed()->findOrFail($id);
        return $model->restore();
    }

    /**
     * Restore multiple trashed records.
     */
    public static function restoreMany(string $modelClass, array $ids): int
    {
        return $modelClass::onlyTrashed()
            ->whereIn('id', $ids)
            ->restore();
    }

    /**
     * Force delete a trashed record permanently.
     */
    public static function forceDelete(string $modelClass, int|string $id): bool
    {
        $model = $modelClass::onlyTrashed()->findOrFail($id);
        return $model->forceDelete();
    }

    /**
     * Force delete multiple trashed records.
     */
    public static function forceDeleteMany(string $modelClass, array $ids): int
    {
        $models = $modelClass::onlyTrashed()->whereIn('id', $ids)->get();
        $count = 0;
        
        foreach ($models as $model) {
            if ($model->forceDelete()) {
                $count++;
            }
        }
        
        return $count;
    }

    /**
     * Empty trash (force delete all trashed records).
     */
    public static function emptyTrash(string $modelClass): int
    {
        $models = $modelClass::onlyTrashed()->get();
        $count = 0;
        
        foreach ($models as $model) {
            if ($model->forceDelete()) {
                $count++;
            }
        }
        
        return $count;
    }

    /**
     * Restore all trashed records.
     */
    public static function restoreAll(string $modelClass): int
    {
        return $modelClass::onlyTrashed()->restore();
    }

    /**
     * Get trash count for a model.
     */
    public static function count(string $modelClass): int
    {
        return $modelClass::onlyTrashed()->count();
    }

    /**
     * Auto-delete old trashed records.
     */
    public static function autoClean(string $modelClass, int $daysOld = 30): int
    {
        $models = $modelClass::onlyTrashed()
            ->where('deleted_at', '<', now()->subDays($daysOld))
            ->get();
        
        $count = 0;
        foreach ($models as $model) {
            if ($model->forceDelete()) {
                $count++;
            }
        }
        
        return $count;
    }
}
`
	dir := filepath.Join(s.ProjectPath, "app/Services")
	os.MkdirAll(dir, 0755)
	return os.WriteFile(filepath.Join(dir, "TrashService.php"), []byte(content), 0644)
}

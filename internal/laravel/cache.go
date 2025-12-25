package laravel

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type CacheSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewCacheSetup(projectPath string, dryRun bool) *CacheSetup {
	return &CacheSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (c *CacheSetup) Setup() error {
	if c.DryRun {
		fmt.Printf("[Dry Run] Would setup caching layer\n")
		return nil
	}

	fmt.Println("💾 Setting up caching layer...")

	// Install predis for Redis support
	fmt.Println("📦 Installing predis/predis...")
	cmd := exec.Command("composer", "require", "predis/predis", "--with-all-dependencies")
	cmd.Dir = c.ProjectPath
	_ = cmd.Run()

	if err := c.createCacheService(); err != nil {
		return err
	}
	if err := c.createCacheableTrait(); err != nil {
		return err
	}

	return nil
}

func (c *CacheSetup) createCacheService() error {
	content := `<?php

namespace App\Services;

use Illuminate\Support\Facades\Cache;
use Closure;

class CacheService
{
    /**
     * Default cache TTL in seconds (1 hour).
     */
    protected static int $defaultTtl = 3600;

    /**
     * Get or set a cached value.
     */
    public static function remember(string $key, Closure $callback, ?int $ttl = null): mixed
    {
        return Cache::remember($key, $ttl ?? self::$defaultTtl, $callback);
    }

    /**
     * Get or set a cached value forever.
     */
    public static function rememberForever(string $key, Closure $callback): mixed
    {
        return Cache::rememberForever($key, $callback);
    }

    /**
     * Get a cached value with tags.
     */
    public static function taggedRemember(array $tags, string $key, Closure $callback, ?int $ttl = null): mixed
    {
        return Cache::tags($tags)->remember($key, $ttl ?? self::$defaultTtl, $callback);
    }

    /**
     * Flush cache by tags.
     */
    public static function flushTags(array $tags): void
    {
        Cache::tags($tags)->flush();
    }

    /**
     * Flush a specific key.
     */
    public static function forget(string $key): bool
    {
        return Cache::forget($key);
    }

    /**
     * Check if a key exists in cache.
     */
    public static function has(string $key): bool
    {
        return Cache::has($key);
    }

    /**
     * Get a value from cache.
     */
    public static function get(string $key, mixed $default = null): mixed
    {
        return Cache::get($key, $default);
    }

    /**
     * Put a value in cache.
     */
    public static function put(string $key, mixed $value, ?int $ttl = null): bool
    {
        return Cache::put($key, $value, $ttl ?? self::$defaultTtl);
    }

    /**
     * Generate a cache key from multiple parts.
     */
    public static function key(string ...$parts): string
    {
        return implode(':', $parts);
    }

    /**
     * Generate a model-specific cache key.
     */
    public static function modelKey(string $model, int|string $id, ?string $suffix = null): string
    {
        $key = strtolower(class_basename($model)) . ':' . $id;
        return $suffix ? $key . ':' . $suffix : $key;
    }

    /**
     * Clear all cache.
     */
    public static function flush(): bool
    {
        return Cache::flush();
    }
}
`
	dir := filepath.Join(c.ProjectPath, "app/Services")
	os.MkdirAll(dir, 0755)
	return os.WriteFile(filepath.Join(dir, "CacheService.php"), []byte(content), 0644)
}

func (c *CacheSetup) createCacheableTrait() error {
	content := `<?php

namespace App\Traits;

use Illuminate\Support\Facades\Cache;

trait Cacheable
{
    /**
     * Cache TTL in seconds.
     */
    protected static int $cacheTtl = 3600;

    /**
     * Boot the cacheable trait.
     */
    public static function bootCacheable(): void
    {
        static::saved(function ($model) {
            $model->flushCache();
        });

        static::deleted(function ($model) {
            $model->flushCache();
        });
    }

    /**
     * Get the cache key for this model.
     */
    public function getCacheKey(?string $suffix = null): string
    {
        $key = strtolower(class_basename($this)) . ':' . $this->getKey();
        return $suffix ? $key . ':' . $suffix : $key;
    }

    /**
     * Get the cache tags for this model.
     */
    public function getCacheTags(): array
    {
        return [strtolower(class_basename($this))];
    }

    /**
     * Cache a value for this model.
     */
    public function cache(string $key, mixed $value, ?int $ttl = null): bool
    {
        return Cache::tags($this->getCacheTags())
            ->put($this->getCacheKey($key), $value, $ttl ?? static::$cacheTtl);
    }

    /**
     * Get a cached value for this model.
     */
    public function cached(string $key, mixed $default = null): mixed
    {
        return Cache::tags($this->getCacheTags())
            ->get($this->getCacheKey($key), $default);
    }

    /**
     * Flush the cache for this model.
     */
    public function flushCache(): void
    {
        Cache::tags($this->getCacheTags())->flush();
    }

    /**
     * Remember a value in cache.
     */
    public function rememberCached(string $key, \Closure $callback, ?int $ttl = null): mixed
    {
        return Cache::tags($this->getCacheTags())
            ->remember($this->getCacheKey($key), $ttl ?? static::$cacheTtl, $callback);
    }

    /**
     * Find a model by ID with caching.
     */
    public static function findCached(int|string $id): ?static
    {
        $key = strtolower(class_basename(static::class)) . ':' . $id;
        
        return Cache::tags([strtolower(class_basename(static::class))])
            ->remember($key, static::$cacheTtl, function () use ($id) {
                return static::find($id);
            });
    }
}
`
	dir := filepath.Join(c.ProjectPath, "app/Traits")
	os.MkdirAll(dir, 0755)
	return os.WriteFile(filepath.Join(dir, "Cacheable.php"), []byte(content), 0644)
}

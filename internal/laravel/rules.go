package laravel

import (
	"fmt"
	"os"
	"path/filepath"
)

type RulesSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewRulesSetup(projectPath string, dryRun bool) *RulesSetup {
	return &RulesSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (r *RulesSetup) Setup() error {
	if r.DryRun {
		fmt.Printf("[Dry Run] Would create custom validation rules\n")
		return nil
	}

	fmt.Println("📏 Creating custom validation rules...")

	rulesDir := filepath.Join(r.ProjectPath, "app/Rules")
	os.MkdirAll(rulesDir, 0755)

	if err := r.createBase64ImageRule(); err != nil {
		return err
	}
	if err := r.createPhoneNumberRule(); err != nil {
		return err
	}
	if err := r.createScopedUniqueRule(); err != nil {
		return err
	}
	if err := r.createTimeFormatRule(); err != nil {
		return err
	}
	if err := r.createStrongPasswordRule(); err != nil {
		return err
	}

	return nil
}

func (r *RulesSetup) createBase64ImageRule() error {
	content := `<?php

namespace App\Rules;

use Closure;
use Illuminate\Contracts\Validation\ValidationRule;

class Base64Image implements ValidationRule
{
    protected array $allowedMimes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];
    protected int $maxSizeKb;

    public function __construct(int $maxSizeKb = 5120)
    {
        $this->maxSizeKb = $maxSizeKb;
    }

    public function validate(string $attribute, mixed $value, Closure $fail): void
    {
        if (!is_string($value)) {
            $fail('The :attribute must be a string.');
            return;
        }

        // Check if it's a valid base64 string
        if (!preg_match('/^data:image\/(\w+);base64,/', $value, $matches)) {
            $fail('The :attribute must be a valid base64 encoded image.');
            return;
        }

        $mimeType = 'image/' . $matches[1];
        if (!in_array($mimeType, $this->allowedMimes)) {
            $fail('The :attribute must be a valid image type (jpeg, png, gif, webp).');
            return;
        }

        // Check file size
        $base64String = preg_replace('/^data:image\/\w+;base64,/', '', $value);
        $decodedSize = strlen(base64_decode($base64String));
        $sizeKb = $decodedSize / 1024;

        if ($sizeKb > $this->maxSizeKb) {
            $fail("The :attribute must not exceed {$this->maxSizeKb}KB.");
        }
    }
}
`
	path := filepath.Join(r.ProjectPath, "app/Rules/Base64Image.php")
	return os.WriteFile(path, []byte(content), 0644)
}

func (r *RulesSetup) createPhoneNumberRule() error {
	content := `<?php

namespace App\Rules;

use Closure;
use Illuminate\Contracts\Validation\ValidationRule;

class PhoneNumber implements ValidationRule
{
    protected ?string $countryCode;

    public function __construct(?string $countryCode = null)
    {
        $this->countryCode = $countryCode;
    }

    public function validate(string $attribute, mixed $value, Closure $fail): void
    {
        if (!is_string($value)) {
            $fail('The :attribute must be a string.');
            return;
        }

        // Remove spaces, dashes, and parentheses
        $cleaned = preg_replace('/[\s\-\(\)]/', '', $value);

        // Check if it starts with + for international format
        if (str_starts_with($cleaned, '+')) {
            // International format: +1234567890 (10-15 digits after +)
            if (!preg_match('/^\+[1-9]\d{9,14}$/', $cleaned)) {
                $fail('The :attribute must be a valid international phone number.');
                return;
            }
        } else {
            // Local format: at least 7 digits, max 15
            if (!preg_match('/^[0-9]{7,15}$/', $cleaned)) {
                $fail('The :attribute must be a valid phone number.');
                return;
            }
        }
    }
}
`
	path := filepath.Join(r.ProjectPath, "app/Rules/PhoneNumber.php")
	return os.WriteFile(path, []byte(content), 0644)
}

func (r *RulesSetup) createScopedUniqueRule() error {
	content := `<?php

namespace App\Rules;

use Closure;
use Illuminate\Contracts\Validation\ValidationRule;
use Illuminate\Support\Facades\DB;

class ScopedUnique implements ValidationRule
{
    protected string $table;
    protected string $column;
    protected array $scopes;
    protected ?int $ignoreId;

    public function __construct(string $table, string $column, array $scopes = [], ?int $ignoreId = null)
    {
        $this->table = $table;
        $this->column = $column;
        $this->scopes = $scopes;
        $this->ignoreId = $ignoreId;
    }

    public function validate(string $attribute, mixed $value, Closure $fail): void
    {
        $query = DB::table($this->table)->where($this->column, $value);

        foreach ($this->scopes as $scopeColumn => $scopeValue) {
            $query->where($scopeColumn, $scopeValue);
        }

        if ($this->ignoreId) {
            $query->where('id', '!=', $this->ignoreId);
        }

        if ($query->exists()) {
            $fail('The :attribute has already been taken.');
        }
    }
}
`
	path := filepath.Join(r.ProjectPath, "app/Rules/ScopedUnique.php")
	return os.WriteFile(path, []byte(content), 0644)
}

func (r *RulesSetup) createTimeFormatRule() error {
	content := `<?php

namespace App\Rules;

use Closure;
use Illuminate\Contracts\Validation\ValidationRule;

class TimeFormat implements ValidationRule
{
    protected string $format;

    public function __construct(string $format = 'H:i')
    {
        $this->format = $format;
    }

    public function validate(string $attribute, mixed $value, Closure $fail): void
    {
        $parsed = \DateTime::createFromFormat($this->format, $value);

        if (!$parsed || $parsed->format($this->format) !== $value) {
            $fail("The :attribute must be a valid time in {$this->format} format.");
        }
    }
}
`
	path := filepath.Join(r.ProjectPath, "app/Rules/TimeFormat.php")
	return os.WriteFile(path, []byte(content), 0644)
}

func (r *RulesSetup) createStrongPasswordRule() error {
	content := `<?php

namespace App\Rules;

use Closure;
use Illuminate\Contracts\Validation\ValidationRule;

class StrongPassword implements ValidationRule
{
    protected int $minLength;
    protected bool $requireUppercase;
    protected bool $requireNumber;
    protected bool $requireSpecial;

    public function __construct(
        int $minLength = 8,
        bool $requireUppercase = true,
        bool $requireNumber = true,
        bool $requireSpecial = true
    ) {
        $this->minLength = $minLength;
        $this->requireUppercase = $requireUppercase;
        $this->requireNumber = $requireNumber;
        $this->requireSpecial = $requireSpecial;
    }

    public function validate(string $attribute, mixed $value, Closure $fail): void
    {
        if (strlen($value) < $this->minLength) {
            $fail("The :attribute must be at least {$this->minLength} characters.");
            return;
        }

        if ($this->requireUppercase && !preg_match('/[A-Z]/', $value)) {
            $fail('The :attribute must contain at least one uppercase letter.');
            return;
        }

        if ($this->requireNumber && !preg_match('/[0-9]/', $value)) {
            $fail('The :attribute must contain at least one number.');
            return;
        }

        if ($this->requireSpecial && !preg_match('/[!@#$%^&*(),.?":{}|<>]/', $value)) {
            $fail('The :attribute must contain at least one special character.');
        }
    }
}
`
	path := filepath.Join(r.ProjectPath, "app/Rules/StrongPassword.php")
	return os.WriteFile(path, []byte(content), 0644)
}

package laravel

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type AuthSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewAuthSetup(projectPath string, dryRun bool) *AuthSetup {
	return &AuthSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (a *AuthSetup) Setup() error {
	if err := a.createAuthController(); err != nil {
		return err
	}
	if err := a.setupRoutes(); err != nil {
		return err
	}
	if err := a.ensureUserHasApiTokens(); err != nil {
		return err
	}
	return nil
}

func (a *AuthSetup) createAuthController() error {
	content := `<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Models\User;
use App\Support\Api\ApiResponse;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Auth;
use Illuminate\Support\Facades\Hash;
use Illuminate\Validation\ValidationException;

class AuthController extends Controller
{
    use ApiResponse;

    public function login(Request $request): JsonResponse
    {
        $request->validate([
            'email' => 'required|email',
            'password' => 'required',
            'device_name' => 'required',
        ]);

        $user = User::where('email', $request->email)->first();

        if (! $user || ! Hash::check($request->password, $user->password)) {
            throw ValidationException::withMessages([
                'email' => ['The provided credentials are incorrect.'],
            ]);
        }

        return $this->ok([
            'token' => $user->createToken($request->device_name)->plainTextToken,
            'user' => $user,
        ], 'Login successful');
    }

    public function logout(Request $request): JsonResponse
    {
        $request->user()->currentAccessToken()->delete();

        return $this->ok(null, 'Logged out successfully');
    }

    public function me(Request $request): JsonResponse
    {
        return $this->ok($request->user());
    }
}
`
	dir := filepath.Join(a.ProjectPath, "app/Http/Controllers/Api")
	if !a.DryRun {
		os.MkdirAll(dir, 0755)
	}
	path := filepath.Join(dir, "AuthController.php")
	if a.DryRun {
		fmt.Printf("[Dry Run] Would create AuthController: %s\n", path)
		return nil
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func (a *AuthSetup) setupRoutes() error {
	path := filepath.Join(a.ProjectPath, "routes/api.php")
	if a.DryRun {
		fmt.Printf("[Dry Run] Would add auth routes to %s\n", path)
		return nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	routes := `
use App\Http\Controllers\Api\AuthController;

Route::post('/login', [AuthController::class, 'login']);

Route::middleware('auth:sanctum')->group(function () {
    Route::get('/me', [AuthController::class, 'me']);
    Route::post('/logout', [AuthController::class, 'logout']);
});
`
	newContent := string(content) + routes
	return os.WriteFile(path, []byte(newContent), 0644)
}

func (a *AuthSetup) ensureUserHasApiTokens() error {
	path := filepath.Join(a.ProjectPath, "app/Models/User.php")
	if a.DryRun {
		fmt.Printf("[Dry Run] Would ensure User model uses HasApiTokens\n")
		return nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	sContent := string(content)
	if !strings.Contains(sContent, "Laravel\\Sanctum\\HasApiTokens") {
		// This is a simple replacement, assuming standard model structure
		sContent = strings.Replace(sContent, "use Illuminate\\Foundation\\Auth\\User as Authenticatable;", "use Illuminate\\Foundation\\Auth\\User as Authenticatable;\nuse Laravel\\Sanctum\\HasApiTokens;", 1)
		sContent = strings.Replace(sContent, "use HasFactory, Notifiable;", "use HasApiTokens, HasFactory, Notifiable;", 1)
	}

	return os.WriteFile(path, []byte(sContent), 0644)
}

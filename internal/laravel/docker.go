package laravel

import (
	"fmt"
	"os"
	"path/filepath"
)

type DockerSetup struct {
	ProjectPath string
	DryRun      bool
}

func NewDockerSetup(projectPath string, dryRun bool) *DockerSetup {
	return &DockerSetup{ProjectPath: projectPath, DryRun: dryRun}
}

func (d *DockerSetup) Setup() error {
	if err := d.createDockerDir(); err != nil {
		return err
	}
	if err := d.createDevDockerfile(); err != nil {
		return err
	}
	if err := d.createProdDockerfile(); err != nil {
		return err
	}
	if err := d.createDockerCompose(); err != nil {
		return err
	}
	return nil
}

func (d *DockerSetup) createDockerDir() error {
	path := filepath.Join(d.ProjectPath, "docker")
	if d.DryRun {
		fmt.Printf("[Dry Run] Would create directory: %s\n", path)
		return nil
	}
	return os.MkdirAll(path, 0755)
}

func (d *DockerSetup) createDevDockerfile() error {
	content := `FROM php:8.3-fpm

# Install system dependencies
RUN apt-get update && apt-get install -y \
    git \
    curl \
    libpng-dev \
    libonig-dev \
    libxml2-dev \
    zip \
    unzip

# Clear cache
RUN apt-get clean && rm -rf /var/lib/apt/lists/*

# Install PHP extensions
RUN docker-php-ext-install pdo_mysql mbstring exif pcntl bcmath gd

# Get latest Composer
COPY --from=composer:latest /usr/bin/composer /usr/bin/composer

# Set working directory
WORKDIR /var/www

USER $user
`
	path := filepath.Join(d.ProjectPath, "docker/Dockerfile")
	if d.DryRun {
		fmt.Printf("[Dry Run] Would create file: %s\n", path)
		return nil
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func (d *DockerSetup) createProdDockerfile() error {
	content := `FROM php:8.3-fpm as build

WORKDIR /var/www

RUN apt-get update && apt-get install -y \
    git \
    unzip \
    libpng-dev \
    libonig-dev \
    libxml2-dev

RUN docker-php-ext-install pdo_mysql mbstring exif pcntl bcmath gd

COPY . .
RUN composer install --no-dev --optimize-autoloader

FROM php:8.3-fpm-alpine

RUN docker-php-ext-install pdo_mysql

COPY --from=build /var/www /var/www

WORKDIR /var/www
`
	path := filepath.Join(d.ProjectPath, "docker/Dockerfile.prod")
	if d.DryRun {
		fmt.Printf("[Dry Run] Would create file: %s\n", path)
		return nil
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func (d *DockerSetup) createDockerCompose() error {
	content := `services:
  app:
    build:
      context: .
      dockerfile: docker/Dockerfile
    image: myapp-app
    container_name: myapp-app
    restart: unless-stopped
    working_dir: /var/www
    volumes:
      - ./:/var/www
    networks:
      - myapp-network

  db:
    image: mysql:8.0
    container_name: myapp-db
    restart: unless-stopped
    environment:
      MYSQL_DATABASE: ${DB_DATABASE}
      MYSQL_ROOT_PASSWORD: ${DB_PASSWORD}
      MYSQL_PASSWORD: ${DB_PASSWORD}
      MYSQL_USER: ${DB_USERNAME}
    volumes:
      - dbdata:/var/lib/mysql
    networks:
      - myapp-network

networks:
  myapp-network:
    driver: bridge

volumes:
  dbdata:
`
	path := filepath.Join(d.ProjectPath, "docker-compose.yml")
	if d.DryRun {
		fmt.Printf("[Dry Run] Would create file: %s\n", path)
		return nil
	}
	return os.WriteFile(path, []byte(content), 0644)
}

.PHONY: dev build generate tailwind clean

# Run app in development
dev: generate
	go run ./cmd/main.go

# Generate templ files
generate:
	templ generate ./templates/...

# Build binary
build: generate
	go build -o ./bin/activity-monitor ./cmd/main.go

# Build tailwind (requires node)
tailwind:
	npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css --minify

# Watch tailwind
tailwind-watch:
	npx tailwindcss -i ./static/css/input.css -o ./static/css/output.css --watch

# Clean generated files
clean:
	find . -name "*_templ.go" -delete
	rm -f ./bin/activity-monitor

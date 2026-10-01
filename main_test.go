package main

import (
	"os"
	"path/filepath"
	"testing"

	"event_explorer/services"
)

func runMainWithStubbedServer(t *testing.T) {
	t.Helper()
	previousApp := services.App
	previousRunServer := runServer
	t.Setenv("GOOGLE_PLACES_API_KEY", "google-test-key")
	t.Setenv("TICKETMASTER_API_KEY", "ticketmaster-test-key")

	serverStarted := false
	runServer = func(...string) {
		serverStarted = true
		if services.App == nil || services.App.Events == nil {
			t.Error("server started before application services were initialized")
		}
	}
	t.Cleanup(func() {
		runServer = previousRunServer
		services.App = previousApp
	})

	main()

	if !serverStarted {
		t.Fatal("main did not start the server")
	}
	if services.App.GooglePlaces == nil || services.App.Ticketmaster == nil || services.App.Cache == nil {
		t.Fatal("main did not initialize the application dependencies")
	}
}

func TestMainInitializesServicesWhenDotEnvIsMissing(t *testing.T) {
	t.Chdir(t.TempDir())
	runMainWithStubbedServer(t)
}

func TestMainInitializesServicesWhenDotEnvExists(t *testing.T) {
	workingDirectory := t.TempDir()
	if err := os.WriteFile(filepath.Join(workingDirectory, ".env"), nil, 0600); err != nil {
		t.Fatalf("create empty .env file: %v", err)
	}
	t.Chdir(workingDirectory)
	runMainWithStubbedServer(t)
}

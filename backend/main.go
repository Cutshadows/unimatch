package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
	"unimatch-back/internal/app"
	"unimatch-back/internal/routes"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 3030, "Port to run the application on")
	flag.Parse()
	app, err := app.NewApplication()
	if err != nil {
		panic(err)
	}

	app.Logger.Println("Application is running")
	r := routes.SetupRoutes(app)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		IdleTimeout: time.Minute, // Idle timeout in seconds
		Handler: r,
		ReadTimeout: 10 * time.Second, // Read timeout in seconds
		WriteTimeout: 30 * time.Second , // Write timeout in seconds
	}
	app.Logger.Printf("Server started successfully on port %d", port)

	err = server.ListenAndServe()
	if err != nil {
		app.Logger.Fatalf("Failed to start server: %v", err)
	} 
		
}


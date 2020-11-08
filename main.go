package main

import (
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"app_api/apis/vehicle"
	gmConnector "app_api/shared/gm"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

//go:generate swagger generate spec -m -o ./swagger/swagger.json

var env *Env
var r *mux.Router

// Env ... export db and router with Env
type Env struct {
	// The struct for storing services with all necessary dependencies.
	// More services will be gradually added for different versions as code refactoring and replacement of functions with services.
	// When refactoring, replace function calls with services and add those services here.
	Services Services
}

// struct for splitting services by versions
type Services struct {
	VehicleService vehicle.Service
}

// Initialize ... initialize the env so we can use it in testing
func Initialize() {
	// init all services
	gmAPIConnector := gmConnector.NewGMAPIConnector()

	vehicleService := vehicle.NewService(gmAPIConnector)

	r = mux.NewRouter()

	env = &Env{
		// init services struct.
		Services: Services{
			VehicleService: vehicleService,
		},
	}
	env.initializeRoutes()
}

func (env *Env) initializeRoutes() {
	r.HandleFunc("/vehicles/{vehicle_id}", env.getVehicle).Methods("GET")
	r.HandleFunc("/vehicles/{vehicle_id}/doors", env.getVehicleDoors).Methods("GET")
	r.HandleFunc("/vehicles/{vehicle_id}/fuel", env.getVehicleFuelStatus).Methods("GET")
	r.HandleFunc("/vehicles/{vehicle_id}/battery", env.getVehicleBatteryStatus).Methods("GET")
	r.HandleFunc("/vehicles/{vehicle_id}/engine", env.actionEngine).Methods("POST")

	r.Use(Logger)
}

func main() {
	// set up logger
	var logFile string = os.Getenv("LOG_FILE")
	var level string = os.Getenv("ENVIRONMENT")
	switch level {
	case "development":
		log.SetReportCaller(false)
		log.SetLevel(log.TraceLevel)
	case "testing":
		log.SetReportCaller(false)
		log.SetLevel(log.DebugLevel)
	case "production":
		log.SetReportCaller(false)
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetReportCaller(false)
		log.SetLevel(log.InfoLevel)
	}
	log.SetOutput(os.Stdout)
	file, err := os.OpenFile(logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	if level == "development" {
		mWriter := io.MultiWriter(os.Stdout, file)
		log.SetOutput(mWriter)
	}
	defer file.Close()

	Initialize()
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8003"
	}
	go func() {
		log.Println("Starting Server")
		if err := http.ListenAndServe(":"+port, r); err != nil {
			log.Fatal("web-server error:", err)
		}
	}()
	// Graceful Shutdown
	waitForShutdown()
}

func waitForShutdown() {
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// Block until we receive our signal.
	<-interruptChan
	os.Exit(0)
}

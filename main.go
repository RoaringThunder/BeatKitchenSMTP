package main

import (
	"fmt"
	"log"
	"net/http"
	controllers "salamander-smtp/controllers/verification"
	"salamander-smtp/database"
	"salamander-smtp/logging"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func init() {
	err := database.InitializeThunderDome()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Welllcomeeee to the Thunderrrrr Dome!")

	err = logging.InitLogging()
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	fmt.Println("SMTP is a go!")

	router := mux.NewRouter()

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "Access-Control-Allow-Credentials"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	origins := handlers.AllowedOrigins([]string{"http://localhost:3000"})
	credentials := handlers.AllowCredentials()

	router.HandleFunc("/api/smtp/verify/send-verification", controllers.SendVerificationEmail).Methods("GET")
	router.HandleFunc("/api/smtp/verify", controllers.VerifyUser).Methods("POST")

	log.Fatal(http.ListenAndServe("127.0.0.1:10042", handlers.CORS(headers, methods, origins, credentials)(router)))
}

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ents-source/gift-cards/api"
	"github.com/ents-source/gift-cards/assets"
	"github.com/ents-source/gift-cards/database"
	"github.com/ents-source/go-amember-api/amember"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	HttpBind string `envconfig:"http_bind" default:"0.0.0.0:8080"`

	AmpApiKey     string `envconfig:"amp_api_key"`
	AmpApiKeyFile string `envconfig:"amp_api_key_file"`
	AmpApiUrl     string `envconfig:"amp_api_url"`

	DataDir string `envconfig:"data_dir" default:"./data"`
}

func main() {
	log.SetOutput(os.Stdout)
	log.SetPrefix("[gift-cards] ")
	log.SetFlags(log.LstdFlags | log.Lshortfile | log.Lmicroseconds)

	var c config
	err := envconfig.Process("gc", &c)
	if err != nil {
		log.Fatal(err)
	}

	devMode := os.Getenv("DEV") == "true"

	webPath := assets.SetupWeb()
	if devMode {
		webPath = "./web"
	}

	migrationsPath := assets.SetupMigrations()
	if devMode {
		migrationsPath = "./migrations"
	}

	paymentsApi := amember.NewClient(c.AmpApiUrl, getPassword(c.AmpApiKey, c.AmpApiKeyFile))
	db, err := database.NewDatabase(c.DataDir, migrationsPath)
	if err != nil {
		log.Fatal(err)
	}
	card, err := db.GetOrCreateCard("test")
	log.Println(card, err)

	wg := api.Start(c.HttpBind, webPath, paymentsApi)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	go func() {
		defer close(stop)
		<-stop

		log.Println("Stopping api...")
		api.Stop()

		log.Println("Cleaning up...")
		_ = os.RemoveAll(webPath)

		log.Println("Done stopping")
	}()

	wg.Add(1)
	wg.Wait()

	log.Println("Goodbye!")
}

func getPassword(in string, f string) string {
	passwd := in
	if f != "" {
		b, err := os.ReadFile(f)
		if err != nil {
			log.Fatal(err)
		}
		passwd = string(b)
	}
	return passwd
}

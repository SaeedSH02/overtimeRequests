package cmd

import (
	"context"
	"shiftdony/config"
	postgres "shiftdony/database"
	"shiftdony/routes"

	log "shiftdony/logs"
	"time"

	"github.com/spf13/cobra"
)

func Start() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "starting....",
		Run: func(cmd *cobra.Command, args []string) {
			start()
		},
	}

	return cmd
}

func start() {
	const op = ("cmd.start")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	// TODO: Read config file in any exists

	// TODO: If there was a redis configuration and that was connectable, connect!
	// If there wasn't any, initialize an ephemeral session
	// session := ephemeral.New()

	db, err := postgres.NewPostgres(config.C.Postgres)
	if err != nil {
		log.Error(op, "cannot connect to postgres", err)
		return
	}
	db.Migrate(ctx)

	
	router := routes.SetupRouter(db.DB()) 
	if err := router.Run(":8080"); err != nil {
		log.Gl.Fatal("Failed to run server: " + err.Error())
	}
	log.Gl.Info("Starting shiftdoni web server...")

}

package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"shiftdony/config"
	log "shiftdony/logs"
	"shiftdony/models"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	bunDebug "github.com/uptrace/bun/extra/bundebug"
)

type Postgres struct {
	db *bun.DB
}

func NewPostgres(cfg config.Postgres) (*Postgres, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DB,
		cfg.SSLMode,
		// cfg.Timezone, // TODO: Timezone should be specified.
	)

	// dsn := "host=localhost user=myuser password=mypass dbname=mydb port=5432 sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())
	if err := sqldb.Ping(); err != nil {
		return nil, fmt.Errorf("cannot connect to database: %w", err)
	}

	db.AddQueryHook(bunDebug.NewQueryHook(
		bunDebug.WithVerbose(true),
		bunDebug.FromEnv("DEBUG"),
		bunDebug.WithWriter(os.Stdout),
	))
	return &Postgres{
		db: db,
	}, nil
}

func (pg *Postgres) Migrate(ctx context.Context) {
	pg.createTables(ctx)
}

func (p *Postgres) DB() *bun.DB {
	return p.db
}
func (pg *Postgres) createTables(ctx context.Context) error {
	// Create User table
	_, err := pg.db.NewCreateTable().
		Model((*models.User)(nil)).
		IfNotExists().
		Exec(ctx)

	if err != nil {
		log.Gl.Fatal(err.Error())
	}

	// Create Team table
	_, err = pg.db.NewCreateTable().
		Model((*models.Team)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		log.Gl.Fatal(err.Error())
	}

	_, err = pg.db.NewCreateTable().
		Model((*models.OvertimeSlot)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		log.Gl.Fatal(err.Error())
	}

	_, err = pg.db.NewCreateTable().
		Model((*models.OvertimeRequest)(nil)).
		IfNotExists().
		Exec(ctx)
	if err != nil {
		log.Gl.Fatal(err.Error())
	}

	return nil
}

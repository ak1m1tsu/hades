package hades

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"time"

	zlog "github.com/rs/zerolog"
	"golang.org/x/sync/errgroup"
)

type Hades struct {
	ctx   context.Context
	errGr *errgroup.Group

	api *API
	log *zlog.Logger
}

func New() *Hades {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	errGr, ctx := errgroup.WithContext(ctx)

	return &Hades{
		ctx:   ctx,
		errGr: errGr,
	}
}

func (a *Hades) Run() error {
	zlog.Ctx(a.ctx).Info().Msg("Starting the service")
	return a.errGr.Wait()
}

func (a *Hades) WithAPI() *Hades {
	a.api = &API{
		Server: &http.Server{
			Addr:              ":3000",
			ReadHeaderTimeout: time.Second * 5,
		},
	}

	a.errGr.Go(func() error {
		zlog.Ctx(a.ctx).Info().Msgf("Starting HTTP API at %s", a.api.Server.Addr)
		if err := a.api.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})
	a.errGr.Go(func() error {
		<-a.ctx.Done()
		zlog.Ctx(a.ctx).Info().Msg("Gracefully shutting down HTTP API")
		return a.api.Server.Shutdown(context.Background())
	})

	return a
}

func (a *Hades) WithLogger() *Hades {
	a.ctx = zlog.New(os.Stdout).With().Timestamp().Caller().Logger().WithContext(a.ctx)

	return a
}

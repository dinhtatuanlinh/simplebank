package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"simplebank/api"
	"simplebank/db"
	sqlc "simplebank/db/sqlc"
	"simplebank/gapi"
	"simplebank/mail"
	"simplebank/pb"
	"simplebank/util"
	"simplebank/worker"
	"syscall"

	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var interruptSignals = []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGINT}

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msgf("cannot load config %s", err)
	}

	if config.Environment == "dev" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	//conn, err := sql.Open(config.DBDriver, config.DBSource)
	//if err != nil {
	//	log.Fatal().Msgf("cannot connect to db: %s", err)
	//}
	connPool, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatal().Msgf("cannot connect to db: %s", err)
	}

	db.RunMigration(&config.MigrationUrl, &config.DBSource)

	store := sqlc.NewStore(connPool)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	waitGroup, ctx := errgroup.WithContext(ctx)

	runTaskProcessor(ctx, waitGroup, config, redisOpt, store)
	runGinServer(ctx, waitGroup, config, store, taskDistributor)
	runGrpcServer(ctx, waitGroup, config, store, taskDistributor)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msgf("wait group finished with error: %s", err)
	}

	//go runTaskProcessor(config, redisOpt, store)
	//go runGinServer(config, store, taskDistributor)
	//runGrpcServer(config, store, taskDistributor)

}

func runGrpcServer(ctx context.Context, waitGroup *errgroup.Group, config util.Config, store sqlc.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msgf("cannot create server: %s", err)
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Msgf("cannot listen: %s", err)
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("grpc server listening on %s", config.GRPCServerAddress)
		err = grpcServer.Serve(listener)
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}

			log.Error().Err(err).Msg("grpc server failed to serve")
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msgf("grpc server gracefully shutting down")

		grpcServer.GracefulStop()
		log.Info().Msgf("grpc server stopped")

		return nil
	})
	//log.Info().Msgf("start gPRC server at %s", config.GRPCServerAddress)
	//if err := grpcServer.Serve(listener); err != nil {
	//	log.Fatal().Msgf("cannot start gRPC server: %s", err)
	//}
}

func runGinServer(ctx context.Context, waitGroup *errgroup.Group, config util.Config, store sqlc.Store, taskDistributor worker.TaskDistributor) {
	server, err := api.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msgf("cannot create server: %s", err)
	}

	httpServer := &http.Server{
		Addr:    config.HTTPServerAddress,
		Handler: server.Router,
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("start HTTP server at %s", config.HTTPServerAddress)
		err = httpServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Error().Err(err).Msg("HTTP server failed to serve")
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msgf("HTTP server gracefully shutting down")


		if err := httpServer.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("failed to shutdown HTTP server gracefully")
			return err
		}

		log.Info().Msgf("HTTP server stopped")
		return nil
	})
}

func runTaskProcessor(ctx context.Context, waitGroup *errgroup.Group, config util.Config, redisOpt asynq.RedisClientOpt, store sqlc.Store) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Msg("cannot start task processor")
	}

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msgf("task processor gracefully shutting down")

		taskProcessor.Shutdown()
		log.Info().Msgf("task processor stopped")

		return nil
	})
}

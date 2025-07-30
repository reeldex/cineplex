package main

import (
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"os"
	service2 "scraper/internal/services/service"
	"time"
)

var (
	scraper *service2.Scraper
	imdb    *service2.IMDB
)

func main() {

	rand.Seed(time.Now().UnixNano())

	lg, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	lg.Info("hello", zap.String("hello", "world"))

	lg.Sync()

	fmt.Println()

	return
	//
	//defer func() {
	//	err = conn.Close()
	//	if err != nil {
	//		logger.Errorf("could not close grpc connection properly %v", err)
	//		ScraperErrors.WithLabelValues("could_not_close_grpc_connection").Inc()
	//	}
	//}()
	//
	//Start(logger)
	//
	//initServices(db, conn, logger)
	//
	//ticker := time.NewTicker(time.Second * 60)
	//defer ticker.Stop()
	//
	//ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	//defer cancel()
	//
	//logger.Info("Scraper started!")
	//
	//process(scraper, imdb, mailer, logger)
	//
	//ScraperHeartbeat.Inc()
	//// syscall.SIGHUP
	//// possibly make ticker change on like the do in caddy and prometheus
	//
	//done := make(chan struct{})
	//
	//go func() {
	//	for {
	//		select {
	//		case <-ctx.Done():
	//			logger.Info("program execution interrupted, exiting")
	//			done <- struct{}{}
	//			return
	//		case <-ticker.C:
	//			ScraperHeartbeat.Inc()
	//			process(scraper, imdb, mailer, logger)
	//			logger.Info("done processing")
	//		}
	//	}
	//}()
	//
	//<-done
	//
	//fmt.Println("Scraper service stopped!")
}

// process works using kind of "pipeline" pattern, or smth close
//func process(scraper *service.Scraper, imdb *service.IMDB, mailer *service.Mailer, logger logrus.FieldLogger) {
//	timer := prometheus.NewTimer(ScraperLatency)
//	defer timer.ObserveDuration()
//
//	ctx, release := context.WithTimeout(context.Background(), time.Second*30)
//	defer release()
//
//	// here
//	err := mailer.Send(imdb.GetFilms(scraper.GetFilms(ctx)))
//	if err != nil {
//		ScraperErrors.WithLabelValues("scraper_error").Inc()
//		logger.Println(err)
//	}
//}

//func initServices(db *sql.DB, grpcConn grpc.ClientConnInterface, logger logrus.FieldLogger) {
//
//	filmRepo := repository.NewFilmStorageRepository(db)
//
//	httpClient, err := client.NewHttpClientWithCookies()
//	if err != nil {
//		telemetry.ScraperErrors.WithLabelValues("unable_to_create_http_client_with_cookies").Inc()
//		logger.Fatalln(err)
//	}
//
//	soupService := fetcher.NewSoupDecorator(httpClient)
//	scraper = service.NewScraper(filmRepo, soupService, logger)
//	if err != nil {
//		telemetry.ScraperErrors.WithLabelValues("unable_to_create_scraper").Inc()
//		logger.Fatalln(err)
//	}
//
//	imdb = service.NewIMDB(proto.NewIMDBClient(grpcConn), logger)
//
//	/*
//		emailRepo := repository.NewSubscriberRepository(logger, db)
//
//			env, err := config.GetEnv()
//			if err != nil {
//				metrics.ScraperErrors.WithLabelValues("incomplete_environment").Inc()
//				logger.Fatalln(err)
//			}
//
//				mailjetClient := mailjet.NewMailjetClient(env.MailjetPubKey, env.MailJetPrivateKey)
//
//				mailer = service.NewMailer(mailjetClient, service.MailerConfig{
//					FromEmail: env.FromEmail,
//					FromName:  env.FromName,
//				}, emailRepo, logger)
//	*/
//}
//
//func initAndMaintainDB() (*sql.DB, error) {
//	db, err := sql.Open("sqlite3", "file:./pirata.db")
//	if err != nil {
//		telemetry.ScraperErrors.WithLabelValues("unable_to_establish_db_connection").Inc()
//		return nil, err
//	}
//
//	sourceDriver, err := iofs.New(migrationFiles, "migrations")
//	if err != nil {
//		telemetry.ScraperErrors.WithLabelValues("fs_error_migration_files").Inc()
//		return nil, err
//	}
//
//	migrationDriver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
//	if err != nil {
//		telemetry.ScraperErrors.WithLabelValues("unable_to_init_migration_driver").Inc()
//		return nil, err
//	}
//
//	// NewWithInstance always returns nil error
//	migration, err := migrate.NewWithInstance("migrations", sourceDriver, "sqlite3", migrationDriver)
//	if err != nil {
//		telemetry.ScraperErrors.WithLabelValues("unable_to_init_migration").Inc()
//		return nil, err
//	}
//
//	err = migration.Up()
//	if err != nil && err != migrate.ErrNoChange {
//		telemetry.ScraperErrors.WithLabelValues("migration_failed").Inc()
//		return nil, err
//	}
//
//	return db, nil
//}
//
//func initLogger(serviceName string) *logrus.Entry {
//	log := logrus.New()
//
//	return log.WithFields(logrus.Fields{"service_name": serviceName})
//}

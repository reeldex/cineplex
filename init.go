package main

import (
	"common/client"
	"common/proto"
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"scraper/internal/services/fetcher"
	"scraper/pkg/telemetry"
	"scraper/service"
	"scraper/service/decorator"

	"google.golang.org/grpc"

	"scraper/storage/repository"

	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/sirupsen/logrus"
)

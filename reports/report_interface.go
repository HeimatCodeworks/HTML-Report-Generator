package reports

import (
	"search-ads-reporter-ui/internal/database"
	"search-ads-reporter-ui/internal/email"
)

type ReportGenerator interface {
	Name() string

	TemplatePath() string

	GenerateEmailData(dbClient *database.DBClient, devMode bool, teamID, mongoURI string) (email.EmailData, error)
}

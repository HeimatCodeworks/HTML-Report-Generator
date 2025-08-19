package main

import (
	"context"
	"fmt"
	"search-ads-reporter-ui/internal/config"
	"search-ads-reporter-ui/internal/database"
	"search-ads-reporter-ui/internal/email"
	"search-ads-reporter-ui/internal/email_sender"
	"search-ads-reporter-ui/reports"
	"sort"
	"strings"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type Company struct {
	TeamID string `json:"teamId"`
	Name   string `json:"name"`
}

type ReportInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (a *App) GetCompanies() ([]Company, error) {
	userTeams, err := database.FetchDummyCompanies("data/dummy_companies.json")
	if err != nil {
		return nil, err
	}

	var companies []Company
	for _, team := range userTeams {
		companies = append(companies, Company{
			TeamID: team.TeamID,
			Name:   team.Company,
		})
	}

	sort.Slice(companies, func(i, j int) bool {
		return companies[i].Name < companies[j].Name
	})

	return companies, nil
}

func (a *App) GetReports() []ReportInfo {
	availableReports := []reports.ReportGenerator{
		&reports.CampaignPerformanceReport{},
		&reports.AccountOverviewReport{},
	}

	var reportInfos []ReportInfo
	for i, r := range availableReports {
		reportInfos = append(reportInfos, ReportInfo{ID: i, Name: r.Name()})
	}
	return reportInfos
}

func (a *App) GenerateReport(companyID string, reportID int) (string, error) {
	devMode := true
	var dbClient *database.DBClient

	availableReports := []reports.ReportGenerator{
		&reports.CampaignPerformanceReport{},
		&reports.AccountOverviewReport{},
	}

	if reportID < 0 || reportID >= len(availableReports) {
		return "", fmt.Errorf("invalid report ID: %d", reportID)
	}

	selectedReport := availableReports[reportID]

	emailData, err := selectedReport.GenerateEmailData(dbClient, devMode, companyID, "")
	if err != nil {
		return "", err
	}

	companies, err := a.GetCompanies()
	if err != nil {
		return "", err
	}
	var companyName = "Selected Company"
	for _, c := range companies {
		if c.TeamID == companyID {
			companyName = c.Name
			break
		}
	}
	emailData["ClientName"] = companyName

	htmlBody, err := email.GenerateEmail(selectedReport.TemplatePath(), emailData)
	if err != nil {
		return "", err
	}

	return htmlBody, nil
}

func (a *App) SendEmail(toRecipients []string, ccRecipients []string, subject string, body string) error {
	cfg, err := config.LoadConfig("configs/config.json")
	if err != nil {
		return fmt.Errorf("could not load config to send email: %w", err)
	}

	if cfg.EmailAddress == "your-email@gmail.com" || cfg.AppPassword == "your-gmail-app-password" {
		return fmt.Errorf("email credentials are not configured in configs/config.json")
	}

	to := cleanEmailList(toRecipients)
	cc := cleanEmailList(ccRecipients)

	if len(to) == 0 {
		return fmt.Errorf("at least one 'To' recipient is required")
	}

	return email_sender.SendEmail(cfg.EmailAddress, cfg.AppPassword, to, cc, subject, body)
}

func cleanEmailList(emails []string) []string {
	var cleaned []string
	for _, email := range emails {
		trimmed := strings.TrimSpace(email)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return cleaned
}

func (a *App) GetUsersForCompany(teamID string) ([]string, error) {
	users, err := database.FetchDummyUsersByTeamID("data/dummy_users.json", teamID)
	if err != nil {
		return nil, err
	}

	var emails []string
	for _, user := range users {
		emails = append(emails, user.Email)
	}

	return emails, nil
}

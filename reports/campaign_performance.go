package reports

import (
	"fmt"
	"search-ads-reporter-ui/internal/database"
	"search-ads-reporter-ui/internal/email"
	"time"
)

type CampaignPerformanceReport struct{}

func (r *CampaignPerformanceReport) Name() string {
	return "Weekly Campaign Performance"
}

func (r *CampaignPerformanceReport) TemplatePath() string {
	return "templates/campaign_performance/template.html"
}

func (r *CampaignPerformanceReport) GenerateEmailData(dbClient *database.DBClient, devMode bool, teamID, mongoURI string) (email.EmailData, error) {
	now := time.Now()
	endDate := now.Truncate(24 * time.Hour)
	startDate := endDate.AddDate(0, 0, -7)
	previousEndDate := startDate
	previousStartDate := previousEndDate.AddDate(0, 0, -7)

	currentWeekDays := getDateRangeStrings(startDate, endDate)
	previousWeekDays := getDateRangeStrings(previousStartDate, previousEndDate)

	var dailyCampaigns, prevDailyCampaigns []database.CampaignReport
	var err error

	if devMode {
		dailyCampaigns, err = database.FetchDummyCampaignData("data/dummy_campaigns.json", teamID, currentWeekDays)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch current week dummy data: %w", err)
		}
		prevDailyCampaigns, err = database.FetchDummyCampaignData("data/dummy_campaigns.json", teamID, previousWeekDays)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch previous week dummy data: %w", err)
		}
	} else {
		if dbClient == nil {
			return nil, fmt.Errorf("database connection is not available for production mode")
		}
		dailyCampaigns, err = dbClient.FetchCampaignData(teamID, currentWeekDays)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch current week campaign data: %w", err)
		}
		prevDailyCampaigns, err = dbClient.FetchCampaignData(teamID, previousWeekDays)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch previous week campaign data: %w", err)
		}
	}

	aggregatedCampaigns := aggregateCampaigns(dailyCampaigns)
	aggregatedPrevCampaigns := aggregateCampaigns(prevDailyCampaigns)

	data := map[string]interface{}{
		"StartDate":         startDate.Format("2006-01-02"),
		"EndDate":           endDate.Add(-24 * time.Hour).Format("2006-01-02"),
		"PreviousStartDate": previousStartDate.Format("2006-01-02"),
		"PreviousEndDate":   previousEndDate.Add(-24 * time.Hour).Format("2006-01-02"),
		"Campaigns":         aggregatedCampaigns,
		"PrevCampaigns":     aggregatedPrevCampaigns,
	}

	return data, nil
}

func getDateRangeStrings(start, end time.Time) []string {
	var dates []string
	for d := start; d.Before(end); d = d.AddDate(0, 0, 1) {
		dates = append(dates, d.Format("2006-01-02"))
	}
	return dates
}

func aggregateCampaigns(dailyReports []database.CampaignReport) []email.CampaignSummary {
	summaryMap := make(map[string]email.CampaignSummary)
	for _, report := range dailyReports {
		campaignName := report.Metadata.CampaignName
		summary := summaryMap[campaignName]
		summary.CampaignName = campaignName
		summary.Impressions += report.Total.Impressions
		summary.Taps += report.Total.Taps
		summary.Installs += report.Total.Installs
		summary.Spend += report.Total.LocalSpend.Amount
		summaryMap[campaignName] = summary
	}
	var result []email.CampaignSummary
	for _, summary := range summaryMap {
		result = append(result, summary)
	}
	return result
}

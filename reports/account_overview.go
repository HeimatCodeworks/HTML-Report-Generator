package reports

import (
	"fmt"
	"search-ads-reporter-ui/internal/database"
	"search-ads-reporter-ui/internal/email"
	"sort"
	"time"
)

type AccountOverviewReport struct{}

func (r *AccountOverviewReport) Name() string {
	return "Monthly Account Overview"
}

func (r *AccountOverviewReport) TemplatePath() string {
	return "templates/account_overview/template.html"
}

type PerfMetrics struct {
	Spend    float64
	Installs int
	CPA      float64
}

type PerformanceChange struct {
	Spend    float64
	Installs float64
	CPA      float64
}

type PerformanceDisplaySimple struct {
	SpendStr    string
	InstallsStr string
	CPAStr      string
}

type PerformanceDisplay struct {
	Metrics           PerfMetrics
	Change            PerformanceChange
	SpendStr          string
	InstallsStr       string
	CPAStr            string
	SpendChangeStr    string
	InstallsChangeStr string
	CPAChangeStr      string
	SpendColor        string
	InstallsColor     string
	CPAColor          string
}

type CampaignPerformance struct {
	CampaignName   string
	CurrentPeriod  PerformanceDisplay
	PreviousPeriod PerformanceDisplaySimple
}

type AppPerformance struct {
	AppName        string
	CurrentPeriod  PerformanceDisplay
	PreviousPeriod PerformanceDisplaySimple
	Campaigns      []CampaignPerformance
}

func (r *AccountOverviewReport) GenerateEmailData(dbClient *database.DBClient, devMode bool, teamID, mongoURI string) (email.EmailData, error) {

	now := time.Now()
	currentPeriodEnd := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Add(-time.Second)
	currentPeriodStart := time.Date(currentPeriodEnd.Year(), currentPeriodEnd.Month(), 1, 0, 0, 0, 0, time.UTC)
	previousPeriodEnd := currentPeriodStart.Add(-time.Second)
	previousPeriodStart := time.Date(previousPeriodEnd.Year(), previousPeriodEnd.Month(), 1, 0, 0, 0, 0, time.UTC)

	currentMonthName := currentPeriodStart.Format("January")
	previousMonthName := previousPeriodStart.Format("January")

	currentPeriodDays := getDateRangeStrings(currentPeriodStart, currentPeriodEnd.AddDate(0, 0, 1))
	previousPeriodDays := getDateRangeStrings(previousPeriodStart, previousPeriodEnd.AddDate(0, 0, 1))

	var allApps []database.App
	var allCampaignReports []database.CampaignReport
	var err error

	if devMode {
		allApps, err = database.FetchDummyApps("data/dummy_apps.json")
		if err != nil {
			return nil, err
		}
		allCampaignReports, err = database.FetchDummyCampaignData("data/dummy_campaigns.json", teamID, append(currentPeriodDays, previousPeriodDays...))
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("production mode for this report is not yet implemented")
	}

	activeApps := filterActiveApps(allApps, teamID)
	sort.Slice(activeApps, func(i, j int) bool { return activeApps[i].AppName < activeApps[j].AppName })

	var allAppsPerformance []AppPerformance
	overallCurrentMetrics := PerfMetrics{}
	overallPreviousMetrics := PerfMetrics{}

	for _, app := range activeApps {
		appCurrentMetrics, appPreviousMetrics, appCampaigns := processCampaignsForApp(allCampaignReports, app.AdamID, currentPeriodDays, previousPeriodDays)

		overallCurrentMetrics.Spend += appCurrentMetrics.Spend
		overallCurrentMetrics.Installs += appCurrentMetrics.Installs
		overallPreviousMetrics.Spend += appPreviousMetrics.Spend
		overallPreviousMetrics.Installs += appPreviousMetrics.Installs

		appData := AppPerformance{
			AppName:        app.AppName,
			CurrentPeriod:  calculatePerformanceDisplay(appCurrentMetrics, appPreviousMetrics),
			PreviousPeriod: formatSimpleDisplay(appPreviousMetrics),
			Campaigns:      appCampaigns,
		}
		allAppsPerformance = append(allAppsPerformance, appData)
	}

	if overallPreviousMetrics.Installs > 0 {
		overallPreviousMetrics.CPA = overallPreviousMetrics.Spend / float64(overallPreviousMetrics.Installs)
	}
	overallPerformance := calculatePerformanceDisplay(overallCurrentMetrics, overallPreviousMetrics)

	data := email.EmailData{
		"ClientName":      "Client Name",
		"CurrentMonth":    currentMonthName,
		"PreviousMonth":   previousMonthName,
		"Date":            now.Format("January 2, 2006"),
		"Overall":         overallPerformance,
		"Apps":            allAppsPerformance,
		"OverallPrevious": formatSimpleDisplay(overallPreviousMetrics),
	}

	return data, nil
}

func filterActiveApps(apps []database.App, teamID string) []database.App {
	var active []database.App
	for _, app := range apps {
		if app.TeamID == teamID && app.IntegrationStatus == "ACTIVE" {
			active = append(active, app)
		}
	}
	return active
}

func processCampaignsForApp(reports []database.CampaignReport, adamID int64, currentDays, previousDays []string) (PerfMetrics, PerfMetrics, []CampaignPerformance) {
	currentMetrics := PerfMetrics{}
	previousMetrics := PerfMetrics{}

	currentDayMap := make(map[string]bool)
	for _, day := range currentDays {
		currentDayMap[day] = true
	}
	previousDayMap := make(map[string]bool)
	for _, day := range previousDays {
		previousDayMap[day] = true
	}

	campaignsCurrent := make(map[string]PerfMetrics)
	campaignsPrevious := make(map[string]PerfMetrics)

	for _, report := range reports {
		meta := report.Metadata
		if meta.App.AdamID != adamID || meta.CampaignStatus != "ENABLED" || meta.DisplayStatus != "RUNNING" {
			continue
		}

		if _, ok := currentDayMap[report.Date]; ok {
			currentMetrics.Spend += report.Total.LocalSpend.Amount
			currentMetrics.Installs += report.Total.Installs

			c := campaignsCurrent[meta.CampaignName]
			c.Spend += report.Total.LocalSpend.Amount
			c.Installs += report.Total.Installs
			campaignsCurrent[meta.CampaignName] = c
		} else if _, ok := previousDayMap[report.Date]; ok {
			previousMetrics.Spend += report.Total.LocalSpend.Amount
			previousMetrics.Installs += report.Total.Installs

			c := campaignsPrevious[meta.CampaignName]
			c.Spend += report.Total.LocalSpend.Amount
			c.Installs += report.Total.Installs
			campaignsPrevious[meta.CampaignName] = c
		}
	}

	if currentMetrics.Installs > 0 {
		currentMetrics.CPA = currentMetrics.Spend / float64(currentMetrics.Installs)
	}
	if previousMetrics.Installs > 0 {
		previousMetrics.CPA = previousMetrics.Spend / float64(previousMetrics.Installs)
	}

	var campaignPerformances []CampaignPerformance
	for name, metrics := range campaignsCurrent {
		prevMetrics := campaignsPrevious[name]
		if prevMetrics.Installs > 0 {
			prevMetrics.CPA = prevMetrics.Spend / float64(prevMetrics.Installs)
		}

		campaignPerformances = append(campaignPerformances, CampaignPerformance{
			CampaignName:   name,
			CurrentPeriod:  calculatePerformanceDisplay(metrics, prevMetrics),
			PreviousPeriod: formatSimpleDisplay(prevMetrics),
		})
	}
	sort.Slice(campaignPerformances, func(i, j int) bool {
		return campaignPerformances[i].CampaignName < campaignPerformances[j].CampaignName
	})

	return currentMetrics, previousMetrics, campaignPerformances
}

func formatSimpleDisplay(metrics PerfMetrics) PerformanceDisplaySimple {
	return PerformanceDisplaySimple{
		SpendStr:    fmt.Sprintf("$%.2f", metrics.Spend),
		InstallsStr: fmt.Sprintf("%d", metrics.Installs),
		CPAStr:      fmt.Sprintf("$%.2f", metrics.CPA),
	}
}

func calculatePerformanceDisplay(current, previous PerfMetrics) PerformanceDisplay {
	if current.Installs > 0 {
		current.CPA = current.Spend / float64(current.Installs)
	}

	change := PerformanceChange{
		Spend:    calculateChange(current.Spend, previous.Spend),
		Installs: calculateChange(float64(current.Installs), float64(previous.Installs)),
		CPA:      calculateChange(current.CPA, previous.CPA),
	}

	return PerformanceDisplay{
		Metrics:           current,
		Change:            change,
		SpendStr:          fmt.Sprintf("$%.2f", current.Spend),
		InstallsStr:       fmt.Sprintf("%d", current.Installs),
		CPAStr:            fmt.Sprintf("$%.2f", current.CPA),
		SpendChangeStr:    formatChange(change.Spend),
		InstallsChangeStr: formatChange(change.Installs),
		CPAChangeStr:      formatChange(change.CPA),
		SpendColor:        formatColor(change.Spend, true),
		InstallsColor:     formatColor(change.Installs, false),
		CPAColor:          formatColor(change.CPA, true),
	}
}

func calculateChange(current, previous float64) float64 {
	if previous == 0 {
		if current == 0 {
			return 0.0
		}
		return 1.0
	}
	return (current - previous) / previous
}

func formatChange(val float64) string {
	if val == 0 {
		return "0.00%"
	}
	return fmt.Sprintf("%.2f%%", val*100)
}

func formatColor(val float64, isCost bool) string {
	green := "#10b981"
	red := "#ef4444"
	if val == 0 {
		return "#6b7280"
	}
	if isCost {
		if val < 0 {
			return green
		}
		return red
	}
	if val > 0 {
		return green
	}
	return red
}

package email

import (
	"bytes"
	"html/template"
)

type EmailData map[string]interface{}

type CampaignSummary struct {
	CampaignName string
	Impressions  int
	Taps         int
	Installs     int
	Spend        float64
}

func GenerateEmail(templatePath string, data EmailData) (string, error) {
	if companyName, ok := data["CompanyName"]; ok {
		data["CompanyName"] = companyName
	}

	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}

	var body bytes.Buffer
	err = tmpl.Execute(&body, data)
	if err != nil {
		return "", err
	}

	return body.String(), nil
}

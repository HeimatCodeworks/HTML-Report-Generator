package database

type UserTeam struct {
	TeamID  string `json:"TeamId"`
	Company string `json:"Company"`
}

type App struct {
	TeamID            string `json:"TeamId"`
	AppName           string `json:"AppName"`
	AdamID            int64  `json:"AdamId"`
	IntegrationStatus string `json:"IntegrationStatus"`
}

type CampaignReport struct {
	TeamID   string `bson:"TeamID"`
	Date     string `bson:"date"`
	Metadata struct {
		CampaignName   string `bson:"campaignName"`
		CampaignStatus string `bson:"campaignStatus"`
		DisplayStatus  string `bson:"displayStatus"`
		App            struct {
			AdamID int64 `bson:"adamId"`
		} `bson:"app"`
	} `bson:"Metadata"`
	Total struct {
		Impressions int `bson:"impressions"`
		Taps        int `bson:"taps"`
		Installs    int `bson:"installs"`
		LocalSpend  struct {
			Amount float64 `bson:"amount"`
		} `bson:"localSpend"`
	} `bson:"total"`
}

type AdGroupReport struct {
	TeamID   string `bson:"TeamID"`
	Date     string `bson:"date"`
	Metadata struct {
		CampaignID    int64  `bson:"campaignId"`
		AdGroupID     int64  `bson:"adGroupId"`
		AdGroupName   string `bson:"adGroupName"`
		AdGroupStatus string `bson:"adGroupStatus"`
	} `bson:"Metadata"`
	Total struct {
		Impressions int `bson:"impressions"`
		Taps        int `bson:"taps"`
		Installs    int `bson:"installs"`
		LocalSpend  struct {
			Amount float64 `bson:"amount"`
		} `bson:"localSpend"`
	} `bson:"total"`
}

type KeywordReport struct {
	TeamID   string `bson:"TeamID"`
	Date     string `bson:"date"`
	Metadata struct {
		CampaignID    int64  `bson:"campaignId"`
		AdGroupID     int64  `bson:"adGroupId"`
		KeywordID     int64  `bson:"keywordId"`
		Keyword       string `bson:"keyword"`
		MatchType     string `bson:"matchType"`
		KeywordStatus string `bson:"keywordStatus"`
	} `bson:"metadata"`
	Total struct {
		Impressions int `bson:"impressions"`
		Taps        int `bson:"taps"`
		Installs    int `bson:"installs"`
		LocalSpend  struct {
			Amount float64 `bson:"amount"`
		} `bson:"localSpend"`
	} `bson:"total"`
}

type MongoUser struct {
	TeamID string `json:"TeamID" bson:"TeamID"`
	Email  string `json:"Email" bson:"Email"`
}

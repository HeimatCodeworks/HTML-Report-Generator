package database

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DBClient struct {
	Client *mongo.Client
}

func Connect(connectionString string) (*DBClient, error) {
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	fmt.Println("Connected to Mongo")
	return &DBClient{Client: client}, nil
}

func (c *DBClient) FetchCampaignData(teamID string, dates []string) ([]CampaignReport, error) {
	var results []CampaignReport
	collection := c.Client.Database("searchadsmaven").Collection("campaignreports")

	filter := bson.M{
		"TeamID": teamID,
		"date": bson.M{
			"$in": dates,
		},
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (c *DBClient) FetchAdGroupData(teamID string, dates []string) ([]AdGroupReport, error) {
	var results []AdGroupReport
	collection := c.Client.Database("searchadsmaven").Collection("adgroupreports")

	filter := bson.M{
		"TeamID": teamID,
		"date": bson.M{
			"$in": dates,
		},
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (c *DBClient) FetchKeywordData(teamID string, dates []string) ([]KeywordReport, error) {
	var results []KeywordReport
	collection := c.Client.Database("searchadsmaven").Collection("targetingkeywordreports")

	filter := bson.M{
		"TeamID": teamID,
		"date": bson.M{
			"$in": dates,
		},
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func FetchDummyCampaignData(path string, teamID string, dates []string) ([]CampaignReport, error) {
	var allReports []CampaignReport
	var results []CampaignReport

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open dummy data file: %w", err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&allReports)
	if err != nil {
		return nil, fmt.Errorf("could not parse dummy data: %w", err)
	}

	datesToFind := make(map[string]bool)
	for _, date := range dates {
		datesToFind[date] = true
	}

	for _, report := range allReports {
		if report.TeamID == teamID {
			if _, found := datesToFind[report.Date]; found {
				results = append(results, report)
			}
		}
	}

	return results, nil
}

func FetchDummyCompanies(path string) ([]UserTeam, error) {
	var results []UserTeam

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open dummy companies file: %w", err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&results)
	if err != nil {
		return nil, fmt.Errorf("could not decode dummy companies data: %w", err)
	}

	return results, nil
}

func FetchDummyApps(path string) ([]App, error) {
	var results []App
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open dummy apps file: %w", err)
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&results)
	if err != nil {
		return nil, fmt.Errorf("could not decode dummy apps data: %w", err)
	}
	return results, nil
}

func (c *DBClient) FetchUsersByTeamID(teamID string) ([]MongoUser, error) {
	var results []MongoUser
	collection := c.Client.Database("searchadsmaven").Collection("mongousers")

	filter := bson.M{
		"TeamID": teamID,
	}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(context.TODO(), &results); err != nil {
		return nil, err
	}

	return results, nil
}

func FetchDummyUsersByTeamID(path string, teamID string) ([]MongoUser, error) {
	var allUsers []MongoUser
	var results []MongoUser

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open dummy users file: %w", err)
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&allUsers)
	if err != nil {
		return nil, fmt.Errorf("could not decode dummy users data: %w", err)
	}

	for _, user := range allUsers {
		if user.TeamID == teamID {
			results = append(results, user)
		}
	}

	return results, nil
}

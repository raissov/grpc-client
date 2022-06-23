package repository

import (
	"bytes"
	"client/internal/models"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

type VerifyTokenRepository struct {
	logger *zap.SugaredLogger
}

func newVerifyTokenRepo(logger *zap.SugaredLogger) *VerifyTokenRepository {
	return &VerifyTokenRepository{
		logger: logger,
	}
}

func (a *VerifyTokenRepository) VerifyToken(token string, url string) error {
	client := &http.Client{}

	connect := &models.Token{
		Token: token,
	}
	a.logger.Infof("VerifyToken: %+v", token)
	jsonBytes, err := json.Marshal(connect)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "http://"+url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	a.logger.Infof("request: %+v", req)
	resp, err := client.Do(req)
	if err != nil {
		a.logger.Infof("Error: %+v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("got %d status code from auth service", resp.StatusCode)
	}

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var r models.VerifyToken
	err = json.Unmarshal(response, &r)
	if err != nil {
		return err
	}
	a.logger.Infof("user public id: %s: ", r.Response.Public)
	return nil
}

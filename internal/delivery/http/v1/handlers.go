package delivery

import (
	"client/config"
	"client/internal/models"
	"client/internal/service"
	"client/pb"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"time"
)

type Handler struct {
	service    *service.Service
	logger     *zap.SugaredLogger
	cfg        *config.Configs
	grpcClient pb.MessageServiceClient
}

type TokenVerify struct {
	Token string `json:"access_token"`
}

//NewHandler - function that creates new handler. It takes services, zap looger and configs as an argument, and returns handler
func NewHandler(services *service.Service, logger *zap.SugaredLogger, cfg *config.Configs, grpcClient pb.MessageServiceClient) *Handler {
	return &Handler{
		service:    services,
		logger:     logger,
		cfg:        cfg,
		grpcClient: grpcClient,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	router.POST("/send-messages", h.SendMessages)
	return router
}

func (h *Handler) SendMessages(c *gin.Context) {
	requestBody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		h.logger.Errorf("Error occurred while reading request body: %s", err.Error())
		c.JSON(401, gin.H{
			"error": models.ErrInvalidInput.Error(),
		})
		return
	}
	var tokenVerify *TokenVerify
	err = json.Unmarshal(requestBody, &tokenVerify)
	if err != nil {
		h.logger.Errorf("Error occurred while unmarshalling request body: %s", err.Error())
		c.JSON(401, gin.H{
			"error": models.ErrInvalidInput.Error(),
		})
		return
	}
	h.logger.Infof("Token: %s", tokenVerify.Token)
	err = h.service.VerifyToken.VerifyToken(tokenVerify.Token, "myapp:3000/verify-token")
	if err != nil {
		h.logger.Errorf("Error occurred while verifying token: %s", err.Error())
		c.JSON(401, gin.H{
			"error": models.ErrInvalidInput.Error(),
		})
		return
	}

	err = h.service.VerifyToken.VerifyToken(tokenVerify.Token, "server-receiver:3002/approve")
	if err != nil {
		h.logger.Errorf("Error occurred while verifying token: %s", err.Error())
		c.JSON(401, gin.H{
			"error": models.ErrInvalidInput.Error(),
		})
		return
	}

	h.logger.Info("Starting to do a Server Streaming RPC...")

	req := &pb.MessageManyTimesRequest{
		Message: &pb.Message{
			FirstName: "SomeName",
			LastName:  "SomeSurname",
		},
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	resStream, err := h.grpcClient.MessageManyTimes(ctx, req)
	if err != nil {
		h.logger.Errorf("Error occurred while doing a Server Streaming RPC: %s", err.Error())
		c.JSON(500, gin.H{
			"error": models.ErrInternalServerError.Error(),
		})
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			if err.Error() == "rpc error: code = Unknown desc = DATABASE_ERROR" {
				break
			}
			h.logger.Errorf("Error occurred while reading response from Server Streaming RPC: %s", err.Error())
			c.JSON(500, gin.H{
				"error": models.ErrInternalServerError.Error(),
			})
			break
		}
		h.logger.Infof("Response from Server Streaming RPC: %s", msg.GetResult())
	}
	h.logger.Infof("Finished Server Streaming RPC")
}

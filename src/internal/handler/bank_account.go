package handler

import (
	"bank_api/src/internal/models/db"
	"bank_api/src/internal/models/dto"
	"bank_api/src/internal/service"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @BasePath /api/v1

type Handler struct {
	ctx     context.Context
	account *service.BankAccountService
}

func NewHandler(ctx context.Context, account *service.BankAccountService) *Handler {
	return &Handler{
		ctx,
		account,
	}
}

// UpdateBalance godoc
// @Summary Пополнение баланса пользователя
// @Tags bank_account
// @Accept  json
// @Produce json
// @Param updateInfo body models.UpdateInfo true "Данные пополнения баланса"
// @Success 200
// @Router /bank_account/update [put]
func (h *Handler) UpdateBalance(c *gin.Context) {
	var info dto.UpdateInfo
	if err := json.NewDecoder(c.Request.Body).Decode(&info); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.account.Update(c, info); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

// Transfer godoc
// @Summary Перевод денежных средств другому пользователю
// @Tags bank_account
// @Accept json
// @Produce json
// @Param transferInfo body models.TransferInfo true "Данные перевода денежных средств"
// @Success 200
// @Router /bank_account/transfer [post]
func (h *Handler) Transfer(c *gin.Context) {
	var info dto.TransferInfo
	if err := json.NewDecoder(c.Request.Body).Decode(&info); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.account.Transfer(c, info); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

// GetLastOperations godoc
// @Summary Получение 10 последних операций пользователя
// @Tags bank_account
// @Accept json
// @Produce json
// @Param   id path int	true "Account ID"
// @Success 200 {array} entities.Operation
// @Router /bank_account/{id}/operations [get]
func (h *Handler) GetLastOperations(c *gin.Context) {
	const COUNT = 10

	var id = c.Param("id")
	var aid, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusBadRequest)
		return
	}

	var result []db.Operation
	result, err = h.account.LastOperations(c, db.User{Id: aid}, COUNT)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
	h.writeJSONResponse(c.Writer, http.StatusOK, result)
}

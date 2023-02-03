package controller

import (
	"go-gin-swagger-test/app/db"
	"go-gin-swagger-test/app/httputil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

const service_name = "gin-swagger-service"

// CreateAccount godoc
//	@Summary Create a new account
//	@Description Add anew account to the persistence
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			account	body		model.AccountDTO	true	"Add account"
//	@Success		200		{object}	model.Account
//	@Failure		400		{object}	httputil.HTTPError
//	@Failure		404		{object}	httputil.HTTPError
//	@Failure		500		{object}	httputil.HTTPError
//	@Router			/accounts [post]
func (c *Controller) CreateAccount(ginContext *gin.Context) {
	newCtx, span := otel.Tracer(service_name).Start(ginContext.Request.Context(), "Create Account")
	defer span.End()

	var dto db.AccountDTO

	if err := ginContext.ShouldBindJSON(&dto); err != nil {
		httputil.NewError(ginContext, http.StatusBadRequest, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	if err := dto.IsNameValid(); err != nil {
		httputil.NewError(ginContext, http.StatusBadRequest, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	Account, err := c.Persistence.InsertAccount(newCtx, dto)

	if err != nil {
		httputil.NewError(ginContext, http.StatusBadRequest, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	span.SetAttributes(attribute.String("request.account", Account.String()))

	ginContext.JSON(http.StatusOK, Account)
}

// 	GetAccounts godoc
//	@Summary Get all accounts
//	@Description Get all accounsts present on the persistence
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Success		200		{array}	model.Account
//	@Failure		400		{object}	httputil.HTTPError
//	@Failure		404		{object}	httputil.HTTPError
//	@Failure		500		{object}	httputil.HTTPError
//	@Router			/accounts [get]
func (c *Controller) GetAccounts(ginContext *gin.Context) {
	newCtx, span := otel.Tracer(service_name).Start(ginContext.Request.Context(), "Get Accounts")

	accounts, err := c.Persistence.GetAllAccounts(newCtx)

	if err != nil {
		httputil.NewError(ginContext, http.StatusNotFound, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	ginContext.JSON(http.StatusOK, accounts)
}

// 	GetAccountById godoc
//	@Summary Get account giving a specific ID
//	@Description Get the account related wiith the id
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Account ID"
//	@Success		200		{object}	model.Account
//	@Failure		400		{object}	httputil.HTTPError
//	@Failure		404		{object}	httputil.HTTPError
//	@Failure		500		{object}	httputil.HTTPError
//	@Router			/accounts/{id} [get]
func (c *Controller) GetAccountById(ginContext *gin.Context) {
	newCtx, span := otel.Tracer(service_name).Start(ginContext.Request.Context(), "Get Account by ID")

	id := ginContext.Param("id")
	convertedId, err := strconv.Atoi(id)
	span.SetAttributes(attribute.String("request.id", id))

	if err != nil {
		httputil.NewError(ginContext, http.StatusBadRequest, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	fetchedAccount, err := c.Persistence.GetAccountById(newCtx, convertedId)
	span.SetAttributes(attribute.String("response.account", fetchedAccount.String()))

	if err != nil {
		httputil.NewError(ginContext, http.StatusBadRequest, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	ginContext.JSON(http.StatusOK, fetchedAccount)
}

// 	UpdateNameById godoc
//	@Summary Update the name of an account
//	@Description  Update the name of an account with a specific ID
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Account ID"
//	@Param			account	body		model.AccountDTO	true	"Update account"
//	@Success		200		{object}	model.Account
//	@Failure		400		{object}	httputil.HTTPError
//	@Failure		404		{object}	httputil.HTTPError
//	@Failure		500		{object}	httputil.HTTPError
//	@Router			/accounts/{id} [put]
func (c *Controller) UpdateNameById(ginContext *gin.Context) {
	newCtx, span := otel.Tracer(service_name).Start(ginContext.Request.Context(), "Update account by ID")

	// Getting ID
	id := ginContext.Param("id")
	convertedId, err := strconv.Atoi(id)
	span.SetAttributes(attribute.String("request.id", id))

	if err != nil {
		httputil.NewError(ginContext, http.StatusBadRequest, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	// Getting body update
	var dto db.AccountDTO

	if err := ginContext.ShouldBindJSON(&dto); err != nil {
		httputil.NewError(ginContext, http.StatusBadRequest, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	if err := dto.IsNameValid(); err != nil {
		httputil.NewError(ginContext, http.StatusBadRequest, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	err = c.Persistence.UpdateAccountById(newCtx, convertedId, dto)

	if err != nil {
		httputil.NewError(ginContext, http.StatusInternalServerError, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	Account, _ := c.Persistence.GetAccountById(newCtx, convertedId)
	span.SetAttributes(attribute.String("response.account", Account.String()))

	ginContext.JSON(http.StatusOK, Account)
}

// 	DeleteAccountById godoc
//	@Summary Delete an account
//	@Description  Delete an account by a speceficic ID
//	@Tags			accounts
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Account ID"
//	@Success		204		{object}	model.Account
//	@Failure		400		{object}	httputil.HTTPError
//	@Failure		404		{object}	httputil.HTTPError
//	@Failure		500		{object}	httputil.HTTPError
//	@Router			/accounts/{id} [delete]
func (c *Controller) DeleteAccountById(ginContext *gin.Context) {
	newCtx, span := otel.Tracer(service_name).Start(ginContext.Request.Context(), "Delete account by ID")

	// Getting ID
	id := ginContext.Param("id")
	convertedId, err := strconv.Atoi(id)
	span.SetAttributes(attribute.String("request.id", id))

	if err != nil {
		httputil.NewError(ginContext, http.StatusBadRequest, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	err = c.Persistence.DeleteAccountById(newCtx, convertedId)

	if err != nil {
		httputil.NewError(ginContext, http.StatusNotFound, err)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return
	}

	ginContext.JSON(http.StatusNoContent, gin.H{})
}

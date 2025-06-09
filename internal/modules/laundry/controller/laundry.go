package controller

import (
	"github.com/audricimanuel/errorutils"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/model"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/modules/laundry/service"
	"github.com/audricimanuel/laundry-routine-tracking-service/utils"
	"github.com/audricimanuel/laundry-routine-tracking-service/utils/constants"
	"github.com/audricimanuel/laundry-routine-tracking-service/utils/httputils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type (
	LaundryController interface {
		GetLaundryList(ctx *gin.Context)
		AddLaundry(ctx *gin.Context)
	}

	LaundryControllerImpl struct {
		laundryService service.LaundryService
	}
)

func NewLaundryController(laundryService service.LaundryService) LaundryController {
	return &LaundryControllerImpl{
		laundryService: laundryService,
	}
}

func (l *LaundryControllerImpl) GetLaundryList(ctx *gin.Context) {
	userDataCtx, ok := ctx.Get(constants.USER_DATA)
	if !ok {
		httputils.InvalidateCookie(ctx, constants.COOKIE_AUTH_TOKEN)
		ctx.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	userData := userDataCtx.(model.UserClaims)

	queryParam := model.LaundryQueryParam{
		CategoryName:    ctx.Query("category_name"),
		LaundryDateFrom: nil,
		LaundryDateTo:   nil,
		Page:            utils.ConvertStrToInt(strings.TrimSpace(ctx.Query("page")), 1),
	}

	if laundryDateFromStr := strings.TrimSpace(ctx.Query("laundry_date_from")); laundryDateFromStr != "" {
		timeObj, err := time.Parse(constants.FORMAT_DATE_DEFAULT, laundryDateFromStr)
		if err == nil {
			queryParam.LaundryDateFrom = &timeObj
		}
	}

	if laundryDateToStr := strings.TrimSpace(ctx.Query("laundry_date_to")); laundryDateToStr != "" {
		timeObj, err := time.Parse(constants.FORMAT_DATE_DEFAULT, laundryDateToStr)
		if err == nil {
			queryParam.LaundryDateTo = &timeObj
		}
	}

	dataHtml := gin.H{}

	result, err := l.laundryService.GetLaundryList(ctx, queryParam, userData.UserId)
	if err != nil {
		dataHtml["error"] = err.Error()
	}
	dataHtml["data"] = result

	ctx.HTML(http.StatusOK, "dashboard.html", dataHtml)
}

func (l *LaundryControllerImpl) AddLaundry(ctx *gin.Context) {
	userDataCtx, ok := ctx.Get(constants.USER_DATA)
	if !ok {
		httputils.InvalidateCookie(ctx, constants.COOKIE_AUTH_TOKEN)
		ctx.Redirect(http.StatusTemporaryRedirect, "/login")
		return
	}

	var request model.AddLaundryRequest
	if err := errorutils.ValidatePayload(ctx.Request, &request); err != nil {
		httputils.SetHttpResponse(ctx, nil, err, nil)
		return
	}

	userData := userDataCtx.(model.UserClaims)

	httputils.SetHttpResponse(ctx, nil, l.laundryService.AddLaundry(ctx, userData.UserId, request), nil)
}

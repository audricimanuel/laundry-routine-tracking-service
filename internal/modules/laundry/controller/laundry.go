package controller

import (
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/model"
	"github.com/audricimanuel/laundry-routine-tracking-service/internal/modules/laundry/service"
	"github.com/audricimanuel/laundry-routine-tracking-service/utils"
	"github.com/audricimanuel/laundry-routine-tracking-service/utils/constants"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type (
	LaundryController interface {
		GetLaundryList(ctx *gin.Context)
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
	userDataCtx, _ := ctx.Get(constants.USER_DATA)
	userData := userDataCtx.(model.UserClaims)

	queryParam := model.LaundryQueryParam{
		CategoryName:    ctx.Query("category_name"),
		LaundryDateFrom: nil,
		LaundryDateTo:   nil,
		DetailNumber:    strings.TrimSpace(ctx.Query("detail_number")),
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
	result = append(result, model.LaundryResponse{
		Id:                "test123",
		DetailNumber:      "25060700001",
		Title:             "Laundry Rukita",
		LaundryDateString: "2025-06-07",
		TotalItems:        5,
		StatusLabel:       "Pending",
	})
	dataHtml["data"] = result

	ctx.HTML(http.StatusOK, "dashboard.html", dataHtml)
}

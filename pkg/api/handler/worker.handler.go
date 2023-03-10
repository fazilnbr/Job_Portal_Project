package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/fazilnbr/project-workey/pkg/common/response"
	"github.com/fazilnbr/project-workey/pkg/domain"
	services "github.com/fazilnbr/project-workey/pkg/usecase/interface"
	"github.com/fazilnbr/project-workey/pkg/utils"
	"github.com/gin-gonic/gin"
)

type WorkerHandler struct {
	workerService services.WorkerUseCase
}

func NewWorkerHandler(workerService services.WorkerUseCase) WorkerHandler {
	return WorkerHandler{
		workerService: workerService,
	}
}

// @Summary Add profile for Worker
// @ID worker add profile
// @Tags Worker Profile Management
// @Security BearerAuth
// @Produce json
// @Param WorkerAddProfile body domain.Profile{} true "Worker Add Profile"
// @Success 200 {object} response.Response{}
// @Failure 422 {object} response.Response{}
// @Router /worker/add-profile [post]
func (cr *WorkerHandler) WorkerAddProfile(c *gin.Context) {
	id, _ := strconv.Atoi(c.Writer.Header().Get("id"))

	fmt.Printf("\n\nidea : %v\n\n", id)
	var userprofile domain.Profile

	err:=c.Bind(&userprofile)
	if err != nil {
		response := response.ErrorResponse("Failed to create user", err.Error(), nil)
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusBadRequest)
		utils.ResponseJSON(*c, response)
		return
	}

	err = cr.workerService.AddProfile(userprofile, id)

	if err != nil {
		response := response.ErrorResponse("Error while adding worker profile", err.Error(), nil)
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)

		utils.ResponseJSON(*c, response)
		return
	}

	response := response.SuccessResponse(true, "SUCCESS", userprofile)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(*c, response)
}

// @Summary Edit profile for Worker
// @ID worker edit profile
// @Tags Worker Profile Management
// @Security BearerAuth
// @Produce json
// @Param WorkerEditProfile body domain.Profile{} true "Worker Edit Profile"
// @Success 200 {object} response.Response{}
// @Failure 422 {object} response.Response{}
// @Router /worker/edit-profile [patch]
func (cr *WorkerHandler) WorkerEditProfile(c *gin.Context) {
	id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	fmt.Printf("\n\n%v\n\n", id)
	var userprofile domain.Profile

	err:=c.Bind(&userprofile)
	if err != nil {
		response := response.ErrorResponse("Failed to create user", err.Error(), nil)
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusBadRequest)
		utils.ResponseJSON(*c, response)
		return
	}

	fmt.Printf("\n\nuser Profile : \n%v\n\n%v\n\n", userprofile, id)

	err = cr.workerService.WorkerEditProfile(userprofile, id)

	if err != nil {
		response := response.ErrorResponse("Error while editing worker profile", err.Error(), nil)
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)

		utils.ResponseJSON(*c, response)
		return
	}

	response := response.SuccessResponse(true, "SUCCESS", userprofile)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(*c, response)
}

// @Summary Change Password for worker
// @ID worker change password
// @Tags Worker Profile Management
// @Security BearerAuth
// @Produce json
// @Param WorekerChangePassword body domain.ChangePassword{} true "Woreker Change Password"
// @Success 200 {object} response.Response{}
// @Failure 422 {object} response.Response{}
// @Router /worker/change-password [patch]
func (cr *WorkerHandler) WorkerChangePassword(c *gin.Context) {

	id, _ := strconv.Atoi(c.Writer.Header().Get("id"))

	fmt.Println("id : ", id)
	var changepassword domain.ChangePassword

	err := c.Bind(&changepassword)
	if err != nil {
		response := response.ErrorResponse("Failed to create user", err.Error(), nil)
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusBadRequest)
		utils.ResponseJSON(*c, response)
		return
	}

	err = cr.workerService.WorkerVerifyPassword(changepassword, id)

	if err != nil {
		response := response.ErrorResponse("Wrong Email id or Password", err.Error(), nil)
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)

		utils.ResponseJSON(*c, response)
		return
	}
	err = cr.workerService.WorkerChangePassword(changepassword.NewPassword, id)

	if err != nil {
		response := response.ErrorResponse("Error while changing Password", err.Error(), nil)
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)

		utils.ResponseJSON(*c, response)
		return
	}

	changepassword.NewPassword = ""
	changepassword.OldPassword = ""
	response := response.SuccessResponse(true, "SUCCESS", changepassword)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(*c, response)
}

// @Summary list all job categories for Worker
// @ID list all job category for worker
// @Tags Worker Job Management
// @Security BearerAuth
// @Produce json
// @Param        page   query      string  true  "Page : "
// @Param        pagesize   query      string  true  "Pagesize : "
// @Success 200 {object} response.Response{}
// @Failure 422 {object} response.Response{}
// @Router /worker/list-job-category [get]
func (cr *WorkerHandler) ListJobCategoryUser(c *gin.Context) {

	page, err := strconv.Atoi(c.Query("page"))

	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	fmt.Printf("\n\nuser : %v\n\nmetea : %v\n\n", page, c.Query("page"))
	log.Println(page, "   ", pageSize)

	pagenation := utils.Filter{
		Page:     page,
		PageSize: pageSize,
	}

	categories, metadata, err := cr.workerService.ListJobCategoryUser(pagenation)

	if err != nil {
		response := response.ErrorResponse("Failed To List Job Category of worker", err.Error(), nil)

		c.Writer.WriteHeader(http.StatusUnprocessableEntity)

		utils.ResponseJSON(*c, response)
		return
	}

	// result := struct {
	// 	jobcategory *[]domain.Category
	// 	Meta        *utils.Metadata
	// }{
	// 	jobcategory: categories,
	// 	Meta:        &metadata,
	// }
	result2 := []interface{}{categories, metadata}

	response := response.SuccessResponse(true, "SUCCESS", result2)
	fmt.Println(response)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(*c, response)
}

// @Summary Add job for Worker
// @ID worker add job
// @Tags Worker Job Management
// @Security BearerAuth
// @Produce json
// @Param WorkerAddJob body domain.Job{} true "Worker Add job"
// @Success 200 {object} response.Response{}
// @Failure 422 {object} response.Response{}
// @Router /worker/add-job [post]
func (cr *WorkerHandler) AddJob(c *gin.Context) {
	id, _ := strconv.Atoi(c.Writer.Header().Get("id"))

	fmt.Printf("\n\nidea : %v\n\n", id)
	var workerjob domain.Job

	// err:=c.Bind(&workerjob)
	err:=c.ShouldBind(&workerjob)
	if err != nil {
		response := response.ErrorResponse("Failed to create user", err.Error(), nil)
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusBadRequest)
		utils.ResponseJSON(*c, response)
		return
	}

	workerjob.IdWorker = id

	_, err = cr.workerService.AddJob(workerjob)

	if err != nil {
		response := response.ErrorResponse("Error while adding worker job", err.Error(), nil)
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)

		utils.ResponseJSON(*c, response)
		return
	}

	response := response.SuccessResponse(true, "SUCCESS", workerjob)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(*c, response)
}

// @Summary list all jobs for worker
// @ID list all job jobs for worker
// @Tags Worker Job Management
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.Response{}
// @Failure 422 {object} response.Response{}
// @Router /worker/view-job [get]
func (cr *WorkerHandler) ViewJob(c *gin.Context) {

	id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	// page, err := strconv.Atoi(c.Query("page"))

	// pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	// fmt.Printf("\n\nuser : %v\n\nmetea : %v\n\n", page, c.Query("page"))
	// log.Println(page, "   ", pageSize)

	// pagenation := utils.Filter{
	// 	Page:     page,
	// 	PageSize: pageSize,
	// }

	jobs, err := cr.workerService.ViewJob(id)

	if err != nil {
		response := response.ErrorResponse("Failed to list workers jobs", err.Error(), nil)

		c.Writer.WriteHeader(http.StatusUnprocessableEntity)

		utils.ResponseJSON(*c, response)
		return
	}

	// result := struct {
	// 	Users *[]domain.UserResponse
	// 	Meta  *utils.Metadata
	// }{
	// 	Users: users,
	// 	Meta:  metadata,
	// }

	response := response.SuccessResponse(true, "SUCCESS", jobs)

	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(*c, response)
}

// @Summary Delete for Worker
// @ID worker delete job
// @Tags Worker Job Management
// @Security BearerAuth
// @Produce json
// @Param        jobid   query      string  true  "Job Id : "
// @Success 200 {object} response.Response{}
// @Failure 422 {object} response.Response{}
// @Router /worker/delete-job [delete]
func (cr *WorkerHandler) DeleteJob(c *gin.Context) {
	// id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	id, _ := strconv.Atoi(c.Query("jobid"))

	// c.Bind(&userprofile)

	fmt.Printf("\n\nuser Profile : \n%v\n\n\n\n", id)

	err := cr.workerService.DeleteJob(id)

	if err != nil {
		response := response.ErrorResponse("Error while deleting worker job", err.Error(), nil)
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)

		utils.ResponseJSON(*c, response)
		return
	}

	response := response.SuccessResponse(true, "SUCCESS", id)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(*c, response)
}

// @Summary list all pending job requset from user for Worker
// @ID list all pending job requset from user for worker
// @Tags Worker User Job Request
// @Security BearerAuth
// @Produce json
// @Param        page   query      string  true  "Page : "
// @Param        pagesize   query      string  true  "Pagesize : "
// @Success 200 {object} response.Response{}
// @Failure 422 {object} response.Response{}
// @Router /worker/list-user-pending-job-request [get]
func (cr *WorkerHandler) ListPendingJobRequsetFromUser(c *gin.Context) {

	id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	// id := 6
	page, err := strconv.Atoi(c.Query("page"))

	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	fmt.Printf("\n\nuser : %v\n\nmetea : %v\n\n", page, c.Query("page"))
	log.Println(page, "   ", pageSize)

	pagenation := utils.Filter{
		Page:     page,
		PageSize: pageSize,
	}

	requests, metadata, err := cr.workerService.ListPendingJobRequsetFromUser(pagenation, id)

	if err != nil {
		response := response.ErrorResponse("Failed To List Pending Job Requests of worker", err.Error(), nil)

		c.Writer.WriteHeader(http.StatusUnprocessableEntity)

		utils.ResponseJSON(*c, response)
		return
	}

	// result := struct {
	// 	jobrequest *[]domain.RequestResponse
	// 	Meta       *utils.Metadata
	// }{
	// 	jobrequest: requests,
	// 	Meta:       metadata,
	// }

	result2 := []interface{}{requests, metadata}

	response := response.SuccessResponse(true, "SUCCESS", result2)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(*c, response)
}

// @Summary list all accepted job requset from user for Worker
// @ID list all accepted job requset from user for worker
// @Tags Worker User Job Request
// @Security BearerAuth
// @Produce json
// @Param        page   query      string  true  "Page : "
// @Param        pagesize   query      string  true  "Pagesize : "
// @Success 200 {object} response.Response{}
// @Failure 422 {object} response.Response{}
// @Router /worker/list-user-accepted-job-request [get]
func (cr *WorkerHandler) ListAcceptedJobRequsetFromUser(c *gin.Context) {

	id, _ := strconv.Atoi(c.Writer.Header().Get("id"))
	// id := 6
	page, err := strconv.Atoi(c.Query("page"))

	pageSize, _ := strconv.Atoi(c.Query("pagesize"))
	fmt.Printf("\n\nuser : %v\n\nmetea : %v\n\n", page, c.Query("page"))
	log.Println(page, "   ", pageSize)

	pagenation := utils.Filter{
		Page:     page,
		PageSize: pageSize,
	}

	requests, metadata, err := cr.workerService.ListAcceptedJobRequsetFromUser(pagenation, id)

	if err != nil {
		response := response.ErrorResponse("Failed To List Accepted Job Requests of worker", err.Error(), nil)

		c.Writer.WriteHeader(http.StatusUnprocessableEntity)

		utils.ResponseJSON(*c, response)
		return
	}

	// result := struct {
	// 	jobrequest *[]domain.RequestResponse
	// 	Meta       *utils.Metadata
	// }{
	// 	jobrequest: requests,
	// 	Meta:       metadata,
	// }

	result2 := []interface{}{requests, metadata}

	response := response.SuccessResponse(true, "SUCCESS", result2)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(*c, response)
}

// @Summary accept job request from user
// @ID worker accept job request from user
// @Tags Worker User Job Request
// @Security BearerAuth
// @Produce json
// @Param        requestId   query      string  true  "RequestId : "
// @Success 200 {object} response.Response{}
// @Failure 422 {object} response.Response{}
// @Router /worker/accept-job-request [patch]
func (cr *WorkerHandler) AcceptJobRequest(c *gin.Context) {

	// id, _ := strconv.Atoi(c.Writer.Header().Get("id"))

	requestId, err := strconv.Atoi(c.Query("requestId"))

	err = cr.workerService.AcceptJobRequest(requestId)

	if err != nil {
		response := response.ErrorResponse("Failed to Accept Job-Request", err.Error(), nil)
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)

		utils.ResponseJSON(*c, response)
		return
	}

	response := response.SuccessResponse(true, "SUCCESS", requestId)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(*c, response)
}

// @Summary reject job request from user
// @ID worker reject job request from user
// @Tags Worker User Job Request
// @Security BearerAuth
// @Produce json
// @Param        requestId   query      string  true  "RequestId : "
// @Success 200 {object} response.Response{}
// @Failure 422 {object} response.Response{}
// @Router /worker/reject-job-request [patch]
func (cr *WorkerHandler) RejectJobRequest(c *gin.Context) {

	// id, _ := strconv.Atoi(c.Writer.Header().Get("id"))

	requestId, err := strconv.Atoi(c.Query("requestId"))

	err = cr.workerService.RejectJobRequest(requestId)

	if err != nil {
		response := response.ErrorResponse("Failed to Accept Job-Request", err.Error(), nil)
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusUnprocessableEntity)

		utils.ResponseJSON(*c, response)
		return
	}

	response := response.SuccessResponse(true, "SUCCESS", requestId)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(http.StatusOK)
	utils.ResponseJSON(*c, response)
}

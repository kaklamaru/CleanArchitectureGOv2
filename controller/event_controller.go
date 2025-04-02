package controller

import (
	"fmt"
	"go-clean-arch/pkg/utility"
	"strconv"

	// "go-clean-arch/structure/entity"
	"go-clean-arch/structure/request"
	"go-clean-arch/usecase"

	"github.com/gofiber/fiber/v2"
)

type EventController struct {
	eventUsecase usecase.EventUsecase
}

func NewEventController(eventUsecase usecase.EventUsecase) *EventController {
	return &EventController{
		eventUsecase: eventUsecase,
	}
}

func (c *EventController) CreateEvent(ctx *fiber.Ctx) error {
	var req request.EventRequest

	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "failed to retrieve claims",
		})
	}
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if err := c.eventUsecase.CreateEvent(&req, claims); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Create successfully",
	})

}

func (c *EventController) GetAllEvent(ctx *fiber.Ctx) error {
	events, err := c.eventUsecase.GetAllEvent()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to retrieve events",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(events)
}

func (c *EventController) GetEventByID(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}
	eventID := uint(id)
	event, err := c.eventUsecase.GetEventByID(eventID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(event)
}

func (c *EventController) ToggleEventStatus(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}
	eventID := uint(id)

	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "failed to retrieve claims",
		})
	}

	status, err := c.eventUsecase.ToggleEventStatus(eventID, claims)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// ตอบกลับค่าที่เปลี่ยนแปลง
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "Event status toggled successfully",
		"event_id":   eventID,
		"new_status": status,
	})
}

func (c *EventController) DeleteEventByID(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}
	eventID := uint(id)

	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "failed to retrieve claims",
		})
	}
	if err := c.eventUsecase.DeleteEventByID(eventID, claims); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"event_id": eventID,
		"massage":  "Event deleted successfully",
	})
}

func (c *EventController) UpdateEventByID(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}
	eventID := uint(id)

	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "failed to retrieve claims",
		})
	}
	var req request.EventRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}
	if err := c.eventUsecase.UpdateEventByID(eventID, claims, req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"event_id": eventID,
		"massage":  "Event updated successfully",
	})
}
func (c *EventController) MyEvent(ctx *fiber.Ctx) error {
	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "failed to retrieve claims",
		})
	}

	events, err := c.eventUsecase.MyEvent(claims)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to retrieve events",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(events)
}

func (c *EventController) AllAllowedEvent(ctx *fiber.Ctx) error {
	events, err := c.eventUsecase.AllAllowedEvent()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to retrieve events",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(events)
}
func (c *EventController) AllCurrentEvent(ctx *fiber.Ctx) error {
	events, err := c.eventUsecase.AllCurrentEvent()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Unable to retrieve events",
		})
	}
	return ctx.Status(fiber.StatusOK).JSON(events)
}
func (c *EventController) MyEventThisYear(ctx *fiber.Ctx) error {
	yearStr := ctx.Params("year")
	yearInt, err := strconv.Atoi(yearStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid year",
		})
	}
	year := uint(yearInt)
	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return err
	}

	insideEvents, outsideEvents, err := c.eventUsecase.MyEventThisYear(claims, year)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"inside_events":  insideEvents,
		"outside_events": outsideEvents,
	})
}


// Inside
func (c *EventController) JoinEvent(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}
	eventID := uint(id)

	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "failed to retrieve claims",
		})
	}

	if err := c.eventUsecase.JoinEvent(eventID, claims); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Joined event successfully",
	})
}

func (c *EventController) UnJoinEvent(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}
	eventID := uint(id)

	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "failed to retrieve claims",
		})
	}

	if err := c.eventUsecase.UnJoinEvent(eventID, claims); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Unjoined event successfully",
	})
}

func (c *EventController) UploadFile(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}
	eventID := uint(id)

	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "failed to retrieve claims",
		})
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "failed to get file",
		})
	}

	if err := c.eventUsecase.UploadFile(eventID, claims,file); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "file uploaded successfully",
	})
}

func (c *EventController) GetFile(ctx *fiber.Ctx) error {
	idStr := ctx.Params("eventid")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}
	eventID := uint(id)

	idStr = ctx.Params("userid")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}
	userID := uint(idInt)

	filePath, err := c.eventUsecase.GetFile(eventID, userID)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return ctx.SendFile(filePath, false)

}

func (c *EventController) MyChecklist(ctx *fiber.Ctx) error {
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}
	eventID := uint(id)

	claims, err := utility.GetClaimsFromContext(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Failed to get user claims",
		})
	}
	
	checklist, err := c.eventUsecase.MyChecklist(eventID, claims)
	if err != nil {
		return err
	}
	return ctx.Status(fiber.StatusOK).JSON(checklist)
}

func (c *EventController) UpdateEventStatusAndComment(ctx *fiber.Ctx) error {
	var req struct {
		Status  bool   `json:"status"`
		Comment string `json:"comment"`
	}
	idStr := ctx.Params("eventid")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}
	eventID := uint(id)

	idStr = ctx.Params("userid")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid UserID",
		})
	}
	userID := uint(idInt)

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	if err := c.eventUsecase.UpdateEventStatusAndComment(eventID, userID, req.Status, req.Comment); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Checking successfully",
	})
}


// Outside
func (c *EventController) CreateFile(ctx *fiber.Ctx) error{
	idStr := ctx.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id format",
		})
	}
	eventID := uint(id)

	data, fileName, err := c.eventUsecase.CreateFile(eventID)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Error: %v", err))
	}

	// กำหนดค่า HTTP headers เพื่อให้ผู้ใช้ดาวน์โหลดไฟล์
	ctx.Set("Content-Type", "application/pdf")
	ctx.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))

	return ctx.Send(data)
}

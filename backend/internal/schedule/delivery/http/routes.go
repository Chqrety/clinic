package http

import (
	"github.com/go-chi/chi/v5"
)

func RegisterPublicScheduleRoutes(r chi.Router, h *ScheduleHandler) {
	r.Get("/doctors", h.GetDoctors)
	r.Get("/doctors/{id}/schedules", h.GetDoctorSchedules)
	r.Get("/doctors/{id}/weekly", h.GetDoctorWeeklySchedules)
	r.Get("/doctors/schedules/summary", h.GetPublicSchedulesSummary)
}

func RegisterDoctorScheduleRoutes(r chi.Router, h *ScheduleHandler) {
	r.Post("/schedules/leave", h.AddDoctorLeave)
	r.Post("/schedules/leave/cancel", h.DeleteDoctorLeave)
	r.Get("/schedules/{id}/weekly", h.GetDoctorWeeklySchedules)
	r.Get("/schedules/{id}/leaves", h.GetDoctorLeaveList)
	r.Post("/schedules/update", h.UpdateDoctorSchedule)
}

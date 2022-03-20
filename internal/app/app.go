package app

import (
    "blog_service2/internal/controller/handler"
    "github.com/julienschmidt/httprouter"
)

func NewRouter(h *handler.Handler) *httprouter.Router {
    r := httprouter.New()
    r.GET("/api/v1/blog", h.GetRecords)
    r.GET("/api/v1/blog/:id", h.GetRecord)
    r.POST("/api/v1/blog", h.AddRecord)
    r.PUT("/api/v1/blog/:id", h.UpdateRecord)
    r.DELETE("/api/v1/blog/:id", h.DeleteRecord)

    return r
}

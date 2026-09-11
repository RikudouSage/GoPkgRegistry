package handlers

import (
	"GoPkgRepository/db"
	"GoPkgRepository/helper"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func GetPackageHandler(writer http.ResponseWriter, req *http.Request, database *db.Client) {
	pkg, err := database.FindPackageByImportPath(chi.URLParam(req, "*"))
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		helper.WriteJSON(writer, map[string]string{
			"error": "internal error",
		})
		return
	}
	if pkg == nil {
		writer.WriteHeader(http.StatusNotFound)
		helper.WriteJSON(writer, map[string]string{
			"error": "package not found",
		})
		return
	}

	helper.WriteJSON(writer, pkg)
}

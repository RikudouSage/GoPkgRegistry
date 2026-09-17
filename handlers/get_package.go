package handlers

import (
	"log"
	"net/http"
	"strconv"

	"go.chrastecky.dev/go-pkg-repository/db"
	"go.chrastecky.dev/go-pkg-repository/helper"

	"github.com/go-chi/chi/v5"
)

// GetPackageHandler writes the package identified by chi's wildcard URL
// parameter as JSON.
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

func GetPackageByIDHandler(writer http.ResponseWriter, request *http.Request, database *db.Client) {
	pkgId, err := strconv.Atoi(chi.URLParam(request, "id"))
	if err != nil {
		helper.WriteJSON(writer, map[string]string{
			"error": "invalid id",
		})
		return
	}

	pkg, err := database.FindPackageByID(pkgId)
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

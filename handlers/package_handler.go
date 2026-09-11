package handlers

import (
	"html/template"
	"log"
	"net/http"

	"go.chrastecky.dev/go-pkg-repository/cfg"
	"go.chrastecky.dev/go-pkg-repository/db"
	"go.chrastecky.dev/go-pkg-repository/helper"
	"go.chrastecky.dev/go-pkg-repository/tpl"

	clone "github.com/huandu/go-clone/generic"
)

// PackageHandler serves go-import metadata for the package named by the
// request URL. If configured, HostOverride is used to derive the import path.
func PackageHandler(
	writer http.ResponseWriter,
	req *http.Request,
	cfg *cfg.GlobalConfig,
	database *db.Client,
) {
	uri := clone.Clone(req.URL)
	if cfg.HostOverride != "" {
		uri.Host = cfg.HostOverride
	} else {
		uri.Host = req.Host
	}

	uri.RawQuery = ""
	uri.Fragment = ""

	packageName := uri.String()[2:]
	pkg, err := database.FindPackageByImportPath(packageName)
	if err != nil {
		log.Println(err)
		writer.WriteHeader(http.StatusInternalServerError)
		helper.WriteJSON(writer, map[string]string{"error": "Internal error"})
		return
	}

	if pkg == nil {
		writer.WriteHeader(http.StatusNotFound)
		helper.WriteJSON(writer, map[string]string{"error": "Package not found"})
		return
	}

	if !pkg.IsValid() {
		writer.WriteHeader(http.StatusInternalServerError)
		helper.WriteJSON(writer, map[string]string{"error": "Package is invalid"})
		return
	}

	writer.Header().Add("Content-Type", "text/html")
	writer.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")

	tmpl := template.Must(template.ParseFS(tpl.Templates, "package.html.tpl"))
	if err := tmpl.Execute(writer, pkg); err != nil {
		log.Println(err)
		helper.WriteJSON(writer, map[string]string{"error": "internal error"})
	}
}

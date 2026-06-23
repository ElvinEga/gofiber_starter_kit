package swaggerui

import (
	"html/template"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v3"
	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/swag"
)

type Config struct {
	URL          string
	InstanceName string
	DocExpansion string
	DomID        string
	DeepLinking  bool
}

func Handler() fiber.Handler {
	config := Config{
		URL:          "doc.json",
		InstanceName: swag.Name,
		DocExpansion: "list",
		DomID:        "swagger-ui",
		DeepLinking:  true,
	}
	index := template.Must(template.New("swagger_index.html").Parse(indexTemplate))

	return func(c fiber.Ctx) error {
		path := strings.TrimPrefix(c.Path(), "/swagger/")
		if path == "/swagger" || path == "" {
			return c.Redirect().Status(fiber.StatusMovedPermanently).To("/swagger/index.html")
		}

		ext := strings.TrimPrefix(filepath.Ext(path), ".")
		switch path {
		case "index.html":
			c.Type("html", "utf-8")
			return index.Execute(c, config)
		case "doc.json":
			doc, err := swag.ReadDoc(config.InstanceName)
			if err != nil {
				return c.Status(http.StatusInternalServerError).SendString(http.StatusText(http.StatusInternalServerError))
			}
			c.Type("json", "utf-8")
			return c.SendString(doc)
		default:
			content, err := swaggerFiles.ReadFile("/" + path)
			if err != nil {
				return c.SendStatus(fiber.StatusNotFound)
			}
			if ext != "" {
				c.Type(ext)
			}
			return c.Send(content)
		}
	}
}

const indexTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Swagger UI</title>
  <link rel="stylesheet" type="text/css" href="./swagger-ui.css">
  <link rel="icon" type="image/png" href="./favicon-32x32.png" sizes="32x32">
  <link rel="icon" type="image/png" href="./favicon-16x16.png" sizes="16x16">
  <style>html{box-sizing:border-box;overflow-y:scroll}*,*:before,*:after{box-sizing:inherit}body{margin:0;background:#fafafa}</style>
</head>
<body>
<div id="swagger-ui"></div>
<script src="./swagger-ui-bundle.js"></script>
<script src="./swagger-ui-standalone-preset.js"></script>
<script>
window.onload = function() {
  window.ui = SwaggerUIBundle({
    url: "{{.URL}}",
    deepLinking: {{.DeepLinking}},
    docExpansion: "{{.DocExpansion}}",
    dom_id: "#{{.DomID}}",
    validatorUrl: null,
    presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
    plugins: [SwaggerUIBundle.plugins.DownloadUrl],
    layout: "StandaloneLayout"
  });
};
</script>
</body>
</html>`

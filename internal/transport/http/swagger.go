package http

import (
	"html/template"
	"net/http"
	"os"
)

// HTML шаблон для Swagger UI
const swaggerUIHTML = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Market API Documentation</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui.css" />
    <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin:0; background: #fafafa; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@4.15.5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
                url: '/api/swagger.yaml',
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout",
                validatorUrl: null
            });
        };
    </script>
</body>
</html>
`

// RegisterSwaggerEndpoints регистрирует эндпоинты для Swagger документации
func RegisterSwaggerEndpoints(router *http.ServeMux) {
	// Swagger UI интерфейс
	router.HandleFunc("GET /docs", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		tmpl := template.Must(template.New("swagger").Parse(swaggerUIHTML))
		tmpl.Execute(w, nil)
	})

	// YAML файл спецификации
	router.HandleFunc("GET /api/swagger.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/x-yaml")
		data, err := os.ReadFile("swagger.yaml")
		if err != nil {
			http.Error(w, "Swagger spec not found", http.StatusNotFound)
			return
		}
		w.Write(data)
	})

	// JSON версия спецификации (для совместимости)
	router.HandleFunc("GET /api/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		// Здесь можно добавить конвертацию YAML в JSON если нужно
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"info": {"title": "Use /api/swagger.yaml for YAML version"}}`))
	})
}

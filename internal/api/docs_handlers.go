package api

import "net/http"

const swaggerHTML = `
	<!DOCTYPE html>
	<html lang="en">
	  <head>
	    <meta charset="UTF-8">
	    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	    <title>Swagger UI</title>
	    <link rel="stylesheet" type="text/css" href="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.0.0/swagger-ui.css" />
	    <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.0.0/swagger-ui-bundle.js"></script>
	    <script src="https://cdnjs.cloudflare.com/ajax/libs/swagger-ui/4.0.0/swagger-ui-standalone-preset.js"></script>
	  </head>
	  <body>
	    <div id="swagger-ui"></div>
	    <script>
	      const ui = SwaggerUIBundle({
	        url: '/docs/swagger.yaml',
	        dom_id: '#swagger-ui',
	        presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
	        layout: "StandaloneLayout"
	      });
	    </script>
	  </body>
	</html>

`

func docsUIHndler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(swaggerHTML))
}

func specHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "openapi3.sepc.yaml")
}

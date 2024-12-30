package pocketsql

import (
	_ "embed"
	"html/template"
)

//go:embed resources/templates/error.gotmpl.html
var errorHtml string

//go:embed resources/templates/exec.gotmpl.html
var execHtml string

//go:embed resources/templates/header.gotmpl.html
var headerHtml string

//go:embed resources/templates/modal.gotmpl.html
var modalHtml string

//go:embed resources/templates/table.gotmpl.html
var tableHtml string

var footerHtml = `</main>
</div><!-- /app -->
` + modalHtml + `
</body>
</html>`

var errorSqlStreamTemplate = template.Must(template.New("pocketsql.error.sql.stream").Parse(`<turbo-stream action="replace" target="error_sql">
  <template>
    <div class="query_error" id="error_sql">{{ .message }}</div>
  </template>
</turbo-stream>`))

var errorTemplate = template.Must(template.New("pocketsql.index").Parse(headerHtml + errorHtml + footerHtml))

var execStreamHtml = `<turbo-stream action="replace" target="main">
  <template>
    <main id="main">
` + execHtml + `
    </div>
  </template>
</turbo-stream>`

var execStreamTemplate = template.Must(template.New("pocketsql.exec.stream").Parse(execStreamHtml))

var tableTemplate = template.Must(template.New("pocketsql.table").Parse(headerHtml + tableHtml + footerHtml))

var tableStreamHtml = `<turbo-stream action="replace" target="main">
  <template>
    <main id="main">
` + tableHtml + `
    </div>
  </template>
</turbo-stream>

<turbo-stream action="replace" target="modals">
  <template>
` + modalHtml + `
  </template>
</turbo-stream>`

var tableStreamTemplate = template.Must(template.New("pocketsql.table.stream").Parse(tableStreamHtml))

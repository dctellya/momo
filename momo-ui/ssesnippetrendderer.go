package main

import (
	"bytes"
	"io"
	"regexp"

	render "github.com/go-echarts/go-echarts/v2/render"
)

type sseSnippetRenderer struct {
	c      interface{}
	before []func()
}

var ssepat = regexp.MustCompile(`(__f__")|("__f__)|(__f__)`)

var sseSnippetTpl = `
<div class="container" style="display: flex;justify-content: center;align-items: center;">
    <div class="item" id="{{ .ChartID }}" style="margin: auto;width:{{ .Initialization.Width }};height:{{ .Initialization.Height }};"></div>
</div>

<script type="text/javascript">
    "use strict";
	addChart("{{ .ChartID | safeJS }}", "{{ .Theme }}", "{{  .Initialization.Renderer }}" , {{ .JSONNotEscaped | safeJS }});
	{{- range  $listener := .EventListeners }}
		{{if .Query  }}
		charts[{{ .ChartID | safeJS }}].on({{ $listener.EventName }}, {{ $listener.Query | safeJS }}, {{ injectInstance $listener.Handler "%MY_ECHARTS%"  $.ChartID | safeJS }});
		{{ else }}
		charts[{{ .ChartID | safeJS }}].on({{ $listener.EventName }}, {{ injectInstance $listener.Handler "%MY_ECHARTS%"  $.ChartID | safeJS }})
		{{ end }}
	{{- end }}
	{{- range .JSFunctions.Fns }}
		{{ injectInstance . "%MY_ECHARTS%"  $.ChartID  | safeJS }}
	{{- end }}

	sseBinding("{{ .ChartID | safeJS }}",getSSENames());
</script>
`

// NewSnippetRenderer returns a render implementation for Chart.
func NewsSSESnippetRenderer(c interface{}, before ...func()) render.Renderer {
	return &sseSnippetRenderer{c: c, before: before}
}

// Render renders the chart into the given io.Writer.
func (r *sseSnippetRenderer) Render(w io.Writer) error {

	content := r.RenderContent()
	_, err := w.Write(content)
	return err
}

func (r *sseSnippetRenderer) RenderContent() []byte {
	const tplName = "chart"
	for _, fn := range r.before {
		fn()
	}

	contents := []string{sseSnippetTpl}
	tpl := render.MustTemplate(tplName, contents)

	var buf bytes.Buffer
	if err := tpl.ExecuteTemplate(&buf, tplName, r.c); err != nil {
		panic(err)
	}

	return ssepat.ReplaceAll(buf.Bytes(), []byte(""))
}

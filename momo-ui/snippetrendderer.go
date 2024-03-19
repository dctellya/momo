package main

import (
	"bytes"
	"io"
	"regexp"

	render "github.com/go-echarts/go-echarts/v2/render"
)

type snippetRenderer struct {
	c      interface{}
	before []func()
}

var pat = regexp.MustCompile(`(__f__")|("__f__)|(__f__)`)

var snippetTpl = `
<div class="container" style="display: flex;justify-content: center;align-items: center;">
    <div class="item" id="{{ .ChartID }}" style="margin: auto;width:{{ .Initialization.Width }};height:{{ .Initialization.Height }};"></div>
</div>

<script type="text/javascript">
    "use strict";
	function init_{{ .ChartID | safeJS }}(){
    	let goecharts_{{ .ChartID | safeJS }} = echarts.init(document.getElementById('{{ .ChartID | safeJS }}'), "{{ .Theme }}", { renderer: "{{  .Initialization.Renderer }}" });
    	let option_{{ .ChartID | safeJS }} = {{ .JSONNotEscaped | safeJS }};
    
		goecharts_{{ .ChartID | safeJS }}.setOption(option_{{ .ChartID | safeJS }});

		{{- range  $listener := .EventListeners }}
			{{if .Query  }}
			goecharts_{{ $.ChartID | safeJS }}.on({{ $listener.EventName }}, {{ $listener.Query | safeJS }}, {{ injectInstance $listener.Handler "%MY_ECHARTS%"  $.ChartID | safeJS }});
			{{ else }}
			goecharts_{{ $.ChartID | safeJS }}.on({{ $listener.EventName }}, {{ injectInstance $listener.Handler "%MY_ECHARTS%"  $.ChartID | safeJS }})
			{{ end }}
		{{- end }}
		{{- range .JSFunctions.Fns }}
			{{ injectInstance . "%MY_ECHARTS%"  $.ChartID  | safeJS }}
		{{- end }}
	};
	init_{{ .ChartID | safeJS }}();
</script>
`

// NewSnippetRenderer returns a render implementation for Chart.
func NewSnippetRenderer(c interface{}, before ...func()) render.Renderer {
	return &snippetRenderer{c: c, before: before}
}

// Render renders the chart into the given io.Writer.
func (r *snippetRenderer) Render(w io.Writer) error {

	content := r.RenderContent()
	_, err := w.Write(content)
	return err
}

func (r *snippetRenderer) RenderContent() []byte {
	const tplName = "chart"
	for _, fn := range r.before {
		fn()
	}

	contents := []string{snippetTpl}
	tpl := render.MustTemplate(tplName, contents)

	var buf bytes.Buffer
	if err := tpl.ExecuteTemplate(&buf, tplName, r.c); err != nil {
		panic(err)
	}

	return pat.ReplaceAll(buf.Bytes(), []byte(""))
}

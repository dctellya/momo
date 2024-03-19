package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var graphNodes = []opts.GraphNode{
	{Name: "Node1"},
	{Name: "Node2"},
	{Name: "Node3"},
	{Name: "Node4"},
	{Name: "Node5"},
	{Name: "Node6"},
	{Name: "Node7"},
	{Name: "Node8"},
}

func genLinks() []opts.GraphLink {
	links := make([]opts.GraphLink, 0)
	for i := 0; i < len(graphNodes); i++ {
		for j := 0; j < len(graphNodes); j++ {
			links = append(links, opts.GraphLink{Source: graphNodes[i].Name, Target: graphNodes[j].Name})
		}
	}
	return links
}

func graphBase() *charts.Graph {
	graph := charts.NewGraph()
	graph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "basic graph example"}),
	)
	graph.AddSeries("graph", graphNodes, genLinks(),
		charts.WithGraphChartOpts(
			opts.GraphChart{Force: &opts.GraphForce{Repulsion: 8000}},
		),
	)
	return graph
}

func graphCircle() *charts.Graph {
	graph := charts.NewGraph()
	graph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Circular layout"}),
	)

	graph.AddSeries("graph", graphNodes, genLinks()).
		SetSeriesOptions(
			charts.WithGraphChartOpts(
				opts.GraphChart{
					Force:  &opts.GraphForce{Repulsion: 8000},
					Layout: "circular",
				}),
			charts.WithLabelOpts(opts.Label{Show: opts.Bool(true), Position: "right"}),
		)
	return graph
}

func graphNpmDep() *charts.Graph {
	graph := charts.NewGraph()
	graph.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "npm dependencies demo",
		}))

	f, err := ioutil.ReadFile("fixtures/npmdepgraph.json")
	if err != nil {
		panic(err)
	}

	type Data struct {
		Nodes []opts.GraphNode
		Links []opts.GraphLink
	}

	var data Data
	if err := json.Unmarshal(f, &data); err != nil {
		fmt.Println(err)
	}

	graph.AddSeries("graph", data.Nodes, data.Links).
		SetSeriesOptions(
			charts.WithGraphChartOpts(opts.GraphChart{
				Layout:             "none",
				Roam:               opts.Bool(true),
				FocusNodeAdjacency: opts.Bool(true),
			}),
			charts.WithEmphasisOpts(opts.Emphasis{
				Label: &opts.Label{
					Show:     opts.Bool(true),
					Color:    "black",
					Position: "left",
				},
			}),
			charts.WithLineStyleOpts(opts.LineStyle{
				Curveness: 0.3,
			}),
		)
	return graph
}

type GraphExamples struct{}

func (GraphExamples) LoadChart() []byte {
	//page := components.NewPage()
	//page.AddCharts(
	//	graphBase(),
	//	graphCircle(),
	//	graphNpmDep(),
	//)
	gb := graphBase()
	gb.Renderer = NewSnippetRenderer(gb, gb.Validate)
	gc := graphCircle()
	gc.Renderer = NewSnippetRenderer(gc, gc.Validate)
	gnd := graphNpmDep()
	gnd.Renderer = NewSnippetRenderer(gnd, gnd.Validate)
	/*
		f, err := os.Create("html/graph.html")
		if err != nil {
			panic(err)

		}
	*/
	buf := new(bytes.Buffer)
	gb.Render(buf)
	gc.Render(buf)
	gnd.Render(buf)

	return buf.Bytes()
}

func gaugeBase() *charts.Gauge {
	gauge := charts.NewGauge()
	gauge.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "basic Gauge example"}),
	)

	gauge.AddSeries("ProjectA", []opts.GaugeData{{Name: "Basic progress", Value: rand.Intn(50)}})
	return gauge
}

func gaugeSSE() *charts.Gauge {
	gauge := charts.NewGauge()
	gauge.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "SSE Gauge example"}),
	)

	gauge.AddSeries("ProjectSSE", []opts.GaugeData{{Name: "SSE progress", Value: rand.Intn(50)}})

	fn := fmt.Sprintf(`function getSSENames() {
		return ["ping","pong"];
	}`)
	fn2 := fmt.Sprintf(`charts["%s"]["%s_ping"]=function(e){
		console.log("ping:",e.data);
	}`, gauge.ChartID, gauge.ChartID)
	fn3 := fmt.Sprintf(`charts["%s"]["%s_pong"]=function(e){
		console.log("pong:",e.data);
	}`, gauge.ChartID, gauge.ChartID)

	gauge.AddJSFuncs(fn)
	gauge.AddJSFuncs(fn2)
	gauge.AddJSFuncs(fn3)

	return gauge
}

func gaugeTimer() *charts.Gauge {
	gauge := charts.NewGauge()
	gauge.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "javascript timer"}),
	)

	gauge.AddSeries("ProjectB", []opts.GaugeData{{Name: "Work progress", Value: rand.Intn(50)}})

	fn := fmt.Sprintf(`setInterval(function () {
			option_%s.series[0].data[0].value = (Math.random() * 100).toFixed(2) - 0;
			goecharts_%s.setOption(option_%s, true);
		}, 2000);`, gauge.ChartID, gauge.ChartID, gauge.ChartID)
	gauge.AddJSFuncs(fn)
	return gauge
}

func (GraphExamples) LoadGaugeChart() []byte {
	//page := components.NewPage()
	//page.AddCharts(
	//	graphBase(),
	//	graphCircle(),
	//	graphNpmDep(),
	//)
	//gb := gaugeBase()
	//gb.Renderer = NewSnippetRenderer(gb, gb.Validate)
	//gt := gaugeTimer()
	//gt.Renderer = NewSnippetRenderer(gt, gt.Validate)
	gsse := gaugeSSE()
	gsse.Renderer = NewsSSESnippetRenderer(gsse, gsse.Validate)
	/*
		f, err := os.Create("html/graph.html")
		if err != nil {
			panic(err)

		}
	*/
	buf := new(bytes.Buffer)
	gsse.Render(buf)

	return buf.Bytes()
}

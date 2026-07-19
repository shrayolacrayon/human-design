package bodygraph

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"humandesign/internal/calculator"
)

// point is an x,y coordinate in the bodygraph SVG coordinate space.
type point struct {
	X, Y float64
}

// gateAnchors maps each of the 64 gates to its position on the edge of its
// center, matching the traditional bodygraph layout.
var gateAnchors = map[int]point{
	// Head (top triangle)
	64: {228, 100}, 61: {260, 100}, 63: {292, 100},
	// Ajna (inverted triangle)
	47: {228, 152}, 24: {260, 152}, 4: {292, 152},
	17: {236, 184}, 43: {260, 196}, 11: {284, 184},
	// Throat (square)
	62: {234, 265}, 23: {260, 265}, 56: {286, 265},
	16: {221, 284}, 20: {221, 308},
	31: {234, 323}, 8: {260, 323}, 33: {283, 323},
	35: {299, 282}, 12: {299, 304}, 45: {300, 328},
	// G (diamond)
	7: {240, 382}, 1: {260, 371}, 13: {280, 382},
	10: {213, 412}, 25: {297, 424},
	15: {234, 443}, 2: {260, 456}, 46: {286, 443},
	// Heart (small triangle)
	21: {356, 437}, 51: {341, 447}, 26: {349, 464}, 40: {381, 440},
	// Spleen (left triangle)
	48: {88, 522}, 57: {104, 536}, 44: {118, 550}, 50: {130, 562},
	32: {118, 580}, 28: {104, 594}, 18: {88, 608},
	// Solar Plexus (right triangle)
	36: {432, 522}, 22: {416, 536}, 37: {402, 550}, 6: {390, 562},
	49: {402, 580}, 55: {416, 594}, 30: {432, 608},
	// Sacral (square)
	5: {234, 561}, 14: {260, 561}, 29: {286, 561},
	34: {221, 582}, 27: {220, 604}, 59: {299, 592},
	42: {236, 614}, 3: {260, 614}, 9: {286, 614},
	// Root (square)
	53: {234, 668}, 60: {260, 668}, 52: {286, 668},
	54: {221, 688}, 38: {221, 708}, 58: {221, 728},
	19: {299, 688}, 39: {299, 708}, 41: {299, 728},
}

// channelCurves gives a quadratic Bezier control point for channels that
// curve around the outside of the body instead of cutting straight across.
var channelCurves = map[[2]int]point{
	{35, 36}: {430, 385}, // Throat right -> Solar Plexus, outer right
	{12, 22}: {402, 404},
	{16, 48}: {90, 385}, // Throat left -> Spleen, outer left
	{20, 57}: {118, 404},
	{26, 44}: {240, 514}, // Heart -> Spleen, sweeping under G
}

// channelPairs lists all 36 channels as gate pairs; each channel is drawn as
// two half-segments colored by that gate's activation.
var channelPairs = [][2]int{
	// Head - Ajna
	{64, 47}, {61, 24}, {63, 4},
	// Ajna - Throat
	{17, 62}, {43, 23}, {11, 56},
	// Throat - G
	{31, 7}, {8, 1}, {33, 13}, {20, 10},
	// Throat - Spleen / Sacral (integration)
	{16, 48}, {20, 57}, {20, 34},
	// Throat - Heart / Solar Plexus
	{45, 21}, {35, 36}, {12, 22},
	// G - Sacral / Heart / Spleen
	{15, 5}, {2, 14}, {46, 29}, {10, 34}, {25, 51}, {10, 57},
	// Heart - Spleen / Solar Plexus
	{26, 44}, {40, 37},
	// Sacral - Spleen / Solar Plexus
	{27, 50}, {34, 57}, {59, 6},
	// Sacral - Root
	{42, 53}, {3, 60}, {9, 52},
	// Spleen - Root
	{18, 58}, {28, 38}, {32, 54},
	// Solar Plexus - Root
	{30, 41}, {55, 39}, {49, 19},
}

// Traditional center colors used when a center is defined.
const (
	colorYellow  = "#f6df4e" // Head, G
	colorGreen   = "#57a05a" // Ajna
	colorBrown   = "#c19a6b" // Throat, Spleen, Solar Plexus, Root
	colorRed     = "#d94f3d" // Heart, Sacral
	colorDesign  = "#b03a2e" // design (red) activations
	colorPers    = "#26241f" // personality (black) activations
	colorTrack   = "#e8e3d7" // undefined channel track
	colorOutline = "#7a7264" // center outlines
)

// centerShape describes how to draw one of the nine centers.
type centerShape struct {
	Name   string
	Kind   string // "poly" or "rect"
	Points string // polygon points
	X, Y   float64
	W, H   float64
	Fill   string // fill when defined
}

var centerShapes = []centerShape{
	{Name: "Head", Kind: "poly", Points: "260,38 205,112 315,112", Fill: colorYellow},
	{Name: "Ajna", Kind: "poly", Points: "205,140 315,140 260,218", Fill: colorGreen},
	{Name: "Throat", Kind: "rect", X: 210, Y: 252, W: 100, H: 84, Fill: colorBrown},
	{Name: "G", Kind: "poly", Points: "260,358 318,412 260,466 202,412", Fill: colorYellow},
	{Name: "Heart", Kind: "poly", Points: "330,432 398,424 360,482", Fill: colorRed},
	{Name: "Spleen", Kind: "poly", Points: "60,490 60,640 150,565", Fill: colorBrown},
	{Name: "SolarPlexus", Kind: "poly", Points: "460,490 460,640 370,565", Fill: colorBrown},
	{Name: "Sacral", Kind: "rect", X: 210, Y: 548, W: 100, H: 78, Fill: colorRed},
	{Name: "Root", Kind: "rect", X: 210, Y: 660, W: 100, H: 84, Fill: colorBrown},
}

// Generator creates SVG body graphs
type Generator struct{}

// NewGenerator creates a new body graph generator
func NewGenerator() *Generator {
	return &Generator{}
}

// GenerateSVG creates an SVG visualization of a Human Design reading in the
// traditional bodygraph style: channels drawn as red (design) / black
// (personality) half-segments, striped when activated by both, with gate
// number badges on each center.
func (g *Generator) GenerateSVG(reading *calculator.Reading) (string, error) {
	pers := map[int]bool{}
	des := map[int]bool{}
	for _, gt := range reading.PersonalityGates {
		pers[gt.Number] = true
	}
	for _, gt := range reading.DesignGates {
		des[gt.Number] = true
	}

	var b strings.Builder
	b.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 520 760" class="bodygraph-svg">`)
	b.WriteString(`<defs><linearGradient id="splitPB" x1="0" y1="0" x2="1" y2="0">` +
		`<stop offset="50%" stop-color="` + colorPers + `"/>` +
		`<stop offset="50%" stop-color="` + colorDesign + `"/>` +
		`</linearGradient></defs>`)

	// Channels (drawn first so center shapes sit on top)
	for _, ch := range channelPairs {
		a1, a2 := gateAnchors[ch[0]], gateAnchors[ch[1]]
		if ctrl, ok := curveFor(ch); ok {
			// Split the quadratic Bezier at t=0.5 so each half can be
			// colored independently.
			c1 := point{(a1.X + ctrl.X) / 2, (a1.Y + ctrl.Y) / 2}
			c2 := point{(a2.X + ctrl.X) / 2, (a2.Y + ctrl.Y) / 2}
			mid := point{(a1.X + 2*ctrl.X + a2.X) / 4, (a1.Y + 2*ctrl.Y + a2.Y) / 4}
			writeCurveHalf(&b, a1, c1, mid, pers[ch[0]], des[ch[0]])
			writeCurveHalf(&b, a2, c2, mid, pers[ch[1]], des[ch[1]])
		} else {
			mid := point{(a1.X + a2.X) / 2, (a1.Y + a2.Y) / 2}
			writeChannelHalf(&b, a1, mid, pers[ch[0]], des[ch[0]])
			writeChannelHalf(&b, a2, mid, pers[ch[1]], des[ch[1]])
		}
	}

	// Centers
	for _, cs := range centerShapes {
		fill := "#ffffff"
		if c, ok := reading.Centers[cs.Name]; ok && c.Defined {
			fill = cs.Fill
		}
		if cs.Kind == "rect" {
			fmt.Fprintf(&b, `<rect x="%.0f" y="%.0f" width="%.0f" height="%.0f" rx="7" fill="%s" stroke="%s" stroke-width="1.5"/>`,
				cs.X, cs.Y, cs.W, cs.H, fill, colorOutline)
		} else {
			fmt.Fprintf(&b, `<polygon points="%s" fill="%s" stroke="%s" stroke-width="1.5"/>`,
				cs.Points, fill, colorOutline)
		}
	}

	// Gate number badges (iterate 1-64 for deterministic output)
	for gate := 1; gate <= 64; gate++ {
		p, ok := gateAnchors[gate]
		if !ok {
			continue
		}
		var fill, textColor, stroke string
		switch {
		case pers[gate] && des[gate]:
			fill, textColor, stroke = "url(#splitPB)", "#ffffff", "none"
		case pers[gate]:
			fill, textColor, stroke = colorPers, "#ffffff", "none"
		case des[gate]:
			fill, textColor, stroke = colorDesign, "#ffffff", "none"
		default:
			fill, textColor, stroke = "#ffffff", "#a49c8e", "#cfc8ba"
		}
		fmt.Fprintf(&b, `<circle cx="%.0f" cy="%.0f" r="8" fill="%s" stroke="%s" stroke-width="1"/>`,
			p.X, p.Y, fill, stroke)
		fmt.Fprintf(&b, `<text x="%.0f" y="%.0f" text-anchor="middle" font-family="Arial, sans-serif" font-size="8.5" font-weight="bold" fill="%s">%d</text>`,
			p.X, p.Y+3, textColor, gate)
	}

	b.WriteString(`</svg>`)
	return b.String(), nil
}

func curveFor(ch [2]int) (point, bool) {
	if p, ok := channelCurves[ch]; ok {
		return p, true
	}
	if p, ok := channelCurves[[2]int{ch[1], ch[0]}]; ok {
		return p, true
	}
	return point{}, false
}

// strokePasses returns the stroke color/dash passes for a half-channel:
// design red, personality black, candy-striped when both, faint track when
// inactive.
func strokePasses(p, d bool) [][2]string {
	switch {
	case p && d:
		return [][2]string{{colorDesign, ""}, {colorPers, ` stroke-dasharray="5,5"`}}
	case p:
		return [][2]string{{colorPers, ""}}
	case d:
		return [][2]string{{colorDesign, ""}}
	default:
		return [][2]string{{colorTrack, ""}}
	}
}

func writeChannelHalf(b *strings.Builder, from, to point, p, d bool) {
	for _, pass := range strokePasses(p, d) {
		fmt.Fprintf(b, `<line x1="%.1f" y1="%.1f" x2="%.1f" y2="%.1f" stroke="%s" stroke-width="6"%s/>`,
			from.X, from.Y, to.X, to.Y, pass[0], pass[1])
	}
}

func writeCurveHalf(b *strings.Builder, from, ctrl, to point, p, d bool) {
	for _, pass := range strokePasses(p, d) {
		fmt.Fprintf(b, `<path d="M %.1f %.1f Q %.1f %.1f %.1f %.1f" fill="none" stroke="%s" stroke-width="6"%s/>`,
			from.X, from.Y, ctrl.X, ctrl.Y, to.X, to.Y, pass[0], pass[1])
	}
}

// planetColumnOrder is the traditional Jovian Archive column ordering.
var planetColumnOrder = []string{
	"Sun", "Earth", "North Node", "South Node", "Moon", "Mercury",
	"Venus", "Mars", "Jupiter", "Saturn", "Uranus", "Neptune", "Pluto",
}

var planetSymbols = map[string]string{
	"Sun": "☉", "Moon": "☽", "Mercury": "☿", "Venus": "♀",
	"Mars": "♂", "Jupiter": "♃", "Saturn": "♄", "Uranus": "♅",
	"Neptune": "♆", "Pluto": "♇", "North Node": "☊",
	"South Node": "☋", "Earth": "⊕",
}

type columnRow struct {
	Symbol string
	Planet string
	Num    string
}

func buildColumn(gates []calculator.Gate) []columnRow {
	byPlanet := map[string]calculator.Gate{}
	for _, g := range gates {
		byPlanet[g.Planet] = g
	}
	rows := []columnRow{}
	for _, p := range planetColumnOrder {
		g, ok := byPlanet[p]
		if !ok {
			continue
		}
		sym := planetSymbols[p]
		if sym == "" {
			sym = "★"
		}
		rows = append(rows, columnRow{
			Symbol: sym,
			Planet: p,
			Num:    fmt.Sprintf("%d.%d", g.Number, g.Line),
		})
	}
	return rows
}

// GenerateHTML creates a complete HTML page with the body graph and reading details
func (g *Generator) GenerateHTML(reading *calculator.Reading) (string, error) {
	svg, err := g.GenerateSVG(reading)
	if err != nil {
		return "", err
	}

	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Human Design Chart</title>
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;
            background: #f1eee7;
            color: #2b2620;
            min-height: 100vh;
            padding: 40px 20px;
        }
        .page { max-width: 1080px; margin: 0 auto; }
        .chart-header { text-align: center; margin-bottom: 30px; }
        .chart-header h1 {
            font-size: 1.6rem;
            font-weight: 300;
            letter-spacing: 6px;
            text-transform: uppercase;
            color: #2b2620;
        }
        .chart-header .birth-info {
            margin-top: 8px;
            font-size: 0.9rem;
            color: #8b8478;
            letter-spacing: 1px;
        }
        .chart-header .type-line {
            margin-top: 12px;
            font-size: 1.05rem;
            letter-spacing: 3px;
            text-transform: uppercase;
            color: #b03a2e;
        }
        .card {
            background: #ffffff;
            border: 1px solid #e3ddd0;
            padding: 30px;
            margin-bottom: 24px;
        }
        .chart-flex {
            display: grid;
            grid-template-columns: 130px 1fr 130px;
            gap: 10px;
            align-items: start;
        }
        .bodygraph-svg { width: 100%; max-width: 470px; height: auto; display: block; margin: 0 auto; }
        .pcol { padding-top: 30px; }
        .pcol-header {
            font-size: 0.72rem;
            text-transform: uppercase;
            letter-spacing: 2.5px;
            padding-bottom: 10px;
            margin-bottom: 12px;
            border-bottom: 2px solid;
            text-align: center;
        }
        .pcol.design .pcol-header { color: #b03a2e; border-color: #b03a2e; }
        .pcol.personality .pcol-header { color: #26241f; border-color: #26241f; }
        .prow {
            display: flex;
            align-items: center;
            justify-content: center;
            gap: 8px;
            font-size: 0.98rem;
            line-height: 2.1;
        }
        .prow .glyph { width: 22px; text-align: center; font-size: 1.05rem; color: #8b8478; }
        .pcol.design .prow .num { color: #b03a2e; font-weight: 600; }
        .pcol.personality .prow .num { color: #26241f; font-weight: 600; }
        .legend {
            display: flex;
            justify-content: center;
            gap: 26px;
            margin-top: 20px;
            font-size: 0.78rem;
            color: #8b8478;
            letter-spacing: 1px;
        }
        .legend .swatch {
            display: inline-block;
            width: 22px; height: 6px;
            margin-right: 6px;
            vertical-align: middle;
        }
        .props {
            display: grid;
            grid-template-columns: 1fr 1fr;
            column-gap: 50px;
        }
        .prop-row {
            display: flex;
            justify-content: space-between;
            align-items: baseline;
            padding: 11px 2px;
            border-bottom: 1px solid #eee8dc;
            font-size: 0.95rem;
        }
        .prop-row .label {
            font-size: 0.72rem;
            text-transform: uppercase;
            letter-spacing: 2px;
            color: #8b8478;
        }
        .prop-row .value { text-align: right; font-weight: 500; }
        .prop-row.wide { grid-column: 1 / -1; }
        .section-title {
            font-size: 0.78rem;
            text-transform: uppercase;
            letter-spacing: 3px;
            color: #8b8478;
            margin-bottom: 16px;
            padding-bottom: 8px;
            border-bottom: 1px solid #e3ddd0;
        }
        .tables-flex { display: grid; grid-template-columns: 1fr 1fr; gap: 30px; }
        .gate-table { width: 100%; border-collapse: collapse; }
        .gate-table td {
            padding: 7px 6px;
            font-size: 0.88rem;
            border-bottom: 1px solid #f2eee5;
        }
        .gate-table tr:last-child td { border-bottom: none; }
        .planet-cell { white-space: nowrap; font-weight: 500; }
        .planet-sym { margin-right: 6px; color: #8b8478; }
        .gate-num-personality, .gate-num-design {
            padding: 2px 9px;
            border-radius: 10px;
            font-size: 0.8rem;
            color: white;
            white-space: nowrap;
        }
        .gate-num-personality { background: #26241f; }
        .gate-num-design { background: #b03a2e; }
        .gate-name-cell { color: #6d675c; }
        .lon-cell { color: #b3aca0; font-size: 0.76rem; text-align: right; }
        .channel-row {
            display: flex;
            align-items: center;
            gap: 14px;
            padding: 10px 2px;
            border-bottom: 1px solid #f2eee5;
            font-size: 0.92rem;
        }
        .channel-row:last-child { border-bottom: none; }
        .channel-key {
            min-width: 64px;
            text-align: center;
            padding: 3px 0;
            font-weight: 600;
            font-size: 0.82rem;
            color: white;
            background: linear-gradient(90deg, #26241f 50%, #b03a2e 50%);
            border-radius: 10px;
        }
        .back-link {
            display: inline-block;
            margin-bottom: 18px;
            font-size: 0.78rem;
            text-transform: uppercase;
            letter-spacing: 2px;
            color: #8b8478;
            text-decoration: none;
        }
        .back-link:hover { color: #b03a2e; }
        @media (max-width: 760px) {
            .chart-flex { grid-template-columns: 74px 1fr 74px; gap: 4px; }
            .prow { font-size: 0.82rem; gap: 4px; }
            .prow .glyph { width: 16px; font-size: 0.9rem; }
            .props, .tables-flex { grid-template-columns: 1fr; }
            .card { padding: 18px; }
        }
    </style>
</head>
<body>
    <div class="page">
        <a class="back-link" href="/">&larr; New Chart</a>
        <div class="chart-header">
            <h1>Human Design</h1>
            {{if .Location}}<div class="birth-info">{{.Location}} &middot; {{.BirthDateTime}}</div>{{end}}
            <div class="type-line">{{.Type}}</div>
        </div>

        <div class="card">
            <div class="chart-flex">
                <div class="pcol design">
                    <div class="pcol-header">Design</div>
                    {{range .DesignCol}}
                    <div class="prow"><span class="glyph">{{.Symbol}}</span><span class="num">{{.Num}}</span></div>
                    {{end}}
                </div>
                <div>{{.SVG}}</div>
                <div class="pcol personality">
                    <div class="pcol-header">Personality</div>
                    {{range .PersonalityCol}}
                    <div class="prow"><span class="glyph">{{.Symbol}}</span><span class="num">{{.Num}}</span></div>
                    {{end}}
                </div>
            </div>
            <div class="legend">
                <span><span class="swatch" style="background:#b03a2e"></span>Design (Unconscious)</span>
                <span><span class="swatch" style="background:#26241f"></span>Personality (Conscious)</span>
                <span><span class="swatch" style="background:repeating-linear-gradient(90deg,#b03a2e 0 5px,#26241f 5px 10px)"></span>Both</span>
            </div>
        </div>

        <div class="card">
            <div class="props">
                <div class="prop-row"><span class="label">Type</span><span class="value">{{.Type}}</span></div>
                <div class="prop-row"><span class="label">Strategy</span><span class="value">{{.Strategy}}</span></div>
                <div class="prop-row"><span class="label">Inner Authority</span><span class="value">{{.Authority}}</span></div>
                <div class="prop-row"><span class="label">Profile</span><span class="value">{{.Profile.Conscious}}/{{.Profile.Unconscious}} {{.Profile.Name}}</span></div>
                <div class="prop-row"><span class="label">Definition</span><span class="value">{{.Definition}}</span></div>
                <div class="prop-row"><span class="label">Not-Self Theme</span><span class="value">{{.NotSelfTheme}}</span></div>
                <div class="prop-row"><span class="label">Signature</span><span class="value">{{.Signature}}</span></div>
                <div class="prop-row"><span class="label">Incarnation Cross</span><span class="value">{{.IncarnationCross}}</span></div>
            </div>
        </div>

        <div class="card">
            <div class="section-title">Activations</div>
            <div class="tables-flex">
                <table class="gate-table">
                    {{range .PersonalityGates}}
                    <tr>
                        <td class="planet-cell"><span class="planet-sym">{{planetSymbol .Planet}}</span>{{.Planet}}</td>
                        <td><span class="gate-num-personality">{{.Number}}.{{.Line}}</span></td>
                        <td class="gate-name-cell">{{.Name}}</td>
                        <td class="lon-cell">{{printf "%.2f" .Longitude}}&deg;</td>
                    </tr>
                    {{end}}
                </table>
                <table class="gate-table">
                    {{range .DesignGates}}
                    <tr>
                        <td class="planet-cell"><span class="planet-sym">{{planetSymbol .Planet}}</span>{{.Planet}}</td>
                        <td><span class="gate-num-design">{{.Number}}.{{.Line}}</span></td>
                        <td class="gate-name-cell">{{.Name}}</td>
                        <td class="lon-cell">{{printf "%.2f" .Longitude}}&deg;</td>
                    </tr>
                    {{end}}
                </table>
            </div>
        </div>

        <div class="card">
            <div class="section-title">Defined Channels</div>
            {{range .Channels}}
            {{if .Defined}}
            <div class="channel-row">
                <span class="channel-key">{{.Gate1}}&ndash;{{.Gate2}}</span>
                <span>{{.Name}}</span>
            </div>
            {{end}}
            {{end}}
        </div>
    </div>
</body>
</html>`

	data := map[string]interface{}{
		"Type":             reading.Type,
		"Strategy":         reading.Strategy,
		"Authority":        reading.Authority,
		"Profile":          reading.Profile,
		"Definition":       reading.Definition,
		"Signature":        reading.Signature,
		"NotSelfTheme":     reading.NotSelfTheme,
		"IncarnationCross": reading.IncarnationCross,
		"PersonalityGates": reading.PersonalityGates,
		"DesignGates":      reading.DesignGates,
		"Channels":         reading.Channels,
		"PersonalityCol":   buildColumn(reading.PersonalityGates),
		"DesignCol":        buildColumn(reading.DesignGates),
		"Location":         reading.BirthData.Location,
		"BirthDateTime":    reading.BirthData.DateTime.Format("Jan 2, 2006 · 15:04 UTC"),
		"SVG":              template.HTML(svg),
	}

	funcMap := template.FuncMap{
		"planetSymbol": func(planet string) string {
			if s, ok := planetSymbols[planet]; ok {
				return s
			}
			return "★"
		},
		"printf": fmt.Sprintf,
	}
	t, err := template.New("reading").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

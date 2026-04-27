package bodygraph

import (
	"bytes"
	"fmt"
	"html/template"
	"humandesign/internal/calculator"
)

// CenterPosition defines the x, y position and size of a center
type CenterPosition struct {
	X      int
	Y      int
	Width  int
	Height int
	Shape  string // "triangle", "square", "diamond"
}

// ChannelPath defines the SVG path for a channel
type ChannelPath struct {
	Gate1   int
	Gate2   int
	Path    string
	Center1 string
	Center2 string
}

// Generator creates SVG body graphs
type Generator struct {
	centerPositions map[string]CenterPosition
	channelPaths    []ChannelPath
}

// NewGenerator creates a new body graph generator
func NewGenerator() *Generator {
	g := &Generator{
		centerPositions: map[string]CenterPosition{
			"Head":        {X: 200, Y: 30, Width: 60, Height: 50, Shape: "triangle"},
			"Ajna":        {X: 200, Y: 100, Width: 60, Height: 50, Shape: "triangle"},
			"Throat":      {X: 200, Y: 180, Width: 70, Height: 50, Shape: "square"},
			"G":           {X: 200, Y: 280, Width: 60, Height: 60, Shape: "diamond"},
			"Heart":       {X: 280, Y: 260, Width: 50, Height: 45, Shape: "triangle"},
			"Spleen":      {X: 100, Y: 340, Width: 50, Height: 60, Shape: "triangle"},
			"SolarPlexus": {X: 300, Y: 340, Width: 50, Height: 60, Shape: "triangle"},
			"Sacral":      {X: 200, Y: 380, Width: 70, Height: 50, Shape: "square"},
			"Root":        {X: 200, Y: 470, Width: 70, Height: 50, Shape: "square"},
		},
	}
	g.initChannelPaths()
	return g
}

func (g *Generator) initChannelPaths() {
	// Define paths connecting centers through their gates
	g.channelPaths = []ChannelPath{
		// Head to Ajna
		{64, 47, "M 200 80 L 200 100", "Head", "Ajna"},
		{63, 4, "M 220 80 L 220 100", "Head", "Ajna"},
		{61, 24, "M 180 80 L 180 100", "Head", "Ajna"},

		// Ajna to Throat
		{43, 23, "M 200 150 L 200 180", "Ajna", "Throat"},
		{17, 62, "M 220 150 L 220 180", "Ajna", "Throat"},
		{11, 56, "M 180 150 L 180 180", "Ajna", "Throat"},

		// Throat to G
		{31, 7, "M 200 230 L 200 280", "Throat", "G"},
		{8, 1, "M 185 230 L 185 280", "Throat", "G"},
		{33, 13, "M 215 230 L 215 280", "Throat", "G"},

		// Throat to Solar Plexus
		{35, 36, "M 235 210 Q 280 280 300 340", "Throat", "SolarPlexus"},
		{12, 22, "M 245 200 Q 290 270 310 340", "Throat", "SolarPlexus"},

		// Throat to Heart
		{45, 21, "M 250 200 L 280 260", "Throat", "Heart"},

		// Throat to Spleen
		{16, 48, "M 165 210 Q 120 280 100 340", "Throat", "Spleen"},
		{20, 57, "M 155 200 Q 110 270 100 330", "Throat", "Spleen"},

		// G to Sacral
		{15, 5, "M 200 340 L 200 380", "G", "Sacral"},
		{2, 14, "M 185 340 L 185 380", "G", "Sacral"},
		{10, 34, "M 175 320 L 170 380", "G", "Sacral"},
		{46, 29, "M 225 320 L 230 380", "G", "Sacral"},

		// G to Heart
		{25, 51, "M 240 290 L 280 270", "G", "Heart"},

		// G to Spleen
		{10, 57, "M 160 300 L 125 340", "G", "Spleen"},

		// Heart to Spleen
		{26, 44, "M 260 285 Q 180 320 125 340", "Heart", "Spleen"},

		// Heart to Solar Plexus
		{37, 40, "M 300 290 L 310 340", "Heart", "SolarPlexus"},

		// Sacral to Spleen
		{50, 27, "M 165 395 L 125 370", "Sacral", "Spleen"},
		{34, 57, "M 155 390 L 120 350", "Sacral", "Spleen"},

		// Sacral to Solar Plexus
		{59, 6, "M 245 395 L 290 370", "Sacral", "SolarPlexus"},

		// Sacral to Root
		{42, 53, "M 185 430 L 185 470", "Sacral", "Root"},
		{3, 60, "M 200 430 L 200 470", "Sacral", "Root"},
		{9, 52, "M 215 430 L 215 470", "Sacral", "Root"},

		// Spleen to Root
		{18, 58, "M 115 400 L 165 470", "Spleen", "Root"},
		{28, 38, "M 100 400 L 155 475", "Spleen", "Root"},
		{32, 54, "M 125 400 L 175 470", "Spleen", "Root"},

		// Solar Plexus to Root
		{36, 35, "M 300 400 L 245 470", "SolarPlexus", "Root"},
		{49, 19, "M 310 400 L 255 475", "SolarPlexus", "Root"},
		{55, 39, "M 290 400 L 235 470", "SolarPlexus", "Root"},
		{30, 41, "M 320 400 L 265 470", "SolarPlexus", "Root"},
	}
}

// GenerateSVG creates an SVG visualization of a Human Design reading
func (g *Generator) GenerateSVG(reading *calculator.Reading) (string, error) {
	tmpl := `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 400 550" width="400" height="550">
  <defs>
    <style>
      .center-undefined { fill: white; stroke: #333; stroke-width: 2; }
      .center-defined { fill: #FFD700; stroke: #333; stroke-width: 2; }
      .channel-undefined { stroke: #ddd; stroke-width: 3; fill: none; }
      .channel-defined { stroke: #8B0000; stroke-width: 4; fill: none; }
      .center-label { font-family: Arial, sans-serif; font-size: 10px; text-anchor: middle; fill: #333; }
      .gate-label { font-family: Arial, sans-serif; font-size: 8px; fill: #666; }
      .title { font-family: Arial, sans-serif; font-size: 14px; font-weight: bold; text-anchor: middle; }
    </style>
  </defs>
  
  <rect width="400" height="550" fill="#f8f9fa"/>
  
  <!-- Channels -->
  {{range .Channels}}
  <path d="{{.Path}}" class="{{if .Defined}}channel-defined{{else}}channel-undefined{{end}}"/>
  {{end}}
  
  <!-- Centers -->
  {{range $name, $center := .Centers}}
  {{$pos := index $.Positions $name}}
  {{if eq $pos.Shape "triangle"}}
  <polygon points="{{$pos.X}},{{minus $pos.Y 25}} {{minus $pos.X 30}},{{plus $pos.Y 25}} {{plus $pos.X 30}},{{plus $pos.Y 25}}" 
           class="{{if $center.Defined}}center-defined{{else}}center-undefined{{end}}"/>
  {{else if eq $pos.Shape "square"}}
  <rect x="{{minus $pos.X 35}}" y="{{minus $pos.Y 25}}" width="70" height="50" rx="5"
        class="{{if $center.Defined}}center-defined{{else}}center-undefined{{end}}"/>
  {{else if eq $pos.Shape "diamond"}}
  <polygon points="{{$pos.X}},{{minus $pos.Y 30}} {{minus $pos.X 30}},{{$pos.Y}} {{$pos.X}},{{plus $pos.Y 30}} {{plus $pos.X 30}},{{$pos.Y}}"
           class="{{if $center.Defined}}center-defined{{else}}center-undefined{{end}}"/>
  {{end}}
  <text x="{{$pos.X}}" y="{{plus $pos.Y 5}}" class="center-label">{{$name}}</text>
  {{end}}
  
</svg>`

	// Create template functions
	funcMap := template.FuncMap{
		"plus": func(a, b int) int {
			return a + b
		},
		"minus": func(a, b int) int {
			return a - b
		},
	}

	// Prepare channel data with defined status
	type ChannelData struct {
		Path    string
		Defined bool
	}
	channelData := []ChannelData{}
	
	definedChannels := make(map[string]bool)
	for _, ch := range reading.Channels {
		if ch.Defined {
			key := fmt.Sprintf("%d-%d", ch.Gate1, ch.Gate2)
			definedChannels[key] = true
			key2 := fmt.Sprintf("%d-%d", ch.Gate2, ch.Gate1)
			definedChannels[key2] = true
		}
	}

	for _, cp := range g.channelPaths {
		key := fmt.Sprintf("%d-%d", cp.Gate1, cp.Gate2)
		channelData = append(channelData, ChannelData{
			Path:    cp.Path,
			Defined: definedChannels[key],
		})
	}

	data := map[string]interface{}{
		"Centers":   reading.Centers,
		"Channels":  channelData,
		"Positions": g.centerPositions,
	}

	t, err := template.New("bodygraph").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("template parse error: %w", err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("template execute error: %w", err)
	}

	return buf.String(), nil
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
    <title>Human Design Reading</title>
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            padding: 20px;
        }
        .container { 
            max-width: 1200px; 
            margin: 0 auto;
            background: white;
            border-radius: 20px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            overflow: hidden;
        }
        .header {
            background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
            color: white;
            padding: 30px;
            text-align: center;
        }
        .header h1 { font-size: 2.5rem; margin-bottom: 10px; }
        .header .type-badge {
            display: inline-block;
            background: #FFD700;
            color: #1a1a2e;
            padding: 8px 20px;
            border-radius: 20px;
            font-weight: bold;
            font-size: 1.2rem;
        }
        .content { display: flex; flex-wrap: wrap; }
        .bodygraph-section {
            flex: 1;
            min-width: 300px;
            padding: 30px;
            display: flex;
            justify-content: center;
            align-items: center;
            background: #f8f9fa;
        }
        .details-section {
            flex: 1;
            min-width: 300px;
            padding: 30px;
        }
        .detail-card {
            background: #f8f9fa;
            border-radius: 10px;
            padding: 20px;
            margin-bottom: 15px;
        }
        .detail-card h3 {
            color: #764ba2;
            margin-bottom: 10px;
            font-size: 0.9rem;
            text-transform: uppercase;
            letter-spacing: 1px;
        }
        .detail-card p {
            font-size: 1.1rem;
            color: #333;
        }
        .gates-section {
            margin-top: 20px;
        }
        .gate-table { width: 100%; border-collapse: collapse; margin-top: 10px; }
        .gate-table td { padding: 6px 8px; font-size: 0.88rem; border-bottom: 1px solid #f0f0f0; }
        .gate-table tr:last-child td { border-bottom: none; }
        .planet-cell { white-space: nowrap; font-weight: 600; }
        .planet-sym { font-size: 1rem; margin-right: 4px; }
        .gate-num-personality {
            background: #333; color: white;
            padding: 2px 8px; border-radius: 10px; font-size: 0.82rem;
            white-space: nowrap;
        }
        .gate-num-design {
            background: #8B0000; color: white;
            padding: 2px 8px; border-radius: 10px; font-size: 0.82rem;
            white-space: nowrap;
        }
        .gate-name-cell { color: #555; font-style: italic; }
        .lon-cell { color: #aaa; font-size: 0.78rem; text-align: right; }
        .channels-section { margin-top: 20px; }
        .channel-item {
            background: linear-gradient(135deg, #FFD700 0%, #FFA500 100%);
            padding: 10px 15px;
            border-radius: 8px;
            margin: 5px 0;
            font-weight: 500;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>Human Design Chart</h1>
            <span class="type-badge">{{.Type}}</span>
        </div>
        <div class="content">
            <div class="bodygraph-section">
                {{.SVG}}
            </div>
            <div class="details-section">
                <div class="detail-card">
                    <h3>Strategy</h3>
                    <p>{{.Strategy}}</p>
                </div>
                <div class="detail-card">
                    <h3>Authority</h3>
                    <p>{{.Authority}}</p>
                </div>
                <div class="detail-card">
                    <h3>Profile</h3>
                    <p>{{.Profile.Conscious}}/{{.Profile.Unconscious}} - {{.Profile.Name}}</p>
                </div>
                <div class="detail-card">
                    <h3>Definition</h3>
                    <p>{{.Definition}}</p>
                </div>
                <div class="detail-card">
                    <h3>Signature & Not-Self Theme</h3>
                    <p>Signature: <strong>{{.Signature}}</strong></p>
                    <p>Not-Self: <strong>{{.NotSelfTheme}}</strong></p>
                </div>
                <div class="detail-card">
                    <h3>Incarnation Cross</h3>
                    <p>{{.IncarnationCross}}</p>
                </div>
                
                <div class="gates-section">
                    <div class="detail-card">
                        <h3>Personality Gates (Conscious &#9679; Black)</h3>
                        <table class="gate-table">
                            {{range .PersonalityGates}}
                            <tr>
                                <td class="planet-cell"><span class="planet-sym">{{planetSymbol .Planet}}</span>{{.Planet}}</td>
                                <td><span class="gate-num-personality">{{.Number}}.{{.Line}}</span></td>
                                <td class="gate-name-cell">{{.Name}}</td>
                                <td class="lon-cell">{{printf "%.2f" .Longitude}}°</td>
                            </tr>
                            {{end}}
                        </table>
                    </div>
                    <div class="detail-card">
                        <h3>Design Gates (Unconscious &#9679; Red)</h3>
                        <table class="gate-table">
                            {{range .DesignGates}}
                            <tr>
                                <td class="planet-cell"><span class="planet-sym">{{planetSymbol .Planet}}</span>{{.Planet}}</td>
                                <td><span class="gate-num-design">{{.Number}}.{{.Line}}</span></td>
                                <td class="gate-name-cell">{{.Name}}</td>
                                <td class="lon-cell">{{printf "%.2f" .Longitude}}°</td>
                            </tr>
                            {{end}}
                        </table>
                    </div>
                </div>

                <div class="channels-section">
                    <div class="detail-card">
                        <h3>Defined Channels</h3>
                        {{range .Channels}}
                        {{if .Defined}}
                        <div class="channel-item">{{.Gate1}}-{{.Gate2}} {{.Name}}</div>
                        {{end}}
                        {{end}}
                    </div>
                </div>
            </div>
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
		"SVG":              template.HTML(svg),
	}

	funcMap := template.FuncMap{
		"planetSymbol": func(planet string) string {
			symbols := map[string]string{
				"Sun": "☉", "Moon": "☽", "Mercury": "☿", "Venus": "♀",
				"Mars": "♂", "Jupiter": "♃", "Saturn": "♄", "Uranus": "♅",
				"Neptune": "♆", "Pluto": "♇", "North Node": "☊",
				"South Node": "☋", "Earth": "⊕",
			}
			if s, ok := symbols[planet]; ok {
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

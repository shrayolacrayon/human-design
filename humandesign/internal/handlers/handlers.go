package handlers

import (
	"encoding/json"
	"html/template"
	"humandesign/internal/bodygraph"
	"humandesign/internal/calculator"
	"log"
	"net/http"
	"time"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	calculator *calculator.Calculator
	bodygraph  *bodygraph.Generator
}

// NewHandler creates a new handler with all dependencies
func NewHandler() *Handler {
	return &Handler{
		calculator: calculator.NewCalculator(),
		bodygraph:  bodygraph.NewGenerator(),
	}
}

// ReadingRequest represents the JSON request for a reading
type ReadingRequest struct {
	DateTime  string  `json:"datetime"`  // ISO 8601 format
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Location  string  `json:"location"`
}

// HomePage serves the main input form
func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	tmpl := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Human Design Calculator</title>
    <style>
        * { box-sizing: border-box; margin: 0; padding: 0; }
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            min-height: 100vh;
            display: flex;
            justify-content: center;
            align-items: center;
            padding: 20px;
        }
        .container {
            background: white;
            padding: 40px;
            border-radius: 20px;
            box-shadow: 0 20px 60px rgba(0,0,0,0.3);
            max-width: 500px;
            width: 100%;
        }
        h1 {
            text-align: center;
            color: #333;
            margin-bottom: 10px;
            font-size: 2rem;
        }
        .subtitle {
            text-align: center;
            color: #666;
            margin-bottom: 30px;
        }
        .form-group {
            margin-bottom: 20px;
        }
        label {
            display: block;
            margin-bottom: 8px;
            color: #333;
            font-weight: 600;
        }
        input[type="date"],
        input[type="time"],
        input[type="text"],
        input[type="number"] {
            width: 100%;
            padding: 12px 15px;
            border: 2px solid #e0e0e0;
            border-radius: 10px;
            font-size: 1rem;
            transition: border-color 0.3s;
        }
        input:focus {
            outline: none;
            border-color: #764ba2;
        }
        .row {
            display: flex;
            gap: 15px;
        }
        .row .form-group {
            flex: 1;
        }
        .coordinates {
            display: flex;
            gap: 15px;
        }
        .coordinates .form-group {
            flex: 1;
        }
        button {
            width: 100%;
            padding: 15px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            color: white;
            border: none;
            border-radius: 10px;
            font-size: 1.1rem;
            font-weight: 600;
            cursor: pointer;
            transition: transform 0.2s, box-shadow 0.2s;
        }
        button:hover {
            transform: translateY(-2px);
            box-shadow: 0 5px 20px rgba(118, 75, 162, 0.4);
        }
        .info {
            margin-top: 20px;
            padding: 15px;
            background: #f8f9fa;
            border-radius: 10px;
            font-size: 0.9rem;
            color: #666;
        }
        .error {
            background: #fee;
            color: #c00;
            padding: 15px;
            border-radius: 10px;
            margin-bottom: 20px;
            display: none;
        }
        .loading {
            display: none;
            text-align: center;
            padding: 20px;
        }
        .spinner {
            border: 3px solid #f3f3f3;
            border-top: 3px solid #764ba2;
            border-radius: 50%;
            width: 40px;
            height: 40px;
            animation: spin 1s linear infinite;
            margin: 0 auto;
        }
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
    </style>
</head>
<body>
    <div class="container">
        <h1>✨ Human Design</h1>
        <p class="subtitle">Discover your unique energetic blueprint</p>
        
        <div class="error" id="error"></div>
        
        <form id="hdForm">
            <div class="row">
                <div class="form-group">
                    <label for="birthdate">Birth Date</label>
                    <input type="date" id="birthdate" name="birthdate" required>
                </div>
                <div class="form-group">
                    <label for="birthtime">Birth Time</label>
                    <input type="time" id="birthtime" name="birthtime" required>
                </div>
            </div>
            
            <div class="form-group">
                <label for="location">Birth Location</label>
                <input type="text" id="location" name="location" placeholder="e.g., New York, NY" required>
            </div>
            
            <div class="coordinates">
                <div class="form-group">
                    <label for="latitude">Latitude</label>
                    <input type="number" id="latitude" name="latitude" step="0.0001" placeholder="40.7128" required>
                </div>
                <div class="form-group">
                    <label for="longitude">Longitude</label>
                    <input type="number" id="longitude" name="longitude" step="0.0001" placeholder="-74.0060" required>
                </div>
            </div>
            
            <button type="submit">Generate My Chart</button>
        </form>
        
        <div class="loading" id="loading">
            <div class="spinner"></div>
            <p>Calculating your chart...</p>
        </div>
        
        <div class="info">
            <strong>💡 Tip:</strong> For accurate results, use your exact birth time from your birth certificate.
            You can look up coordinates for your birth city online.
        </div>
    </div>
    
    <script>
        document.getElementById('hdForm').addEventListener('submit', async function(e) {
            e.preventDefault();
            
            const form = e.target;
            const loading = document.getElementById('loading');
            const error = document.getElementById('error');
            
            error.style.display = 'none';
            form.style.display = 'none';
            loading.style.display = 'block';
            
            const birthdate = document.getElementById('birthdate').value;
            const birthtime = document.getElementById('birthtime').value;
            const datetime = birthdate + 'T' + birthtime + ':00Z';
            
            const data = {
                datetime: datetime,
                latitude: parseFloat(document.getElementById('latitude').value),
                longitude: parseFloat(document.getElementById('longitude').value),
                location: document.getElementById('location').value
            };
            
            try {
                const response = await fetch('/api/reading', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(data)
                });
                
                if (!response.ok) {
                    throw new Error('Failed to generate chart');
                }
                
                const html = await response.text();
                document.open();
                document.write(html);
                document.close();
            } catch (err) {
                loading.style.display = 'none';
                form.style.display = 'block';
                error.textContent = 'Error: ' + err.message;
                error.style.display = 'block';
            }
        });
    </script>
</body>
</html>`

	t, err := template.New("home").Parse(tmpl)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	t.Execute(w, nil)
}

// GenerateReading handles the API request to generate a reading
func (h *Handler) GenerateReading(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ReadingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Parse the datetime
	dt, err := time.Parse(time.RFC3339, req.DateTime)
	if err != nil {
		// Try alternate format
		dt, err = time.Parse("2006-01-02T15:04:05Z", req.DateTime)
		if err != nil {
			http.Error(w, "Invalid datetime format: "+err.Error(), http.StatusBadRequest)
			return
		}
	}

	birthData := calculator.BirthData{
		DateTime:  dt,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Location:  req.Location,
	}

	// Calculate the reading
	reading, err := h.calculator.Calculate(birthData)
	if err != nil {
		http.Error(w, "Calculation error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Generate the HTML response
	html, err := h.bodygraph.GenerateHTML(reading)
	if err != nil {
		http.Error(w, "Generation error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// GetReadingJSON returns the reading as JSON (for API clients)
func (h *Handler) GetReadingJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ReadingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	dt, err := time.Parse(time.RFC3339, req.DateTime)
	if err != nil {
		dt, err = time.Parse("2006-01-02T15:04:05Z", req.DateTime)
		if err != nil {
			http.Error(w, "Invalid datetime format", http.StatusBadRequest)
			return
		}
	}

	birthData := calculator.BirthData{
		DateTime:  dt,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
		Location:  req.Location,
	}

	reading, err := h.calculator.Calculate(birthData)
	if err != nil {
		log.Printf("Calculation error: %v", err)
		http.Error(w, "Calculation error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reading)
}

package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"humandesign/internal/astrocartography"
	"humandesign/internal/astrology"
	"humandesign/internal/bodygraph"
	"humandesign/internal/calculator"
	"humandesign/internal/cities"
	"humandesign/internal/database"
	"log"
	"net/http"
	"strings"
	"time"
)

// Handler holds dependencies for HTTP handlers
type Handler struct {
	calculator     *calculator.Calculator
	bodygraph      *bodygraph.Generator
	astroCalc      *astrology.Calculator
	astrocartoCalc *astrocartography.Calculator
	db             *database.Database
}

// NewHandler creates a new handler with all dependencies
func NewHandler(db *database.Database) *Handler {
	return &Handler{
		calculator:     calculator.NewCalculator(),
		bodygraph:      bodygraph.NewGenerator(),
		astroCalc:      astrology.NewCalculator(),
		astrocartoCalc: astrocartography.NewCalculator(),
		db:             db,
	}
}

// ReadingRequest represents the JSON request for a reading
type ReadingRequest struct {
	DateTime  string  `json:"datetime"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Location  string  `json:"location"`
}

// --- Shared layout / navigation ---

const navCSS = `
* { box-sizing: border-box; margin: 0; padding: 0; }
body {
    font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;
    background: #f1eee7;
    color: #2b2620;
    min-height: 100vh;
}
nav {
    background: #ffffff;
    border-bottom: 1px solid #e3ddd0;
    padding: 0 30px;
    display: flex;
    align-items: center;
    gap: 0;
    position: sticky;
    top: 0;
    z-index: 100;
}
nav .brand {
    color: #2b2620;
    font-size: 0.95rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 3px;
    margin-right: 30px;
    padding: 18px 0;
}
nav a {
    color: #8b8478;
    text-decoration: none;
    padding: 18px 18px;
    font-size: 0.85rem;
    font-weight: 500;
    text-transform: uppercase;
    letter-spacing: 1.5px;
    transition: color 0.2s, border-bottom 0.2s;
    border-bottom: 2px solid transparent;
}
nav a:hover { color: #2b2620; }
nav a.active { color: #b03a2e; border-bottom-color: #b03a2e; }
.page-container {
    max-width: 900px;
    margin: 36px auto;
    padding: 0 20px;
}
.card {
    background: #ffffff;
    border: 1px solid #e3ddd0;
    padding: 30px;
    margin-bottom: 24px;
}
.card h2 {
    color: #2b2620; margin-bottom: 20px; font-size: 1.05rem;
    font-weight: 600; text-transform: uppercase; letter-spacing: 2.5px;
    padding-bottom: 12px; border-bottom: 1px solid #e3ddd0;
}
.card h3 {
    color: #8b8478; margin: 18px 0 8px; font-size: 0.72rem;
    text-transform: uppercase; letter-spacing: 2px;
}
.form-group { margin-bottom: 18px; }
label {
    display: block; margin-bottom: 6px; color: #8b8478; font-size: 0.72rem;
    font-weight: 600; text-transform: uppercase; letter-spacing: 2px;
}
input[type="date"], input[type="time"], input[type="text"], input[type="number"] {
    width: 100%; padding: 11px 14px; border: 1px solid #d8d1c3; border-radius: 2px;
    background: #fdfcfa; font-size: 1rem; color: #2b2620; transition: border-color 0.2s;
}
input:focus { outline: none; border-color: #b03a2e; }
.row { display: flex; gap: 15px; }
.row .form-group { flex: 1; }
button, .btn {
    padding: 12px 30px; background: #26241f;
    color: white; border: none; border-radius: 2px; font-size: 0.85rem; font-weight: 600;
    text-transform: uppercase; letter-spacing: 2px;
    cursor: pointer; transition: background 0.2s; text-decoration: none; display: inline-block;
}
button:hover, .btn:hover { background: #b03a2e; }
.btn-sm { padding: 8px 16px; font-size: 0.72rem; }
.btn-danger { background: #b03a2e; }
.btn-danger:hover { background: #8e2f25; }
.btn-success { background: #57a05a; }
.btn-success:hover { background: #468148; }
table { width: 100%; border-collapse: collapse; margin-top: 15px; }
th, td { padding: 11px 14px; text-align: left; border-bottom: 1px solid #eee8dc; }
th {
    background: #faf8f3; color: #8b8478; font-size: 0.72rem;
    text-transform: uppercase; letter-spacing: 1.5px;
}
tr:hover { background: #f7f4ec; }
.tag { display: inline-block; padding: 4px 12px; border-radius: 10px; font-size: 0.78rem; font-weight: 600; margin: 2px; }
.tag-fire { background: #f6e3d5; color: #b45e26; }
.tag-earth { background: #e2ecdd; color: #468148; }
.tag-air { background: #dde8ef; color: #3d6e8f; }
.tag-water { background: #e6e0eb; color: #6d5488; }
.tag-harmonious { background: #e2ecdd; color: #468148; }
.tag-challenging { background: #f2ddda; color: #b03a2e; }
.tag-neutral { background: #f1ead9; color: #a07b2e; }
.tag-personality { background: #26241f; color: white; }
.tag-design { background: #b03a2e; color: white; }
.strength-very-strong { color: #b03a2e; font-weight: bold; }
.strength-strong { color: #b45e26; font-weight: 600; }
.strength-moderate { color: #3d6e8f; }
.strength-weak { color: #a49c8e; }
.dropdown-wrap { position: relative; }
.dropdown-list {
    display: none; position: absolute; top: 100%; left: 0; right: 0; max-height: 250px;
    overflow-y: auto; background: white; border: 1px solid #d8d1c3; border-top: none;
    z-index: 50; box-shadow: 0 8px 20px rgba(43,38,32,0.1);
}
.dropdown-list.open { display: block; }
.dropdown-item {
    padding: 10px 14px; cursor: pointer; font-size: 0.95rem; border-bottom: 1px solid #f2eee5;
}
.dropdown-item:hover, .dropdown-item.highlighted { background: #f4f0e6; }
.dropdown-item .country { color: #a49c8e; font-size: 0.8rem; margin-left: 6px; }
.error { background: #f2ddda; color: #8e2f25; padding: 15px; border: 1px solid #e4c2bd; margin-bottom: 20px; display: none; }
.spinner { border: 3px solid #eee8dc; border-top: 3px solid #b03a2e; border-radius: 50%; width: 40px; height: 40px; animation: spin 1s linear infinite; margin: 10px auto; }
@keyframes spin { 0%{transform:rotate(0)} 100%{transform:rotate(360deg)} }
.loading { display: none; text-align: center; padding: 20px; }
.info-box { background: #faf8f3; border: 1px solid #eee8dc; padding: 15px; font-size: 0.9rem; color: #6d675c; margin-top: 15px; }
`

func navHTML(active string) string {
	items := []struct{ href, label, id string }{
		{"/", "Human Design", "hd"},
		{"/astrology", "Astrology", "astrology"},
		{"/astrocartography", "Astrocartography", "astrocarto"},
		{"/people", "People", "people"},
	}
	var sb strings.Builder
	sb.WriteString(`<nav><span class="brand">Cosmic Blueprint</span>`)
	for _, item := range items {
		cls := ""
		if item.id == active {
			cls = ` class="active"`
		}
		sb.WriteString(fmt.Sprintf(`<a href="%s"%s>%s</a>`, item.href, cls, item.label))
	}
	sb.WriteString(`</nav>`)
	return sb.String()
}

func pageHead(title string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>%s - Cosmic Blueprint</title>
<style>%s</style>
</head>
<body>`, title, navCSS)
}

const pageFoot = `</body></html>`

// birthFormJS returns JS for the birth-data form used on multiple pages
func birthFormJS(endpoint string) string {
	return fmt.Sprintf(`
(function() {
    let citiesData = [];
    let highlightIdx = -1;

    fetch('/api/cities').then(r => r.json()).then(data => { citiesData = data; });

    const searchInput = document.getElementById('locationSearch');
    const dropdown = document.getElementById('dropdownList');
    const latInput = document.getElementById('latitude');
    const lonInput = document.getElementById('longitude');
    const locInput = document.getElementById('location');
    const tzInput  = document.getElementById('timezone');

    // Convert a local date+time in an IANA timezone to a UTC ISO string.
    function localToUTC(dateStr, timeStr, timezone) {
        const target = dateStr + 'T' + timeStr + ':00';
        let guess = new Date(target + 'Z'); // initial guess: treat as UTC
        for (let i = 0; i < 5; i++) {
            // Format guess in the target timezone using en-CA (gives YYYY-MM-DD, HH:MM)
            const inTZ = guess.toLocaleString('en-CA', {
                timeZone: timezone,
                year: 'numeric', month: '2-digit', day: '2-digit',
                hour: '2-digit', minute: '2-digit', hour12: false
            }).replace(',', '');
            const parts = inTZ.trim().split(' ');
            const tzLocal = parts[0] + 'T' + parts[1] + ':00Z';
            const diff = new Date(target + 'Z') - new Date(tzLocal);
            guess = new Date(guess.getTime() + diff);
            if (Math.abs(diff) < 60000) break;
        }
        return guess.toISOString();
    }

    function renderDropdown(matches) {
        dropdown.innerHTML = '';
        highlightIdx = -1;
        if (matches.length === 0) { dropdown.classList.remove('open'); return; }
        matches.forEach(function(city, i) {
            const div = document.createElement('div');
            div.className = 'dropdown-item';
            div.innerHTML = city.name + '<span class="country">' + city.country + '</span>';
            div.addEventListener('mousedown', function(e) { e.preventDefault(); selectCity(city); });
            dropdown.appendChild(div);
        });
        dropdown.classList.add('open');
    }

    function selectCity(city) {
        searchInput.value = city.name + ', ' + city.country;
        latInput.value = city.latitude;
        lonInput.value = city.longitude;
        locInput.value = city.name + ', ' + city.country;
        tzInput.value  = city.timezone || '';
        dropdown.classList.remove('open');
    }

    searchInput.addEventListener('input', function() {
        const q = this.value.toLowerCase().trim();
        if (q.length < 1) { dropdown.classList.remove('open'); return; }
        const matches = citiesData.filter(function(c) {
            return c.name.toLowerCase().indexOf(q) !== -1 || c.country.toLowerCase().indexOf(q) !== -1;
        }).slice(0, 20);
        renderDropdown(matches);
    });

    searchInput.addEventListener('keydown', function(e) {
        const items = dropdown.querySelectorAll('.dropdown-item');
        if (e.key === 'ArrowDown') {
            e.preventDefault();
            highlightIdx = Math.min(highlightIdx + 1, items.length - 1);
            items.forEach(function(el, i) { el.classList.toggle('highlighted', i === highlightIdx); });
            if (items[highlightIdx]) items[highlightIdx].scrollIntoView({block:'nearest'});
        } else if (e.key === 'ArrowUp') {
            e.preventDefault();
            highlightIdx = Math.max(highlightIdx - 1, 0);
            items.forEach(function(el, i) { el.classList.toggle('highlighted', i === highlightIdx); });
            if (items[highlightIdx]) items[highlightIdx].scrollIntoView({block:'nearest'});
        } else if (e.key === 'Enter' && highlightIdx >= 0 && dropdown.classList.contains('open')) {
            e.preventDefault();
            const q = searchInput.value.toLowerCase().trim();
            const matches = citiesData.filter(function(c) {
                return c.name.toLowerCase().indexOf(q) !== -1 || c.country.toLowerCase().indexOf(q) !== -1;
            }).slice(0, 20);
            if (matches[highlightIdx]) selectCity(matches[highlightIdx]);
        }
    });

    searchInput.addEventListener('blur', function() { setTimeout(function(){ dropdown.classList.remove('open'); }, 200); });
    searchInput.addEventListener('focus', function() { if (this.value.length >= 1) this.dispatchEvent(new Event('input')); });

    document.getElementById('birthForm').addEventListener('submit', async function(e) {
        e.preventDefault();
        if (!latInput.value || !lonInput.value) {
            const error = document.getElementById('error');
            error.textContent = 'Please select a city from the dropdown.';
            error.style.display = 'block';
            return;
        }
        const form = e.target;
        const loading = document.getElementById('loading');
        const error = document.getElementById('error');
        error.style.display = 'none';
        form.style.display = 'none';
        loading.style.display = 'block';
        const birthdate = document.getElementById('birthdate').value;
        const birthtime = document.getElementById('birthtime').value;
        const tz = tzInput.value;
        const datetime = tz ? localToUTC(birthdate, birthtime, tz) : birthdate + 'T' + birthtime + ':00Z';
        const data = {
            datetime: datetime,
            latitude: parseFloat(latInput.value),
            longitude: parseFloat(lonInput.value),
            location: locInput.value
        };
        try {
            const response = await fetch('%s', {
                method: 'POST',
                headers: {'Content-Type': 'application/json'},
                body: JSON.stringify(data)
            });
            if (!response.ok) throw new Error(await response.text());
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
})();`, endpoint)
}

func birthFormHTML() string {
	return `
<div class="error" id="error"></div>
<form id="birthForm">
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
        <label for="locationSearch">Birth Location</label>
        <div class="dropdown-wrap">
            <input type="text" id="locationSearch" autocomplete="off" placeholder="Type to search cities..." required>
            <div class="dropdown-list" id="dropdownList"></div>
        </div>
        <input type="hidden" id="latitude" name="latitude">
        <input type="hidden" id="longitude" name="longitude">
        <input type="hidden" id="location" name="location">
        <input type="hidden" id="timezone" name="timezone">
    </div>
    <button type="submit">Calculate</button>
</form>
<div class="loading" id="loading"><div class="spinner"></div><p>Calculating...</p></div>
<div class="info-box"><strong>Tip:</strong> Use your exact birth time from your birth certificate. Select the city closest to your birth location.</div>`
}

// ========================
// HOME PAGE - Human Design
// ========================

func (h *Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	html := pageHead("Human Design") + navHTML("hd") + `
<div class="page-container">
<div class="card">
    <h2>Human Design Calculator</h2>
    <p style="color:#8b8478;margin-bottom:20px;">Discover your unique energetic blueprint</p>` +
		birthFormHTML() + `
</div>
</div>
<script>` + birthFormJS("/api/reading") + `</script>` + pageFoot

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// ========================
// ASTROLOGY PAGE
// ========================

func (h *Handler) AstrologyPage(w http.ResponseWriter, r *http.Request) {
	html := pageHead("Astrology") + navHTML("astrology") + `
<div class="page-container">
<div class="card">
    <h2>Natal Astrology Chart</h2>
    <p style="color:#8b8478;margin-bottom:20px;">Calculate your Western astrology natal chart with zodiac placements, houses, and aspects</p>` +
		birthFormHTML() + `
</div>
</div>
<script>` + birthFormJS("/api/astrology") + `</script>` + pageFoot

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// ========================
// ASTROCARTOGRAPHY PAGE
// ========================

func (h *Handler) AstrocartographyPage(w http.ResponseWriter, r *http.Request) {
	html := pageHead("Astrocartography") + navHTML("astrocarto") + `
<div class="page-container">
<div class="card">
    <h2>Astrocartography</h2>
    <p style="color:#8b8478;margin-bottom:20px;">Discover how planetary energies manifest at different locations on Earth</p>` +
		birthFormHTML() + `
</div>
</div>
<script>` + birthFormJS("/api/astrocartography") + `</script>` + pageFoot

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// ========================
// PEOPLE DATABASE PAGE
// ========================

func (h *Handler) PeoplePage(w http.ResponseWriter, r *http.Request) {
	people := h.db.List()

	var rows strings.Builder
	for _, p := range people {
		rows.WriteString(fmt.Sprintf(`<tr>
<td>%s</td><td>%s</td><td>%s</td><td>%s</td>
<td>
  <a href="/people/%s/hd" class="btn btn-sm">HD</a>
  <a href="/people/%s/astrology" class="btn btn-sm">Astro</a>
  <a href="/people/%s/astrocartography" class="btn btn-sm">AstroC</a>
</td>
<td><button onclick="deletePerson('%s')" class="btn btn-sm btn-danger">Delete</button></td>
</tr>`, template.HTMLEscapeString(p.Name), p.BirthDate, p.BirthTime, template.HTMLEscapeString(p.Location),
			p.ID, p.ID, p.ID, p.ID))
	}

	html := pageHead("People") + navHTML("people") + `
<div class="page-container">
<div class="card">
    <h2>People Database</h2>
    <p style="color:#8b8478;margin-bottom:20px;">Store birth data for quick chart generation</p>
    <table>
        <thead><tr><th>Name</th><th>Birth Date</th><th>Birth Time</th><th>Location</th><th>Charts</th><th></th></tr></thead>
        <tbody>` + rows.String() + `</tbody>
    </table>
</div>
<div class="card">
    <h2>Add Person</h2>
    <form id="addPersonForm">
        <div class="form-group">
            <label for="name">Name</label>
            <input type="text" id="name" name="name" placeholder="Full name" required>
        </div>
        <div class="row">
            <div class="form-group">
                <label for="p_birthdate">Birth Date</label>
                <input type="date" id="p_birthdate" name="birthdate" required>
            </div>
            <div class="form-group">
                <label for="p_birthtime">Birth Time</label>
                <input type="time" id="p_birthtime" name="birthtime" required>
            </div>
        </div>
        <div class="form-group">
            <label for="p_locationSearch">Birth Location</label>
            <div class="dropdown-wrap">
                <input type="text" id="p_locationSearch" autocomplete="off" placeholder="Type to search cities..." required>
                <div class="dropdown-list" id="p_dropdownList"></div>
            </div>
            <input type="hidden" id="p_latitude" name="latitude">
            <input type="hidden" id="p_longitude" name="longitude">
            <input type="hidden" id="p_location" name="location">
            <input type="hidden" id="p_timezone" name="timezone">
        </div>
        <div class="form-group">
            <label for="notes">Notes (optional)</label>
            <input type="text" id="notes" name="notes" placeholder="Any notes...">
        </div>
        <button type="submit">Add Person</button>
    </form>
    <div id="addMsg" style="margin-top:15px;display:none;padding:10px;border-radius:8px;"></div>
</div>
</div>
<script>
(function() {
    let citiesData = [];
    let highlightIdx = -1;
    fetch('/api/cities').then(r => r.json()).then(data => { citiesData = data; });

    const searchInput = document.getElementById('p_locationSearch');
    const dropdown = document.getElementById('p_dropdownList');
    const latInput = document.getElementById('p_latitude');
    const lonInput = document.getElementById('p_longitude');
    const locInput = document.getElementById('p_location');
    const tzInput = document.getElementById('p_timezone');

    function renderDropdown(matches) {
        dropdown.innerHTML = '';
        highlightIdx = -1;
        if (matches.length === 0) { dropdown.classList.remove('open'); return; }
        matches.forEach(function(city) {
            const div = document.createElement('div');
            div.className = 'dropdown-item';
            div.innerHTML = city.name + '<span class="country">' + city.country + '</span>';
            div.addEventListener('mousedown', function(e) { e.preventDefault(); selectCity(city); });
            dropdown.appendChild(div);
        });
        dropdown.classList.add('open');
    }

    function selectCity(city) {
        searchInput.value = city.name + ', ' + city.country;
        latInput.value = city.latitude;
        lonInput.value = city.longitude;
        locInput.value = city.name + ', ' + city.country;
        tzInput.value = city.timezone || '';
        dropdown.classList.remove('open');
    }

    searchInput.addEventListener('input', function() {
        const q = this.value.toLowerCase().trim();
        if (q.length < 1) { dropdown.classList.remove('open'); return; }
        const matches = citiesData.filter(function(c) {
            return c.name.toLowerCase().indexOf(q) !== -1 || c.country.toLowerCase().indexOf(q) !== -1;
        }).slice(0, 20);
        renderDropdown(matches);
    });

    searchInput.addEventListener('keydown', function(e) {
        const items = dropdown.querySelectorAll('.dropdown-item');
        if (e.key === 'ArrowDown') { e.preventDefault(); highlightIdx = Math.min(highlightIdx+1, items.length-1); items.forEach(function(el,i){el.classList.toggle('highlighted',i===highlightIdx);}); }
        else if (e.key === 'ArrowUp') { e.preventDefault(); highlightIdx = Math.max(highlightIdx-1, 0); items.forEach(function(el,i){el.classList.toggle('highlighted',i===highlightIdx);}); }
        else if (e.key === 'Enter' && highlightIdx >= 0 && dropdown.classList.contains('open')) {
            e.preventDefault();
            const q = searchInput.value.toLowerCase().trim();
            const matches = citiesData.filter(function(c) { return c.name.toLowerCase().indexOf(q)!==-1||c.country.toLowerCase().indexOf(q)!==-1; }).slice(0,20);
            if (matches[highlightIdx]) selectCity(matches[highlightIdx]);
        }
    });

    searchInput.addEventListener('blur', function() { setTimeout(function(){ dropdown.classList.remove('open'); }, 200); });
    searchInput.addEventListener('focus', function() { if (this.value.length >= 1) this.dispatchEvent(new Event('input')); });

    document.getElementById('addPersonForm').addEventListener('submit', async function(e) {
        e.preventDefault();
        if (!latInput.value || !lonInput.value) {
            const msg = document.getElementById('addMsg');
            msg.textContent = 'Please select a city from the dropdown.';
            msg.style.background = '#fadbd8'; msg.style.color = '#c0392b'; msg.style.display = 'block';
            return;
        }
        const msg = document.getElementById('addMsg');
        const data = {
            name: document.getElementById('name').value,
            birth_date: document.getElementById('p_birthdate').value,
            birth_time: document.getElementById('p_birthtime').value,
            timezone: tzInput.value,
            location: locInput.value,
            latitude: parseFloat(latInput.value),
            longitude: parseFloat(lonInput.value),
            notes: document.getElementById('notes').value
        };
        try {
            const resp = await fetch('/api/people', { method:'POST', headers:{'Content-Type':'application/json'}, body: JSON.stringify(data) });
            if (!resp.ok) throw new Error(await resp.text());
            msg.textContent = 'Person added successfully!';
            msg.style.background = '#d5f5e3'; msg.style.color = '#27ae60'; msg.style.display = 'block';
            setTimeout(() => window.location.reload(), 800);
        } catch(err) {
            msg.textContent = 'Error: ' + err.message;
            msg.style.background = '#fadbd8'; msg.style.color = '#c0392b'; msg.style.display = 'block';
        }
    });
})();
async function deletePerson(id) {
    if (!confirm('Delete this person?')) return;
    await fetch('/api/people/' + id, { method: 'DELETE' });
    window.location.reload();
}
</script>` + pageFoot

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(html))
}

// ========================
// API: Human Design Reading
// ========================

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

	dt, err := parseDateTime(req.DateTime)
	if err != nil {
		http.Error(w, "Invalid datetime format: "+err.Error(), http.StatusBadRequest)
		return
	}

	birthData := calculator.BirthData{DateTime: dt, Latitude: req.Latitude, Longitude: req.Longitude, Location: req.Location}
	reading, err := h.calculator.Calculate(birthData)
	if err != nil {
		http.Error(w, "Calculation error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	logPositions(reading, req.Location, req.DateTime)

	htmlOut, err := h.bodygraph.GenerateHTML(reading)
	if err != nil {
		http.Error(w, "Generation error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlOut))
}

// logPositions prints all planetary positions as JSON to the terminal.
func logPositions(reading *calculator.Reading, location, datetimeUTC string) {
	type gateLog struct {
		Planet    string  `json:"planet"`
		Gate      int     `json:"gate"`
		Line      int     `json:"line"`
		Longitude float64 `json:"longitude"`
	}
	toLog := func(gates []calculator.Gate) []gateLog {
		out := make([]gateLog, len(gates))
		for i, g := range gates {
			out[i] = gateLog{g.Planet, g.Number, g.Line, g.Longitude}
		}
		return out
	}
	payload := map[string]interface{}{
		"location":          location,
		"datetime_utc":      datetimeUTC,
		"type":              reading.Type,
		"personality_gates": toLog(reading.PersonalityGates),
		"design_gates":      toLog(reading.DesignGates),
	}
	if b, err := json.MarshalIndent(payload, "", "  "); err == nil {
		log.Printf("=== HD Chart Positions ===\n%s", string(b))
	}
}

// GetReadingJSON returns the reading as JSON
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

	dt, err := parseDateTime(req.DateTime)
	if err != nil {
		http.Error(w, "Invalid datetime format", http.StatusBadRequest)
		return
	}

	birthData := calculator.BirthData{DateTime: dt, Latitude: req.Latitude, Longitude: req.Longitude, Location: req.Location}
	reading, err := h.calculator.Calculate(birthData)
	if err != nil {
		log.Printf("Calculation error: %v", err)
		http.Error(w, "Calculation error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reading)
}

// ========================
// API: Astrology
// ========================

func (h *Handler) GenerateAstrology(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ReadingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	dt, err := parseDateTime(req.DateTime)
	if err != nil {
		http.Error(w, "Invalid datetime format: "+err.Error(), http.StatusBadRequest)
		return
	}

	chart, err := h.astroCalc.CalculateChart(dt, req.Latitude, req.Longitude)
	if err != nil {
		http.Error(w, "Calculation error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	htmlOut := renderAstrologyHTML(chart, req.Location)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlOut))
}

func renderAstrologyHTML(chart *astrology.NatalChart, location string) string {
	var planets strings.Builder
	for _, p := range chart.Planets {
		planets.WriteString(fmt.Sprintf(`<tr><td><strong>%s</strong></td><td>%s %s</td><td>%d°%d'</td><td>%d</td></tr>`,
			p.Planet, p.SignSymbol, p.Sign, int(p.DegreeInSign), int((p.DegreeInSign-float64(int(p.DegreeInSign)))*60), p.House))
	}

	var aspects strings.Builder
	for _, a := range chart.Aspects {
		cls := "tag-neutral"
		if a.Harmony == "harmonious" {
			cls = "tag-harmonious"
		} else if a.Harmony == "challenging" {
			cls = "tag-challenging"
		}
		aspects.WriteString(fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td><span class="tag %s">%s</span></td><td>%.1f°</td></tr>`,
			a.Planet1, a.Planet2, cls, a.Type, a.Orb))
	}

	var houses strings.Builder
	for _, house := range chart.Houses {
		houses.WriteString(fmt.Sprintf(`<tr><td>House %d</td><td>%s</td><td>%.1f°</td></tr>`, house.Number, house.Sign, house.Degree))
	}

	var elements strings.Builder
	for el, count := range chart.Elements {
		cls := "tag-" + strings.ToLower(el)
		elements.WriteString(fmt.Sprintf(`<span class="tag %s">%s: %d</span> `, cls, el, count))
	}

	var modalities strings.Builder
	for mod, count := range chart.Modalities {
		modalities.WriteString(fmt.Sprintf(`<span class="tag tag-neutral">%s: %d</span> `, mod, count))
	}

	return pageHead("Astrology Chart") + navHTML("astrology") + fmt.Sprintf(`
<div class="page-container">
<div class="card">
    <h2>Natal Chart - %s</h2>
    <div class="row" style="gap:30px;flex-wrap:wrap;">
        <div style="flex:1;min-width:200px;">
            <h3>Sun Sign</h3><p style="font-size:1.3rem;">%s</p>
        </div>
        <div style="flex:1;min-width:200px;">
            <h3>Moon Sign</h3><p style="font-size:1.3rem;">%s</p>
        </div>
        <div style="flex:1;min-width:200px;">
            <h3>Rising Sign</h3><p style="font-size:1.3rem;">%s</p>
        </div>
    </div>
    <div style="margin-top:15px;">
        <h3>Midheaven (MC)</h3><p>%s</p>
    </div>
</div>

<div class="card">
    <h2>Planetary Placements</h2>
    <table><thead><tr><th>Planet</th><th>Sign</th><th>Degree</th><th>House</th></tr></thead>
    <tbody>%s</tbody></table>
</div>

<div class="card">
    <h2>Elements & Modalities</h2>
    <h3>Elements</h3><div>%s</div>
    <h3>Modalities</h3><div>%s</div>
</div>

<div class="card">
    <h2>Aspects</h2>
    <table><thead><tr><th>Planet 1</th><th>Planet 2</th><th>Aspect</th><th>Orb</th></tr></thead>
    <tbody>%s</tbody></table>
</div>

<div class="card">
    <h2>Houses</h2>
    <table><thead><tr><th>House</th><th>Sign</th><th>Cusp</th></tr></thead>
    <tbody>%s</tbody></table>
</div>
</div>`,
		template.HTMLEscapeString(location),
		chart.SunSign, chart.MoonSign, chart.RisingSign,
		chart.MCSign,
		planets.String(),
		elements.String(), modalities.String(),
		aspects.String(),
		houses.String(),
	) + pageFoot
}

// ========================
// API: Astrocartography
// ========================

func (h *Handler) GenerateAstrocartography(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ReadingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	dt, err := parseDateTime(req.DateTime)
	if err != nil {
		http.Error(w, "Invalid datetime format: "+err.Error(), http.StatusBadRequest)
		return
	}

	chart, err := h.astrocartoCalc.Calculate(dt, req.Latitude, req.Longitude)
	if err != nil {
		http.Error(w, "Calculation error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	htmlOut := renderAstrocartographyHTML(chart, req.Location)
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(htmlOut))
}

func renderAstrocartographyHTML(chart *astrocartography.AstrocartoChart, location string) string {
	var influences strings.Builder
	if chart.Location != nil && len(chart.Location.Influences) > 0 {
		for _, inf := range chart.Location.Influences {
			strengthClass := "strength-" + strings.ReplaceAll(inf.Strength, " ", "-")
			influences.WriteString(fmt.Sprintf(`<tr>
<td><strong>%s</strong></td>
<td>%s</td>
<td class="%s">%s</td>
<td>%.1f°</td>
<td>%s</td>
</tr>`, inf.Planet, inf.LineType, strengthClass, inf.Strength, inf.Distance, inf.Meaning))
		}
	}

	noInfluences := ""
	if chart.Location == nil || len(chart.Location.Influences) == 0 {
		noInfluences = `<p style="color:#a49c8e;margin-top:10px;">No strong planetary lines near this location. Try exploring other locations using the lines listed below.</p>`
	}

	// Group lines by planet for the summary
	planetLines := make(map[string][]string)
	for _, line := range chart.Lines {
		key := line.Planet
		planetLines[key] = append(planetLines[key], fmt.Sprintf("%s", line.LineType))
	}

	var linesSummary strings.Builder
	for planet, types := range planetLines {
		linesSummary.WriteString(fmt.Sprintf(`<tr><td><strong>%s</strong></td><td>%s</td><td>%s</td></tr>`,
			planet, strings.Join(types, ", "),
			getMeaningForPlanet(planet)))
	}

	return pageHead("Astrocartography") + navHTML("astrocarto") + fmt.Sprintf(`
<div class="page-container">
<div class="card">
    <h2>Astrocartography - %s</h2>
    <p style="color:#8b8478;margin-bottom:10px;">Showing planetary line influences at your specified location</p>
    <h3>Planetary Influences at This Location</h3>
    %s
    <table><thead><tr><th>Planet</th><th>Line</th><th>Strength</th><th>Distance</th><th>Meaning</th></tr></thead>
    <tbody>%s</tbody></table>
</div>

<div class="card">
    <h2>All Planetary Lines</h2>
    <p style="color:#8b8478;margin-bottom:10px;">Each planet creates 4 lines across the globe (MC, IC, ASC, DSC)</p>
    <table><thead><tr><th>Planet</th><th>Lines</th><th>General Theme</th></tr></thead>
    <tbody>%s</tbody></table>
</div>

<div class="card">
    <h2>Understanding Astrocartography Lines</h2>
    <div class="info-box">
        <p><strong>MC (Midheaven):</strong> Career, public life, reputation - how you're seen in the world</p>
        <p><strong>IC (Imum Coeli):</strong> Home, family, roots - where you feel a sense of belonging</p>
        <p><strong>ASC (Ascendant):</strong> Self-expression, identity - where your personality shines</p>
        <p><strong>DSC (Descendant):</strong> Relationships, partnerships - where you attract significant connections</p>
    </div>
</div>
</div>`,
		template.HTMLEscapeString(location),
		noInfluences,
		influences.String(),
		linesSummary.String(),
	) + pageFoot
}

func getMeaningForPlanet(planet string) string {
	meanings := map[string]string{
		"Sun":        "Vitality, identity, recognition",
		"Moon":       "Emotions, comfort, intuition",
		"Mercury":    "Communication, learning, connection",
		"Venus":      "Love, beauty, harmony, pleasure",
		"Mars":       "Energy, drive, passion, courage",
		"Jupiter":    "Expansion, luck, growth, wisdom",
		"Saturn":     "Discipline, structure, achievement",
		"Uranus":     "Innovation, change, liberation",
		"Neptune":    "Spirituality, dreams, creativity",
		"Pluto":      "Transformation, power, rebirth",
		"North Node": "Destiny, karmic growth, life purpose",
	}
	if m, ok := meanings[planet]; ok {
		return m
	}
	return "Planetary influence"
}

// ========================
// API: People CRUD
// ========================

func (h *Handler) HandlePeople(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(h.db.List())
	case http.MethodPost:
		var p database.Person
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "Invalid body: "+err.Error(), http.StatusBadRequest)
			return
		}
		if err := h.db.Add(p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(p)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) HandlePerson(w http.ResponseWriter, r *http.Request) {
	// Extract ID from path: /api/people/{id}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/people/"), "/")
	if len(parts) == 0 || parts[0] == "" {
		http.Error(w, "Missing person ID", http.StatusBadRequest)
		return
	}
	id := parts[0]

	switch r.Method {
	case http.MethodGet:
		p, err := h.db.Get(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	case http.MethodDelete:
		if err := h.db.Delete(id); err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// HandlePersonChart generates a chart for a person from the DB
func (h *Handler) HandlePersonChart(w http.ResponseWriter, r *http.Request) {
	// Path: /people/{id}/{chartType}
	path := strings.TrimPrefix(r.URL.Path, "/people/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	id := parts[0]
	chartType := parts[1]

	person, err := h.db.Get(id)
	if err != nil {
		http.Error(w, "Person not found: "+err.Error(), http.StatusNotFound)
		return
	}

	dt, err := personBirthUTC(person)
	if err != nil {
		http.Error(w, "Invalid date/time for person: "+err.Error(), http.StatusInternalServerError)
		return
	}

	switch chartType {
	case "hd":
		birthData := calculator.BirthData{DateTime: dt, Latitude: person.Latitude, Longitude: person.Longitude, Location: person.Location}
		reading, err := h.calculator.Calculate(birthData)
		if err != nil {
			http.Error(w, "Calculation error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		htmlOut, err := h.bodygraph.GenerateHTML(reading)
		if err != nil {
			http.Error(w, "Generation error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(htmlOut))

	case "astrology":
		chart, err := h.astroCalc.CalculateChart(dt, person.Latitude, person.Longitude)
		if err != nil {
			http.Error(w, "Calculation error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(renderAstrologyHTML(chart, person.Location+" ("+person.Name+")")))

	case "astrocartography":
		chart, err := h.astrocartoCalc.Calculate(dt, person.Latitude, person.Longitude)
		if err != nil {
			http.Error(w, "Calculation error: "+err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(renderAstrocartographyHTML(chart, person.Location+" ("+person.Name+")")))

	default:
		http.Error(w, "Unknown chart type: "+chartType, http.StatusBadRequest)
	}
}

// ========================
// API: Cities
// ========================

func (h *Handler) CitiesAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cities.All)
}

// ========================
// Helpers
// ========================

func parseDateTime(s string) (time.Time, error) {
	dt, err := time.Parse(time.RFC3339, s)
	if err != nil {
		dt, err = time.Parse("2006-01-02T15:04:05Z", s)
	}
	return dt, err
}

// personBirthUTC converts a stored person's local birth date/time into the
// correct UTC instant used by the ephemeris.
//
// The birth_time is a wall-clock time at the birth location, so it must be
// interpreted in that location's timezone before being converted to UTC.
// The timezone is taken from the person's record when present, otherwise it
// is resolved from the stored location/coordinates. Historically this was
// parsed directly as UTC, which shifted every non-UTC birth by the zone's
// offset and produced incorrect gates, profiles and types.
func personBirthUTC(person *database.Person) (time.Time, error) {
	tzName := person.Timezone
	if tzName == "" {
		tzName = cities.TimezoneFor(person.Location, person.Latitude, person.Longitude)
	}

	loc, err := time.LoadLocation(tzName)
	if err != nil {
		// Fall back to UTC if the zone database is unavailable/unknown.
		loc = time.UTC
	}

	local, err := time.ParseInLocation("2006-01-02T15:04", person.BirthDate+"T"+person.BirthTime, loc)
	if err != nil {
		return time.Time{}, err
	}
	return local.UTC(), nil
}

# CallSheet

**Turn sales call recordings into structured CRM rows.** Upload or record audio in the browser; CallSheet uses Google Gemini to extract summary, action items, sentiment, and urgency, then appends everything to a Google Sheet.

---

## What it does

- **Upload or record** — Drop an audio file (WAV, MP3, WebM) or record directly in the browser
- **AI analysis** — Gemini transcribes and analyzes the call, returning summary, sentiment (positive/neutral/negative), urgency (1–10), action items, and client name
- **One-click to Sheet** — Each analysis is appended as a row: Timestamp, Client, Summary, Sentiment, Urgency, and an **Urgent** flag (Yes when urgency > 7)

No manual note-taking; the sheet becomes your call log.

---

## Tech stack

| Layer    | Tech |
|----------|------|
| Backend  | Go, [Gin](https://github.com/gin-gonic/gin) |
| AI       | [Google Gemini 2.5 Flash](https://ai.google.dev/) (audio → structured JSON) |
| Storage  | Google Sheets API (service account) |
| Frontend | Vanilla HTML, CSS, JS — upload form + record button + results |

---

## Getting started

### Prerequisites

- **Go 1.23+** — [Install Go](https://go.dev/dl/)
- **Gemini API key** — [Google AI Studio](https://aistudio.google.com/apikey)
- **Google Sheet** + **Service account** with Sheets API enabled (see below)

### 1. Clone and run

```bash
git clone https://github.com/YOUR_USERNAME/CallSheet.git
cd CallSheet
go mod tidy
go run main.go
```

Server runs at **http://localhost:8080**.

### 2. Configure environment

Create a `.env` file in the project root (gitignored — never commit it):

| Variable | Description |
|----------|-------------|
| `GEMINI_API_KEY` | Your API key from [Google AI Studio](https://aistudio.google.com/apikey) |
| `SPREADSHEET_ID` | The ID in your sheet URL: `https://docs.google.com/spreadsheets/d/<SPREADSHEET_ID>/edit` |
| `GOOGLE_APPLICATION_CREDENTIALS` | Absolute path to your service account JSON key file |

### 3. Google Cloud setup

1. Create a [Google Cloud project](https://console.cloud.google.com/) (or use an existing one).
2. Enable the **Google Sheets API**.
3. Create a **service account** (IAM & Admin → Service Accounts → Create), then download its JSON key.
4. Create a Google Sheet and add a header row: **Timestamp** | **Client** | **Summary** | **Sentiment** | **Urgency** | **Urgent**
5. **Share the sheet** with the service account email (from the JSON `client_email`) and give it **Editor** access.

### 4. Use the app

Open http://localhost:8080. Choose **Upload** or **Record**, add your audio, then click **Analyze**. The result appears on the page and a new row is added to your sheet.

---

## Project structure

```
CallSheet/
├── main.go                 # Entrypoint, server, route wiring
├── web/
│   └── index.html          # Single-page UI: upload, record, results
└── internal/
    ├── config/             # Loads .env (GEMINI_API_KEY, SPREADSHEET_ID, etc.)
    ├── handlers/           # POST /voice-to-crm — validate file, call services, respond
    ├── models/             # VoiceAnalysis struct (summary, sentiment, urgency, …)
    └── services/
        ├── gemini.go       # Upload audio to Gemini, prompt for JSON, parse response
        └── sheets.go       # Append one row per analysis to the sheet
```

---

## API

**POST** `/voice-to-crm`

- **Content-Type:** `multipart/form-data`
- **Body:** field name `file` — audio file (`.wav`, `.mp3`, or `.webm`, max 25 MB)
- **Response (200):** JSON with `summary`, `action_items`, `sentiment`, `urgency_score`, `client_name`

Example:

```bash
curl -X POST http://localhost:8080/voice-to-crm -F "file=@call.mp3"
```

---

## License

MIT.

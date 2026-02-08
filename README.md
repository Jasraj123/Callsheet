# VoiceLine — Intelligent Sales Routing Engine

A Go backend that turns voice into actionable CRM data: it analyzes sales audio with **Google Gemini 1.5 Flash**, logs every call to **Google Sheets**, and triggers **instant alerts** for high-urgency leads via a configurable webhook (e.g. Slack).

---

## Why VoiceLine?

### Smart Routing as a Competitive Advantage

Most voice-to-CRM pipelines stop at “log the call.” VoiceLine adds **intelligent routing** so that:

- **Every call is captured** in a single sheet (Timestamp, Client, Summary, Sentiment, Urgency).
- **Urgent deals don’t slip:** when urgency &gt; 7, the engine immediately POSTs to your `ALERT_WEBHOOK_URL` with a clear message (e.g. “URGENT: [Client] needs immediate follow up!”), so sales and support can act in real time.
- **AI does the triage:** Gemini extracts summary, action items, sentiment, and a 1–10 urgency score so your team focuses on what matters instead of re-listening to every recording.

That combination—**always log, conditionally alert**—makes VoiceLine a routing engine, not just a logger, and gives you a clear edge in response time and pipeline visibility.

---

## Tech Stack

| Layer        | Technology                          |
|-------------|--------------------------------------|
| Framework   | **Gin** (Go)                         |
| AI          | **Google Gemini 1.5 Flash** (`google.golang.org/genai`) |
| Integrations| **Google Sheets API**, generic **Webhook** (Slack-compatible) |
| Config      | **.env** (e.g. `godotenv`)           |

---

## Project Structure (Clean Architecture)

```
voiceline/
├── main.go                 # Entrypoint, wiring
├── .env.example            # Env template
├── go.mod / go.sum
├── internal/
│   ├── config/             # Load GEMINI_API_KEY, SPREADSHEET_ID, ALERT_WEBHOOK_URL
│   ├── handlers/           # HTTP: POST /voice-to-crm
│   ├── models/             # VoiceAnalysis (summary, action_items, sentiment, urgency_score, client_name)
│   └── services/           # Gemini, Sheets, Webhook (with context timeouts)
```

- **Handlers** — API surface and validation.
- **Services** — External APIs (Gemini, Sheets, Webhook); all use `context` timeouts and proper error handling.

---

## Requirements

- **Go 1.21+**
- **GEMINI_API_KEY** — [Google AI Studio](https://aistudio.google.com/apikey)
- **SPREADSHEET_ID** — Target Google Sheet (service account must have edit access).
- **ALERT_WEBHOOK_URL** — Webhook URL for urgent alerts (e.g. Slack incoming webhook).
- **Google Sheets auth** — Service account JSON path in `GOOGLE_APPLICATION_CREDENTIALS`, or application default credentials.

---

## Setup

1. **Clone and install**

   ```bash
   cd voiceline
   go mod tidy
   go build -o voiceline .
   ```

2. **Configure environment**

   ```bash
   cp .env.example .env
   # Edit .env: GEMINI_API_KEY, SPREADSHEET_ID, ALERT_WEBHOOK_URL
   # Optional: GOOGLE_APPLICATION_CREDENTIALS=/path/to/service-account.json
   ```

3. **Google Sheet**

   - Create a sheet with (or leave empty for auto-append):  
     `Timestamp` | `Client` | `Summary` | `Sentiment` | `Urgency`
   - Share the sheet with the **service account email** (e.g. `xxx@yyy.iam.gserviceaccount.com`) with “Editor” access.

4. **Run**

   ```bash
   ./voiceline
   # Server listens on :8080
   ```

---

## API

### `POST /voice-to-crm`

- **Content-Type:** `multipart/form-data`
- **Field:** `file` — audio file (`.wav` or `.mp3`, max 25 MB).

**Flow:**

1. Upload is validated (type, size).
2. Audio is sent to **Gemini 1.5 Flash** with a Sales Assistant prompt; the model returns JSON: `summary`, `action_items`, `sentiment`, `urgency_score`, `client_name`.
3. **Always:** one row is appended to the Google Sheet: `[Timestamp, Client, Summary, Sentiment, Urgency]`.
4. **If urgency_score &gt; 7:** a POST is sent to `ALERT_WEBHOOK_URL` with:  
   `"🚨 URGENT: [Client Name] needs immediate follow up! Summary: [Summary]"` (e.g. in a JSON body like `{"text": "..."}` for Slack).

**Example (curl):**

```bash
curl -X POST http://localhost:8080/voice-to-crm \
  -F "file=@/path/to/call.mp3"
```

**Example response (200):**

```json
{
  "summary": "Client asked about enterprise pricing and timeline.",
  "action_items": ["Send proposal by EOW", "Schedule technical call"],
  "sentiment": "positive",
  "urgency_score": 8,
  "client_name": "Acme Corp"
}
```

---

## Smart Routing Logic (Summary)

| Condition        | Action |
|------------------|--------|
| Every request    | Append row to Google Sheet (Timestamp, Client, Summary, Sentiment, Urgency). |
| `urgency_score > 7` | POST to `ALERT_WEBHOOK_URL` with the urgent message above. |

---

## License

MIT.

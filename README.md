# CallSheet

CallSheet turns sales call audio into structured CRM data. Upload or record a call in the browser; the app uses **Google Gemini** to extract summary, action items, sentiment, and urgency, then appends a row to **Google Sheets**.

**Tech:** Go, Gin, Google Gemini 2.5 Flash, Google Sheets API, vanilla HTML/JS (upload + record UI).

---

## Run locally

**Prerequisites:** Go 1.23+, a Gemini API key, a Google Sheet, and a Google Cloud service account with Sheets API access.

1. **Clone and build**

   ```bash
   git clone https://github.com/YOUR_USERNAME/CallSheet.git
   cd CallSheet
   go mod tidy
   go run main.go
   ```

2. **Add your own config**

   Create a `.env` file in the project root (this file is gitignored; never commit it):

   - `GEMINI_API_KEY` — get one at [Google AI Studio](https://aistudio.google.com/apikey)
   - `SPREADSHEET_ID` — from your sheet’s URL: `.../d/SPREADSHEET_ID/edit`
   - `GOOGLE_APPLICATION_CREDENTIALS` — path to your service account JSON key

   Create a service account in Google Cloud, enable the Sheets API, download the JSON key, and **share the sheet** with the service account email (Editor).

3. **Use the app**

   Open **http://localhost:8080**. Upload or record audio, then click Analyze. Rows are appended to your sheet with: Timestamp, Client, Summary, Sentiment, Urgency, Urgent (Yes/No).

---

## Project layout

```
├── main.go              # Server, routes, wiring
├── web/index.html       # Frontend: upload / record + results
└── internal/
    ├── config/          # Load env (no secrets in repo)
    ├── handlers/        # POST /voice-to-crm
    ├── models/          # VoiceAnalysis struct
    └── services/        # Gemini client, Sheets client
```

---

## API

- **POST /voice-to-crm** — `multipart/form-data`, field `file` (`.wav`, `.mp3`, or `.webm`, max 25MB). Returns JSON with `summary`, `action_items`, `sentiment`, `urgency_score`, `client_name`.

---

## License

MIT.

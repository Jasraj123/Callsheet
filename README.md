# VoiceLine — Intelligent Sales Routing Engine

VoiceLine turns sales call audio into structured data on Google sheets. Upload or record a call, and it analyzes the conversation with **Google Gemini**, extracts summary, action items, sentiment, and urgency, then logs everything to **Google Sheets**.

---

## Quick Start — How to Spin This Up

### 1. Prerequisites

- **Go 1.23+** — [Install Go](https://go.dev/dl/)
- **Gemini API key** — [Google AI Studio](https://aistudio.google.com/apikey)
- **Google Cloud project** with Sheets API enabled and a service account

### 2. Clone & Install

```bash
git clone <repo-url>
cd Voiceline
go mod tidy
go build -o voiceline .
```

### 3. Environment Variables

Create a `.env` file in the project root:

```env
GEMINI_API_KEY=your_gemini_api_key
SPREADSHEET_ID=your_google_sheet_id
GOOGLE_APPLICATION_CREDENTIALS=/path/to/your-service-account.json
```

- **GEMINI_API_KEY** — From [aistudio.google.com/apikey](https://aistudio.google.com/apikey)
- **SPREADSHEET_ID** — The ID from your Google Sheet URL:  
  `https://docs.google.com/spreadsheets/d/SPREADSHEET_ID/edit`
- **GOOGLE_APPLICATION_CREDENTIALS** — Path to the service account JSON key file

### 4. Google Sheet Setup

1. Create a Google Sheet (or use an existing one).
2. Add a header row: `Timestamp` | `Client` | `Summary` | `Sentiment` | `Urgency` | `Urgent`
3. Enable **Google Sheets API** for your GCP project.
4. Create a service account in [Google Cloud Console](https://console.cloud.google.com/) → IAM & Admin → Service Accounts → Create.
5. Download the JSON key and set its path in `GOOGLE_APPLICATION_CREDENTIALS`.
6. **Share the sheet** with the service account email (from the JSON `client_email`) with **Editor** access.

### 5. Run

```bash
go run main.go
```

Server listens on **http://localhost:8080**.

- **Web UI:** Open http://localhost:8080 — upload or record audio, then click Analyze.
- **API:** `POST /voice-to-crm` with multipart `file` (`.wav`, `.mp3`, or `.webm`).

### 6. Test with curl

```bash
curl -X POST http://localhost:8080/voice-to-crm -F "file=@/path/to/call.mp3"
```

---

## Tech Stack (Current)

| Layer      | Technology                           |
|-----------|--------------------------------------|
| Backend   | Go, Gin                              |
| AI        | Google Gemini 2.5 Flash              |
| Storage   | Google Sheets API                    |
| Frontend  | Vanilla HTML/JS                      |

---

## Project Structure

```
Voiceline/
├── main.go
├── .env                    # Your secrets 
├── web/
│   └── index.html         # Upload/Record + results UI
└── internal/
    ├── config/            # Load env vars
    ├── handlers/          # POST /voice-to-crm
    ├── models/            # VoiceAnalysis struct
    └── services/          # Gemini, Sheets
```

---

## If You Built VoiceLine From Scratch — Architecture & Tech Choices

If you were designing VoiceLine as a full product (mobile app, backend, infra, AI) from scratch, here’s a practical stack and how the pieces would talk to each other.

### 1. Building Blocks

| Block         | Role |
|---------------|------|
| **Mobile app**| Capture and upload voice; show results and CRM history. |
| **Backend API** | Receive audio, call AI, persist to DB/sheets, manage auth. |
| **AI service**  | Transcribe + analyze audio (summary, sentiment, urgency, etc.). |
| **Storage**     | CRM records, user data, call history. |
| **Auth / identity** | Who can record, which CRM to write to. |

### 2. Tech Choices

**Mobile app**

- **React Native** or **Flutter** for iOS + Android with shared logic.
- **Expo** (React Native) for quicker iteration and OTA updates.
- Record with native modules (e.g. `expo-av`, `react-native-audio-recorder`).

**Backend**

- **Go** (as in current VoiceLine) or **Node/TypeScript** with a framework like **Fastify**.
- REST for normal flows; **WebSockets** if you add real-time features (live transcription, notifications).
- **PostgreSQL** for users, orgs, CRM metadata; optionally **Redis** for caching and queues.

**AI**

- **Google Gemini** (or **Whisper + GPT**) for transcription + structured analysis.
- Use **structured output** (e.g. JSON schema) so parsing is reliable.

**Infra**

- **Cloud Run** or **AWS Lambda** for the API (auto-scale, pay-per-use).
- **GCP** if you lean on Gemini and Vertex AI.
- **Terraform** or **Pulumi** for infra-as-code.
- **GitHub Actions** or **Cloud Build** for CI/CD.

**Storage**

- **PostgreSQL** (users, orgs, call metadata).
- **Google Sheets** or **Airtable** as “CRM” if you want spreadsheets; otherwise a proper DB with an ORM.
- **Cloud Storage (GCS/S3)** for raw audio if you need long-term storage or replay.

### 3. Communication Between Blocks

```
┌─────────────┐     HTTPS      ┌─────────────┐     API calls     ┌─────────────┐
│  Mobile App │ ◄────────────► │  Backend    │ ◄───────────────► │  Gemini AI  │
│  (RN/Flutter)│   REST/WS     │  (Go/Node)  │   (REST/GRPC)     │             │
└─────────────┘                └──────┬──────┘                   └─────────────┘
                                     │
                                     │ SQL / Sheets API
                                     ▼
                              ┌─────────────┐
                              │  PostgreSQL │
                              │  or Sheets  │
                              └─────────────┘
```

**Flow**

1. **Mobile → Backend:**  
   - Auth (JWT or session cookie).  
   - `POST /voice-to-crm` with multipart audio.  
   - Optional: WebSocket for live transcription.

2. **Backend → AI:**  
   - Upload audio to Gemini (or Whisper) via REST.  
   - Request structured JSON (summary, sentiment, urgency, etc.).  
   - Parse and validate response.

3. **Backend → Storage:**  
   - Write analysis rows to PostgreSQL or Google Sheets.  
   - Optionally store raw audio in GCS/S3 and link by ID.

4. **Backend → Mobile:**  
   - Return analysis JSON.  
   - Optionally push updates (e.g. via WebSocket or push notifications) for urgent alerts.

**Auth**

- **OAuth 2.0** (Google, Microsoft) for sign-in.  
- **JWT** for API auth; refresh tokens for long-lived sessions.  
- Backend checks token and maps user → org/CRM before writing.

### 4. Minimal Viable Stack

For an MVP:

- **Expo / React Native** app, or simple web UI (like current VoiceLine).
- **Go** backend with **Gin**, deployed on **Cloud Run**.
- **Gemini** for audio analysis.
- **Google Sheets** as CRM (or Postgres if you want a real DB).
- **Firebase Auth** or **Supabase** for auth if you need it quickly.

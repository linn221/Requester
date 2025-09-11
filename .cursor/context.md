# Project Context & Guidelines

This project is written in **Go**, **Templ**, **Alpine.js**, **HTMX**, and **GORM**.  
It is designed to be used in **Bug Hunting** workflows — specifically to record, analyze, and monitor HTTP requests.  

---

## Goals
- Keep detailed records of HTTP requests.
- Enable later analysis using:
  - Entropy scanning
  - Sensitive comment detection
  - Wordlist generation
  - Daily endpoint change detection
- Maintain **clean code** and **clarity of design**.
- Strong preference for **server-side rendering (SSR) with Templ** over JavaScript.

---

## Build & Workflow Rules
- **Always regenerate Templ files** whenever a template is changed.  
- **Always build the Go binary** once a coding prompt / task is finished.  

---

## Technology Preferences
### Templ
- Default to **Templ** for UI rendering.
- Treat Templ as the **preferred layer** over JavaScript.  

### HTMX
- Prefer **HTMX** for fetching, interactivity, and UI updates.
- Follow and appreciate HTMX patterns from the official examples:  
  👉 https://htmx.org/examples/  
- Whenever convenient, choose an HTMX pattern over writing custom JS.  

### Alpine.js
- Use **Alpine.js inline styles** instead of `<script>` tags.
- Alpine should only be used where interactivity is lightweight and unavoidable.  

### JavaScript
- Avoid JavaScript unless **absolutely necessary**.
- Never introduce front-end frameworks beyond HTMX + Alpine.  

---

## UI/UX Guidelines
- **Action buttons** should appear on the **right side**.  
- **Search bar** should take the **full width**.  
- **Text suggestions** should behave like Google Search autocomplete.  
- Strive for a **minimal, clean, and functional UI**.  

---

## General Rules
- Do **not** refactor the **User Interface** unless explicitly asked.  
- Code should always follow **clean code principles**.  
- Prioritize **server-side rendering**.  

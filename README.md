Gemini Reporter: Automated Apple Search Ads Report Generator
Gemini Reporter is a cross-platform desktop application built to automate the generation of client-facing performance reports for Apple Search Ads. It connects to a MongoDB data source (or uses local JSON data in dev mode), processes ad performance data, and populates professional HTML email templates.

The application provides a user-friendly desktop UI for generating, previewing, editing, and sending these reports. It also integrates with Large Language Models (LLMs) by generating pre-written prompts that include the data-filled HTML, allowing for automated qualitative analysis.

Key Features
Modular Reporting System: Easily add new report types by implementing a simple Go interface.

Live Editing & Preview: A split-screen view allows for direct editing of the generated HTML with a live-reloading <iframe> preview.

Integrated AI Prompting: A "Copy AI Prompt" feature constructs a detailed prompt, wrapping the generated HTML with instructions for an LLM to perform analysis.

Direct Email Sending: Connects to a Gmail account using an App Password to send the final HTML report to a list of recipients.

Offline Development Mode: Run the application entirely offline by using local dummy JSON data instead of a live database connection.

Dynamic Recipient Loading: The email composer automatically fetches a list of default recipients based on the selected company.

Technology Stack
The application is built using the Wails (v2) framework, which combines a Go backend with a modern web-based frontend.

Role	Technology	Description
Backend	Go (Golang)	Handles all business logic, data fetching and processing, template rendering, and email sending.
Frontend	Svelte & Vite	Manages the entire User Interface (UI) and user experience with a modern, component-based approach.
Database	MongoDB	Primary data source for ad performance metrics and user data. Can be swapped for local JSON files.

Export to Sheets
Architecture Overview
Backend (Go)
main.go: Wails application entry point. Configures the main application window.

app.go: The bridge between the frontend and backend. Methods on the App struct are exposed to the Svelte UI.

/reports: Contains the modular system for different report types. Each report implements the ReportGenerator interface.

/internal/database: The data layer, containing data models and functions for fetching data from MongoDB or local JSON files.

/internal/email: A utility to execute Go's html/template engine.

/internal/email_sender: A dedicated package for sending emails via SMTP.

/internal/config: Manages loading settings from configs/config.json.

Frontend (Svelte)
App.svelte: The main component that acts as a router, controlling which view is visible.

/lib/CompanySelect.svelte: Displays a searchable list of clients.

/lib/ReportSelect.svelte: Displays a list of available report types.

/lib/ReportView.svelte: The core view with a split-screen layout for HTML editing and live preview.

/lib/EmailSender.svelte: A modal for composing and sending the final email report.

Workflow Example
Select Company: The user starts by selecting a client from a searchable list.

Select Report: The user then chooses one of the available report types (e.g., "Monthly Account Overview").

Generate Report: The Go backend fetches the required data, performs aggregations and calculations (like month-over-month changes), and populates the corresponding HTML template.

Preview & Edit: The final HTML is displayed in a split-screen view where the user can preview the rendered report and make live edits to the source code if needed.

AI Analysis (Optional): The user can copy a pre-formatted prompt to their clipboard, paste it into an LLM, and receive a qualitative analysis to insert back into the report.

Send Email: The user opens the email sender, which is pre-populated with default recipients, and sends the final HTML report directly from the application.
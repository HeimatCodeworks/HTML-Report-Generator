Project Technical Summary: Automated Apple Search Ads Reporter
1. Core Purpose
The project is a cross-platform desktop application designed to automate the generation of client-facing performance reports for Apple Search Ads. Its primary function is to query advertising performance data from a MongoDB database, process and aggregate this data according to specific business logic, and inject it into predefined HTML email templates. The application features a user interface for selecting the client and report type, and it provides a preview of the generated HTML, which can then be used for client communication.

A key feature is its integration with Large Language Models (LLMs). The generated HTML reports contain structured prompts within summary sections, allowing a user to paste the entire output into an LLM to automatically generate qualitative analyses of the quantitative data.

2. Architecture & Technology Stack
The application is built using the Wails (v2) framework, which combines a Go backend with a modern web-based frontend.

Backend: Go (Golang)

Responsibilities: All business logic, data fetching, data processing, calculations, and template rendering.

Key Packages:

go.mongodb.org/mongo-driver for database connectivity.

Standard libraries: html/template, encoding/json, time.

Frontend: Svelte with Vite

Responsibilities: All user interface elements, state management, and user interaction.

Styling: Modern dark theme (dark gray and blue) with custom CSS, implemented using global CSS variables in Svelte.

Platform: Cross-platform desktop application (builds to native executables for Windows, macOS, etc.).

3. Backend Breakdown
The Go backend is organized into several distinct packages for modularity and maintainability.

main.go & app.go (Root)

main.go: The main entry point of the Wails application. It configures and runs the app, setting window properties like size, transparency, and the application icon.

app.go: Defines the App struct, which serves as the bridge between the Go backend and the Svelte frontend. All methods on this struct (e.g., GetCompanies(), GenerateReport()) are automatically bound and callable from the JavaScript frontend.

/reports Package (Modular Reporting Core)

This is the heart of the business logic. It's designed to be extensible for new report types.

report_interface.go: Defines the ReportGenerator interface, which mandates that every report type must have Name(), TemplatePath(), and GenerateEmailData() methods.

Implemented Reports:

campaign_performance.go: A simple weekly report showing campaign performance.

account_overview.go: A complex monthly report that performs multi-level aggregation. It fetches data for the previous two full calendar months, aggregates campaign data up to the app level and the overall account level, calculates month-over-month changes, and formats all data for its template.

/internal/database Package (Data Layer)

Handles all interactions with data sources.

Dev Mode: A key feature is a "dev mode" (controlled by config.json) that reads data from local JSON files (/data/*.json) instead of connecting to a live MongoDB instance. This allows for offline development and testing.

models.go: Defines the Go structs that map to the MongoDB collections and dummy JSON files (e.g., CampaignReport, App, UserTeam).

database.go: Contains functions for connecting to MongoDB and for fetching data (e.g., FetchDummyCompanies, FetchDummyCampaignData).

/internal/email Package (Template Engine)

A lightweight utility package responsible for parsing and executing Go's html/template files. The GenerateEmail function takes a template path and a data map and returns the final, populated HTML string.

4. Frontend Breakdown (Svelte)
The frontend is a single-page application built with Svelte, managing the UI through a component-based, state-driven workflow.

App.svelte (Main Component/Router)

Manages the current view state (home, company, report, result).

Conditionally renders the appropriate child component based on the current state.

Handles events dispatched from child components to transition between states (e.g., listens for on:companySelected to move from the company list to the report list).

Contains the global CSS theme definitions (<style global>).

View Components (/lib/*.svelte)

CompanySelect.svelte: Fetches the list of companies from the Go backend (wailsjs/go/main/App.GetCompanies), displays them alphabetically, and includes a real-time search/filter input.

ReportSelect.svelte: After a company is selected, this view fetches the available report types from Go and displays them. Includes a "Back" button for navigation.

ReportView.svelte: The final view, featuring a split-screen layout:

Left Panel: Displays the raw, generated HTML source code in a <pre> block, with a custom inlaid "Copy" button.

Right Panel: Renders a live preview of the HTML inside a sandboxed <iframe> using the srcdoc attribute.

Contains action buttons ("Generate New Report," "Generate Another Report for this Company") for restarting the workflow.

5. Data Flow Example (Account Overview Report)
User clicks "Generate New Report" in the UI.

App.svelte changes its state to 'company'.

CompanySelect.svelte mounts and calls the GetCompanies() method in app.go.

app.go calls database.FetchDummyCompanies(), which reads and parses /data/dummy_companies.json, returning a list of companies.

The user selects a company. CompanySelect.svelte dispatches a companySelected event.

App.svelte catches the event and changes state to 'report'.

ReportSelect.svelte mounts, displaying the available reports. The user selects "Monthly Account Overview" (ID 1).

ReportSelect.svelte calls the GenerateReport(companyID, 1) method in app.go.

app.go selects the AccountOverviewReport generator.

The generator's GenerateEmailData() method executes:

It calculates the required date ranges (previous and prior full months).

It fetches all necessary app and campaign data from the dummy JSON files.

It filters, aggregates, and performs calculations (MoM changes, CPAs).

It structures the final data into a large map.

app.go receives this data map, calls email.GenerateEmail(), which parses /templates/account_overview/template.html and populates it with the data.

The final HTML string is returned to the ReportSelect.svelte component.

ReportSelect.svelte dispatches a reportGenerated event containing the HTML.

App.svelte catches the event, stores the HTML, and changes state to 'result'.

ReportView.svelte mounts and displays the HTML in the split-screen preview.
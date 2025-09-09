# Templates Organization

This directory contains organized template files for the Requester application.

## File Structure

### Core Templates
- **`layout.templ`** - Main layout template with navbar, CSS, and common structure
- **`helpers.go`** - Helper functions for formatting and data processing

### Page Templates
- **`home.templ`** - Home page template
- **`request_listing.templ`** - Requests listing page with Vue.js interactivity
- **`request_details.templ`** - Request details page with Vue.js expand/collapse functionality
- **`import.templ`** - Import form and result pages

### Common Components
- **`common.templ`** - Shared components like ErrorBox and JobStatus

## Usage

All templates are re-exported through the main `views/templates.templ` file, which maintains the same API as before but with better organization.

## Features

- **Vue.js 3** - Modern reactive framework for client-side interactivity
- **Bootstrap 5** - Responsive UI components
- **HTMX** - Server-side interactions without page reloads
- **Full-width View More buttons** - Better UX for expandable content
- **Organized structure** - Each page has its own template file

## Generated Files

Each `.templ` file generates a corresponding `_templ.go` file with the compiled template code. These generated files should not be edited manually.

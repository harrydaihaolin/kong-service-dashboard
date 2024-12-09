#!/bin/bash

# Setup environment variables
export SERVICE_DASHBOARD_DB_PORT=5432
export SERVICE_DASHBOARD_DB_USER=postgres
export SERVICE_DASHBOARD_DB_PASSWORD=example
export SERVICE_DASHBOARD_DB_NAME=postgres
export SERVICE_DASHBOARD_DB_HOST=localhost

# Set variables for filenames
COVERAGE_DIR="coverage"
COVERAGE_PROFILE="$COVERAGE_DIR/coverage.out"
HTML_REPORT="$COVERAGE_DIR/coverage.html"
SUMMARY_REPORT="$COVERAGE_DIR/coverage_summary.txt"

# Create coverage directory if it doesn't exist
mkdir -p $COVERAGE_DIR

# Step 1: Run tests and generate coverage profile
echo "Running tests and generating coverage profile..."
go test ./cmd/... -coverprofile=$COVERAGE_PROFILE

# Step 2: Generate HTML coverage report
echo "Generating HTML coverage report..."
go tool cover -html=$COVERAGE_PROFILE -o $HTML_REPORT

# Step 3: Extract percentage data
echo "Extracting percentage data..."
go tool cover -func=$COVERAGE_PROFILE > $SUMMARY_REPORT
TOTAL_COVERAGE=$(grep total $SUMMARY_REPORT | awk '{print $3}')

# Step 4: Display results
echo "Coverage Summary:"
cat $SUMMARY_REPORT
echo
echo "Total Coverage: $TOTAL_COVERAGE"

# Step 5: Notify about HTML report
echo "HTML report generated: $HTML_REPORT"

# Exit
exit 0

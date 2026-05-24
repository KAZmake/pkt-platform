#!/usr/bin/env bash
# Run k6 load tests against staging or local environment.
#
# Usage:
#   ./run.sh public                         # public endpoints, local
#   ./run.sh borrower                       # borrower flow, local (needs BORROWER_TOKEN)
#   ./run.sh employee                       # employee flow, local (needs EMPLOYEE_TOKEN)
#   ./run.sh spike                          # spike test
#   ./run.sh all                            # all scenarios sequentially
#
#   BASE_URL=https://staging.pkt.kz ./run.sh public   # against staging
#
# Prerequisites: k6 installed — https://k6.io/docs/get-started/installation/

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

BASE_URL="${BASE_URL:-http://localhost:8080}"
WEB_URL="${WEB_URL:-http://localhost:3100}"
EXPERTISE_URL="${EXPERTISE_URL:-http://localhost:8081}"
SYNC_URL="${SYNC_URL:-http://localhost:8082}"
ASSISTANT_URL="${ASSISTANT_URL:-http://localhost:8083}"
BORROWER_TOKEN="${BORROWER_TOKEN:-}"
EMPLOYEE_TOKEN="${EMPLOYEE_TOKEN:-}"

# Output dir for reports
REPORT_DIR="${SCRIPT_DIR}/reports"
mkdir -p "$REPORT_DIR"
TIMESTAMP="$(date +%Y%m%d_%H%M%S)"

run_scenario() {
  local name="$1"
  local file="${SCRIPT_DIR}/scenarios/${name}.js"
  local report="${REPORT_DIR}/${name}_${TIMESTAMP}.json"

  echo ""
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
  echo "  Running scenario: ${name}"
  echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

  k6 run \
    --out "json=${report}" \
    -e "BASE_URL=${BASE_URL}" \
    -e "WEB_URL=${WEB_URL}" \
    -e "EXPERTISE_URL=${EXPERTISE_URL}" \
    -e "SYNC_URL=${SYNC_URL}" \
    -e "ASSISTANT_URL=${ASSISTANT_URL}" \
    -e "BORROWER_TOKEN=${BORROWER_TOKEN}" \
    -e "EMPLOYEE_TOKEN=${EMPLOYEE_TOKEN}" \
    "$file"

  echo ""
  echo "  Report saved: ${report}"
}

SCENARIO="${1:-public}"

case "$SCENARIO" in
  public|borrower|employee|spike)
    run_scenario "$SCENARIO"
    ;;
  all)
    run_scenario public
    run_scenario borrower
    run_scenario employee
    ;;
  *)
    echo "Unknown scenario: $SCENARIO"
    echo "Available: public, borrower, employee, spike, all"
    exit 1
    ;;
esac

echo ""
echo "Done."

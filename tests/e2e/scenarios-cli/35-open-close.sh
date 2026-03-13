#!/bin/bash
# 35-open-close.sh — CLI open/close commands

source "$(dirname "$0")/common.sh"

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab open <url>"

pt_ok open "${FIXTURES_URL}/index.html"
assert_output_json
assert_output_contains "tabId" "returns tab ID"
assert_output_contains "title" "returns page title"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab open --block-images <url>"

pt_ok open --block-images "${FIXTURES_URL}/index.html"
assert_output_json
assert_output_contains "tabId" "navigated with images blocked"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab open (no args → error)"

pt_fail open
assert_output_contains "requires" "shows usage error"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab close <tabId>"

# Open a tab to close
pt_ok open --new-tab "${FIXTURES_URL}/index.html"
TAB_ID=$(echo "$PT_OUT" | jq -r '.tabId // empty')

if [ -n "$TAB_ID" ]; then
  pt_ok close "$TAB_ID"
  echo -e "  ${GREEN}✓${NC} closed tab $TAB_ID"
  ((ASSERTIONS_PASSED++)) || true
else
  echo -e "  ${RED}✗${NC} could not get tab ID to close"
  ((ASSERTIONS_FAILED++)) || true
fi

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab close --tab <tabId>"

pt_ok open --new-tab "${FIXTURES_URL}/index.html"
TAB_ID=$(echo "$PT_OUT" | jq -r '.tabId // empty')

if [ -n "$TAB_ID" ]; then
  pt_ok close --tab "$TAB_ID"
  echo -e "  ${GREEN}✓${NC} closed tab via --tab flag"
  ((ASSERTIONS_PASSED++)) || true
else
  echo -e "  ${RED}✗${NC} could not get tab ID to close"
  ((ASSERTIONS_FAILED++)) || true
fi

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab close (no args, no --tab → error)"

pt_fail close
assert_output_contains "specify a tab ID" "requires tab ID"

end_test

# ─────────────────────────────────────────────────────────────────
# ─────────────────────────────────────────────────────────────────
start_test "pinchtab goto <url> (alias for open)"

pt_ok goto "${FIXTURES_URL}/index.html"
assert_output_json
assert_output_contains "tabId" "goto works as alias"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab navigate <url> (alias for open)"

pt_ok navigate "${FIXTURES_URL}/index.html"
assert_output_json
assert_output_contains "tabId" "navigate works as alias"

end_test

# ─────────────────────────────────────────────────────────────────
start_test "pinchtab nav shows deprecation notice"

pt_ok nav "${FIXTURES_URL}/index.html"
# Cobra prints deprecation to stderr, but output should still work
assert_output_json
assert_output_contains "tabId" "nav still works as alias"

end_test

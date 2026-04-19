package cli

// Doctor usage documentation.
const doctorDoc = `Usage: envchain doctor

Checks the health of your envchain setup:
  - Verifies the store directory exists
  - Checks store directory permissions (warns if group/world readable)
  - Reports the number of stored projects

Exits with a non-zero status if any issues are found.
`

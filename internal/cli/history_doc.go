package cli

// history subcommand documentation.
//
// Usage:
//
//	envchain history <project>
//
// Prints a list of all snapshots for the given project in chronological order,
// showing each snapshot name alongside its creation timestamp.
//
// Snapshots are created with `envchain snapshot <project> <label>`.
const historyDoc = `
Usage: envchain history <project>

List all snapshots for <project> in chronological order.

Each line shows the snapshot name and the UTC timestamp at which it was taken.
`

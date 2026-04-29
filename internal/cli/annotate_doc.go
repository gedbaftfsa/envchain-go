package cli

// annotate command documentation.
//
// The annotate command lets you attach a short human-readable description to
// any project stored in envchain.  Annotations are encrypted at rest together
// with the environment variable set and are never written in plain text.
//
// Examples:
//
//	# Display the current annotation
//	envchain annotate myproject
//
//	# Set (or overwrite) an annotation
//	envchain annotate myproject "Production database credentials"
//
//	# Clear an annotation by setting it to empty
//	envchain annotate myproject ""
const annotateDoc = `annotate <project> [text...]  show or set a project annotation`

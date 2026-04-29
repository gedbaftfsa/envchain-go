package cli

// shell and unshell commands
//
// Usage:
//
//	envchain shell <project>
//		Print export statements for all variables in <project>.
//		Intended to be eval'd in the current shell:
//			eval "$(envchain shell myproject)"
//
//	envchain unshell <project>
//		Print unset statements for all variables in <project>.
//		Intended to be eval'd in the current shell:
//			eval "$(envchain unshell myproject)"

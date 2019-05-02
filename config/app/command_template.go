package app

// Command Help Template
// Purpose: Every command has a unique usage message with flags and their descriptions
//          This template maps the text for each command and outputs it in a help message
//          for convenience.
var CommandHelpTemplate = `NAME: 
   {{.HelpName}}{{if .Description}} - {{.Description}}{{end}}

USAGE:
   {{.HelpName}}{{if .UsageText}} {{.UsageText}}{{end}}{{if .ArgsUsage}} {{.ArgsUsage}}{{end}}{{if .Flags}}

OPTIONS:
   {{range .Flags}}{{.}}
   {{end}}{{end}}
`

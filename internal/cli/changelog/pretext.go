package changelog

// pretext is the placeholder text for the input file
const pretext = `// New release for {{.Name}} v{{.Version}}
// All lines starting with '//' will be excluded from final changelog
// Insert changelog below this comment. An example format has been given:

## v{{.Version}} ({{.Date}})

FEATURES:

* **Feature Name**: Description about new feature this release

IMPROVEMENTS:

* **Improvement Name**: Description about new improvement this release

BUG FIXES:

* topic: Description of the bug. Example below [bug#]
* api: Fix Go API using lease revocation via URL instead of body [GH-7777]
`

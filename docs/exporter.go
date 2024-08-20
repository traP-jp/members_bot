package docs

import (
	_ "embed"
)

//go:embed help.md
var HelpTemplate string

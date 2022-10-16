package mrpkg

import (
	// import cobra explicitly so go would include
	// it in go.mod file, then there are no errors
	// when executing 'go run "github.com/Boyux/mrpkg/loadc"'
	_ "github.com/spf13/cobra"
)

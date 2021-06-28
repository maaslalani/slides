package version

// Set at build time with:
//   -ldflags "-X github.com/maaslalani/slides/internal/version.Version=<v> -X ..."
var (
	Version         string = "dev"
	ShortCommit     string = "unset"
	CommitDate      string = "date unset"
	VersionTemplate string = `   _______
  | # ... |  slides %s
  | ..... |  Commit %s (%s)
  |_______|  Copyright (C) 2021 Maas Lalani
     /|\
`
)

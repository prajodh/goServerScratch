module github.com/prajodh/goServerScratch

go 1.22.5

replace github.com/prajodh/goServerScratch/database v0.0.0 => ./database

require (
	github.com/prajodh/goServerScratch/database v0.0.0
	golang.org/x/crypto v0.26.0
)

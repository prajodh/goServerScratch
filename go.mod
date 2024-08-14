module github.com/prajodh/goServerScratch

go 1.22.5

replace github.com/prajodh/goServerScratch/database v0.0.0 => ./database

require (
	github.com/prajodh/goServerScratch/database v0.0.0
	golang.org/x/crypto v0.26.0
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.1 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
)

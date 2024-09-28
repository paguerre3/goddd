# goddd
DDD Onion project


---
### Project structure

```plaintext
padel-tournament/
│
├── cmd/
│   └── main.go                              # Entry point of the application
│
├── internal/modules/
│            ├── player-couple/                       # Player couple module
│            │   ├── api/
│            │   │   └── player_couple_handler.go     # REST handlers for player couple
│            │   ├── application/
│            │   │   └── player_couple_service.go     # Service layer for player couple
│            │   ├── domain/
│            │   │   └── player_couple.go             # Player couple domain entities
│            │   │   └── i_player_couple_repo.go      # Player couple repository interface
│            │   └── infrastructure/
│            │       └── mongo/
│            │           └── player_couple_repo.go    # MongoDB repository for player couple
│            │
│            ├── tournament/                          # Tournament module
│            │   ├── api/
│            │   │   └── tournament_handler.go        # REST handlers for tournament
│            │   ├── application/
│            │   │   └── tournament_service.go        # Service layer for tournament
│            │   ├── domain/
│            │   │   ├── tournament.go                # Tournament domain entities
│            │   │   └── i_tournament_repo.go         # Tournament repository interface
│            │   └── infrastructure/
│            │       └── mongo/
│            │           └── tournament_repo.go       # MongoDB repository for tournament
│            │
│            └── common/                              # Shared common utilities
│                ├── mongo/
│                │   └── mongo_client.go              # MongoDB client setup
│                └── utils/
│                    └── id_generator.go              # ID generation utility
│
├── docker-compose.yml                       # Docker Compose configuration
├── Dockerfile                               # Dockerfile for Go app
└── go.mod                                   # Go modules
```


---
### Requirements
1. ⚠️Docker must be running before executing Application.
2. <code>docker-compose -f docker-compose.yml up -d</code> for running tests and application. 
3. [Mongo Express URI](http://localhost:8081/)

***Optional***: Running under WSL needs allowing traffic through the firewall, i.e. 
using PS <code>New-NetFirewallRule -DisplayName "Allow MongoDB" -Direction Inbound -LocalPort 27017 -Protocol TCP -Action Allow</code>
and <code>New-NetFirewallRule -DisplayName "Allow MongoExpress" -Direction Inbound -LocalPort 8081 -Protocol TCP -Action Allow</code>.  


---
### DDD reading

[Mastering DDD Repository Design Patterns in Go](https://medium.com/@yohata/mastering-ddd-repository-design-patterns-in-go-2034486c82b3)

[DDD site reference](https://www.domainlanguage.com/ddd/reference/)

[DDD PDF reference](docs/DDD_Reference_2015-03.pdf)

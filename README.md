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
├── modules/
│   ├── player-couple/                       # Player couple module
│   │   ├── api/
│   │   │   └── player_couple_handler.go     # REST handlers for player couple
│   │   ├── application/
│   │   │   └── player_couple_service.go     # Service layer for player couple
│   │   ├── domain/
│   │   │   └── player_couple.go             # Player couple domain entities
│   │   │   └── i_player_couple_repo.go      # Player couple repository interface
│   │   └── infrastructure/
│   │       └── mongo/
│   │           └── player_couple_repo.go    # MongoDB repository for player couple
│   │
│   ├── tournament/                          # Tournament module
│   │   ├── api/
│   │   │   └── tournament_handler.go        # REST handlers for tournament
│   │   ├── application/
│   │   │   └── tournament_service.go        # Service layer for tournament
│   │   ├── domain/
│   │   │   ├── tournament.go                # Tournament domain entities
│   │   │   └── i_tournament_repo.go         # Tournament repository interface
│   │   └── infrastructure/
│   │       └── mongo/
│   │           └── tournament_repo.go       # MongoDB repository for tournament
│   │
│   └── common/                              # Shared common utilities
│       ├── mongo/
│       │   └── mongo_client.go              # MongoDB client setup
│       └── utils/
│           └── id_generator.go              # ID generation utility
│
├── docker-compose.yml                       # Docker Compose configuration
├── Dockerfile                               # Dockerfile for Go app
└── go.mod                                   # Go modules
```


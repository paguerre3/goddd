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

**Alternative 1: Using Docker isolated**

0. [Docker install](docs/0_docker-install-in-wsl.txt)
1. ⚠️Docker must be running before executing Application.
2. <code>docker-compose -f docker-compose.yml up -d</code> for running tests and application. 
3. [Mongo Express URI](http://localhost:8081/)

**Alternative 2: Using Docker, Kubectl and Minikube (K8s)**

0. [Docker install](docs/0_docker-install-in-wsl.txt)
1. [Kubectl and Minikube install](docs/1_minikube-install.txt)
2. [Build Docker image and publish it to Dockerhub *(Already done)*](docs/2_build_docker_image_and_publish_it.txt)
3. K8s deployment, i.e. <code>kubectl apply</code> *in order*:
Namespace 1st, then deployments
```bash
kubectl apply -f ./deployments/k8s/goddd-namespace.yaml
```
Mongodb deployment
```bash
kubectl apply -f ./deployments/k8s/mongodb-secret.yaml ./deployments/k8s/mongodb-deployment.yaml
```
Mongo-express deployment
```bash
kubectl apply -f ./deployments/k8s/mongo-express-deployment.yaml
```
Padel-place deployment and ingress
```bash
kubectl apply -f ./deployments/k8s/padel-place-*.yaml
```
4. [Mongo Express URI](http://localhost:8081/)

***Optional***: Running under WSL needs allowing traffic through the firewall, i.e. 
using PS <code>New-NetFirewallRule -DisplayName "Allow MongoDB" -Direction Inbound -LocalPort 27017 -Protocol TCP -Action Allow</code>
and <code>New-NetFirewallRule -DisplayName "Allow MongoExpress" -Direction Inbound -LocalPort 8081 -Protocol TCP -Action Allow</code>.  



---
### DDD reading

[Mastering DDD Repository Design Patterns in Go](https://medium.com/@yohata/mastering-ddd-repository-design-patterns-in-go-2034486c82b3)

[DDD site reference](https://www.domainlanguage.com/ddd/reference/)

[DDD PDF reference](docs/DDD_Reference_2015-03.pdf)

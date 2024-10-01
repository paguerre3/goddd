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

    **3.1.** Namespace 1st, then deployments
    ```bash
    kubectl apply -f ./deployments/k8s/goddd-namespace.yaml
    ```
    **3.2.** Create Secret for the Namespace
    ```bash
    kubectl apply -f ./deployments/k8s/mongodb-secret.yaml --namespace=goddd
    ```
    **3.3.** Mongodb deployment
    ```bash
    kubectl apply -f ./deployments/k8s/mongodb-deployment.yaml --namespace=goddd
    ```
    **3.4.** Mongo-express deployment *(check **⚠️4**)*
    ```bash
    kubectl apply -f ./deployments/k8s/mongo-express-deployment.yaml --namespace=goddd
    ```
    **3.5.** Padel-place deployment & ingress *(check Pre-requisite **⚠️5**)*
    ```bash
    kubectl apply -f ./deployments/k8s/padel-place-*.yaml --namespace=goddd
    ```
4. ***⚠️ Only for Minikube in case of Mongo-express external access)***: It shows "pending" EXTERNAL IP because of the usage of Minikube (using k8s directly should display external IP right away). Is needed to additionally execute "manually" <code>minikube service mongo-express-service</code> so Minikube assigns the external IP to the ExternalService of mongoexpress already defined, 
e.g. using docker driver and tunneling
    ```bash	
    minikube service mongo-express-service --namespace=goddd
    ```    
    [Mongo Express URI (port range from 30000 -it will be assigned by Minikube)](http://127.0.0.1:30000/)

5. ***⚠️ Only for Minikube in case of Padel-place  external access -Pre-requisite before applying [padle-place-ingress.yaml](deployments/k8s/padel-place-ingress.yaml)***: `Ingress` is used in real production environments where having `ExternalService` for exposing "external" IP isn't adequate, i.e. normally an application is accessed setting its domain name and secured port through the client browser. `Ingres` configuration is of `kind`: `Ingress`, and then `spec` contains a section rules where "routing rules" have defined -host domain addresses that receive requests and forwards to `http: paths: -backend: serviceName and servicePort`, i.e. Routing rules forward requests to `InternalService/s`. In other words in `Ingress` its defined a mapping that forwards requests from `Host` to `InternalService`. ⚠️*Warning, Ingress = `spe: rules: http` doesn't correspond to the "external communication protocol" that public URL uses, e.g. HTTPS or HTTP of my-domain, and instead it belongs to the "internal protocol" being used for forwarding the requests to `InternalService`*. `host` present in routes of `Ingress` should be a valid domain address as it maps domain name to a Node's IP address which is considered the "entry point" OR, alternatively, host maps domain name to a server outside of k8s cluster that acts like a Proxy or Secured Gateway that behaves as "entry point", i.e. Ingress will receive request from the internal or external "entry point"/host and then it will forward to `InternalService`.

    In order to work `Ingress` needs an "implementation" of it which is an `IngressController` Pod or set of Pods, i.e. `IngressController` runs on Pod or a Set of Pods of a Node in k8s cluster and does the "evaluation and processing" of `Ingress rules`. In other words `IngressController` is the "entry point" of k8s cluster that evaluates all rules and manages redirections. There are many third-party implementations of `IngressControllers` and k8s also offers its own implementation which is **"Nginx"** `IngressController` therefore it needs to be installed so Ingress can function
    ![Ingress Controller Implementation](https://github.com/paguerre3/kubeops/blob/master/support/22-ingress-controller.PNG)

    *if k8s cluster runs under a Cloud Service Provider like AWS or Google Cloud that have out-of-the-box kubernetes solutions or that use their own virtualized load balancer then there is normally a Cloud Provider "Load Balancer" placed in front of k8s cluster that behaves as a Secured Load Balancer "entry point" that receives and forwards requests to the `IngresController` of k8s, e.g.*
    ![AWS Ingress Contrtoller Implementation](https://github.com/paguerre3/kubeops/blob/master/support/23-ingress-controller-cloud-provider.PNG)

    **5.1.** Install `IngressController` in Minikube so `Ingress` can work, i.e. it automatically starts the Nginx implementation of `IngressController`
    ```bash
    minikube addons enable ingress
    * After the addon is enabled, please run "minikube tunnel" and your ingress resources would be available at "127.0.0.1"
    * Verifying ingress addon...
    * The 'ingress' addon is enabled
    ```
    "At the end of all steps" enable tunneling if needed for testing purposes `minikube tunnel`.
    
    Check Nginx `IngressController` is running under "ingress-nginx" NS
    ```bash
    kubectl get pod -n ingress-nginx
    NAME                                       READY   STATUS      RESTARTS   AGE
    ingress-nginx-admission-create-k9vr6       0/1     Completed   0          5m55s
    ingress-nginx-admission-patch-q8zwq        0/1     Completed   1          5m55s
    ingress-nginx-controller-bc57996ff-ncnnn   1/1     Running     0          5m55s 
    ```    
    **5.2.** Enable k8s dashboard and metrics-server (dependency of dashboard) in Minikube to do a Demo of `Ingress` configuration, i.e. execute `minikube addons enable dashboard` and then `minikube addons enable metrics-server`. To check the list of Minikube enabled `addons minikube addons list`
    ```bash
    minikube addons enable dashboard
    ```
    ```bash
    minikube addons enable metrics-server
    ```
    ```bash
    minikube addons list
    ``` 
    k8s dashboard has `InternalService` and Pod already configured but it doesn't have a `Ingress/IngressController` enabled.

    Check NS in order to visualize the one associated to the dashboard
    ```bash
    kubectl get namespace
    NAME                   STATUS   AGE
    default                Active   3d5h
    goddd                  Active   25h
    kube-node-lease        Active   3d5h
    kube-public            Active   3d5h
    kube-system            Active   3d5h
    kubernetes-dashboard   Active   23m	
    ```	
    Check configurations filtered by dashboard NS (`Ingress` "rule" configuration isn't present).
    ```bash	
    kubectl get all -n kubernetes-dashboard
    NAME                                            READY   STATUS    RESTARTS   AGE
    pod/dashboard-metrics-scraper-c95fcf479-kk692   1/1     Running   0          25m
    pod/kubernetes-dashboard-6cff4c7c4f-9vpvg       1/1     Running   0          25m
    NAME                                TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)    AGE
    service/dashboard-metrics-scraper   ClusterIP   10.100.171.241   none          8000/TCP   25m
    service/kubernetes-dashboard        ClusterIP   10.99.44.121     none          80/TCP     25m
    NAME                                        READY   UP-TO-DATE   AVAILABLE   AGE
    deployment.apps/dashboard-metrics-scraper   1/1     1            1           25m
    deployment.apps/kubernetes-dashboard        1/1     1            1           25m
    NAME                                                  DESIRED   CURRENT   READY   AGE
    replicaset.apps/dashboard-metrics-scraper-c95fcf479   1         1         1       25m
    replicaset.apps/kubernetes-dashboard-6cff4c7c4f       1         1         1       25m
    ```

    **5.3.** Create `Ingress` "rule" resource for k8s dashboard, i.e. [dashboard Ingress](deployments/k8s/dashboard-ingress.yaml) and execute
    ```bash
    kubectl apply -f ./deployments/k8s/dashboard-ingress.yaml
    ```
    Check Ingress rule creation
    ```bash
    kubectl get ingress -n kubernetes-dashboard
    NAME                CLASS    HOSTS           ADDRESS        PORTS   AGE
    dashboard-ingress   none     dashboard.com   192.168.49.2   80      3m7s
    ```

    **5.4.** Emulate "entry point" that behaves as a Proxy in front of `IngressController` outside k8s cluster so `IngressController` can use "dashboard" Ingress rule to evaluate and manage redirection (forwarding requests to "dashboard" `InternalService`), i.e. go to "hosts" file of os and create dns rule that matches with HOST and IP address of dashboard-ingress, e.g.
    
    ![hosts dns update](https://github.com/paguerre3/kubeops/blob/master/support/25-hosts-as-proxy.PNG)

    Or editing `/etc/hosts` file with IP address of dashboard-ingress and HOST name of "dashboard.com" under Linux distro OS.

    **5.5.** ***⚠️ Only if running on embedded Linux distro: Enable port forwarding using `Kubctl`:***
    ```bash    
    kubectl port-forward svc/kubernetes-dashboard -n kubernetes-dashboard 8001:80
    ```
    ***In this case edit host file having dashboard pointing to `127.0.0.1:8001` and execute `ipconfig /flushdns`***

    **5.6.** Open browser, write domain "dashboard.com" *(or "dashboard.com:8001" in case port forwarding)* and check k8s dashboard.

    **5.7.** Repeat steps 5.3 to 5.6 for "padel-place.com" using/applying [padel-place-ingress.yaml](/deployments/k8s/padel-place-ingress.yaml)
    ```bash
    kubectl apply -f ./deployments/k8s/padel-place-ingress.yaml --namespace=goddd
    ```
    ```bash
    kubectl get ingress -n goddd
    ```
    ```bash    
    kubectl port-forward svc/padel-place-service -n goddd 8001:8080
    ```
    APIs will be available under "padel-place.com" *(or "padel-place.com:8001" in case port forwarding)*.

***Optional***: Running under WSL needs allowing traffic through the firewall, i.e. 
using PS <code>New-NetFirewallRule -DisplayName "Allow MongoDB" -Direction Inbound -LocalPort 27017 -Protocol TCP -Action Allow</code>
and <code>New-NetFirewallRule -DisplayName "Allow MongoExpress" -Direction Inbound -LocalPort 8081 -Protocol TCP -Action Allow</code>.


---
### DDD reading

[Mastering DDD Repository Design Patterns in Go](https://medium.com/@yohata/mastering-ddd-repository-design-patterns-in-go-2034486c82b3)

[DDD site reference](https://www.domainlanguage.com/ddd/reference/)

[DDD PDF reference](docs/DDD_Reference_2015-03.pdf)

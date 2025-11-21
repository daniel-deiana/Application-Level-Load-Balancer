#### Description
This project implements a simple application level load balancer. A backend can register to the load balancer to be added to the pool of backends to forward the api traffic to them. The backends expose an health API so that the Lb can monitor their state and act if one or more of them goes down.

#### To test
in the terminal run `docker compose build` to build the containter images of the services described in the compose. After this `docker compose up -d ` to create the newtork of containers and get them running. To see if client request are correctly being server get inside the test containter and run the test script.

#### Things to note
If one tries to scale the backend service of the compose (by using `docker compose up scale backends=n`) the docker engine will return an error code telling us we are trying to bind a port that is already used. This is because we are actually exposing the backend port `8081` to the outside (in this case the outside is the host networking stack). In fact in the compose we specified: 

    ports: 
      - 8081:8081

This means actually that we are telling the host routing table to add a mapping between host_ip:8081 to containter_ip:8081, when the second replica of the backend service is started we fall into this error because the mapping is then ambigous, to resolve. Since the backend service does not need to be exposed to the host network we just do not expose the port 8081 to the outside, so that we have no conflicts.
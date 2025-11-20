#### Description
This project implements a simple application level load balancer. A backend can register to the load balancer to be added to the pool of backends to forward the api traffic to them. The backends expose an health API so that the Lb can monitor their state and act if one or more of them goes down.

#### To test
in the terminal run `docker compose build` to build the containter images of the services described in the compose. After this `docker compose up -d ` to create the newtork of containers and get them running. To see if client request are correctly being server get inside the test containter and run the test script.

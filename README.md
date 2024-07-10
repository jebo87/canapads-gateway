
# canapads-gateway
This repository is part of Canapads.ca.  This is a website being created to provide access to rental properties in Canada. 

![canapads_polygon](https://github.com/jebo87/canapads-gateway/assets/7273200/8247c1a9-bd63-4d57-bada-f3c1cb98348c)

This go microservice acts as the api gateway for all online services. Receives the requests from the web client and retreives the information from the backend API. Provides middleware to validate requests.

Frontend is being built using ReactJS and the backend is composed of microservices written in Go and NodeJS

We also use ElasticSearch for searching, Postgres as a Database, Mongo DB for user storage, ORY Hydra as authorization server 

This is a work in progress in early stages and I am currently working in what you see in the diagram below, however more micro services will be added in the future.



![canapads architecture](https://github.com/jebo87/canapads-gateway/assets/7273200/ef55d5f0-77c2-4dac-95df-ffbfbbd1a691)

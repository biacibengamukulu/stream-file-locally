
# My Spec
I want the user to be able to upload files to the server using a POST request with a JSON body containing the base64 encoded content of the file. The server should respond with a JSON object containing the file information upon successful upload.
please provide two ways to save file:
1. Save file to disk (file system) this will run on the docker container so make sure to use a volume and it is running localhost provide default directory
2. Save file to database (cassandra ) the database is running default on safer.easipath.com but you will use the config from the environment variables, make sure to user migration as i dont have access to the server cqlsh

Please follow the best practices for the project, and make sure to use the best tools for the job.
Also i provide the DDD approach for the project if you can follow it.
Please complete my code i did provide the basic structure.
Create also the public stream handler which will receive the url and stream the file to the user.

Please use docker-compose to run the project. and docker username is: 010309  so the conatiner name will be:  010309/stream-file-locally:latest

Please create api spec and swagger documentation for the project. and api document for integration and usage.

for deployment server user "ssh safer" and put the docker-compose.yml file in the  directory /apps/docker-compose-script/stream-file-locally

Please use the best practices for the project.
Happy coding!
# morgan
morgan is a suite of CLI automation commands. The name morgan comes from a TV series called "Chuck". morgan is the best friend of the main character "Chuck".

This tool is intented to be light weight automation tool. It's focused on two goals. #1 provide value to the user. #2 try not to reinvent the wheel. Having said that, we can take a look at the following AWS command as an example.

AWS CLI command for creating ECS service goes like this:
```
aws ecs create-service \
    --cluster MyCluster \
    --service-name ecs-simple-service \
    --task-definition sleep360:2 \
    --desired-count 1
```

The value morgan is providing in the following is creating and registering a minimal task definition. Using sensible defaults and naming convention can speed up deployment of services during development. 
```
morgan aws ecs create-service \
    ecs-simple-service small 8080 nginx:latest \    
    --cluster MyCluster \
    --desired-count 1
```

## Installation
```
go get github.com/7onetella/morgan
```

Here is more examples

![ECS Opearation](/asset/ecs-crud.png)
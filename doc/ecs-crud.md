morgan aws ec2 start ecs

morgan aws ecs create-service hello-world small 80 nginx:latest --cluster Development

morgan aws ecs describe-services hello-world --cluster Development

morgan aws ecs update-service hello-world latest --cluster Development

morgan aws ecs stop-service hello-world

morgan aws ecs describe-services hello-world --cluster Development

morgan aws ecs delete-service hello-world

morgan aws ec2 stop ecs
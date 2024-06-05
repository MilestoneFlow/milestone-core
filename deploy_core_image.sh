aws ecr-public get-login-password --region us-east-1 | docker login --username AWS --password-stdin public.ecr.aws/h4j3w3q7
docker build -t milestone-test .
docker tag milestone-test:latest public.ecr.aws/h4j3w3q7/milestone-test:latest
docker push public.ecr.aws/h4j3w3q7/milestone-test:latest
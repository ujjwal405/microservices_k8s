This project has Authentication, Mail and Product service.

Authentication service: which authenticates the user and generates OTP and stores the otp of specific user for 30 sec of its creation. It then adds the user email and otp to the RabbitMQ.

Mail service: which receives message from the RabbitMQ and sends Email to the respective user.

Product service: which helps to add  and display the product for admin role  but displays product for the user.


Both Authentication and Product service has their own MongoDB as a Database.



To run this project successfully please update the mail service under which update  

FROM_EMAIL_ADDRESS=  // which the email address of the mailer.
FROM_PASSWORD_EMAIL=  // which is the password for the mail service to send email.

in docker-compose file.



// For kubernetes

steps:
1 kubectl create ns app

2 kubectl create ns db

3 helm repo add sealed-secrets https://bitnami-labs.github.io/sealed-secrets

4 helm repo update

5 helm install sealed-secrets-controller --namespace kube-system  sealed-secrets/sealed-secrets

6 create  k8s secret for auth, product, mail, mongodb and rabbitmq  in db namespace inside secret folder in root directory.

7 Get access public key from the controller using below command

   kubeseal --fetch-cert --controller-name <name of controller> --controller-namespace kube-system
   
8 

9  kubeseal -f ./secrets/ -w ./kubernetes/sealed_secret/

10  kubectl apply -f ./kubernetes/sealed_secret/

11 kubectl apply -f ./kubernetes/authentication_mongodb/ -n db

12 kubectl apply -f ./kubernetes/product_mongodb/ -n db

13 kubectl apply -f ./kubernetes/rabbitmq/ -n db

14 helm install my-microservice ./kubernetes/helm/my-product -n app

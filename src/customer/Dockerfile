FROM golang:1.8.1

ENV GOPATH /go

ENV GOENV local

ENV PATH $GOPATH/bin:$PATH

RUN mkdir -p /go/src/customer

ADD . /go/src/customer

WORKDIR /go/src/customer

RUN apt-get install wget

RUN wget https://raw.githubusercontent.com/pote/gpm/v1.4.0/bin/gpm

RUN chmod +x gpm

RUN mv gpm /usr/local/bin

RUN gpm install

RUN go get github.com/beego/bee

# RUN go build -o /go/src/customer/customer /go/src/customer/main.go

# RUN bee run customer

# EXPOSE 8081

CMD ["/go/bin/bee","run","customer"]


# flow at minikube
# 1. minikube start
# 2. eval $(minikube docker-env)
# 3. docker build -t customer:v1 .
# 4. kubectl run customer --image=customer:v1 --port=8081 --replicas=4
# 5. kubectl create -f ./app.yaml
# 6. minikube service customer

# access to pods => kubectl exec -t -i <namepods> /bin/bash

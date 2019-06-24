pipeline {
    agent none
    stages {
        stage('Build') {
            agent {
                docker { image 'obraun/vss-protoactor-jenkins' }
            }
            steps {
                sh 'go get github.com/golang/protobuf/protoc-gen-go'
                sh 'go get github.com/micro/protoc-gen-micro'
                sh 'cd cinema_hall/proto && protoc --micro_out=. --go_out=. cinema_hall.proto'
                sh 'cd movie/proto && protoc --micro_out=. --go_out=. movie.proto'
                sh 'cd cinema_showing/proto && protoc --micro_out=. --go_out=. cinema_showing.proto'
                sh 'cd user/proto && protoc --micro_out=. --go_out=. user.proto'
                sh 'cd reservation/proto && protoc --micro_out=. --go_out=. reservation.proto'
                sh 'cd cinema_hall && go build main.go'
                sh 'cd movie && go build main.go'
                sh 'cd cinema_showing && go build main.go'
                sh 'cd user && go build main.go'
                sh 'cd reservation && go build main.go'
                sh 'cd client && go build main.go'
            }
        }
        stage('Lint') {
            agent {
                docker { image 'obraun/vss-protoactor-jenkins' }
            }   
            steps {
                sh 'golangci-lint run --deadline 20m --enable-all'
            }
        }
        stage('Test') {
            agent {
                docker { image 'obraun/vss-protoactor-jenkins' }
            }
            steps {
                sh 'go test ./...'
            }
        }
        stage('Build Docker Image') {
                    agent any
                    steps {
                        sh "docker-build-and-push -b ${BRANCH_NAME} -s cinema_hall -f cinema_hall/cinema_hall.dockerfile"
                        sh "docker-build-and-push -b ${BRANCH_NAME} -s movie -f movie/movie.dockerfile"
                        sh "docker-build-and-push -b ${BRANCH_NAME} -s cinema_showing -f cinema_showing/cinema_showing.dockerfile"
                        sh "docker-build-and-push -b ${BRANCH_NAME} -s user -f user/user.dockerfile"
                        sh "docker-build-and-push -b ${BRANCH_NAME} -s reservation -f reservation/reservation.dockerfile"
                        sh "docker-build-and-push -b ${BRANCH_NAME} -s client -f client/client.dockerfile"
                    }
                }
    }
}

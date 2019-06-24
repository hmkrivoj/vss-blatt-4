pipeline {
    agent none
    stages {
        stage('Build') {
            agent {
                docker { image 'obraun/vss-protoactor-jenkins' }
            }
            steps {
                sh 'go get github.com/micro/protoc-gen-micro'
                sh 'cd cinema_hall/proto && protoc --micro_out=. --go_out=. cinema_hall.proto'
                sh 'cd cinema_hall && go build main.go'
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
        stage('Build Docker Image') {
                    agent any
                    steps {
                        sh "docker-build-and-push -b ${BRANCH_NAME} -s cinema_hall -f cinema_hall/cinema_hall.dockerfile"
                    }
                }
    }
}

pipeline {
    agent any
    tools{
        go '1.21.2'
    }
    stages {
        stage('Github-Clone') {
            steps {
                git branch: 'dev',
                credentialsId: 'github-token',
                url: 'https://github.com/off-chain-storage/GoSphere.git'
            }
        }
        stage('Build') {
            steps {
                sh 'go version'
                sh 'go build -o GoSphere ./cmd/GoSphere/main.go'
            }
        }
        stage('Build Docker Image') {
            steps {
                sh 'docker image build -t jinbum99/gosphere .'
            }
        }
        stage('Push Docker Image') {
            steps {
                withCredentials([usernamePassword(credentialsId: 'DOCKER_REGISTRY_CREDENTIAL_ID', usernameVariable: 'DOCKER_USERNAME', passwordVariable: 'DOCKER_PASSWORD')]) {
                    sh "echo $DOCKER_PASSWORD | docker login --username $DOCKER_USERNAME --password-stdin"
                    sh 'docker push jinbum99/gosphere'
                }
            }
        }
        /*
        stage('Deploy') {
            steps {
                sh 'cd var/jenkins_home/projecty/deploy'
                sh 'chmod 777 ./deploy.sh'
                sh './deploy.sh'
            }
        }
        */
    }
}

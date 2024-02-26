pipeline {
    agent any
    tools{
        go 'go1.21.2'
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
                script {
                    sh """
                        sh 'go version'
                        sh 'go build -o GoSphere ./cmd/GoSphere/main.go'
                    """
                }
            }
        }
        stage('Build Docker Image') {
            steps {
                script {
                    sh """
                        sh 'docker image build -t jinbum99/GoSphere .'
                    """
                }
            }
        }
        stage('Push Docker Image') {
            steps {
                script {
                    withCredentials([usernamePassword(credentialsId: 'DOCKER_REGISTRY_CREDENTIALS_ID', usernameVariable: 'DOCKER_USERNAME', passwordVariable: 'DOCKER_PASSWORD')]) {
                        sh """
                            echo $DOCKER_PASSWORD | docker login --username $DOCKER_USERNAME --password-stdin
                            docker push jinbum99/GoSphere
                        """
                    }
                }
            }
        }
        /*
        stage('Deploy') {
            steps {
                script {
                    sh """
                        cd var/jenkins_home/projecty/deploy
                        chmod 777 ./deploy.sh
                        ./deploy.sh
                    """
                }
            }
        }
        */
    }
}
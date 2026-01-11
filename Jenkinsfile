pipeline {
    agent {
        node {
            label 'docker-agent-alpine'
        }
    }
    triggers {
        pollSCM '* * * * *'
    }
    stages {
        stage('Build')  {
            steps {
                echo "Building..."
                sh '''
                    make build
                '''
            }
        }
        stage('Test') {
            steps {
                echo "Testing..."
                sh '''
                    make test
                '''
            }
        }
        stage("Delivery") {
            steps {
                echo "Delivery"
                sh '''
                    make docker-run
                '''

            }
        }
    }
}

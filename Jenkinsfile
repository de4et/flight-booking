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
                    echo "doing build stuff"; sleep 5;
                '''
            }
        }
        stage('Test') {
            steps {
                echo "Testing..."
            }
        }
        stage("Delivery") {
            steps {
                echo "Delivery"
            }
        }
    }
}

pipeline {
    agent any
    options { timestamps () }

    stages {
        stage('Test') {
            tools {
                go '1.25.5'
            }
            steps {
                sh 'go test ./... -v'
            }
        }

        stage('Build') {
            steps {
                script {
                    docker.withRegistry('', 'docker-hub') {
                        def image = docker.build("de4et/flight-booking:${GIT_COMMIT}")
                        image.push()
                    }
                }
            }
        }

        stage('Deploy') {
            steps {
                withCredentials([string(credentialsId: 'HOME_IP', variable: 'HOME_IP')]) {
                    sshagent(['deploy-key']) {
                        sh """
                            ssh -o StrictHostKeyChecking=no deploy@$HOME_IP "
                                cd deploy
                                docker pull de4et/flight-booking:${GIT_COMMIT}
                                docker stop app
                                docker run \
                                    -d \
                                    --rm \
                                    --network flight-booking_blueprint \
                                    --name app \
                                    --env-file .env \
                                    -p 8081:8080 \
                                    de4et/flight-booking:${GIT_COMMIT}
                                if [ \"\$?\" = \"1\" ]; then echo \"Started successfully\"; else echo \"Error occured $?\"; fi
                            "
                        """
                    }
                }
            }
        }
    }
}

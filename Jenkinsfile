pipeline {
    agent any

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
            docker.withRegistry('registry.hub.docker.com', 'docker-hub') {
                def image = docker.build("de4et/flight-booking:${GIT_COMMIT}")
                image.push()
            }
        }

        stage('Deploy') {
            steps {
                sshagent(['deploy-key']) {
                    sh '''
                        ssh deploy@${HOME_IP} "
                            cd deploy
                            docker pull de4et/flight-booking:${GIT_COMMIT}
                            docker run \
                                -d \
                                --network flight-booking_blueprint \
                                --name app \
                                --env-file .env \
                                -p 8081:8080 \
                                de4et/flight-booking:${GIT_COMMIT}
                        "
                    '''
                }
            }
        }
    }
}

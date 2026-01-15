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
            tools {
                go '1.25.5'
            }
            steps {
                sh 'docker build -t de4et/flight-booking:${GIT_COMMIT} .'
                sh 'docker push de4et/flight-booking:${GIT_COMMIT}'
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

pipeline {
    agent none

    stages {
        stage('Build') {
            agent {
                node { label 'docker-agent-golang' }
            }
            tools {
                go '1.25.5'
            }
            steps {
                sh 'go build -o main cmd/api/main.go'
            }
        }

        stage('Test') {
            agent {
                node { label 'docker-agent-golang' }
            }
            tools {
                go '1.25.5'
            }
            steps {
                sh 'go test ./... -v'
            }
        }

        stage('Deploy') {
            agent {node {
                label 'master'
            }}
            steps {
                script {
                    sh '''
                        cat > .env.production << 'EOF'
                        PORT=8080
                        APP_ENV=production
                        BLUEPRINT_DB_HOST=psql_bp
                        BLUEPRINT_DB_PORT=5432
                        BLUEPRINT_DB_DATABASE=blueprint
                        BLUEPRINT_DB_USERNAME=melkey
                        BLUEPRINT_DB_PASSWORD=password1234
                        BLUEPRINT_DB_SCHEMA=public
                        REDIS_HOST=redis
                        REDIS_PORT=6379
                        REDIS_PASSWORD=hkjchzcxvysdafas2345345akljkjkbz
                        GZIP_LEVEL=6
                        EOF
                    '''

                    // Stop and remove old container
                    sh '''
                        export DOCKER_HOST=unix:///var/run/docker.sock
                        docker build -t my-app:latest .
                        docker stop my-app || true
                        docker rm my-app || true
                    '''

                    // Run new container
                    sh '''
                            docker run -d \
                            --name my-app \
                            --network blueprint \
                            --env-file .env.production \
                            -p 8080:8080 \
                            my-app:latest
                    '''
                }
            }
        }
    }
}

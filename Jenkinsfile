pipeline {
    agent none

    stages {
        // stage('Build') {
        //     agent {
        //         node { label 'docker-agent-golang' }
        //     }
        //     tools {
        //         go '1.25.5'
        //     }
        //     steps {
        //         sh 'go build -o main cmd/api/main.go'
        //     }
        // }
        //
        // stage('Test') {
        //     agent {
        //         node { label 'docker-agent-golang' }
        //     }
        //     tools {
        //         go '1.25.5'
        //     }
        //     steps {
        //         sh 'go test ./... -v'
        //     }
        // }

        stage('Deploy') {
            agent {node {
                label 'master'
            }}
            // environment {
            //     DOCKER_HOST = 'tcp://localhost:2375'
            //     // DOCKER_TLS_VERIFY = '0'
            // }
            steps {
                script {
                    sh '''
                        cat > .env << 'EOF'
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
                    // sh '''
                    //     export DOCKER_HOST=unix:///var/run/docker.sock
                    //     docker build -t my-app:latest .
                    //     docker stop my-app || true
                    //     docker rm my-app || true
                    // '''

                    // sh '''
                    //     docker compose build app
                    //     docker compose up -d --no-deps app
                    // '''
                    sh '''
                        echo "=== Docker Debug ==="
                        echo "DOCKER_HOST: $DOCKER_HOST"
                        echo "Checking docker socket..."
                        ls -la /var/run/docker.sock 2>/dev/null || echo "No socket"

                        # Try different approaches
                        echo "Trying direct docker command..."
                        /usr/bin/docker ps 2>&1 | head -5
                    '''

                    // Try multiple approaches
                    sh '''
                        # Approach 1: Clear DOCKER_HOST
                        unset DOCKER_HOST

                        # Approach 2: If Approach 1 fails, use full path
                        if ! docker compose build app; then
                            echo "Trying with full path..."
                            /usr/bin/docker compose build app
                            /usr/bin/docker compose up -d --no-deps app
                        else
                            docker compose up -d --no-deps app
                        fi
                    '''


                    // Run new container
                    // sh '''
                    //         docker run -d \
                    //         --name my-app \
                    //         --network blueprint \
                    //         --env-file .env.production \
                    //         -p 8080:8080 \
                    //         my-app:latest
                    // '''
                }
            }
        }
    }
}

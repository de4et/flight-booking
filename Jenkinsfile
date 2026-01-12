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
                        echo "=== Setting up environment ==="

                        # 1. Clear Docker environment variables
                        unset DOCKER_HOST
                        unset DOCKER_TLS_VERIFY

                        # 2. Check if Docker is available
                        echo "Testing Docker..."
                        if command -v docker >/dev/null 2>&1; then
                            echo "Docker found at: $(which docker)"
                        else
                            echo "ERROR: Docker not found!"
                            exit 1
                        fi

                        # 3. Try to use Docker
                        if docker ps >/dev/null 2>&1; then
                            echo "Docker is accessible!"
                        else
                            echo "ERROR: Cannot connect to Docker daemon"
                            echo "Trying different methods..."

                            # Try Unix socket (might be different location)
                            export DOCKER_HOST=unix:///var/run/docker.sock
                            docker ps 2>&1 | head -5 || echo "Unix socket failed"

                            # Try Docker Desktop socket
                            export DOCKER_HOST=unix:///var/run/docker-desktop/docker.sock
                            docker ps 2>&1 | head -5 || echo "Docker Desktop socket failed"

                            # Try socat if you have it
                            export DOCKER_HOST=tcp://localhost:2375
                            docker ps 2>&1 | head -5 || echo "TCP socket failed"

                            exit 1
                        fi

                        # 4. Now build and deploy
                        echo "Building and deploying..."
                        # Build and run
                        docker compose build app
                        docker compose up -d --no-deps app

                        echo "✅ Deployment complete!"

                    '''
sh '''
            # Test HTTP connection to socat
            echo "Testing socat connection..."
            curl -s http://docker:2375/version | jq -r '.Version'

            # Test Docker with HTTP
            DOCKER_HOST=tcp://docker:2375 docker ps

            # Test Docker with HTTP explicitly
            DOCKER_HOST=tcp://docker:2375 DOCKER_TLS_VERIFY=0 docker ps
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

pipeline {
    agent any
    environment {
        GO_VERSION = '1.24.4' // Specify your Go version
        ARTIFACTS_DIR = 'artifacts'
        BINARY_NAME = 'gobs' // Replace with your binary name
        PATH = "/usr/local/bin/go/bin:$PATH"
    }
    stages {
        stage('Checkout') {
            steps {
                checkout scm
            }
        }
        stage('Code Check and Test') {
            when {
                changeRequest target: 'main'
            }
            steps {
                sh 'go version'
                sh 'go fmt ./...'
                sh 'go vet ./...'
                sh 'go test ./... -cover'
            }
        }
        stage('Build and Export Artifacts') {
            when {
                anyOf {
                    branch 'main'
                    expression { return env.GIT_BRANCH =~ /^release\/.*/ }
                    expression { return env.GIT_TAG_NAME =~ /^v[0-9]+\.[0-9]+\.[0-9]+$/ }
                }
            }
            steps {
                sh 'mkdir -p ${ARTIFACTS_DIR}'
                sh 'GOOS=linux GOARCH=amd64 go build -o ${ARTIFACTS_DIR}/${BINARY_NAME} ./cmd/main.go' // Adjust path to main.go as needed
                archiveArtifacts artifacts: "${ARTIFACTS_DIR}/${BINARY_NAME}", fingerprint: true
            }
        }
    }
    post {
        always {
            cleanWs()
        }
        success {
            echo 'Pipeline completed successfully!'
        }
        failure {
            echo 'Pipeline failed!'
        }
    }
}

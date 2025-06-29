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

pipeline {
    environment {
        imageName = "araji/ghpr"
        dockerImage = ''
    }
    agent any
    stages {
      stage('Docker Build') {
         steps {
             script {
                echo "build image"
                dockerImage  = docker.build imageName
                }   
            }
        }

      stage('Security Scanning") {
         parallel {
             steps {
                script {
                    echo "running trivy"
                    sh '''
                        trivy -d -i dockerImage
                    '''
                }   
            steps {
                script {
                    echo "Running Anchore"
                    sh '''
                        trivy -d -i dockerImage
                    '''
                }
            }
        }
  
        stage('Push Container') {
            steps {
                echo "workspace is $WORKSPACE"
                dir("$WORKSPACE") {
                    script {
                        docker.withRegistry('https://index.docker.io/v1/','DockerHub') {
                            dockerImage.push("$BRANCH_NAME-$BUILD_NUMBER")
                            dockerImage.push("latest")
                        }
                    }
                }
            }`
        }
    }
}
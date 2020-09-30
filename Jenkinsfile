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
        stage('Push Container') {
            steps {
                echo "workspace is $WORKSPACE"
                dir("$WORKSPACE") {
                    script {
                        docker.withRegistry('https://index.docker.io/v1/','DockerHub') {
                            dockerImage.push("$BUILD_NUMBER")
                            dockerImage.push("latest")
                        }
                    }
                }
            }
        }
    }
}
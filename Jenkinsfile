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

      stage('Security Scanning') {
          parallel {
              stage('Trivy Analysis') {
                  steps {
                      echo "running trivy"
                      sh '''
                        sleep 30
                        '''
                    }
              }   
                stage('Anchore Analysis') {
                  steps {
                      echo "running trivy"
                      sh '''
                        sleep 30
                        '''
                    }
                }
          }
      }
        stage('Push Container') {
            steps {
                echo "workspace is $WORKSPACE"
                dir("$WORKSPACE") {
                    script {
                        docker.withRegistry('https://index.docker.io/v1/','DockerHub') {
                            dockerImage.push("$env.BRANCH_NAME-$BUILD_NUMBER")
                            dockerImage.push("latest")
                        }
                    }
                }
            }
        }
    }
}
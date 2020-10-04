pipeline {
    environment {
        imageName = "araji/ghpr"
        testImage = "araji/ghpr-test"
        dockerImage = ''
        dockerTestImage = ''
    }
    agent any
    stages {
       
      stage('Docker Build') {
         steps {
             script {
                echo "build prod image"
                dockerImage  = docker.build imageName
                }  
             script {
                echo "build test image"
                def mytest = docker.build(testImage, "--target builder .")
                }  
        }
      }
    stage('Run Tests'){
          agent {
              docker {image 'araji/ghpr-test'}
          }
          steps{
              sh 'cd /build;GOCACHE=/tmp go test  -v . '

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
                      echo "running Anchore"
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


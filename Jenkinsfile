pipeline {
   agent any

   stages {
      stage('Docker Build') {
         steps {
            sh(script: """
               docker build -t araji/ghpr .
            """)
      }
    }
    stage('Push Container') {
        steps {
            echo "workspace is $WORKSPACE"
            dir("$WORKSPACE") {
                script {
                    docker.withRegistry('https://index.docker.io/v1/','DockerHub') {
                        def image = docker.build('araji/ghpr:latest')
                        image.push()
                    }
                }
            }
        }
    }
   }
}
      
pipeline {
  agent {
    docker {
      image 'golangci/build-runner:latest'
    }

  }
  stages {
    stage('检出') {
      steps {
        checkout([$class: 'GitSCM', branches: [[name: env.GIT_BUILD_REF]], 
                                                                                                                                    userRemoteConfigs: [[url: env.GIT_REPO_URL, credentialsId: env.CREDENTIALS_ID]]])
        sh 'go version'
      }
    }
    stage('并行阶段 2') {
      parallel {
        stage('构建game') {
          steps {
            echo '构建game中...'
            sh '''cd apps/game/
make build'''
            archiveArtifacts(artifacts: 'apps/game/game', fingerprint: true)
            echo '构建game完成.'
          }
        }
        stage('构建gate') {
          steps {
            echo '构建gate中...'
            sh '''cd apps/gate/
make build'''
            archiveArtifacts(artifacts: 'apps/gate/gate', fingerprint: true)
            echo '构建gate完成'
          }
        }
      }
    }
    stage('测试') {
      parallel {
        stage('game单元测试') {
          steps {
            echo '单元测试中...'
            sh '''cd apps/game
make test'''
            echo '单元测试完成.'
          }
        }
        stage('gate单元测试') {
          steps {
            echo '单元测试中...'
            sh '''cd apps/gate
make test'''
            echo '单元测试完成.'
          }
        }
        stage('接口测试') {
          steps {
            echo '接口测试中...'
            echo '接口测试完成.'
          }
        }
      }
    }
  }
  environment {
    GO111MODULE = 'on'
    GOPROXY = 'https://goproxy.cn,direct'
  }
}

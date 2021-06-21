pipeline {
  agent any
  stages {
    stage('检出') {
      steps {
        checkout([$class: 'GitSCM', branches: [[name: env.GIT_BUILD_REF]], 
        userRemoteConfigs: [[url: env.GIT_REPO_URL, credentialsId: env.CREDENTIALS_ID]]])
      }
    }
    stage('构建并打包推送') {
      agent {
        docker {
          reuseNode 'true'
          registryUrl 'https://mmstudio-docker.pkg.coding.net'
          registryCredentialsId "${env.DOCKER_REGISTRY_CREDENTIALS_ID}"
          image 'blade/server/ci-building-base:latest'
          // 挂载外部虚拟机的 docker socket
          args '-v /var/run/docker.sock:/var/run/docker.sock -v /usr/bin/docker:/usr/bin/docker'
        }
      }
      
      steps {
        echo '构建中...'
        sh 'git config --global user.name "coding-ci"'
        sh 'git config --global user.email "hellodudu86@gmail.com"'
        sh 'git config --global url."https://${CODING_USERNAME}:${CODING_PASSWORD}@e.coding.net".insteadOf "https://e.coding.net"'
        sh 'make build'
        sh 'make docker'
        sh 'make push_coding'
        echo '构建成功...'
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
    GOPRIVATE = 'e.coding.net'
    GOPROXY = 'https://goproxy.cn,direct'
  }
}
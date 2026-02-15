const { defineConfig } = require('@tarojs/cli')
const path = require('path')
const devConfig = require('./dev')
const prodConfig = require('./prod')

module.exports = defineConfig({
  projectName: 'whatsdot-aibuddy-taro',
  date: '2026-02-15',
  designWidth: 390,
  deviceRatio: {
    390: 2,
    640: 2.34 / 2,
    750: 1,
    828: 1.81 / 2
  },
  sourceRoot: 'src',
  outputRoot: 'dist',
  framework: 'react',
  plugins: [],
  alias: {
    '@': path.resolve(__dirname, '..', 'src')
  },
  compiler: {
    type: 'webpack5',
    prebundle: {
      enable: false
    }
  },
  defineConstants: {
    'process.env.TARO_APP_API_BASE': JSON.stringify(process.env.TARO_APP_API_BASE || 'http://127.0.0.1:8080')
  },
  copy: {
    patterns: [],
    options: {}
  },
  cache: { enable: true },
  mini: {
    postcss: {
      pxtransform: { enable: true, config: {} },
      cssModules: { enable: false }
    }
  },
  h5: {
    publicPath: '/',
    staticDirectory: 'static',
    router: {
      mode: 'hash',
      customRoutes: {
        '/': 'pages/home/index'
      }
    },
    devServer: {
      host: '127.0.0.1',
      port: 10086
    },
    webpackChain(chain) {
      chain.merge({
        ignoreWarnings: [
          { module: /taro-video-core\.js/ },
          /webpackExports/
        ]
      })
    }
  }
}, (merge) => process.env.NODE_ENV === 'development' ? merge({}, devConfig) : merge({}, prodConfig))

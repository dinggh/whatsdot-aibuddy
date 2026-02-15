module.exports = {
  env: {
    NODE_ENV: 'development'
  },
  defineConstants: {
    'process.env.TARO_APP_API_BASE': JSON.stringify(process.env.TARO_APP_API_BASE || 'http://127.0.0.1:8080')
  },
  mini: {},
  h5: {}
}

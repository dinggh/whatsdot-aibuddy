module.exports = {
  entryPagePath: 'pages/home/index',
  router: {
    customRoutes: {
      '/': 'pages/home/index'
    }
  },
  pages: [
    'pages/home/index',
    'pages/loading/index',
    'pages/result/index',
    'pages/mode/index',
    'pages/history/index',
    'pages/profile/index'
  ],
  window: {
    navigationStyle: 'custom',
    backgroundTextStyle: 'dark',
    backgroundColor: '#F8F5F0'
  }
}

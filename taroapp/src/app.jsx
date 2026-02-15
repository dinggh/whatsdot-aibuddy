import './app.scss'

if (process.env.TARO_ENV === 'h5' && typeof window !== 'undefined') {
  const hash = window.location.hash || ''
  const pathname = window.location.pathname || ''
  const isPagePath = /^\/pages\/[^/]+\/index\/?$/.test(pathname)
  if (!hash && (pathname === '/' || isPagePath)) {
    window.location.replace('/#/pages/home/index')
  }
}

function App(props) {
  return props.children
}

export default App

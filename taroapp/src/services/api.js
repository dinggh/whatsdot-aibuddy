import Taro from '@tarojs/taro'

const API_BASE = process.env.TARO_APP_API_BASE || 'http://127.0.0.1:8080'
const TOKEN_KEY = 'auth_token'
const USER_KEY = 'auth_user'

function getToken() {
  return Taro.getStorageSync(TOKEN_KEY) || ''
}

function setSession(token, user) {
  Taro.setStorageSync(TOKEN_KEY, token || '')
  Taro.setStorageSync(USER_KEY, user || null)
}

function request(path, options = {}) {
  const { method = 'GET', data = null, auth = false } = options
  const header = { 'content-type': 'application/json' }
  if (auth) {
    const token = getToken()
    if (token) {
      header.Authorization = `Bearer ${token}`
    }
  }

  return Taro.request({
    url: `${API_BASE}${path}`,
    method,
    data,
    header
  }).then((res) => {
    if (res.statusCode >= 200 && res.statusCode < 300) {
      return res.data || {}
    }
    const msg = (res.data && res.data.error) || `HTTP ${res.statusCode}`
    throw new Error(msg)
  })
}

function isWeApp() {
  return Taro.getEnv() === Taro.ENV_TYPE.WEAPP
}

export async function ensureLogin() {
  const token = getToken()
  if (token) return token

  if (!isWeApp()) {
    throw new Error('H5 不支持微信登录，请在小程序环境测试')
  }

  const loginRes = await Taro.login()
  const code = loginRes && loginRes.code
  if (!code) {
    throw new Error('wx.login 未返回 code')
  }

  const data = await request('/api/v1/auth/wechat/login', {
    method: 'POST',
    data: { code }
  })

  if (!data.token) {
    throw new Error('后端未返回 token')
  }
  setSession(data.token, data.user || null)
  return data.token
}

export async function fetchMe() {
  await ensureLogin()
  const data = await request('/api/v1/me', { method: 'GET', auth: true })
  if (data.user) setSession(getToken(), data.user)
  return data.user || null
}

export async function updateProfile(nickName, avatarUrl) {
  await ensureLogin()
  const data = await request('/api/v1/auth/wechat/profile', {
    method: 'POST',
    auth: true,
    data: { nickName, avatarUrl }
  })
  if (data.user) setSession(getToken(), data.user)
  return data.user || null
}

export async function bindPhoneByCode(code) {
  await ensureLogin()
  const data = await request('/api/v1/auth/wechat/phone', {
    method: 'POST',
    auth: true,
    data: { code }
  })
  if (data.user) setSession(getToken(), data.user)
  return data.user || null
}

export async function fetchHistory() {
  await ensureLogin()
  const data = await request('/api/v1/history', { method: 'GET', auth: true })
  return data.items || []
}

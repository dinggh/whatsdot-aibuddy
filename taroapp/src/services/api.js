import Taro from '@tarojs/taro'

const API_BASE = process.env.TARO_APP_API_BASE || 'http://127.0.0.1:8080'
const DEVICE_KEY = 'device_id'
const MODE_KEY = 'analysis_mode'

const MODE_OPTIONS = [
  { key: 'guided', label: '引导思考', desc: '引导孩子一步一步思考（默认推荐）' },
  { key: 'detailed', label: '详细讲解', desc: '完整讲解解题过程和知识点' },
  { key: 'noanswer', label: '不给答案模式', desc: '只给思路和提示，不出现答案' },
  { key: 'quick', label: '快速提示', desc: '快速给出关键提示，节省时间' }
]

function ensureDeviceId() {
  let id = Taro.getStorageSync(DEVICE_KEY)
  if (!id) {
    id = `dev_${Date.now()}_${Math.random().toString(36).slice(2, 10)}`
    Taro.setStorageSync(DEVICE_KEY, id)
  }
  return id
}

function parseResp(payload, statusCode) {
  if (statusCode < 200 || statusCode >= 300) {
    throw new Error((payload && payload.message) || `HTTP ${statusCode}`)
  }
  if (!payload || typeof payload !== 'object') {
    throw new Error('empty response')
  }
  if (payload.code !== 0) {
    throw new Error(payload.message || 'request failed')
  }
  return payload.data || {}
}

function request(path, options = {}) {
  const { method = 'GET', data = null } = options
  const deviceId = ensureDeviceId()

  return Taro.request({
    url: `${API_BASE}${path}`,
    method,
    data,
    header: {
      'content-type': 'application/json',
      'X-Device-Id': deviceId
    }
  }).then((res) => parseResp(res.data, res.statusCode))
}

export function getModeOptions() {
  return MODE_OPTIONS
}

export function getCurrentMode() {
  const v = Taro.getStorageSync(MODE_KEY)
  if (MODE_OPTIONS.find((it) => it.key === v)) {
    return v
  }
  return 'guided'
}

export function setCurrentMode(mode) {
  const m = MODE_OPTIONS.find((it) => it.key === mode) ? mode : 'guided'
  Taro.setStorageSync(MODE_KEY, m)
  return m
}

export function modeLabel(mode) {
  const found = MODE_OPTIONS.find((it) => it.key === mode)
  return found ? found.label : '引导思考'
}

export function uploadHomework(imagePath, mode) {
  const m = setCurrentMode(mode)
  const deviceId = ensureDeviceId()

  return new Promise((resolve, reject) => {
    Taro.uploadFile({
      url: `${API_BASE}/api/v1/homework/analyze`,
      filePath: imagePath,
      name: 'image',
      formData: { mode: m },
      header: { 'X-Device-Id': deviceId },
      success: (res) => {
        try {
          const payload = JSON.parse(res.data || '{}')
          resolve(parseResp(payload, res.statusCode || 200))
        } catch (err) {
          reject(new Error('invalid upload response'))
        }
      },
      fail: (err) => reject(new Error(err.errMsg || 'upload failed'))
    })
  })
}

export async function fetchHistory() {
  const data = await request('/api/v1/history', { method: 'GET' })
  return data.items || []
}

export async function fetchHistoryDetail(id) {
  const data = await request(`/api/v1/history/${id}`, { method: 'GET' })
  return data.record || null
}

export async function regenerateHomework(id, mode) {
  const m = setCurrentMode(mode)
  const data = await request(`/api/v1/homework/${id}/regenerate?mode=${encodeURIComponent(m)}`, {
    method: 'POST',
    data: { mode: m }
  })
  return data.record || null
}

export function formatTime(input) {
  const d = input ? new Date(input) : new Date()
  const yyyy = d.getFullYear()
  const mm = String(d.getMonth() + 1).padStart(2, '0')
  const dd = String(d.getDate()).padStart(2, '0')
  const hh = String(d.getHours()).padStart(2, '0')
  const mi = String(d.getMinutes()).padStart(2, '0')
  return `${yyyy}-${mm}-${dd} ${hh}:${mi}`
}

export function buildAssetURL(path) {
  const p = String(path || '').trim()
  if (!p) return ''
  if (p.startsWith('http://') || p.startsWith('https://')) return p
  return `${API_BASE}${p.startsWith('/') ? '' : '/'}${p}`
}

// Backward compatible no-op APIs for profile page.
export async function ensureLogin() {
  return ensureDeviceId()
}

export async function fetchMe() {
  return {
    nickName: '家长用户',
    avatarUrl: '',
    phoneNumber: ''
  }
}

export async function updateProfile(nickName, avatarUrl) {
  return { nickName: nickName || '家长用户', avatarUrl: avatarUrl || '', phoneNumber: '' }
}

export async function bindPhoneByCode() {
  return { nickName: '家长用户', avatarUrl: '', phoneNumber: '138****0000' }
}

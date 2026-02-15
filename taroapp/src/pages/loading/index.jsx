import React from 'react'
import { View, Text, Image, Button } from '@tarojs/components'
import Taro, { useLoad } from '@tarojs/taro'

import { StatusBar } from '@/components/layout'
import { uploadHomework } from '@/services/api'
import '@/styles/common.scss'
import './index.scss'

const h = React.createElement

export default function LoadingPage() {
  const [preview, setPreview] = React.useState('')
  const [mode, setMode] = React.useState('guided')
  const [loading, setLoading] = React.useState(true)
  const [error, setError] = React.useState('')

  const submit = (path, selectedMode) => {
    setLoading(true)
    setError('')
    uploadHomework(path, selectedMode)
      .then((data) => {
        const id = data && data.record && data.record.id
        if (!id) throw new Error('后端未返回记录ID')
        Taro.redirectTo({ url: `/pages/result/index?id=${id}` })
      })
      .catch((err) => {
        setError(err.message || '上传失败')
      })
      .finally(() => setLoading(false))
  }

  useLoad((query) => {
    const imagePath = decodeURIComponent((query && query.imagePath) || '')
    const selectedMode = decodeURIComponent((query && query.mode) || 'guided')
    setPreview(imagePath)
    setMode(selectedMode)
    if (!imagePath) {
      setError('图片路径丢失，请返回重试')
      setLoading(false)
      return
    }
    submit(imagePath, selectedMode)
  })

  return h(View, { className: 'screen' },
    h(StatusBar),
    h(View, { className: 'nav-back', onClick: () => Taro.navigateBack() }, '←'),
    h(View, { className: 'loading-content' },
      preview ? h(Image, { className: 'loading-preview', mode: 'aspectFill', src: preview }) : null,
      h(View, { className: 'loading-ring' }, loading ? '✶' : '!'),
      h(Text, { className: 'loading-title' }, loading ? '正在识题目...' : '处理结束'),
      h(Text, { className: 'loading-sub' }, loading ? 'AI正在整理讲解方式...' : (error || '完成')),
      h(View, { className: 'loading-tip' }, '🔍  小贴士：引导孩子自己思考比直接告诉答案更有效')
    ),
    !loading && error ? h(Button, { className: 'loading-next', onClick: () => submit(preview, mode) }, '重试上传') : null
  )
}

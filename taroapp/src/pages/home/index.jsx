import React from 'react'
import { View, Text, Button } from '@tarojs/components'
import Taro, { useDidShow } from '@tarojs/taro'

import { BottomTabBar, StatusBar } from '@/components/layout'
import { ensureLogin } from '@/services/api'
import '@/styles/common.scss'
import './index.scss'

const h = React.createElement

export default function HomePage() {
  useDidShow(() => {
    ensureLogin().catch((err) => {
      Taro.showToast({ title: err.message || '登录失败', icon: 'none' })
    })
  })

  return h(View, { className: 'screen' },
    h(StatusBar),
    h(View, { className: 'home-content' },
      h(View, { className: 'home-header' },
        h(Text, { className: 'h1' }, '微点辅导助手'),
        h(Text, { className: 'sub' }, '让辅导作业变得简单轻松')
      ),
      h(Button, { className: 'camera-btn', onClick: () => Taro.navigateTo({ url: '/pages/loading/index' }) }, '拍作业'),
      h(View, { className: 'panel album-btn', onClick: () => Taro.navigateTo({ url: '/pages/loading/index' }) }, '从相册上传'),
      h(View, { className: 'home-tip' }, '帮助家长引导孩子思考，而不是直接给答案')
    ),
    h(BottomTabBar, { active: 'home' })
  )
}

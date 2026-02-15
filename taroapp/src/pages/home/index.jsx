import React from 'react'
import { View, Text, Button } from '@tarojs/components'
import Taro from '@tarojs/taro'

import { BottomTabBar, StatusBar } from '@/components/layout'
import { getCurrentMode } from '@/services/api'
import '@/styles/common.scss'
import './index.scss'

const h = React.createElement

export default function HomePage() {
  const pickImage = (sourceType) => {
    const mode = getCurrentMode()
    Taro.chooseImage({
      count: 1,
      sizeType: ['compressed'],
      sourceType
    }).then((res) => {
      const filePath = (res.tempFilePaths && res.tempFilePaths[0]) || ''
      if (!filePath) {
        Taro.showToast({ title: '未选择图片', icon: 'none' })
        return
      }
      const url = `/pages/loading/index?imagePath=${encodeURIComponent(filePath)}&mode=${encodeURIComponent(mode)}`
      Taro.navigateTo({ url })
    }).catch((err) => {
      if (!String(err.errMsg || '').includes('cancel')) {
        Taro.showToast({ title: '选择图片失败', icon: 'none' })
      }
    })
  }

  return h(View, { className: 'screen' },
    h(StatusBar),
    h(View, { className: 'home-content' },
      h(View, { className: 'home-header' },
        h(Text, { className: 'h1 home-title' }, '微点辅导助手'),
        h(Text, { className: 'sub home-sub' }, '让辅导作业变得简单轻松')
      ),
      h(Button, { className: 'camera-btn', onClick: () => pickImage(['camera']) },
        h(Text, { className: 'camera-icon' }, '◉'),
        h(Text, { className: 'camera-text' }, '拍作业')
      ),
      h(View, { className: 'panel album-btn', onClick: () => pickImage(['album']) }, '◩  从相册上传'),
      h(View, { className: 'home-tip' }, '💡  帮助家长引导孩子思考，而不是直接给答案')
    ),
    h(BottomTabBar, { active: 'home' })
  )
}

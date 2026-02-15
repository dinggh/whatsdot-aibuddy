import React from 'react'
import { View, Text, Image, Button } from '@tarojs/components'
import Taro from '@tarojs/taro'

import { StatusBar } from '@/components/layout'
import '@/styles/common.scss'
import './index.scss'
import preview from '@/assets/generated-1771138856711.png'

const h = React.createElement

export default function LoadingPage() {
  return h(View, { className: 'screen' },
    h(StatusBar),
    h(View, { className: 'nav-back', onClick: () => Taro.navigateBack() }, '←'),
    h(View, { className: 'loading-content' },
      h(Image, { className: 'loading-preview', mode: 'aspectFill', src: preview }),
      h(View, { className: 'loading-ring' }, '✶'),
      h(Text, { className: 'loading-title' }, '正在识别题目...'),
      h(Text, { className: 'loading-sub' }, 'AI正在整理讲解方式...'),
      h(View, { className: 'loading-tip' }, '小贴士：引导孩子自己思考比直接告诉答案更有效哦'),
      h(Button, { className: 'loading-next', onClick: () => Taro.navigateTo({ url: '/pages/result/index' }) }, '查看结果')
    )
  )
}

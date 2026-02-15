import React from 'react'
import { View, Text, Image, ScrollView, Button } from '@tarojs/components'
import Taro from '@tarojs/taro'

import { StatusBar } from '@/components/layout'
import '@/styles/common.scss'
import './index.scss'
import qimg from '@/assets/generated-1771138893602.png'

const h = React.createElement

function Card(title, body, cls = 'panel result-card') {
  return h(View, { className: cls },
    h(Text, { className: 'result-h' }, title),
    h(Text, { className: cls.includes('result-strong') ? 'result-p result-dark' : 'result-p' }, body)
  )
}

export default function ResultPage() {
  return h(View, { className: 'screen' },
    h(StatusBar),
    h(View, { className: 'result-topbar' },
      h(Text, { className: 'result-back', onClick: () => Taro.navigateBack() }, '←'),
      h(Text, { className: 'result-title' }, '讲解结果'),
      h(View, { className: 'result-mode', onClick: () => Taro.navigateTo({ url: '/pages/mode/index' }) }, '引导思考')
    ),
    h(ScrollView, { className: 'result-content', scrollY: true },
      h(View, { className: 'panel result-card' },
        h(Text, { className: 'result-card-title' }, '题目 · 三年级'),
        h(Image, { className: 'result-qimg', mode: 'aspectFill', src: qimg }),
        h(Text, { className: 'result-qtext' }, '24 x 15 = ?')
      ),
      Card('解题思路（给家长看）', '这道题考查的是两位数乘法。可用竖式计算：1) 24x5=120；2) 24x10=240；3) 相加=360。'),
      Card('讲给孩子听（语气简单）', '把 15 拆成 10 和 5，再分别和 24 相乘，最后把结果加起来。'),
      Card('家长可以这样引导（重点）', '先问怎么拆 15，再引导算 24x5 和 24x10，最后让孩子自己说答案。', 'panel result-card result-strong')
    ),
    h(View, { className: 'result-actions' },
      h(Button, { className: 'result-btn result-btn-primary', onClick: () => Taro.navigateTo({ url: '/pages/history/index' }) }, '保存'),
      h(Button, { className: 'result-btn' }, '分享'),
      h(Button, { className: 'result-btn' }, '复制')
    )
  )
}

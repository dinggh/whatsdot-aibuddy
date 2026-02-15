import React from 'react'
import { View, Text, Image, ScrollView } from '@tarojs/components'

import { BottomTabBar, StatusBar } from '@/components/layout'
import '@/styles/common.scss'
import './index.scss'
import img1 from '@/assets/generated-1771139016204.png'
import img2 from '@/assets/generated-1771138856711.png'
import img3 from '@/assets/generated-1771138893602.png'

const h = React.createElement

function Item(props) {
  return h(View, { className: 'panel history-item' },
    h(Image, { className: 'history-thumb', mode: 'aspectFill', src: props.src }),
    h(View, { className: 'history-meta' },
      h(Text, { className: 'history-title' }, props.title),
      h(Text, { className: 'history-sub' }, props.sub)
    ),
    h(Text, { className: 'history-arrow' }, '›')
  )
}

export default function HistoryPage() {
  return h(View, { className: 'screen' },
    h(StatusBar),
    h(View, { className: 'history-header' },
      h(Text, { className: 'history-h1' }, '历史记录'),
      h(Text, { className: 'history-search' }, '⌕')
    ),
    h(ScrollView, { className: 'history-list', scrollY: true },
      h(Text, { className: 'history-date' }, '今天'),
      h(Item, { src: img1, title: '24 x 15 = ?', sub: '三年级 · 今天19:50' }),
      h(Item, { src: img2, title: '阅读理解：小蝌蚪找妈妈', sub: '四年级 · 今天18:15' }),
      h(Text, { className: 'history-date' }, '昨天'),
      h(Item, { src: img3, title: '长方形面积计算', sub: '三年级 · 昨天20:10' }),
      h(Item, { src: img1, title: '比喻句仿写', sub: '五年级 · 昨天19:45' })
    ),
    h(BottomTabBar, { active: 'history' })
  )
}

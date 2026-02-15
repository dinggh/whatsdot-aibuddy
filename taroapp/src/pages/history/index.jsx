import React from 'react'
import { View, Text, Image, ScrollView } from '@tarojs/components'
import Taro, { useDidShow } from '@tarojs/taro'

import { BottomTabBar, StatusBar } from '@/components/layout'
import { fetchHistory } from '@/services/api'
import '@/styles/common.scss'
import './index.scss'
import img1 from '@/assets/generated-1771139016204.png'
import img2 from '@/assets/generated-1771138856711.png'
import img3 from '@/assets/generated-1771138893602.png'

const h = React.createElement

const fallbackImages = [img1, img2, img3]

function formatTime(input) {
  const d = input ? new Date(input) : new Date()
  const yyyy = d.getFullYear()
  const mm = String(d.getMonth() + 1).padStart(2, '0')
  const dd = String(d.getDate()).padStart(2, '0')
  const hh = String(d.getHours()).padStart(2, '0')
  const mi = String(d.getMinutes()).padStart(2, '0')
  return `${yyyy}-${mm}-${dd} ${hh}:${mi}`
}

function normalizeItems(items) {
  return (items || []).map((it, idx) => ({
    id: it.id || idx,
    title: it.title || '题目',
    sub: `${it.grade || '-'} · ${formatTime(it.solvedAt)}`,
    src: fallbackImages[idx % fallbackImages.length]
  }))
}

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
  const [items, setItems] = React.useState([])
  const [loading, setLoading] = React.useState(false)

  useDidShow(() => {
    setLoading(true)
    fetchHistory()
      .then((list) => setItems(normalizeItems(list)))
      .catch((err) => {
        Taro.showToast({ title: err.message || '加载失败', icon: 'none' })
      })
      .finally(() => setLoading(false))
  })

  return h(View, { className: 'screen' },
    h(StatusBar),
    h(View, { className: 'history-header' },
      h(Text, { className: 'history-h1' }, '历史记录'),
      h(Text, { className: 'history-search' }, '⌕')
    ),
    h(ScrollView, { className: 'history-list', scrollY: true },
      h(View, { className: 'history-list-inner' },
        loading ? h(Text, { className: 'history-date' }, '加载中...') : null,
        items.map((item) => h(Item, { key: item.id, src: item.src, title: item.title, sub: item.sub })),
        !loading && items.length === 0 ? h(Text, { className: 'history-date' }, '暂无记录') : null
      )
    ),
    h(BottomTabBar, { active: 'history' })
  )
}

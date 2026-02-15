import React from 'react'
import { View, Text, Image, ScrollView } from '@tarojs/components'
import Taro, { useDidShow } from '@tarojs/taro'

import { BottomTabBar, StatusBar } from '@/components/layout'
import { fetchHistory, formatTime, modeLabel, buildAssetURL } from '@/services/api'
import '@/styles/common.scss'
import './index.scss'

const h = React.createElement

function dayLabel(input) {
  const d = input ? new Date(input) : new Date()
  const now = new Date()
  const a = new Date(d.getFullYear(), d.getMonth(), d.getDate()).getTime()
  const b = new Date(now.getFullYear(), now.getMonth(), now.getDate()).getTime()
  const diff = Math.round((b - a) / 86400000)
  if (diff === 0) return '今天'
  if (diff === 1) return '昨天'
  return `${d.getMonth() + 1}月${d.getDate()}日`
}

function Item(props) {
  return h(View, {
    className: 'panel history-item',
    onClick: () => Taro.navigateTo({ url: `/pages/result/index?id=${props.id}` })
  },
  h(Image, { className: 'history-thumb', mode: 'aspectFill', src: props.src }),
  h(View, { className: 'history-meta' },
    h(Text, { className: 'history-title' }, props.title),
    h(View, { className: 'history-sub-row' },
      h(Text, { className: 'history-tag' }, modeLabel(props.mode)),
      h(Text, { className: 'history-sub' }, props.time)
    )
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
      .then((list) => {
        const normalized = (list || []).map((it) => ({
          id: it.id,
          mode: it.mode,
          title: it.title || '题目',
          src: buildAssetURL(it.thumbUrl || ''),
          solvedAt: it.solvedAt,
          day: dayLabel(it.solvedAt),
          time: formatTime(it.solvedAt).slice(11)
        }))
        setItems(normalized)
      })
      .catch((err) => {
        Taro.showToast({ title: err.message || '加载失败', icon: 'none' })
      })
      .finally(() => setLoading(false))
  })

  let lastDay = ''
  const rows = []
  items.forEach((it) => {
    if (it.day !== lastDay) {
      lastDay = it.day
      rows.push(h(Text, { key: `day-${it.day}-${it.id}`, className: 'history-day' }, it.day))
    }
    rows.push(h(Item, { key: it.id, ...it }))
  })

  return h(View, { className: 'screen' },
    h(StatusBar),
    h(View, { className: 'history-header' },
      h(Text, { className: 'history-h1' }, '历史记录'),
      h(View, { className: 'history-search' }, '⌕')
    ),
    h(ScrollView, { className: 'history-list', scrollY: true },
      h(View, { className: 'history-list-inner' },
        loading ? h(Text, { className: 'history-day' }, '加载中...') : null,
        ...rows,
        !loading && items.length === 0 ? h(Text, { className: 'history-day' }, '暂无记录') : null
      )
    ),
    h(BottomTabBar, { active: 'history' })
  )
}

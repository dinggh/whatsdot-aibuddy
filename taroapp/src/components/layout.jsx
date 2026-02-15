import React from 'react'
import { View, Text } from '@tarojs/components'
import Taro from '@tarojs/taro'

const h = React.createElement

export function StatusBar() {
  return h(View, { className: 'status-bar' }, h(Text, null, '9:41'), h(Text, null, '⋯'))
}

export function BottomTabBar(props) {
  const { active } = props
  const tab = (key, label, url) => h(View, {
    className: `tab ${active === key ? 'active' : ''}`,
    onClick: () => url && Taro.redirectTo({ url })
  }, label)

  return h(View, { className: 'tabbar' },
    tab('home', '首页', '/pages/home/index'),
    tab('history', '历史记录', '/pages/history/index'),
    tab('learn', '家长课堂', ''),
    tab('profile', '我的', '/pages/profile/index')
  )
}

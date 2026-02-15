import React from 'react'
import { View, Text } from '@tarojs/components'

import { BottomTabBar, StatusBar } from '@/components/layout'
import '@/styles/common.scss'
import './index.scss'

const h = React.createElement

function MenuItem(props) {
  return h(View, { className: 'profile-menu-item' }, h(Text, null, props.text), h(Text, null, '›'))
}

export default function ProfilePage() {
  return h(View, { className: 'screen' },
    h(StatusBar),
    h(View, { className: 'profile-content' },
      h(View, { className: 'profile-user' },
        h(View, { className: 'profile-avatar' }, '👤'),
        h(View, null,
          h(Text, { className: 'profile-name' }, '张妈妈'),
          h(Text, { className: 'profile-child' }, '小明 · 三年级')
        )
      ),
      h(View, { className: 'profile-stats' },
        h(View, { className: 'panel profile-stat' }, h(Text, { className: 'profile-num green' }, '47'), h(Text, { className: 'profile-label' }, '已用题量')),
        h(View, { className: 'panel profile-stat' }, h(Text, { className: 'profile-num orange' }, '53'), h(Text, { className: 'profile-label' }, '剩余次数'))
      ),
      h(View, { className: 'profile-buy' },
        h(Text, { className: 'profile-buy-title' }, '购买题包'),
        h(Text, { className: 'profile-buy-sub' }, '100次 / 月 · 限时优惠')
      ),
      h(View, { className: 'panel profile-menu' },
        h(MenuItem, { text: '家长成长指南' }),
        h(MenuItem, { text: '辅导设置' }),
        h(MenuItem, { text: '意见反馈' }),
        h(MenuItem, { text: '关于我们' })
      )
    ),
    h(BottomTabBar, { active: 'profile' })
  )
}

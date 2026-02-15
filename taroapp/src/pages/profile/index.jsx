import React from 'react'
import { View, Text, Button, Image } from '@tarojs/components'
import Taro, { useDidShow } from '@tarojs/taro'

import { BottomTabBar, StatusBar } from '@/components/layout'
import { bindPhoneByCode, fetchMe, updateProfile } from '@/services/api'
import '@/styles/common.scss'
import './index.scss'

const h = React.createElement

function MenuItem(props) {
  return h(View, { className: 'profile-menu-item' }, h(Text, null, props.text), h(Text, null, '›'))
}

function normalizeUser(u) {
  const user = u || {}
  return {
    nickName: user.nickName || '未设置昵称',
    avatarUrl: user.avatarUrl || '',
    phoneNumber: user.phoneNumber || '',
    usedCount: Number(user.usedCount || 0),
    remainingCount: Number(user.remainingCount || 0)
  }
}

export default function ProfilePage() {
  const [loading, setLoading] = React.useState(false)
  const [user, setUser] = React.useState(normalizeUser(null))
  const isWeApp = Taro.getEnv() === Taro.ENV_TYPE.WEAPP

  const loadUser = React.useCallback(() => {
    setLoading(true)
    fetchMe()
      .then((u) => setUser(normalizeUser(u)))
      .catch((err) => Taro.showToast({ title: err.message || '加载失败', icon: 'none' }))
      .finally(() => setLoading(false))
  }, [])

  useDidShow(() => {
    loadUser()
  })

  const onSyncProfile = () => {
    if (!isWeApp) {
      Taro.showToast({ title: '请在微信小程序中使用', icon: 'none' })
      return
    }

    Taro.getUserProfile({
      desc: '用于展示头像和昵称',
      success: async (res) => {
        try {
          const info = (res && res.userInfo) || {}
          const updated = await updateProfile(info.nickName || '微信用户', info.avatarUrl || '')
          setUser(normalizeUser(updated))
          Taro.showToast({ title: '昵称已同步', icon: 'success' })
        } catch (err) {
          Taro.showToast({ title: err.message || '同步失败', icon: 'none' })
        }
      },
      fail: () => {
        Taro.showToast({ title: '你取消了授权', icon: 'none' })
      }
    })
  }

  const onGetPhoneNumber = async (e) => {
    const code = e && e.detail && e.detail.code
    if (!code) {
      Taro.showToast({ title: '未获取到手机号授权码', icon: 'none' })
      return
    }

    try {
      const updated = await bindPhoneByCode(code)
      setUser(normalizeUser(updated))
      Taro.showToast({ title: '手机号已绑定', icon: 'success' })
    } catch (err) {
      Taro.showToast({ title: err.message || '绑定失败', icon: 'none' })
    }
  }

  return h(View, { className: 'screen' },
    h(StatusBar),
    h(View, { className: 'profile-content' },
      h(View, { className: 'profile-user' },
        user.avatarUrl
          ? h(Image, { className: 'profile-avatar profile-avatar-image', mode: 'aspectFill', src: user.avatarUrl })
          : h(View, { className: 'profile-avatar' }, '👤'),
        h(View, null,
          h(Text, { className: 'profile-name' }, user.nickName),
          h(Text, { className: 'profile-child' }, user.phoneNumber ? `手机号：${user.phoneNumber}` : '未绑定手机号')
        )
      ),
      h(View, { className: 'profile-actions' },
        h(Button, { className: 'panel profile-auth-btn', onClick: onSyncProfile, loading }, '同步昵称头像'),
        isWeApp
          ? h(Button, { className: 'panel profile-auth-btn', openType: 'getPhoneNumber', onGetPhoneNumber, loading }, '绑定手机号')
          : h(Button, { className: 'panel profile-auth-btn', onClick: () => Taro.showToast({ title: '请在微信小程序中使用', icon: 'none' }) }, '绑定手机号')
      ),
      h(View, { className: 'profile-stats' },
        h(View, { className: 'panel profile-stat' }, h(Text, { className: 'profile-num green' }, String(user.usedCount)), h(Text, { className: 'profile-label' }, '已用题量')),
        h(View, { className: 'panel profile-stat' }, h(Text, { className: 'profile-num orange' }, String(user.remainingCount)), h(Text, { className: 'profile-label' }, '剩余次数'))
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

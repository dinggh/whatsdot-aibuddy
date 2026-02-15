import React from 'react'
import { View, Text } from '@tarojs/components'
import Taro, { useLoad } from '@tarojs/taro'

import { getModeOptions, getCurrentMode, setCurrentMode } from '@/services/api'
import './index.scss'

const h = React.createElement
const PENDING_MODE_KEY = 'pending_mode_change'
const modeIcon = { guided: '💡', detailed: '📖', noanswer: '🙈', quick: '⚡' }

export default function ModePage() {
  const [mode, setMode] = React.useState(getCurrentMode())
  const [recordId, setRecordId] = React.useState(0)

  useLoad((query) => {
    if (query && query.mode) setMode(query.mode)
    if (query && query.recordId) setRecordId(Number(query.recordId))
  })

  const chooseMode = (nextMode) => {
    setCurrentMode(nextMode)
    if (recordId) {
      Taro.setStorageSync(PENDING_MODE_KEY, { id: recordId, mode: nextMode, ts: Date.now() })
    }
    Taro.navigateBack()
  }

  return h(View, { className: 'mode-mask' },
    h(View, { className: 'mode-sheet' },
      h(View, { className: 'mode-handle' }),
      h(View, { className: 'mode-header' },
        h(Text, { className: 'mode-title' }, '选择讲解模式'),
        h(Text, { className: 'mode-close', onClick: () => Taro.navigateBack() }, '×')
      ),
      ...getModeOptions().map((it) => h(View, {
        key: it.key,
        className: `mode-item ${mode === it.key ? 'active' : ''}`,
        onClick: () => chooseMode(it.key)
      },
      h(View, { className: 'mode-left' },
        h(View, { className: 'mode-icon' }, modeIcon[it.key] || '•'),
        h(View, null,
          h(Text, { className: 'mode-item-title' }, it.label),
          h(Text, { className: 'mode-item-desc' }, it.desc)
        )
      ),
      mode === it.key ? h(Text, { className: 'mode-check' }, '✓') : h(Text, { className: 'mode-check mode-check-empty' }, '○')
      ))
    )
  )
}

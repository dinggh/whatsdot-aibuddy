import React from 'react'
import { View, Text } from '@tarojs/components'
import Taro from '@tarojs/taro'

import './index.scss'

const h = React.createElement

const options = [
  { title: '引导思考', desc: '引导孩子一步一步思考（默认推荐）', active: true },
  { title: '详细讲解', desc: '完整讲解解题过程和知识点' },
  { title: '不给答案模式', desc: '只给思路和提示，不出现答案' },
  { title: '快速提示', desc: '快速给出关键提示，节省时间' }
]

export default function ModePage() {
  return h(View, { className: 'mode-mask' },
    h(View, { className: 'mode-sheet' },
      h(View, { className: 'mode-handle' }),
      h(View, { className: 'mode-header' },
        h(Text, { className: 'mode-title' }, '选择讲解模式'),
        h(Text, { className: 'mode-close', onClick: () => Taro.navigateBack() }, '×')
      ),
      ...options.map((it) => h(View, { key: it.title, className: `mode-item ${it.active ? 'active' : ''}` },
        h(View, null,
          h(Text, { className: 'mode-item-title' }, it.title),
          h(Text, { className: 'mode-item-desc' }, it.desc)
        ),
        it.active ? h(Text, { className: 'mode-check' }, '✓') : null
      ))
    )
  )
}

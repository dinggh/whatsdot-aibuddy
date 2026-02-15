import React from 'react'
import { View, Text, Image, ScrollView, Button } from '@tarojs/components'
import Taro, { useDidShow, useLoad } from '@tarojs/taro'

import { StatusBar } from '@/components/layout'
import {
  fetchHistoryDetail,
  regenerateHomework,
  getCurrentMode,
  modeLabel,
  buildAssetURL
} from '@/services/api'
import '@/styles/common.scss'
import './index.scss'

const h = React.createElement
const PENDING_MODE_KEY = 'pending_mode_change'

function Card(title, badge, body, strong = false) {
  return h(View, { className: `panel result-card ${strong ? 'result-card-strong' : ''}` },
    h(View, { className: 'result-card-head' },
      h(Text, { className: 'result-card-head-title' }, title),
      badge ? h(Text, { className: `result-badge ${strong ? 'result-badge-strong' : ''}` }, badge) : null
    ),
    h(Text, { className: `result-card-text ${strong ? 'result-card-text-strong' : ''}` }, body)
  )
}

function listText(items, fallback) {
  if (!Array.isArray(items) || items.length === 0) return fallback
  return items.filter(Boolean).map((it, idx) => `${idx + 1}. ${it}`).join('\n')
}

export default function ResultPage() {
  const [recordId, setRecordId] = React.useState(0)
  const [record, setRecord] = React.useState(null)
  const [loading, setLoading] = React.useState(false)
  const [mode, setMode] = React.useState(getCurrentMode())

  const refresh = (id) => {
    setLoading(true)
    return fetchHistoryDetail(id)
      .then((rec) => {
        setRecord(rec)
        setMode(rec && rec.mode ? rec.mode : getCurrentMode())
      })
      .catch((err) => {
        Taro.showToast({ title: err.message || '加载失败', icon: 'none' })
      })
      .finally(() => setLoading(false))
  }

  useLoad((query) => {
    const id = Number((query && query.id) || 0)
    setRecordId(id)
    if (!id) {
      Taro.showToast({ title: '记录ID缺失', icon: 'none' })
      return
    }
    refresh(id)
  })

  useDidShow(() => {
    if (!recordId) return
    const pending = Taro.getStorageSync(PENDING_MODE_KEY)
    if (pending && Number(pending.id) === Number(recordId) && pending.mode) {
      Taro.removeStorageSync(PENDING_MODE_KEY)
      setLoading(true)
      regenerateHomework(recordId, pending.mode)
        .then((nextRec) => {
          setRecord(nextRec)
          setMode(nextRec.mode)
        })
        .catch((err) => Taro.showToast({ title: err.message || '再生成失败', icon: 'none' }))
        .finally(() => setLoading(false))
    }
  })

  const result = (record && record.result) || {}
  const parentGuidance = listText(result.parent_guidance || result.parentGuidance, '1. 先让孩子说思路\n2. 用追问带孩子拆步骤\n3. 鼓励孩子自己验证答案')
  const stuck = listText(result.child_stuck_points || result.childStuckPoints, '1. 可能漏步骤\n2. 计算过程易出错')
  const knowledge = listText(result.knowledge_points || result.knowledgePoints, '1. 基础运算')
  const thinkingText = `${result.solution_thoughts || result.solutionThoughts || '暂无内容'}\n\n孩子可能卡点：\n${stuck}`
  const childText = `${result.explain_to_child || result.explainToChild || '暂无内容'}\n\n知识点：\n${knowledge}`

  return h(View, { className: 'screen' },
    h(StatusBar),
    h(View, { className: 'result-topbar' },
      h(Text, { className: 'result-back', onClick: () => Taro.navigateBack() }, '←'),
      h(Text, { className: 'result-title' }, '讲解结果'),
      h(View, {
        className: 'result-mode',
        onClick: () => Taro.navigateTo({ url: `/pages/mode/index?recordId=${recordId}&mode=${mode}` })
      }, modeLabel(mode))
    ),
    h(ScrollView, { className: 'result-content', scrollY: true },
      h(View, { className: 'result-content-inner' },
        h(View, { className: 'panel result-question' },
          h(View, { className: 'result-question-head' },
            h(Text, { className: 'result-question-title' }, '题目'),
            h(Text, { className: 'result-grade' }, record ? (record.suggestedGrade || '待识别') : '-')
          ),
          record
            ? h(Image, {
              className: 'result-qimg',
              mode: 'aspectFill',
              src: buildAssetURL(record.sourceImageUrl || '')
            })
            : null,
          h(Text, { className: 'result-qtext' }, loading ? '识别中...' : (record && record.questionText) || '暂无题干')
        ),
        Card('解题思路', '给家长看', thinkingText, false),
        Card('讲给孩子听', '简单说', childText, false),
        Card('家长可以这样引导', '重点', parentGuidance, true)
      )
    ),
    h(View, { className: 'result-actions' },
      h(Button, {
        className: 'result-btn result-btn-primary',
        onClick: () => {
          if (!recordId) return
          setLoading(true)
          regenerateHomework(recordId, mode)
            .then((nextRec) => setRecord(nextRec))
            .catch((err) => Taro.showToast({ title: err.message || '再生成失败', icon: 'none' }))
            .finally(() => setLoading(false))
        }
      }, '保存'),
      h(Button, { className: 'result-btn' }, '分享'),
      h(Button, { className: 'result-btn', onClick: () => Taro.setClipboardData({ data: parentGuidance }) }, '复制')
    )
  )
}

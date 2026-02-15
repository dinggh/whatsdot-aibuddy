const api = require('../../utils/api')

Page({
  data: {
    items: [],
    loading: false
  },

  async onShow() {
    this.setData({ loading: true });
    try {
      const items = await api.fetchHistory();
      this.setData({ items: (items || []).map(normalizeItem) });
    } catch (err) {
      wx.showToast({ title: err.message || '加载失败', icon: 'none' });
    } finally {
      this.setData({ loading: false });
    }
  },

  goHome() {
    wx.redirectTo({ url: '/pages/home/home' });
  },
  goProfile() {
    wx.navigateTo({ url: '/pages/profile/profile' });
  }
});

function normalizeItem(item) {
  const it = item || {};
  const date = it.solvedAt ? new Date(it.solvedAt) : new Date();
  const yyyy = date.getFullYear();
  const mm = String(date.getMonth() + 1).padStart(2, '0');
  const dd = String(date.getDate()).padStart(2, '0');
  const hh = String(date.getHours()).padStart(2, '0');
  const min = String(date.getMinutes()).padStart(2, '0');

  return {
    id: it.id,
    title: it.title || '题目',
    grade: it.grade || '-',
    thumbUrl: it.thumbUrl || '/images/generated-1771139016204.png',
    displayTime: `${yyyy}-${mm}-${dd} ${hh}:${min}`
  };
}

const api = require('../../utils/api')

Page({
  async onShow() {
    try {
      await api.ensureLogin();
    } catch (err) {
      wx.showToast({ title: '登录失败，请检查后端', icon: 'none' });
    }
  },

  goLoading() {
    wx.navigateTo({ url: '/pages/loading/loading' });
  },
  goHistory() {
    wx.redirectTo({ url: '/pages/history/history' });
  },
  goProfile() {
    wx.redirectTo({ url: '/pages/profile/profile' });
  }
});

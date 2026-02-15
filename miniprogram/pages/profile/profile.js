const api = require('../../utils/api')

Page({
  data: {
    loading: false,
    user: {
      nickName: '未设置昵称',
      avatarUrl: '',
      phoneNumber: '',
      usedCount: 0,
      remainingCount: 0
    }
  },

  async onShow() {
    await this.loadUser();
  },

  async loadUser() {
    this.setData({ loading: true });
    try {
      const user = await api.fetchMe();
      this.setData({ user: normalizeUser(user) });
    } catch (err) {
      wx.showToast({ title: err.message || '加载失败', icon: 'none' });
    } finally {
      this.setData({ loading: false });
    }
  },

  onTapSyncProfile() {
    wx.getUserProfile({
      desc: '用于展示头像和昵称',
      success: async (res) => {
        const info = res.userInfo || {};
        try {
          const user = await api.updateProfile(info.nickName || '微信用户', info.avatarUrl || '');
          this.setData({ user: normalizeUser(user) });
          wx.showToast({ title: '昵称已同步', icon: 'success' });
        } catch (err) {
          wx.showToast({ title: err.message || '同步失败', icon: 'none' });
        }
      },
      fail: () => {
        wx.showToast({ title: '你取消了授权', icon: 'none' });
      }
    });
  },

  async onGetPhoneNumber(e) {
    const code = e.detail && e.detail.code;
    if (!code) {
      wx.showToast({ title: '未获取到手机号授权码', icon: 'none' });
      return;
    }
    try {
      const user = await api.bindPhoneByCode(code);
      this.setData({ user: normalizeUser(user) });
      wx.showToast({ title: '手机号已绑定', icon: 'success' });
    } catch (err) {
      wx.showToast({ title: err.message || '绑定失败', icon: 'none' });
    }
  },

  goHome() {
    wx.redirectTo({ url: '/pages/home/home' });
  },
  goHistory() {
    wx.redirectTo({ url: '/pages/history/history' });
  }
});

function normalizeUser(user) {
  const u = user || {};
  return {
    nickName: u.nickName || '未设置昵称',
    avatarUrl: u.avatarUrl || '',
    phoneNumber: u.phoneNumber || '',
    usedCount: Number(u.usedCount || 0),
    remainingCount: Number(u.remainingCount || 0)
  };
}

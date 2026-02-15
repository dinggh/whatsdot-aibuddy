Page({
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

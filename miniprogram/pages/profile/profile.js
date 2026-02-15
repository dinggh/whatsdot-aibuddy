Page({
  goHome() {
    wx.redirectTo({ url: '/pages/home/home' });
  },
  goHistory() {
    wx.redirectTo({ url: '/pages/history/history' });
  }
});

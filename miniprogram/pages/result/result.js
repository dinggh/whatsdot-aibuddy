Page({
  goBack() {
    wx.navigateBack({ delta: 1 });
  },
  goMode() {
    wx.navigateTo({ url: '/pages/mode/mode' });
  },
  goHistory() {
    wx.navigateTo({ url: '/pages/history/history' });
  }
});

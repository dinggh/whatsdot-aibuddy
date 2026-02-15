Page({
  goBack() {
    wx.navigateBack({ delta: 1 });
  },
  goResult() {
    wx.navigateTo({ url: '/pages/result/result' });
  }
});

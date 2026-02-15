Page({
  goHome() {
    wx.redirectTo({ url: '/pages/home/home' });
  },
  goProfile() {
    wx.navigateTo({ url: '/pages/profile/profile' });
  }
});

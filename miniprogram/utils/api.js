function getAppSafe() {
  try {
    return getApp();
  } catch (e) {
    return null;
  }
}

function getBaseUrl() {
  const app = getAppSafe();
  return (app && app.globalData && app.globalData.apiBaseUrl) || 'http://127.0.0.1:8080';
}

function getToken() {
  const app = getAppSafe();
  return (app && app.globalData && app.globalData.token) || wx.getStorageSync('token') || '';
}

function saveSession(token, user) {
  const app = getAppSafe();
  if (app && app.globalData) {
    app.globalData.token = token || '';
    app.globalData.user = user || null;
  }
  wx.setStorageSync('token', token || '');
  wx.setStorageSync('user', user || null);
}

function request(options) {
  const {
    path,
    method = 'GET',
    data = {},
    auth = false
  } = options;

  return new Promise((resolve, reject) => {
    const headers = { 'content-type': 'application/json' };
    if (auth) {
      const token = getToken();
      if (token) {
        headers.Authorization = `Bearer ${token}`;
      }
    }

    wx.request({
      url: `${getBaseUrl()}${path}`,
      method,
      data,
      header: headers,
      success(res) {
        if (res.statusCode >= 200 && res.statusCode < 300) {
          resolve(res.data || {});
          return;
        }
        const message = (res.data && res.data.error) || `HTTP ${res.statusCode}`;
        reject(new Error(message));
      },
      fail(err) {
        reject(err);
      }
    });
  });
}

function wxLoginCode() {
  return new Promise((resolve, reject) => {
    wx.login({
      success: (res) => {
        if (!res.code) {
          reject(new Error('wx.login 没有返回 code'));
          return;
        }
        resolve(res.code);
      },
      fail: reject
    });
  });
}

async function loginWithWechat() {
  const code = await wxLoginCode();
  const data = await request({
    path: '/api/v1/auth/wechat/login',
    method: 'POST',
    data: { code }
  });
  if (!data.token) {
    throw new Error('登录失败：后端未返回 token');
  }
  saveSession(data.token, data.user || null);
  return data;
}

async function ensureLogin() {
  const token = getToken();
  if (token) {
    return token;
  }
  const data = await loginWithWechat();
  return data.token;
}

async function fetchMe() {
  await ensureLogin();
  const data = await request({
    path: '/api/v1/me',
    method: 'GET',
    auth: true
  });
  if (data.user) {
    saveSession(getToken(), data.user);
  }
  return data.user;
}

async function updateProfile(nickName, avatarUrl) {
  await ensureLogin();
  const data = await request({
    path: '/api/v1/auth/wechat/profile',
    method: 'POST',
    auth: true,
    data: { nickName, avatarUrl }
  });
  if (data.user) {
    saveSession(getToken(), data.user);
  }
  return data.user;
}

async function bindPhoneByCode(code) {
  await ensureLogin();
  const data = await request({
    path: '/api/v1/auth/wechat/phone',
    method: 'POST',
    auth: true,
    data: { code }
  });
  if (data.user) {
    saveSession(getToken(), data.user);
  }
  return data.user;
}

async function fetchHistory() {
  await ensureLogin();
  const data = await request({
    path: '/api/v1/history',
    method: 'GET',
    auth: true
  });
  return data.items || [];
}

module.exports = {
  request,
  ensureLogin,
  loginWithWechat,
  fetchMe,
  updateProfile,
  bindPhoneByCode,
  fetchHistory
};

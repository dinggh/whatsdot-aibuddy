"use client";

import {
  ArrowLeft,
  BatteryFull,
  BookOpen,
  Brain,
  Camera,
  CheckCircle2,
  ChevronRight,
  CircleHelp,
  Clock3,
  Copy,
  EyeOff,
  FileText,
  Gift,
  GraduationCap,
  Home,
  Image as ImageIcon,
  Info,
  Lightbulb,
  Loader,
  MessageCircle,
  MessageSquare,
  RefreshCw,
  Search,
  Settings,
  Share2,
  Signal,
  Star,
  User,
  UserRound,
  Wifi,
  X,
  Zap
} from "lucide-react";
import { type ReactNode, useMemo, useState } from "react";

import { Card } from "@/components/ui/card";
import { cn } from "@/lib/utils";

type ScreenKey = "home" | "loading" | "result" | "mode" | "history" | "profile";

const screenOptions: Array<{ key: ScreenKey; label: string }> = [
  { key: "home", label: "首页" },
  { key: "loading", label: "识别中" },
  { key: "result", label: "讲解结果" },
  { key: "mode", label: "模式弹窗" },
  { key: "history", label: "历史记录" },
  { key: "profile", label: "我的" }
];

function StatusBar() {
  return (
    <div className="flex h-[54px] items-center justify-between px-6 pt-3.5">
      <span className="text-base font-semibold text-[#2D2A26]">9:41</span>
      <div className="flex items-center gap-1.5 text-[#2D2A26]">
        <Signal className="h-4 w-4" strokeWidth={2} />
        <Wifi className="h-4 w-4" strokeWidth={2} />
        <BatteryFull className="h-4 w-4" strokeWidth={2} />
      </div>
    </div>
  );
}

function MiniShell({ children }: { children: ReactNode }) {
  return (
    <div className="mx-auto flex h-[844px] w-[390px] flex-col overflow-hidden rounded-[32px] border border-[#E8E4DE] bg-[#F8F5F0] shadow-[0_24px_60px_rgba(39,38,36,0.14)]">
      {children}
    </div>
  );
}

function TabBar({ active }: { active: "home" | "history" | "profile" | "learn" }) {
  const iconCls = "h-[22px] w-[22px]";
  const tab = (name: string) => (active === name ? "text-[#4A9B6E]" : "text-[#9C9892]");

  return (
    <div className="flex h-[82px] items-start justify-around border-t border-[#E8E4DE] bg-white pb-7 pt-2.5">
      <div className="flex flex-1 flex-col items-center gap-1">
        <Home className={cn(iconCls, tab("home"))} />
        <span className={cn("text-[10px]", active === "home" ? "font-semibold text-[#4A9B6E]" : "font-medium text-[#9C9892]")}>首页</span>
      </div>
      <div className="flex flex-1 flex-col items-center gap-1">
        <Clock3 className={cn(iconCls, tab("history"))} />
        <span className={cn("text-[10px]", active === "history" ? "font-semibold text-[#4A9B6E]" : "font-medium text-[#9C9892]")}>历史记录</span>
      </div>
      <div className="flex flex-1 flex-col items-center gap-1">
        <BookOpen className={cn(iconCls, tab("learn"))} />
        <span className={cn("text-[10px]", active === "learn" ? "font-semibold text-[#4A9B6E]" : "font-medium text-[#9C9892]")}>家长课堂</span>
      </div>
      <div className="flex flex-1 flex-col items-center gap-1">
        <UserRound className={cn(iconCls, tab("profile"))} />
        <span className={cn("text-[10px]", active === "profile" ? "font-semibold text-[#4A9B6E]" : "font-medium text-[#9C9892]")}>我的</span>
      </div>
    </div>
  );
}

function HomeScreen() {
  return (
    <MiniShell>
      <StatusBar />
      <div className="flex flex-1 flex-col gap-6 px-6">
        <header className="space-y-1 pt-4">
          <h1 className="text-[26px] font-semibold leading-none text-[#2D2A26]">微点辅导助手</h1>
          <p className="text-sm text-[#6D6A65]">让辅导作业变得简单轻松</p>
        </header>

        <section className="flex flex-col items-center gap-5 py-5">
          <button
            className="flex h-[200px] w-[200px] flex-col items-center justify-center gap-3 rounded-full bg-[#4A9B6E] text-white shadow-[0_4px_20px_rgba(74,155,110,0.25)] transition hover:bg-[#3A7D58]"
            type="button"
          >
            <Camera className="h-12 w-12" strokeWidth={2.25} />
            <span className="text-xl font-semibold">拍作业</span>
          </button>

          <Card className="flex h-[52px] w-full items-center justify-center gap-2 rounded-2xl border-[#E8E4DE] bg-white text-[#6D6A65] shadow-[0_2px_8px_rgba(26,25,24,0.03)]">
            <ImageIcon className="h-5 w-5" />
            <span className="text-[15px] font-medium">从相册上传</span>
          </Card>
        </section>

        <section className="flex items-center gap-2.5 rounded-2xl bg-[#FFF8E1] px-5 py-4">
          <Lightbulb className="h-5 w-5 shrink-0 text-[#C4A43A]" />
          <p className="text-[13px] leading-[1.5] text-[#6D6A65]">帮助家长引导孩子思考，而不是直接给答案</p>
        </section>

        <div className="flex-1" />
      </div>
      <TabBar active="home" />
    </MiniShell>
  );
}

function LoadingScreen() {
  return (
    <MiniShell>
      <StatusBar />
      <div className="flex h-11 items-center px-4">
        <ArrowLeft className="h-6 w-6 text-[#2D2A26]" />
      </div>

      <div className="flex flex-1 flex-col items-center justify-center gap-8 px-10">
        <img alt="作业图片" className="h-40 w-[220px] rounded-2xl border border-[#E8E4DE] object-cover shadow-[0_2px_12px_rgba(26,25,24,0.03)]" src="/images/generated-1771138856711.png" />

        <div className="flex flex-col items-center gap-4">
          <div className="flex h-12 w-12 items-center justify-center rounded-full border-[3px] border-[#4A9B6E] bg-[#E8F5ED]">
            <Loader className="h-6 w-6 animate-spin text-[#4A9B6E]" />
          </div>
          <div className="space-y-1 text-center">
            <p className="text-[18px] font-semibold leading-tight text-[#2D2A26]">正在识别题目...</p>
            <p className="text-sm text-[#9C9892]">AI正在整理讲解方式...</p>
          </div>
        </div>

        <div className="flex w-full items-start gap-2.5 rounded-[14px] bg-[#E6F0F8] px-[18px] py-[14px]">
          <Info className="mt-0.5 h-[18px] w-[18px] shrink-0 text-[#6BA3C7]" />
          <p className="text-xs leading-[1.5] text-[#6D6A65]">小贴士：引导孩子自己思考比直接告诉答案更有效哦</p>
        </div>
      </div>
    </MiniShell>
  );
}

function ResultScreen() {
  return (
    <MiniShell>
      <StatusBar />

      <div className="flex h-11 items-center justify-between px-4">
        <ArrowLeft className="h-6 w-6 text-[#2D2A26]" />
        <span className="text-[17px] font-semibold text-[#2D2A26]">讲解结果</span>
        <div className="flex items-center gap-1 rounded-full bg-[#E8F5ED] px-2.5 py-1.5 text-[#4A9B6E]">
          <Settings className="h-3.5 w-3.5" />
          <span className="text-xs font-medium">引导思考</span>
        </div>
      </div>

      <div className="min-h-0 flex-1 space-y-4 overflow-auto px-4 pb-4">
        <Card className="space-y-3 rounded-2xl border-[#E8E4DE] bg-white p-4 shadow-[0_2px_8px_rgba(26,25,24,0.03)]">
          <div className="flex items-center gap-1.5">
            <FileText className="h-4 w-4 text-[#4A9B6E]" />
            <span className="text-[13px] font-semibold text-[#4A9B6E]">题目</span>
            <span className="rounded-full bg-[#E6F0F8] px-2 py-0.5 text-[11px] font-medium text-[#6BA3C7]">三年级</span>
          </div>
          <img alt="题目" className="h-[100px] w-full rounded-xl object-cover" src="/images/generated-1771138893602.png" />
          <p className="text-[15px] font-medium text-[#2D2A26]">24 x 15 = ?</p>
        </Card>

        <Card className="space-y-3 rounded-2xl border-[#E8E4DE] bg-white p-4 shadow-[0_2px_8px_rgba(26,25,24,0.03)]">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-1.5">
              <Brain className="h-[18px] w-[18px] text-[#4A9B6E]" />
              <span className="text-base font-semibold text-[#2D2A26]">解题思路</span>
              <span className="rounded-full bg-[#E8F5ED] px-2 py-0.5 text-[11px] font-medium text-[#4A9B6E]">给家长看</span>
            </div>
            <RefreshCw className="h-4 w-4 text-[#9C9892]" />
          </div>
          <p className="whitespace-pre-line text-sm leading-[1.6] text-[#6D6A65]">这道题考查的是两位数乘法。可以用竖式计算法：{"\n\n"}1. 先算 24 x 5 = 120{"\n"}2. 再算 24 x 10 = 240{"\n"}3. 最后把两个结果相加：120 + 240 = 360{"\n\n"}核心知识点：乘法分配律</p>
        </Card>

        <Card className="space-y-3 rounded-2xl border-[#E8E4DE] bg-white p-4 shadow-[0_2px_8px_rgba(26,25,24,0.03)]">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-1.5">
              <MessageCircle className="h-[18px] w-[18px] text-[#E8A87C]" />
              <span className="text-base font-semibold text-[#2D2A26]">讲给孩子听</span>
              <span className="rounded-full bg-[#FFF0E6] px-2 py-0.5 text-[11px] font-medium text-[#E8A87C]">语气简单</span>
            </div>
            <RefreshCw className="h-4 w-4 text-[#9C9892]" />
          </div>
          <p className="whitespace-pre-line text-sm leading-[1.6] text-[#6D6A65]">宝贝，我们来想想 24 x 15 怎么算：{"\n\n"}想象一下，你有 24 颗糖果，要分成 15 份。{"\n\n"}我们可以把 15 拆成 10 和 5：{"\n"}- 24 个 5 是多少呢？就是 120{"\n"}- 24 个 10 呢？就是 240{"\n"}- 加起来就是 360 啦！</p>
        </Card>

        <Card className="space-y-3 rounded-2xl border-[1.5px] border-[#4A9B6E] bg-[#E8F5ED] p-4 shadow-[0_4px_16px_rgba(74,155,110,0.12)]">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-1.5">
              <Star className="h-[18px] w-[18px] text-[#4A9B6E]" />
              <span className="text-base font-semibold text-[#3A7D58]">家长可以这样引导</span>
              <span className="rounded-full bg-[#4A9B6E] px-2 py-0.5 text-[11px] font-semibold text-white">重点</span>
            </div>
            <RefreshCw className="h-4 w-4 text-[#4A9B6E]" />
          </div>
          <p className="whitespace-pre-line text-sm leading-[1.6] text-[#2D2A26]">第一步：先问孩子{"\n"}「24 x 15，你觉得可以怎么拆？」{"\n\n"}第二步：如果孩子不知道，提示{"\n"}「15 可以拆成哪两个数？10 和 5 对不对？」{"\n\n"}第三步：引导计算{"\n"}「那 24 x 5 是多少？24 x 10 呢？」{"\n\n"}第四步：总结{"\n"}「把两个结果加起来，就是答案了！」</p>
        </Card>
      </div>

      <div className="flex h-[84px] items-center gap-2.5 border-t border-[#E8E4DE] bg-white px-4 pb-7 pt-3">
        <button
          className="flex h-11 flex-1 items-center justify-center gap-1.5 rounded-xl bg-[#4A9B6E] text-sm font-medium text-white transition hover:bg-[#3A7D58]"
          type="button"
        >
          <Star className="h-4 w-4" />
          保存
        </button>
        <button
          className="flex h-11 flex-1 items-center justify-center gap-1.5 rounded-xl border border-[#E8E4DE] bg-white text-sm font-medium text-[#6D6A65] transition hover:bg-[#F8F5F0]"
          type="button"
        >
          <Share2 className="h-4 w-4" />
          分享
        </button>
        <button
          className="flex h-11 flex-1 items-center justify-center gap-1.5 rounded-xl border border-[#E8E4DE] bg-white text-sm font-medium text-[#6D6A65] transition hover:bg-[#F8F5F0]"
          type="button"
        >
          <Copy className="h-4 w-4" />
          复制
        </button>
      </div>
    </MiniShell>
  );
}

function ModeSelectScreen() {
  const options = [
    { title: "引导思考", desc: "引导孩子一步一步思考（默认推荐）", icon: Lightbulb, bg: "bg-[#4A9B6E]", active: true, iconColor: "text-white" },
    { title: "详细讲解", desc: "完整讲解解题过程和知识点", icon: BookOpen, bg: "bg-[#FFF0E6]", iconColor: "text-[#E8A87C]" },
    { title: "不给答案模式", desc: "只给思路和提示，不出现答案", icon: EyeOff, bg: "bg-[#E6F0F8]", iconColor: "text-[#6BA3C7]" },
    { title: "快速提示", desc: "快速给出关键提示，节省时间", icon: Zap, bg: "bg-[#FFF8E1]", iconColor: "text-[#C4A43A]" }
  ] as const;

  return (
    <MiniShell>
      <div className="flex flex-1 flex-col justify-end bg-[#00000060]">
        <div className="space-y-5 rounded-t-[24px] bg-white px-5 pb-9 pt-6">
          <div className="flex justify-center">
            <div className="h-1 w-10 rounded-full bg-[#E8E4DE]" />
          </div>
          <div className="flex items-center justify-between">
            <h2 className="text-[20px] font-semibold text-[#2D2A26]">选择讲解模式</h2>
            <div className="flex h-8 w-8 items-center justify-center rounded-full bg-[#F8F5F0]">
              <X className="h-[18px] w-[18px] text-[#6D6A65]" />
            </div>
          </div>

          <div className="space-y-2.5">
            {options.map((item) => (
              <div
                className={cn(
                  "flex items-center gap-3.5 rounded-2xl border bg-white px-[18px] py-4",
                  item.active ? "border-2 border-[#4A9B6E] bg-[#E8F5ED]" : "border border-[#E8E4DE]"
                )}
                key={item.title}
              >
                <div className={cn("flex h-10 w-10 items-center justify-center rounded-xl", item.bg)}>
                  <item.icon className={cn("h-5 w-5", item.iconColor)} />
                </div>
                <div className="flex-1 space-y-0.5">
                  <p className={cn("text-[15px]", item.active ? "font-semibold" : "font-medium")}>{item.title}</p>
                  <p className="text-xs text-[#6D6A65]">{item.desc}</p>
                </div>
                {item.active ? <CheckCircle2 className="h-[22px] w-[22px] text-[#4A9B6E]" /> : null}
              </div>
            ))}
          </div>
        </div>
      </div>
    </MiniShell>
  );
}

function HistoryItem({
  title,
  tag,
  tagColor,
  time,
  image
}: {
  title: string;
  tag: string;
  tagColor: "green" | "blue" | "orange";
  time: string;
  image: string;
}) {
  const tagMap = {
    green: "bg-[#E8F5ED] text-[#4A9B6E]",
    blue: "bg-[#E6F0F8] text-[#6BA3C7]",
    orange: "bg-[#FFF0E6] text-[#E8A87C]"
  } as const;

  return (
    <Card className="flex items-center gap-3 rounded-2xl border-[#E8E4DE] bg-white px-3.5 py-3.5 shadow-[0_1px_6px_rgba(26,25,24,0.03)]">
      <img alt={title} className="h-14 w-14 rounded-[10px] object-cover" src={image} />
      <div className="min-w-0 flex-1 space-y-1">
        <p className="truncate text-[15px] font-medium text-[#2D2A26]">{title}</p>
        <div className="flex items-center gap-2">
          <span className={cn("rounded-full px-1.5 py-0.5 text-[11px] font-medium", tagMap[tagColor])}>{tag}</span>
          <span className="text-[11px] text-[#9C9892]">{time}</span>
        </div>
      </div>
      <ChevronRight className="h-[18px] w-[18px] text-[#9C9892]" />
    </Card>
  );
}

function HistoryScreen() {
  return (
    <MiniShell>
      <StatusBar />
      <div className="flex h-[52px] items-center justify-between px-5">
        <h2 className="text-[22px] font-semibold text-[#2D2A26]">历史记录</h2>
        <div className="flex h-9 w-9 items-center justify-center rounded-full border border-[#E8E4DE] bg-white">
          <Search className="h-[18px] w-[18px] text-[#6D6A65]" />
        </div>
      </div>

      <div className="min-h-0 flex-1 space-y-3 overflow-auto px-4 py-2">
        <p className="text-[13px] font-medium text-[#9C9892]">今天</p>
        <HistoryItem image="/images/generated-1771139016204.png" tag="三年级" tagColor="green" time="今天19:50" title="24 x 15 = ?" />
        <HistoryItem image="https://images.unsplash.com/photo-1660287082054-6a8842ee4695?auto=format&fit=crop&w=200&q=80" tag="四年级" tagColor="blue" time="今天18:15" title="阅读理解：小蝌蚪找妈妈" />

        <p className="pt-1 text-[13px] font-medium text-[#9C9892]">昨天</p>
        <HistoryItem image="https://images.unsplash.com/photo-1676911809779-5ce408f0cf26?auto=format&fit=crop&w=200&q=80" tag="三年级" tagColor="green" time="昨天20:10" title="长方形面积计算" />
        <HistoryItem image="https://images.unsplash.com/photo-1770240366367-9bdc759e4445?auto=format&fit=crop&w=200&q=80" tag="五年级" tagColor="orange" time="昨天19:45" title="比喻句仿写" />
      </div>

      <TabBar active="history" />
    </MiniShell>
  );
}

function ProfileMenuItem({ icon, title, colorClass }: { icon: ReactNode; title: string; colorClass: string }) {
  return (
    <div className="flex items-center gap-3 border-b border-[#F0EDE8] px-[18px] py-4 last:border-0">
      <div className={cn("flex h-8 w-8 items-center justify-center rounded-lg", colorClass)}>{icon}</div>
      <span className="flex-1 text-[15px] font-medium text-[#2D2A26]">{title}</span>
      <ChevronRight className="h-[18px] w-[18px] text-[#9C9892]" />
    </div>
  );
}

function ProfileScreen() {
  return (
    <MiniShell>
      <StatusBar />
      <div className="flex min-h-0 flex-1 flex-col gap-5 px-5">
        <div className="flex items-center gap-3.5 pt-4">
          <div className="flex h-14 w-14 items-center justify-center rounded-full bg-[#E8F5ED]">
            <User className="h-7 w-7 text-[#4A9B6E]" />
          </div>
          <div>
            <p className="text-[30px] font-semibold leading-none text-[#2D2A26]">张妈妈</p>
            <p className="mt-1 text-sm text-[#6D6A65]">小明 · 三年级</p>
          </div>
        </div>

        <div className="flex gap-3">
          <Card className="flex flex-1 flex-col items-center gap-1 rounded-2xl border-[#E8E4DE] bg-white py-3.5">
            <p className="text-[28px] font-bold leading-none text-[#4A9B6E]">47</p>
            <p className="text-xs text-[#6D6A65]">已用题量</p>
          </Card>
          <Card className="flex flex-1 flex-col items-center gap-1 rounded-2xl border-[#E8E4DE] bg-white py-3.5">
            <p className="text-[28px] font-bold leading-none text-[#E8A87C]">53</p>
            <p className="text-xs text-[#6D6A65]">剩余次数</p>
          </Card>
        </div>

        <div className="flex items-center gap-3 rounded-2xl bg-gradient-to-r from-[#4A9B6E] to-[#3A7D58] px-[18px] py-4 text-white">
          <div className="flex h-10 w-10 items-center justify-center rounded-xl bg-[#FFFFFF20]">
            <Gift className="h-[22px] w-[22px]" />
          </div>
          <div className="flex-1">
            <p className="text-base font-semibold">购买题包</p>
            <p className="text-xs text-white/70">100次 / 月 · 限时优惠</p>
          </div>
          <ChevronRight className="h-5 w-5" />
        </div>

        <Card className="overflow-hidden rounded-2xl border-[#E8E4DE] bg-white shadow-[0_1px_6px_rgba(26,25,24,0.03)]">
          <ProfileMenuItem colorClass="bg-[#FFF8E1]" icon={<GraduationCap className="h-[18px] w-[18px] text-[#C4A43A]" />} title="家长成长指南" />
          <ProfileMenuItem colorClass="bg-[#E6F0F8]" icon={<Settings className="h-[18px] w-[18px] text-[#6BA3C7]" />} title="辅导设置" />
          <ProfileMenuItem colorClass="bg-[#E8F5ED]" icon={<MessageSquare className="h-[18px] w-[18px] text-[#4A9B6E]" />} title="意见反馈" />
          <div className="flex items-center gap-3 px-[18px] py-4">
            <div className="flex h-8 w-8 items-center justify-center rounded-lg bg-[#FFF0E6]">
              <CircleHelp className="h-[18px] w-[18px] text-[#E8A87C]" />
            </div>
            <span className="flex-1 text-[15px] font-medium text-[#2D2A26]">关于我们</span>
            <ChevronRight className="h-[18px] w-[18px] text-[#9C9892]" />
          </div>
        </Card>

        <div className="flex-1" />
      </div>
      <TabBar active="profile" />
    </MiniShell>
  );
}

export function MpMiniProgram() {
  const [screen, setScreen] = useState<ScreenKey>("home");

  const content = useMemo(() => {
    if (screen === "home") return <HomeScreen />;
    if (screen === "loading") return <LoadingScreen />;
    if (screen === "result") return <ResultScreen />;
    if (screen === "mode") return <ModeSelectScreen />;
    if (screen === "history") return <HistoryScreen />;
    return <ProfileScreen />;
  }, [screen]);

  return (
    <div className="mx-auto min-h-screen bg-[radial-gradient(circle_at_top,_#ffffff_0%,_#f8f5f0_55%,_#efe9de_100%)] px-4 py-8">
      <div className="mx-auto mb-6 flex w-full max-w-[860px] flex-wrap justify-center gap-2">
        {screenOptions.map((item) => (
          <button
            className={cn(
              "rounded-full border px-4 py-2 text-sm font-medium transition",
              item.key === screen
                ? "border-[#4A9B6E] bg-[#E8F5ED] text-[#3A7D58]"
                : "border-[#E8E4DE] bg-white text-[#6D6A65] hover:bg-[#F8F5F0]"
            )}
            key={item.key}
            onClick={() => setScreen(item.key)}
            type="button"
          >
            {item.label}
          </button>
        ))}
      </div>
      {content}
    </div>
  );
}

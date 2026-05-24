# AI 智联论坛 · Android APK

WebView 壳工程，加载 `assets/dist/` 中的 Vue 构建产物。

## 方式一：Android Studio（推荐）

1. 同步前端资源（每次改完 Vue 后执行）：

   ```powershell
   cd ..\..
   npm run build
   .\scripts\sync-android-assets.ps1
   ```

2. Android Studio → **File → Open** → 选择本目录 `frontend-standalone/android`

3. 等待 Gradle Sync 完成

4. 菜单 **Build → Build Bundle(s) / APK(s) → Build APK(s)**  
   或运行配置选择 `app` 后点击 Run

5. Debug APK 输出路径：

   `app/build/outputs/apk/debug/app-debug.apk`

### 发布 Release APK

**Build → Generate Signed Bundle / APK → APK**，创建或选择 keystore 后构建 `release`。

命令行（需自行配置签名）：

```powershell
cd ..
.\scripts\build-apk.ps1 -Variant release
```

## 方式二：命令行一键打包

在项目根目录 `frontend-standalone`：

```powershell
.\scripts\build-apk.ps1
```

需要：

- 已安装 Android SDK（`local.properties` 中 `sdk.dir`）
- JDK 17+（可用 Android Studio 自带 JBR）

若 Gradle 拉取依赖失败，可在 Android Studio 中打开工程同步一次（会使用 IDE 代理/镜像）。

## 安装到手机

```powershell
adb install -r android\app\build\outputs\apk\debug\app-debug.apk
```

## 说明

- 使用 `WebViewAssetLoader` 加载 `https://appassets.androidplatform.net/assets/dist/index.html`
- 与前端 Hash 路由、`base: './'` 配套
- 返回键优先关闭导航抽屉

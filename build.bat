@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion
echo ========================================
echo           自动构建脚本
echo ========================================
set FRONTEND_DIR=web
set BACKEND_DIR=%CD%
set BACKEND_OUTPUT=app.exe
go version >nul 2>&1 || (echo 错误: Go未安装 && pause && exit /b 1)
bun --version >nul 2>&1 || (echo 错误: bun未安装 && pause && exit /b 1)
echo.
echo ========================================
echo           开始构建前端
echo ========================================
cd /d "%FRONTEND_DIR%"
if not exist "node_modules" (
    echo 正在安装依赖...
    bun install --frozen-lockfile
) else (
    echo 依赖已存在，跳过安装阶段...
)
echo 正在构建前端...
bun run build || (echo 错误: 前端构建失败！ && cd /d "%BACKEND_DIR%" && pause && exit /b 1)

echo ✅ 前端构建成功！

cd /d "%BACKEND_DIR%"

echo.
echo ========================================
echo           开始构建后端
echo ========================================

echo 正在整理依赖...
go mod tidy

echo 正在构建Go项目...
go build -tags=go_json -ldflags="-s -w" -o "%BACKEND_OUTPUT%" . || (echo 错误: 后端构建失败！ && pause && exit /b 1)

echo ✅ 后端构建成功！
echo.
echo ========================================
echo           构建完成
echo ========================================
pause

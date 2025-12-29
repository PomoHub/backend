<div align="center">

# <img src="public/logos/pomohub-logo-white.svg" alt="PomoHub Logo" width="300" height="80" />

**Focus. Build Habits. Achieve Goals.**

[![License](https://img.shields.io/github/license/PomoHub/PomoHub?style=flat-square)](LICENSE)
[![GitHub package.json version](https://img.shields.io/github/package-json/v/PomoHub/PomoHub?style=flat-square)](package.json)
[![GitHub Downloads](https://img.shields.io/github/downloads/PomoHub/PomoHub/total?style=flat-square)](https://github.com/PomoHub/PomoHub/releases)
![GitHub repo size](https://img.shields.io/github/repo-size/PomoHub/PomoHub?style=flat-square)

<br/>

![Tauri](https://img.shields.io/badge/Tauri-FFC131?style=for-the-badge&logo=tauri&logoColor=black)
![React](https://img.shields.io/badge/React-20232A?style=for-the-badge&logo=react&logoColor=61DAFB)
![TypeScript](https://img.shields.io/badge/TypeScript-007ACC?style=for-the-badge&logo=typescript&logoColor=white)
![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-316192?style=for-the-badge&logo=postgresql&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2CA5E0?style=for-the-badge&logo=docker&logoColor=white)

</div>

<br/>

**PomoHub** is a modern, all-in-one productivity application built with **Tauri**, **React**, and **TypeScript**. It combines the power of the Pomodoro technique with habit tracking, goal setting, and task management in a beautiful, customizable interface. With a robust backend built in **Go**, it supports real-time social features to help you stay accountable with friends.

---

## 🛠️ Tech Stack

### Frontend (Desktop & Mobile)
- **Framework:** ![React](https://img.shields.io/badge/-React-black?style=flat-square&logo=react) ![Vite](https://img.shields.io/badge/-Vite-black?style=flat-square&logo=vite)
- **Core:** ![Tauri](https://img.shields.io/badge/-Tauri-black?style=flat-square&logo=tauri) ![Rust](https://img.shields.io/badge/-Rust-black?style=flat-square&logo=rust)
- **Language:** ![TypeScript](https://img.shields.io/badge/-TypeScript-black?style=flat-square&logo=typescript)
- **Styling:** ![TailwindCSS](https://img.shields.io/badge/-TailwindCSS-black?style=flat-square&logo=tailwindcss)
- **State Management:** Zustand
- **Motion:** Framer Motion

### Backend (API & Microservices)
- **Language:** ![Go](https://img.shields.io/badge/-Go-black?style=flat-square&logo=go)
- **Framework:** Fiber
- **Database:** ![PostgreSQL](https://img.shields.io/badge/-PostgreSQL-black?style=flat-square&logo=postgresql)
- **Caching:** Valkey (Redis compatible)
- **Infrastructure:** ![Docker](https://img.shields.io/badge/-Docker-black?style=flat-square&logo=docker)

### Website
- **Framework:** ![Next.js](https://img.shields.io/badge/-Next.js-black?style=flat-square&logo=next.js)

---

## ✨ Features

### 🍅 Pomodoro Timer
- **Customizable Modes:** Switch between Work, Short Break, and Long Break.
- **Visual Progress:** Elegant circular progress indicator.
- **Session Logging:** Automatically logs work sessions to track your focus time.
- **Configurable Settings:** Adjust durations for work and breaks to suit your workflow.

### ✅ Habit Tracker
- **Daily Tracking:** Easily mark habits as completed for the day.
- **Color Coding:** Assign custom colors to different habits for better visualization.
- **Streak Logic:** Visual cues for completed habits (strikethrough and color fill).
- **Persistent Data:** All habits and logs are stored locally.

### 📅 Calendar
- **Monthly Overview:** View your activity across the entire month.
- **Daily Insights:** Click on any day to see detailed stats:
  - Habits completed
  - Total focus time (Pomodoro minutes)
  - Tasks finished
- **Visual Indicators:** Color-coded dots on calendar days show activity types.

### 📝 Todo List
- **Task Management:** Create, read, update, and delete tasks.
- **Due Dates:** Assign due dates to keep track of deadlines.
- **Smart Sorting:** Uncompleted tasks always appear at the top.

### 🎯 Goals
- **Long-term Tracking:** Set numeric goals (e.g., "Read 10 Books").
- **Progress Bars:** Visual progress tracking with percentage indicators.
- **Target Dates:** Set deadlines for your goals.

### ⚙️ Customization
- **Theme Support:** Switch between **Light**, **Dark**, or **System** themes.
- **Custom Backgrounds:** Choose any image from your computer as the application background.
- **Privacy First:** All data is stored locally on your device using SQLite.

---

## 🚀 Getting Started

### Prerequisites
- [Node.js](https://nodejs.org/) (v18+)
- [Rust](https://www.rust-lang.org/tools/install) (for Tauri)
- [Go](https://go.dev/) (v1.23+, for Backend)
- [Docker](https://www.docker.com/) (optional, for Backend services)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/PomoHub/PomoHub.git
   cd PomoHub
   ```

2. **Frontend (Desktop App)**
   ```bash
   # Install dependencies
   npm install

   # Run development mode
   npm run tauri dev
   ```

3. **Backend (Server)**
   ```bash
   cd backend

   # Run with Docker (Recommended)
   docker-compose up -d

   # OR Run manually
   go run cmd/server/main.go
   ```

---

## 🗺️ Roadmap

### v0.1.0 (Current)
- [x] Basic Application Infrastructure (Tauri + React)
- [x] SQLite Database Integration
- [x] Pomodoro Timer Feature
- [x] Habit Tracker Feature
- [x] Calendar & Stats Feature
- [x] Todo List Feature
- [x] Goals Feature
- [x] Settings (Theme & Background)

### v0.1.1 (Bug Fixes & Improvements)
- **Fix:** Resolved database connection issues preventing creation of Habits, Todos, and Goals.
- **Fix:** Fixed theme switching not applying correctly.
- **Fix:** Corrected Calendar stats aggregation.
- **New:** Added configurable auto-transition between Pomodoro modes.
- **New:** Added "Long Break Interval" setting to customize session count before long break.
- **New:** Added sound notification (placeholder) on timer completion.

### v0.1.2
- **New:** User Onboarding & Profile System (Name, Birthday).
- **New:** Birthday Celebration with confetti.
- **New:** Seasonal Winter Snowfall effect (Dec-Feb).
- **New:** Profile Page with comprehensive statistics (Total Time, Streaks).
- **New:** Achievements System with unlock notifications.
- **New:** Customizable Notification Sounds (Presets).
- **New:** Android Support (Tauri Mobile initialized).

### v0.1.3
- **New:** Daily Task Reminders (Notifications).
- **Improvement:** Reliable Background Timer (Timestamp-based logic).
- **Improvement:** Native Notifications for Desktop & Android.
- **Improvement:** Custom Notification Sounds (File selection).
- **Fix:** Android background execution issues.

### v0.1.4 (Current)
- **New:** Desktop Window Title Timer (Visible countdown in title bar).
- **New:** Note Taking System (Text, Drawings, Attachments).
- **New:** Smart Macros (Create Todos/Tasks directly from notes with `@create-todo` syntax).
- **Improvement:** Persistent Mobile Notifications (Lock screen timer updates).
- **Improvement:** Notification Sound Reliability (Fallback beep mechanism).
- **Fix:** Desktop Notification Spam (Optimized notification frequency).

### v0.2.0 (Planned)
- [x] **Social Spaces:** Create private rooms, invite friends via code, and focus together.
- [x] **Real-time Chat:** Chat with friends in your space.
- [x] **Friend System:** Add friends, block users, and see what they are up to.

---

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

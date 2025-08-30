# 🚀 Open-Same - Same.AI Replica Platform

[![CI/CD Pipeline](https://github.com/you112ef/Open-same/workflows/CI/CD%20Pipeline/badge.svg)](https://github.com/you112ef/Open-same/actions)
[![Auto Deploy](https://github.com/you112ef/Open-same/workflows/Deploy%20to%20Cloudflare%20Pages/badge.svg)](https://github.com/you112ef/Open-same/actions)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**Open-Same** هو منصة مفتوحة المصدر تحاكي Same.AI مع ميزات متقدمة للنشر التلقائي على Cloudflare Pages.

## ✨ الميزات الرئيسية

### 🤖 Same.AI Replica Features
- **AI-Powered Content Generation** - OpenAI GPT-4 & Anthropic Claude
- **Real-time Collaboration** - Live editing with WebSocket support
- **User Authentication** - Secure login and role-based access
- **Content Management** - Version control and content organization
- **Mobile-Responsive Design** - Optimized for all devices

### 🚀 Auto-Deployment Features
- **GitHub Actions CI/CD** - Automated testing and deployment
- **Cloudflare Pages Integration** - Instant global deployment
- **Multi-Environment Support** - Production & Staging
- **Release Management** - Automatic deployment on tags
- **Quick Deploy** - Manual deployment when needed

## 🏗️ البنية التقنية

```
Open-Same/
├── frontend/          # React + TypeScript application
├── backend/           # Go + PostgreSQL + Redis API
├── sdk/              # Professional SDK for integration
├── docs/             # Comprehensive documentation
├── scripts/          # Deployment and utility scripts
├── .github/          # GitHub Actions workflows
└── monitoring/       # Observability and monitoring
```

## 🚀 النشر التلقائي

### ✅ ما يعمل تلقائياً

1. **Push إلى `main`** → نشر تلقائي إلى Production
2. **Push إلى `develop`** → نشر تلقائي إلى Staging  
3. **إنشاء Tag جديد** → نشر تلقائي للإصدار
4. **Pull Request** → اختبار تلقائي

### 🔧 الإعداد السريع

```bash
# تشغيل سكريبت الإعداد
./scripts/setup-auto-deploy.sh

# أو اتبع الدليل اليدوي
docs/auto-deployment.md
```

### 📋 المتطلبات

- [Cloudflare API Token](https://dash.cloudflare.com/profile/api-tokens)
- [GitHub Repository Secrets](https://github.com/you112ef/Open-same/settings/secrets/actions)

## 🛠️ التثبيت والتشغيل

### المتطلبات الأساسية

- Node.js 18+
- Go 1.21+
- PostgreSQL 15+
- Redis 7+
- Docker & Docker Compose

### التثبيت السريع

```bash
# Clone repository
git clone https://github.com/you112ef/Open-same.git
cd Open-same

# Setup with Docker
docker-compose up -d

# Or setup manually
./scripts/setup.sh
```

### التشغيل المحلي

```bash
# Frontend
cd frontend
npm install
npm run dev

# Backend
cd backend
go mod download
go run cmd/server/main.go
```

## 📚 الوثائق

- **[النشر التلقائي](docs/auto-deployment.md)** - دليل شامل للنشر التلقائي
- **[بنية المشروع](PROJECT_STRUCTURE.md)** - تفاصيل البنية التقنية
- **[دليل المساهمة](CONTRIBUTING.md)** - كيفية المساهمة في المشروع
- **[النشر على Cloudflare](docs/deployment.md)** - دليل النشر المفصل

## 🔄 Workflows

### 1. النشر التلقائي (`deploy-cloudflare.yml`)
- بناء Frontend و Backend
- اختبار شامل
- نشر تلقائي إلى Staging/Production
- نشر Backend إلى Cloudflare Workers

### 2. النشر السريع (`quick-deploy.yml`)
- نشر فوري عند الطلب
- اختيار البيئة
- تقارير فورية

### 3. نشر الإصدارات (`release-deploy.yml`)
- نشر تلقائي عند إنشاء tags
- إدارة الإصدارات
- تقارير مفصلة

## 🌐 النشر

### Cloudflare Pages
- **Production:** `https://open-same.pages.dev`
- **Staging:** `https://open-same-staging.pages.dev`

### GitHub Actions
- [CI/CD Pipeline](https://github.com/you112ef/Open-same/actions/workflows/ci.yml)
- [Auto Deploy](https://github.com/you112ef/Open-same/actions/workflows/deploy-cloudflare.yml)
- [Quick Deploy](https://github.com/you112ef/Open-same/actions/workflows/quick-deploy.yml)

## 🤝 المساهمة

نرحب بمساهماتكم! راجع [دليل المساهمة](CONTRIBUTING.md) للبدء.

### كيفية المساهمة

1. Fork المشروع
2. أنشئ feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit التغييرات (`git commit -m 'Add AmazingFeature'`)
4. Push إلى branch (`git push origin feature/AmazingFeature`)
5. أنشئ Pull Request

## 📄 الترخيص

هذا المشروع مرخص تحت [MIT License](LICENSE).

## 🙏 الشكر

- [Same.AI](https://same.ai) - للمفهوم الأصلي
- [Cloudflare](https://cloudflare.com) - للنشر السريع
- [GitHub Actions](https://github.com/features/actions) - للنشر التلقائي

## 📞 الدعم

- **Issues:** [GitHub Issues](https://github.com/you112ef/Open-same/issues)
- **Discussions:** [GitHub Discussions](https://github.com/you112ef/Open-same/discussions)
- **Documentation:** [docs/](docs/)

---

## 🎉 النشر التلقائي يعمل الآن!

**Open-Same** جاهز للنشر التلقائي! كل push إلى `main` أو `develop` سيؤدي إلى نشر تلقائي على Cloudflare Pages.

### 🚀 الخطوات التالية:

1. **أضف Cloudflare API Token** إلى GitHub Secrets
2. **أنشئ Environments** (production & staging)
3. **اختبر النشر** باستخدام Quick Deploy
4. **استمتع بالنشر التلقائي!**

---

**⭐ إذا أعجبك المشروع، لا تنس إعطاءه نجمة!**

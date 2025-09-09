# 🚀 إعداد النشر التلقائي - Open-Same

## ✅ ما تم إعداده

- [x] التحقق من المتطلبات الأساسية
- [x] إنشاء branch develop
- [x] التحقق من ملفات Workflow
- [x] دليل إعداد Environment Variables

## 🔑 المطلوب منك

1. **إعداد GitHub Environments:**
   - production
   - staging

2. **إضافة GitHub Secret:**
   - CLOUDFLARE_API_TOKEN

3. **الحصول على Cloudflare API Token**

## 📁 ملفات Workflow

- `deploy-cloudflare.yml` - النشر التلقائي
- `quick-deploy.yml` - النشر السريع
- `release-deploy.yml` - نشر الإصدارات

## 🚀 كيفية الاستخدام

1. **النشر التلقائي:** push إلى main/develop
2. **النشر السريع:** من Actions tab
3. **نشر الإصدارات:** إنشاء tags جديدة

## 📚 الدليل الكامل

راجع `docs/auto-deployment.md` للحصول على دليل مفصل.

---
*تم إنشاء هذا الملف بواسطة setup-auto-deploy.sh*

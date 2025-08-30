# 🚀 النشر التلقائي - Open-Same

دليل شامل لإعداد النشر التلقائي لمشروع Open-Same على Cloudflare Pages باستخدام GitHub Actions.

## 📋 نظرة عامة

تم إعداد ثلاثة أنواع من workflows للنشر التلقائي:

1. **النشر التلقائي** (`deploy-cloudflare.yml`) - يعمل عند push إلى main/develop
2. **النشر السريع** (`quick-deploy.yml`) - يعمل يدوياً عند الطلب
3. **نشر الإصدارات** (`release-deploy.yml`) - يعمل عند إنشاء tags جديدة

## ⚙️ المتطلبات الأساسية

### 1. Cloudflare API Token

1. اذهب إلى [Cloudflare Dashboard](https://dash.cloudflare.com/)
2. انتقل إلى **My Profile** → **API Tokens**
3. انقر على **Create Token**
4. اختر **Custom token**
5. أضف الأذونات التالية:
   - **Zone:Zone:Read**
   - **Zone:Zone Settings:Edit**
   - **Account:Cloudflare Pages:Edit**
   - **Account:Cloudflare Workers:Edit**

### 2. إضافة Secrets إلى GitHub

1. اذهب إلى repository الخاص بك على GitHub
2. انتقل إلى **Settings** → **Secrets and variables** → **Actions**
3. أضف secret جديد:
   - **Name:** `CLOUDFLARE_API_TOKEN`
   - **Value:** API Token من Cloudflare

## 🔧 إعداد البيئات

### إعداد Environment: Production

1. في GitHub repository، اذهب إلى **Settings** → **Environments**
2. انقر على **New environment**
3. أدخل **Environment name:** `production`
4. أضف **Environment protection rules**:
   - **Required reviewers:** أضف نفسك أو فريقك
   - **Wait timer:** 0 (أو حسب الحاجة)

### إعداد Environment: Staging

1. أنشئ environment جديد باسم `staging`
2. لا تحتاج إلى protection rules للإعدادات الأولية

## 📁 ملفات Workflow

### 1. `deploy-cloudflare.yml`

**الميزات:**
- ✅ بناء Frontend و Backend
- ✅ تشغيل الاختبارات
- ✅ نشر تلقائي إلى Staging (develop branch)
- ✅ نشر تلقائي إلى Production (main branch)
- ✅ نشر Backend إلى Cloudflare Workers
- ✅ تقارير مفصلة عن النشر

**متى يعمل:**
- عند push إلى `main` أو `develop`
- عند إنشاء Pull Request إلى `main`

### 2. `quick-deploy.yml`

**الميزات:**
- ✅ نشر سريع عند الطلب
- ✅ اختيار البيئة (Production/Staging)
- ✅ خيار إعادة البناء القسري
- ✅ تقارير فورية

**متى يعمل:**
- يدوياً من GitHub Actions tab
- مفيد للنشر السريع أو إصلاح المشاكل

### 3. `release-deploy.yml`

**الميزات:**
- ✅ نشر تلقائي عند إنشاء tags
- ✅ نشر الإصدارات الجديدة
- ✅ تقارير مفصلة عن الإصدار

**متى يعمل:**
- عند إنشاء tag جديد (مثل `v1.2.0`)
- عند نشر release على GitHub

## 🚀 كيفية الاستخدام

### النشر التلقائي

1. **Push إلى main branch:**
   ```bash
   git add .
   git commit -m "Update features"
   git push origin main
   ```
   - سيتم البناء والنشر تلقائياً
   - انتظر 2-3 دقائق لإكمال العملية

2. **Push إلى develop branch:**
   ```bash
   git checkout develop
   git add .
   git commit -m "New feature"
   git push origin develop
   ```
   - سيتم النشر إلى بيئة Staging

### النشر السريع

1. اذهب إلى **Actions** tab في GitHub
2. اختر **Quick Deploy to Cloudflare Pages**
3. انقر على **Run workflow**
4. اختر البيئة المطلوبة
5. انقر على **Run workflow**

### نشر الإصدارات

1. **إنشاء tag جديد:**
   ```bash
   git tag -a v1.3.0 -m "New features and improvements"
   git push origin v1.3.0
   ```

2. **أو إنشاء Release على GitHub:**
   - اذهب إلى **Releases**
   - انقر على **Create a new release**
   - اختر tag موجود أو أنشئ واحد جديد
   - اكتب وصف الإصدار
   - انشر Release

## 📊 مراقبة النشر

### 1. GitHub Actions

- اذهب إلى **Actions** tab
- راقب حالة workflows
- راجع logs للأخطاء

### 2. Cloudflare Dashboard

- اذهب إلى **Pages** في Cloudflare
- راقب deployments
- تحقق من URLs

### 3. تقارير النشر

كل workflow ينشئ تقرير مفصل يتضمن:
- ✅ ما تم نشره
- 🔗 الروابط المهمة
- 📋 الخطوات التالية
- ⏰ وقت النشر

## 🔍 استكشاف الأخطاء

### مشاكل شائعة

1. **خطأ في API Token:**
   - تأكد من صحة `CLOUDFLARE_API_TOKEN`
   - تحقق من صلاحيات Token

2. **فشل في البناء:**
   - راجع logs في GitHub Actions
   - تحقق من dependencies

3. **فشل في النشر:**
   - تحقق من إعدادات Cloudflare Pages
   - راجع project name في wrangler.toml

### نصائح للتشخيص

1. **راجع logs بالتفصيل:**
   - انقر على job فاشل في Actions
   - راجع كل step بالتفصيل

2. **اختبر محلياً:**
   ```bash
   cd frontend
   npm run build
   wrangler pages deploy dist
   ```

3. **تحقق من Secrets:**
   - تأكد من وجود `CLOUDFLARE_API_TOKEN`
   - تحقق من صحة القيمة

## 🎯 أفضل الممارسات

### 1. إدارة Branches

- **main:** للإنتاج (deploy تلقائي)
- **develop:** للتطوير (deploy إلى staging)
- **feature branches:** للتطوير المحلي

### 2. إدارة الإصدارات

- استخدم semantic versioning (v1.2.3)
- أنشئ tags لكل إصدار مهم
- اكتب release notes مفصلة

### 3. المراقبة

- راقب deployments بانتظام
- تحقق من الأداء بعد كل نشر
- احتفظ بسجل التغييرات

## 🔗 روابط مفيدة

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [Cloudflare Pages Documentation](https://developers.cloudflare.com/pages/)
- [Wrangler CLI Documentation](https://developers.cloudflare.com/workers/wrangler/)
- [Open-Same Repository](https://github.com/you112ef/Open-same)

## 📞 الدعم

إذا واجهت أي مشاكل:

1. راجع logs في GitHub Actions
2. تحقق من إعدادات Cloudflare
3. راجع هذا الدليل
4. أنشئ issue في repository

---

**🎉 تهانينا! الآن لديك نظام نشر تلقائي كامل يعمل مع كل push و release!**
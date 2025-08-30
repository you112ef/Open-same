#!/bin/bash

# 🚀 Open-Same Auto-Deployment Setup Script
# هذا السكريبت يساعدك في إعداد النشر التلقائي

set -e

echo "🚀 مرحباً بك في إعداد النشر التلقائي لـ Open-Same!"
echo "=================================================="
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# Check if we're in the right directory
if [ ! -f "wrangler.toml" ]; then
    print_error "يجب تشغيل هذا السكريبت من مجلد المشروع الرئيسي"
    exit 1
fi

print_step "1. التحقق من المتطلبات الأساسية..."

# Check if git is available
if ! command -v git &> /dev/null; then
    print_error "Git غير مثبت. يرجى تثبيته أولاً."
    exit 1
fi

# Check if node is available
if ! command -v node &> /dev/null; then
    print_error "Node.js غير مثبت. يرجى تثبيته أولاً."
    exit 1
fi

# Check if npm is available
if ! command -v npm &> /dev/null; then
    print_error "npm غير مثبت. يرجى تثبيته أولاً."
    exit 1
fi

print_status "✅ جميع المتطلبات الأساسية متوفرة"

print_step "2. التحقق من حالة Git repository..."

# Check git status
if [ -z "$(git status --porcelain)" ]; then
    print_status "✅ Working directory نظيف"
else
    print_warning "⚠️  هناك تغييرات غير محفوظة. يرجى commit أو stash التغييرات أولاً."
    git status --short
    echo ""
    read -p "هل تريد المتابعة؟ (y/N): " -n 1 -r
    echo ""
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_error "تم إلغاء العملية"
        exit 1
    fi
fi

print_step "3. التحقق من GitHub repository..."

# Get remote origin
REMOTE_URL=$(git remote get-url origin 2>/dev/null || echo "")

if [ -z "$REMOTE_URL" ]; then
    print_error "لا يوجد remote origin. يرجى إضافة GitHub repository أولاً."
    echo "مثال: git remote add origin https://github.com/username/repository.git"
    exit 1
fi

print_status "✅ Remote origin: $REMOTE_URL"

print_step "4. إنشاء branch develop (إذا لم يكن موجوداً)..."

# Check if develop branch exists
if ! git show-ref --verify --quiet refs/heads/develop; then
    print_status "إنشاء branch develop..."
    git checkout -b develop
    git push -u origin develop
    print_status "✅ تم إنشاء branch develop"
else
    print_status "✅ branch develop موجود بالفعل"
    git checkout develop
fi

print_step "5. التحقق من ملفات Workflow..."

# Check if workflow files exist
WORKFLOW_DIR=".github/workflows"
if [ ! -d "$WORKFLOW_DIR" ]; then
    print_error "مجلد workflows غير موجود: $WORKFLOW_DIR"
    exit 1
fi

# List workflow files
echo "📁 ملفات Workflow الموجودة:"
ls -la "$WORKFLOW_DIR"/*.yml 2>/dev/null || print_warning "لا توجد ملفات workflow"

print_step "6. إعداد Environment Variables..."

echo ""
echo "🔑 الآن تحتاج إلى إعداد Environment Variables في GitHub:"
echo ""
echo "1. اذهب إلى: https://github.com/$(git remote get-url origin | sed 's/.*github.com[:/]\([^/]*\/[^/]*\).*/\1/')/settings/environments"
echo ""
echo "2. أنشئ Environment جديد باسم 'production':"
echo "   - اضغط على 'New environment'"
echo "   - أدخل الاسم: production"
echo "   - أضف protection rules إذا أردت"
echo ""
echo "3. أنشئ Environment جديد باسم 'staging':"
echo "   - اضغط على 'New environment'"
echo "   - أدخل الاسم: staging"
echo ""
echo "4. أضف Secret جديد:"
echo "   - اذهب إلى Settings → Secrets and variables → Actions"
echo "   - اضغط على 'New repository secret'"
echo "   - الاسم: CLOUDFLARE_API_TOKEN"
echo "   - القيمة: API Token من Cloudflare"
echo ""

print_step "7. إعداد Cloudflare API Token..."

echo "🔑 للحصول على Cloudflare API Token:"
echo ""
echo "1. اذهب إلى: https://dash.cloudflare.com/profile/api-tokens"
echo "2. اضغط على 'Create Token'"
echo "3. اختر 'Custom token'"
echo "4. أضف الأذونات التالية:"
echo "   - Zone:Zone:Read"
echo "   - Zone:Zone Settings:Edit"
echo "   - Account:Cloudflare Pages:Edit"
echo "   - Account:Cloudflare Workers:Edit"
echo "5. انسخ Token وأضفه إلى GitHub Secrets"
echo ""

print_step "8. اختبار الإعداد..."

echo "🧪 لاختبار الإعداد:"
echo ""
echo "1. تأكد من إضافة CLOUDFLARE_API_TOKEN إلى GitHub Secrets"
echo "2. اذهب إلى Actions tab في GitHub"
echo "3. اختر 'Quick Deploy to Cloudflare Pages'"
echo "4. اضغط على 'Run workflow'"
echo "5. اختر البيئة: staging"
echo "6. اضغط على 'Run workflow'"
echo ""

print_step "9. النشر التلقائي..."

echo "🚀 بعد إكمال الإعداد، سيعمل النشر التلقائي:"
echo ""
echo "✅ عند push إلى main → نشر تلقائي إلى Production"
echo "✅ عند push إلى develop → نشر تلقائي إلى Staging"
echo "✅ عند إنشاء tag جديد → نشر تلقائي للإصدار"
echo "✅ نشر سريع عند الطلب من Actions tab"
echo ""

print_step "10. مراقبة النشر..."

echo "📊 لمراقبة النشر:"
echo ""
echo "1. GitHub Actions tab → لمراقبة workflows"
echo "2. Cloudflare Dashboard → Pages → لمراقبة deployments"
echo "3. تقارير مفصلة في كل workflow"
echo ""

# Create a summary file
SUMMARY_FILE="AUTO_DEPLOYMENT_SETUP.md"
cat > "$SUMMARY_FILE" << EOF
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

- \`deploy-cloudflare.yml\` - النشر التلقائي
- \`quick-deploy.yml\` - النشر السريع
- \`release-deploy.yml\` - نشر الإصدارات

## 🚀 كيفية الاستخدام

1. **النشر التلقائي:** push إلى main/develop
2. **النشر السريع:** من Actions tab
3. **نشر الإصدارات:** إنشاء tags جديدة

## 📚 الدليل الكامل

راجع \`docs/auto-deployment.md\` للحصول على دليل مفصل.

---
*تم إنشاء هذا الملف بواسطة setup-auto-deploy.sh*
EOF

print_status "✅ تم إنشاء ملف ملخص: $SUMMARY_FILE"

echo ""
echo "🎉 تم إعداد النشر التلقائي بنجاح!"
echo ""
echo "📋 الخطوات التالية:"
echo "1. اتبع التعليمات أعلاه لإعداد Environment Variables"
echo "2. أضف CLOUDFLARE_API_TOKEN إلى GitHub Secrets"
echo "3. اختبر النشر باستخدام Quick Deploy"
echo "4. راجع $SUMMARY_FILE للمراجعة"
echo ""
echo "🔗 روابط مفيدة:"
echo "- دليل النشر التلقائي: docs/auto-deployment.md"
echo "- GitHub Actions: .github/workflows/"
echo "- Cloudflare Pages: wrangler.toml"
echo ""
echo "🚀 الآن يمكنك النشر تلقائياً مع كل push!"
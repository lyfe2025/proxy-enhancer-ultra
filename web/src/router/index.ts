import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { ElMessage } from 'element-plus'
import {
  Connection,
  SetUp,
  DataAnalysis,
  Setting,
  Odometer,
  PieChart,
  Monitor,
  TrendCharts
} from '@element-plus/icons-vue'

// 导入页面组件
import Login from '@/views/auth/Login.vue'
import AdminLayout from '@/views/layout/AdminLayout.vue'
import Dashboard from '@/views/dashboard/Dashboard.vue'
import ProxyManagement from '@/views/proxy/ProxyManagement.vue'
import RulesManagement from '@/views/rules/RulesManagement.vue'
import DataCollection from '@/views/data/DataCollection.vue'
import SystemSettings from '@/views/system/SystemSettings.vue'
import NotFound from '@/views/error/NotFound.vue'

// 定义路由
const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: Login,
    meta: {
      title: '登录',
      requireAuth: false,
      hidden: true
    }
  },
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/admin',
    name: 'Admin',
    component: AdminLayout,
    redirect: '/dashboard',
    meta: {
      title: '管理后台',
      requireAuth: true
    },
    children: [
      {
        path: '/dashboard',
        name: 'Dashboard',
        component: Dashboard,
        meta: {
          title: '仪表盘',
          icon: 'Odometer',
          requireAuth: true,
          affix: true
        }
      },
      {
        path: '/proxy',
        name: 'ProxyManagement',
        component: ProxyManagement,
        meta: {
          title: '代理管理',
          icon: 'Connection',
          requireAuth: true,
          permissions: ['proxy:view']
        }
      },
      {
        path: '/rules',
        name: 'RulesManagement',
        component: RulesManagement,
        meta: {
          title: '规则配置',
          icon: 'Setting',
          requireAuth: true,
          permissions: ['rules:view']
        }
      },
      {
        path: '/data',
        name: 'DataCollection',
        component: DataCollection,
        meta: {
          title: '数据收集',
          icon: 'DataAnalysis',
          requireAuth: true,
          permissions: ['data:view']
        }
      },
      {
        path: '/settings',
        name: 'SystemSettings',
        component: SystemSettings,
        meta: {
          title: '系统设置',
          icon: 'Tools',
          requireAuth: true,
          permissions: ['system:view']
        }
      }
    ]
  },
  {
    path: '/404',
    name: 'NotFound',
    component: NotFound,
    meta: {
      title: '页面不存在',
      hidden: true
    }
  },
  // 404 页面
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFoundCatch',
    component: NotFound,
    meta: {
      title: '页面未找到',
      requireAuth: false
    }
  }
]

// 创建路由实例
const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// 路由守卫
router.beforeEach(async (to, from, next) => {
  // 设置页面标题
  if (to.meta?.title) {
    document.title = `${to.meta.title} - 代理增强器`
  }

  const authStore = useAuthStore()
  
  // 检查是否需要认证
  if (to.meta?.requireAuth) {
    // 检查是否已登录
    if (!authStore.isAuthenticated) {
      ElMessage.warning('请先登录')
      next({
        path: '/login',
        query: { redirect: to.fullPath }
      })
      return
    }

    // 检查权限
    if (to.meta?.permissions && Array.isArray(to.meta.permissions)) {
      const hasPermission = to.meta.permissions.some(permission => 
        authStore.hasPermission(permission)
      )
      
      if (!hasPermission) {
        ElMessage.error('您没有访问该页面的权限')
        next('/dashboard')
        return
      }
    }

    // 检查角色
    if (to.meta?.roles && Array.isArray(to.meta.roles)) {
      const hasRole = to.meta.roles.some(role => 
        authStore.hasRole(role)
      )
      
      if (!hasRole) {
        ElMessage.error('您没有访问该页面的权限')
        next('/dashboard')
        return
      }
    }
  }

  // 如果已登录且访问登录页，重定向到仪表盘
  if (to.path === '/login' && authStore.isAuthenticated) {
    next('/dashboard')
    return
  }

  next()
})

// 路由错误处理
router.onError((error) => {
  console.error('路由错误:', error)
  ElMessage.error('页面加载失败，请刷新重试')
})

export default router

// 导出路由配置用于菜单生成
export const menuRoutes = routes.find(route => route.path === '/admin')?.children || []
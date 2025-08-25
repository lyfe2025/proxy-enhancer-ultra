import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import { useAuthStore } from '@/store/auth'

// 路由配置
const routes: RouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: {
      title: '登录',
      requiresAuth: false,
      hideInMenu: true
    }
  },
  {
    path: '/',
    redirect: '/dashboard'
  },
  {
    path: '/admin',
    name: 'Admin',
    component: () => import('@/views/layout/AdminLayout.vue'),
    meta: {
      title: '管理后台',
      requiresAuth: true
    },
    children: [
      {
        path: '/dashboard',
        name: 'Dashboard',
        component: () => import('@/views/dashboard/Dashboard.vue'),
        meta: {
          title: '仪表盘',
          icon: 'Dashboard',
          requiresAuth: true
        }
      },
      {
        path: '/proxy',
        name: 'ProxyManagement',
        component: () => import('@/views/proxy/ProxyManagement.vue'),
        meta: {
          title: '代理管理',
          icon: 'Connection',
          requiresAuth: true,
          permissions: ['proxy:read']
        }
      },
      {
        path: '/proxy/create',
        name: 'ProxyCreate',
        component: () => import('@/views/proxy/ProxyForm.vue'),
        meta: {
          title: '创建代理',
          hideInMenu: true,
          requiresAuth: true,
          permissions: ['proxy:create']
        }
      },
      {
        path: '/proxy/edit/:id',
        name: 'ProxyEdit',
        component: () => import('@/views/proxy/ProxyForm.vue'),
        meta: {
          title: '编辑代理',
          hideInMenu: true,
          requiresAuth: true,
          permissions: ['proxy:update']
        }
      },
      {
        path: '/rules',
        name: 'RulesManagement',
        component: () => import('@/views/rules/RulesManagement.vue'),
        meta: {
          title: '规则配置',
          icon: 'Setting',
          requiresAuth: true,
          permissions: ['rule:read']
        }
      },
      {
        path: '/rules/create',
        name: 'RuleCreate',
        component: () => import('@/views/rules/RuleForm.vue'),
        meta: {
          title: '创建规则',
          hideInMenu: true,
          requiresAuth: true,
          permissions: ['rule:create']
        }
      },
      {
        path: '/rules/edit/:id',
        name: 'RuleEdit',
        component: () => import('@/views/rules/RuleForm.vue'),
        meta: {
          title: '编辑规则',
          hideInMenu: true,
          requiresAuth: true,
          permissions: ['rule:update']
        }
      },
      {
        path: '/data-collection',
        name: 'DataCollection',
        component: () => import('@/views/data/DataCollection.vue'),
        meta: {
          title: '数据收集',
          icon: 'DataAnalysis',
          requiresAuth: true,
          permissions: ['popup:read']
        }
      },
      {
        path: '/data-collection/popup/create',
        name: 'PopupCreate',
        component: () => import('@/views/data/PopupForm.vue'),
        meta: {
          title: '创建弹窗',
          hideInMenu: true,
          requiresAuth: true,
          permissions: ['popup:create']
        }
      },
      {
        path: '/data-collection/popup/edit/:id',
        name: 'PopupEdit',
        component: () => import('@/views/data/PopupForm.vue'),
        meta: {
          title: '编辑弹窗',
          hideInMenu: true,
          requiresAuth: true,
          permissions: ['popup:update']
        }
      },
      {
        path: '/data-collection/submissions',
        name: 'Submissions',
        component: () => import('@/views/data/Submissions.vue'),
        meta: {
          title: '提交数据',
          hideInMenu: true,
          requiresAuth: true,
          permissions: ['submission:read']
        }
      },
      {
        path: '/system',
        name: 'SystemSettings',
        component: () => import('@/views/system/SystemSettings.vue'),
        meta: {
          title: '系统设置',
          icon: 'Tools',
          requiresAuth: true,
          permissions: ['system:read']
        }
      },
      {
        path: '/system/users',
        name: 'UserManagement',
        component: () => import('@/views/system/UserManagement.vue'),
        meta: {
          title: '用户管理',
          hideInMenu: true,
          requiresAuth: true,
          permissions: ['user:read']
        }
      },
      {
        path: '/system/roles',
        name: 'RoleManagement',
        component: () => import('@/views/system/RoleManagement.vue'),
        meta: {
          title: '角色管理',
          hideInMenu: true,
          requiresAuth: true,
          permissions: ['role:read']
        }
      },
      {
        path: '/system/permissions',
        name: 'PermissionManagement',
        component: () => import('@/views/system/PermissionManagement.vue'),
        meta: {
          title: '权限管理',
          hideInMenu: true,
          requiresAuth: true,
          permissions: ['permission:read']
        }
      },
      {
        path: '/system/logs',
        name: 'SystemLogs',
        component: () => import('@/views/system/SystemLogs.vue'),
        meta: {
          title: '系统日志',
          hideInMenu: true,
          requiresAuth: true,
          permissions: ['log:read']
        }
      }
    ]
  },
  {
    path: '/404',
    name: 'NotFound',
    component: () => import('@/views/error/404.vue'),
    meta: {
      title: '页面不存在',
      hideInMenu: true
    }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/404'
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
  const authStore = useAuthStore()
  
  // 设置页面标题
  if (to.meta.title) {
    document.title = `${to.meta.title} - 代理增强器`
  }
  
  // 检查是否需要认证
  if (to.meta.requiresAuth) {
    // 检查是否已登录
    if (!authStore.isAuthenticated) {
      // 尝试从本地存储恢复认证状态
      await authStore.checkAuth()
      
      if (!authStore.isAuthenticated) {
        next({
          path: '/login',
          query: { redirect: to.fullPath }
        })
        return
      }
    }
    
    // 检查权限
    if (to.meta.permissions && Array.isArray(to.meta.permissions)) {
      const hasPermission = to.meta.permissions.some(permission => 
        authStore.hasPermission(permission)
      )
      
      if (!hasPermission) {
        // 权限不足，跳转到403页面或首页
        next('/dashboard')
        return
      }
    }
  } else if (to.path === '/login' && authStore.isAuthenticated) {
    // 已登录用户访问登录页，重定向到首页
    next('/dashboard')
    return
  }
  
  next()
})

export default router
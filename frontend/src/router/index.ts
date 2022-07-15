import {createRouter, createWebHistory} from 'vue-router'
import MockView from '@/views/mock/MockView.vue'
import SettingView from '@/views/SettingView.vue'
import RoutesView from '@/views/mock/RoutesView.vue'
import ProxyView from '@/views/mock/ProxyView.vue'
import LogView from '@/views/mock/LogView.vue'
import SettingsView from '@/views/mock/SettingsView.vue'

const router = createRouter({
    history: createWebHistory(import.meta.env.BASE_URL),
    routes: [
        {
            path: '/',
            name: 'root',
            component: MockView,
            children: [
                {path: "", component: RoutesView},
            ]
        },
        {
            path: '/mocks',
            name: 'mocks-view',
            component: MockView,
            children: [
                {path: "", component: RoutesView},
            ]
        },
        {
            path: '/mocks/:id',
            name: 'mock-view',
            component: MockView,
            children: [
                {path: "", component: RoutesView},
                {
                    path: "routes", name: "routes-view", component: RoutesView, children: [
                        {
                            path: ":routeId",
                            name: "route-view",
                            component: RoutesView
                        }
                    ]
                },
                {path: "proxy", name: "route-proxy-view", component: ProxyView},
                {path: "log", name: "route-log-view", component: LogView},
                {path: "settings", name: "route-settings-view", component: SettingsView},
            ]
        },
        {
            path: '/setting',
            name: 'setting',
            component: SettingView
        }
    ]
})

export default router

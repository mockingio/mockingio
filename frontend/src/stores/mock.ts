import {defineStore} from 'pinia'
import axios from 'axios'

export interface MockData {
    id: string
    name: string
    description: string
    routes: Route[]
    url: string
}

export interface MockState {
    url: string
    status: string
    mock_id: string
}

export interface Mock {
    data: MockData
    state: MockState
}

export interface Route {
    id: string
    method: string
    path: string
    description: string
    responses: Response[]
}

export interface Response {
    id: string
    status: number
    body: string
}

export const useMockStore = defineStore({
    id: "mocks",
    state: () => ({
        activeId: null as string | null,
        activeRouteId: null as string | null,
        mocks: [] as Mock[],
        error: null as any,
    }),
    getters: {
        getMockByID(state) {
            return (id: string): Mock | undefined => state.mocks.find(m => m.data.id === id)
        },
        activeRoute(state) {
            const mock = state.mocks.find(m => m.data.id === state.activeId)
            if (!mock) return null
            return mock.data.routes.find(r => r.id === state.activeRouteId)
        },
        activeMock(state): Mock | undefined {
            if (!!state.activeId) {
                return this.getMockByID(state.activeId)
            }
        }
    },
    actions: {
        async stopMockServer(id: string) {
            const mock = this.getMockByID(id)
            if (!mock) return

            const state = await axios.delete<MockState>(getUrl(`/mocks/${id}/stop`))
            this.updateMockState(id, state.data)
        },
        async startMockServer(id: string) {
            const mock = this.getMockByID(id)
            if (!mock) return

            const state = await axios.post<MockState>(getUrl(`/mocks/${id}/start`))
            this.updateMockState(id, state.data)
        },
        updateMockState(id: string, state: MockState) {
            const mock = this.getMockByID(id)
            if (!mock) return

            const idx = this.mocks.findIndex(m => m.data.id === id)
            console.log({state})
            this.mocks[idx] = {...mock, state}
        },
        setActiveMock(id: string) {
            if (id === this.activeId) {
                return
            }
            this.activeId = id
            const mock = this.getMockByID(id)
            if (!!mock) {
                this.setActiveRoute(mock.data.routes[0].id)
            }
        },
        setDefaultActiveMock() {
            if (this.mocks.length === 0) {
                return undefined
            }
            this.setActiveMock(this.mocks[0].data.id)
        },
        setActiveRoute(id: string) {
            if (!id) {
                this.setDefaultActiveRoute()
                return
            }
            this.activeRouteId = id
        },
        setDefaultActiveRoute() {
            if (!this.activeId) {
                return
            }
            const mock = this.getMockByID(this.activeId)
            if (!mock) {
                return
            }
            if (mock.data.routes.length === 0) {
                return
            }

            this.activeRouteId = mock.data.routes[0].id
        },
        async fetchMocks() {
            try {
                const [mockData, stateData] = await Promise.all([
                    axios.get<MockData[]>(getUrl("/mocks")),
                    axios.get<{ [key: string]: MockState }>(getUrl("/mocks/states")),
                ])

                this.mocks = mockData.data.map(mock => ({
                    data: mock,
                    state: stateData.data[mock.id]
                }))
            } catch (error) {
                this.error = {
                    error,
                }
            }
        },
    }
})

const getUrl = (path: string): string => {
    return import.meta.env.VITE_API_ENDPOINT + path
}
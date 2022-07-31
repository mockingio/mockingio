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
    status: string
    body: string
    headers: { [key: string]: string }
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
        async patchRoute(mockId: string, routeId: string, data: { [key: string]: any }) {
            const mockIdx = this.mocks.findIndex(m => m.data.id === mockId)
            if (mockIdx === undefined) {
                return
            }

            const mock = this.mocks[mockIdx]

            const routeIdx = mock.data.routes.findIndex(r => r.id === routeId)
            if (routeIdx === undefined) {
                return
            }

            mock.data.routes[routeIdx] = {...mock.data.routes[routeIdx], ...data}

            this.mocks[mockIdx] = {...mock, data: {...mock.data, routes: [...mock.data.routes]}}
            // send data to server
            return axios.patch(getUrl(`/mocks/${mockId}/routes/${routeId}`), data)
        },
        async patchResponse(mockId: string, routeId: string, responseId: string, data: { [key: string]: any }) {
            const mockIdx = this.mocks.findIndex(m => m.data.id === mockId)
            if (mockIdx === undefined) {
                return
            }
            const mock = this.mocks[mockIdx]

            const routeIdx = mock.data.routes.findIndex(r => r.id === routeId)
            if (routeIdx === undefined) {
                return
            }

            const responseIdx = mock.data.routes[routeIdx].responses.findIndex(r => r.id === responseId)
            if (responseIdx === undefined) {
                return
            }

            mock.data.routes[routeIdx].responses[responseIdx] = {...mock.data.routes[routeIdx].responses[responseIdx], ...data}

            this.mocks[mockIdx] = {...mock, data: {...mock.data, routes: [...mock.data.routes]}}
            // send data to server
            return axios.patch(getUrl(`/mocks/${mockId}/routes/${routeId}/responses/${responseId}`), data)
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
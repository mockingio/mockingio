import {defineStore} from 'pinia'
import axios from 'axios'

export enum Status {
    STOPPED = "stopped",
    RUNNING = "running"
}

export interface MockData {
    id: string
    name: string
    description: string
    routes: Route[]
}

export interface MockState {
    url?: string
    status: Status
}

export interface Mock {
    data: MockData
    state: MockState
}

export interface Route {
    request: string
    responses: Response[]
}

export interface Response {
    status: number
    body: string
}

export const useMockStore = defineStore({
    id: "mocks",
    state: () => ({
        activeId: null as string | null,
        mocks: [] as Mock[],
        error: null as any,
    }),
    getters: {
        getMockByID(state) {
            return (id: string): Mock | undefined => state.mocks.find(m => m.data.id === id)
        },
        activeMock(state): Mock | undefined {
            if (!!state.activeId) {
                return this.getMockByID(state.activeId)
            }

            if (state.mocks.length === 0) {
                return undefined
            }

            this.activeId = state.mocks[0].data.id
            return this.getMockByID(this.activeId)
        }
    },
    actions: {
        setActiveMock(id: string) {
            this.activeId = id
        },
        async fetchMocks(selectedMockId: string) {
            try {
                const {data} = await axios.get<MockData[]>(getUrl("/mocks"))
                this.mocks = data.map(mock => ({
                    data: mock,
                    state: {
                        status: Status.STOPPED
                    }
                }))
                if (!!selectedMockId) {
                    this.setActiveMock(selectedMockId)
                }
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
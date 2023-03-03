import { configureStore, createSlice, PayloadAction } from '@reduxjs/toolkit'

const authSlice = createSlice({
    name: 'auth',
    initialState: {
        token: "",
    },
    reducers: {
        setToken(state, action: PayloadAction<string>) {
            state.token = action.payload
        }
    }
})

export const store = configureStore({
    reducer: authSlice.reducer,
})

export const { setToken } = authSlice.actions

// Infer the `RootState` and `AppDispatch` types from the store itself
export type RootState = ReturnType<typeof store.getState>
// Inferred type: {posts: PostsState, comments: CommentsState, users: UsersState}
export type AppDispatch = typeof store.dispatch

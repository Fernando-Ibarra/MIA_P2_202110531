import { createSlice } from '@reduxjs/toolkit';

export const userSlice = createSlice({
    name: 'user',
    initialState: {
        user: '',
        pass: '',
        status: 'not-authenticated', // 'checking', 'not-authenticated', 'authenticated'
        users: []
    },
    reducers: {
        login: (state, action) => {
            state.user = action.payload.user;
            state.pass = action.payload.pass;
            state.status = 'authenticated';
        },
        checking: (state) => {
            state.status = 'checking';
        },
        logout: (state) => {
            state.user = '';
            state.pass = '';
            state.status = 'not-authenticated';
        },
        setUsers: (state, action) => {
            state.users = action.payload;
        }
    }
});

export const { login, logout, setUsers, checking } = userSlice.actions;
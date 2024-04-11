import { createSlice } from '@reduxjs/toolkit';

export const userSlice = createSlice({
    name: 'user',
    initialState: {
        user: '',
        pass: '',
        grupo: ''
    },
    reducers: {
        login: (state, action) => {
            state.user = action.payload.user;
            state.pass = action.payload.pass;
            state.grupo = action.payload.grupo;
        },
        logout: (state) => {
            state.user = '';
            state.pass = '';
            state.grupo = '';
        }
    }
});

export const { login, logout } = userSlice.actions;
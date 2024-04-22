import { createSlice } from '@reduxjs/toolkit';

export const fileSlice = createSlice({
    name: 'file',
    initialState: {
        file: '',
    },
    reducers: {
        setFile: ( state, { payload }) => {
            state.file = payload
        },
    }
});

export const { setFile } = fileSlice.actions;
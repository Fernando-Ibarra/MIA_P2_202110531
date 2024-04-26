import { createSlice } from '@reduxjs/toolkit';

export const fileSlice = createSlice({
    name: 'file',
    initialState: {
        file: '',
        currentDisk: {},
        currentPartition: {},
    },
    reducers: {
        setFile: ( state, { payload }) => {
            state.file = payload
        },
        setCurrentDisk: (state, { payload }) => {
            state.currentDisk = payload
        },
        setCurrentPartition: (state, { payload }) => {
            state.currentPartition = payload
        }
    }
});

export const { setFile, setCurrentDisk, setCurrentPartition } = fileSlice.actions;
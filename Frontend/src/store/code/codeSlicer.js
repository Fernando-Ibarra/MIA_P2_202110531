import { createSlice } from '@reduxjs/toolkit';

export const codeSlice = createSlice({
    name: 'code',
    initialState: {
        code: '',
        output: '',
    },
    reducers: {
        setCode: ( state, { payload }) => {
            state.code = payload
            state.output = ''
        },
        setOutput: (state, { payload }) => {
            state.output = payload
        },
    }
});

export const { setCode, setOutput } = codeSlice.actions;
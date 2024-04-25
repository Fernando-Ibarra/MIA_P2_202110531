import { createSlice } from '@reduxjs/toolkit';

export const reporteSlice = createSlice({
    name: 'report',
    initialState: {
        report: [],
        currentReport: {},
    },
    reducers: {
        setReport: ( state, { payload }) => {
            state.report = payload
        },
        setCurrentReport: ( state, { payload }) => {
            state.currentReport = payload
        },
    }
});

export const { setReport, setCurrentReport } = reporteSlice.actions;

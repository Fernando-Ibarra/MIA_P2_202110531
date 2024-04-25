import { configureStore } from '@reduxjs/toolkit';
import { userSlice, codeSlice, fileSlice, reporteSlice }  from './';

export const store = configureStore({
    reducer: {
      user: userSlice.reducer,
      file: fileSlice.reducer,
      report: reporteSlice.reducer,
      code: codeSlice.reducer,
    }
});
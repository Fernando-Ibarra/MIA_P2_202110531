import { configureStore } from '@reduxjs/toolkit';
import { userSlice, codeSlice, fileSlice }  from './';

export const store = configureStore({
    reducer: {
      user: userSlice.reducer,
      file: fileSlice.reducer,
      code: codeSlice.reducer,
    }
});
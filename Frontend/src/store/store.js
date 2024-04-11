import { configureStore } from '@reduxjs/toolkit';
import { userSlice}  from './';

export const store = configureStore({
    reducer: {
      user: userSlice.reducer
    }
});
import { Route, Routes } from 'react-router-dom';
import { Console, FileSystem, Reports } from '../dashboard';

export const AppRouter = () => {
    return (
        <Routes>
            <Route path="/" element={<Console />} />
            <Route path="/file-system" element={<FileSystem />} />
            <Route path="/reports" element={<Reports />} />
        </Routes>
    )
}

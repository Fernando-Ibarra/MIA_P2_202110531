import { Route, Routes } from 'react-router-dom';
import { Console, FileSystem, Reports, ViewFile, ViewReport } from '../dashboard';

export const AppRouter = () => {
    return (
        <Routes>
            <Route path="/" element={<Console />} />
            <Route path="/file-system" element={<FileSystem />} />
            <Route path="/reports" element={<Reports />} />
            <Route path="/report" element={<ViewReport />} />
            <Route path="/file" element={<ViewFile />} />
        </Routes>
    )
}

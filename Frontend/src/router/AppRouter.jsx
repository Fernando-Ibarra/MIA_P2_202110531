import { Route, Routes } from 'react-router-dom';
import { Console, FileSystem, Reports, ViewFile, ViewReport } from '../dashboard';
import { Disk, Partition } from '../components';

export const AppRouter = () => {
    return (
        <Routes>
            <Route path="/" element={<Console />} />
            <Route path="/file-system" element={<FileSystem />} />
            <Route path="/reports" element={<Reports />} />
            <Route path="/report" element={<ViewReport />} />
            <Route path="/file" element={<ViewFile />} />
            <Route path="/disk" element={<Disk />} />
            <Route path="/partition" element={<Partition />} />
        </Routes>
    )
}

import { useSelector, useDispatch } from 'react-redux';
import { appApi } from '../api';
import { setReport, setCurrentReport } from '../store';

export const useReport = () => {

    const dispatch = useDispatch();
    const { report, currentReport } = useSelector(state => state.report);

    const makeReport = async (data) => {
        const { data: response } = await appApi.post('/reports', data);
        dispatch(setReport(response));
    }

    const onCurrentReport = (data) => {
        dispatch(setCurrentReport(data));
    }

    return {
        report,
        currentReport,
        makeReport,
        onCurrentReport
    }

}
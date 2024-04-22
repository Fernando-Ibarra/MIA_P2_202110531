import { useSelector, useDispatch } from 'react-redux';

import { setFile } from '../store';
import { appApi } from '../api';

export const useFile = () => {

    const dispatch = useDispatch();
    const { file } = useSelector(state => state.file);

    const makefileSystem = async () => {
        const { data } = await appApi.post('/file-system', {});
        const fixedStr = data.replace(/,(\s*})/g, '$1').replace(/,(\s*\])/g, '$1');
        console.log(fixedStr);
        dispatch(setFile(JSON.parse(fixedStr)));
    }

    return {
        file,
        makefileSystem
    }

}
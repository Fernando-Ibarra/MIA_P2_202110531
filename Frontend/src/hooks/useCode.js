import { useSelector, useDispatch } from 'react-redux';

import { setCode, setOutput } from '../store';
import { appApi } from '../api';

export const useCode = () => {

    const dispatch = useDispatch();
    const { code, output } = useSelector(state => state.code);

    const setActiveCode = (code) => {
        dispatch(setCode(code));
    }

    const setCodeOutput = async () => {
        const dataSend = {
            comands_req: code
        }
        const { data } = await appApi.post('/makeMagic', dataSend);
        const codeOutputString = data.replace("\n", "\n");
        dispatch(setOutput(codeOutputString));
    }

    return {
        code,
        output,
        setActiveCode,
        setCodeOutput
    }
}
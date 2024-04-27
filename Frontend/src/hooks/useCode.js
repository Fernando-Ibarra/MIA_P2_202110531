import { useSelector, useDispatch } from 'react-redux';
import Swal from 'sweetalert2'

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

    const deleteAll = async () => {
        const { data } = await appApi.post('/delete', {});
        Swal.fire({
            position: "bottom-left",
            icon: "success",
            title: `${data}`,
            showConfirmButton: false,
            timer: 1500
        });
        
    }

    return {
        code,
        output,
        setActiveCode,
        setCodeOutput,
        deleteAll
    }
}
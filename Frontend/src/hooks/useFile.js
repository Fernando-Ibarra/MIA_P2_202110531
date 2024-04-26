import { useSelector, useDispatch } from 'react-redux';
import { setFile, setCurrentDisk, setCurrentPartition } from '../store';
import { appApi } from '../api';

export const useFile = () => {

    const dispatch = useDispatch();
    const { file, currentDisk, currentPartition } = useSelector(state => state.file);

    const makefileSystem = async () => {
        const { data } = await appApi.post('/file-system', {});
        const fixedStr = data.replace(/,(\s*})/g, '$1').replace(/,(\s*\])/g, '$1');
        dispatch(setFile(JSON.parse(fixedStr)));
    }

    const onCurrentDisk = (disk) => {
        dispatch(setCurrentDisk(disk));
    }

    const onCurrentPartition = (partition) => {
        dispatch(setCurrentPartition(partition));
    }

    return {
        file,
        currentDisk,
        currentPartition,
        makefileSystem,
        onCurrentDisk,
        onCurrentPartition
    }

}
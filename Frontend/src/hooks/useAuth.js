import { useSelector, useDispatch } from 'react-redux';
import { login, logout, checking, setUsers } from '../store';

export const useAuth = () => {

    const dispatch = useDispatch();
    const { user, pass, status, users } = useSelector(state => state.user);


    const onSetUsers = (users) => {
        dispatch(setUsers(users));
    }

    const onLogin = ({ user, pass }) => {
        let founded = false;
        
        users.forEach((usr) => {
            if (usr.user === user && usr.pass === pass) {
                founded = true;
            }
        });

        if (!founded) {
            alert('User not found');
            return;
        }

        console.log('Login', user, pass);

        dispatch(login({ user, pass }));
    }

    const onLogout = () => {
        dispatch(logout());
    }

    const onChecking = () => {
        dispatch(checking());
    }

    return {
        user,
        pass,
        status,
        users,
        onSetUsers,
        onLogin,
        onLogout,
        onChecking

    }

}
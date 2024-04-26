import { Grid, Stack, Typography, IconButton, TextField, Button, Box, Fab } from '@mui/material';
import { useNavigate } from 'react-router-dom';
import SourceOutlinedIcon from '@mui/icons-material/SourceOutlined';
import ArrowLeftIcon from '@mui/icons-material/ArrowLeft';

import { useAuth, useFile, useForm } from '../../hooks';
import { useEffect } from 'react';
import { Folder } from '../';
import { DrawerComponent } from '../../Layout';

const loginFormFields = {
    loginUserName: '',
    loginPassword: ''
}

export const Partition = () => {

    const { currentPartition } = useFile();
    const { status, onSetUsers, onLogin, onLogout  } = useAuth();
    const { loginUserName, loginPassword, onInputChange } = useForm( loginFormFields );

    const navigate = useNavigate();

    useEffect(() => {
        console.log('Current Partition', currentPartition);
        onSetUsers(currentPartition.users);
    }, []);


    const loginPartition = ( event ) => {
        event.preventDefault();
        const login = {
            user: loginUserName,
            pass: loginPassword
        }
        onLogin(login);
    }

    const backReports = () => {
        navigate('/disk');
    }


    // AUTHENTICATION
    if( status !== 'authenticated' ) {
        return (
            <DrawerComponent
                name="Login"
            >
                <Fab color="primary" aria-label="back"
                    onClick={backReports}
                    sx={{
                        position: 'absolute',
                        bottom: 32,
                        left: 32,
                        backgroundColor: '#ffc300',
                    }}
                >
                    <ArrowLeftIcon />
                </Fab>
                <Box 
                    sx={{
                        borderRadius: '10px',
                        display: 'flex',
                        justifyContent: 'center',
                        alignItems: 'center',
                        padding: '0',
                        margin: '0',
                    }}
                    lg={4}
                    spacing={2}
                >

                    <Grid container
                        sx={{
                            margin: '0',
                            display: 'flex',
                            justifyContent: 'center',
                            alignItems: 'center',
                            width: '50%',
                        }}
                        lg={12}
                        spacing={2}
                    >
                        <Grid item lg={12}
                            
                        >
                            <Stack
                                spacing={2}
                                direction="column"
                                sx={{
                                    justifyContent: 'center',
                                    alignItems: 'center',
                                    marginTop: '90px',
                                }}
                            >
                                <Typography
                                    sx={{
                                        color: '#001d3d',
                                        fontSize: '2.5rem',
                                        fontWeight: 'bold',
                                        padding: '20px',
                                    }}
                                >
                                    Sign In
                                </Typography>
                            </Stack>
                        </Grid>

                        <Grid item lg={12}
                            sx={{
                                display: 'flex',
                                justifyContent: 'center',
                                alignItems: 'center',
                            }}
                        >
                            <Stack
                                spacing={2}
                                direction="column"
                                sx={{
                                    
                                }}
                            >
                                <TextField
                                    required
                                    id="user"
                                    label="user"
                                    defaultValue="user"
                                    variant="outlined"
                                    name="loginUserName"
                                    value={loginUserName}
                                    onChange={onInputChange}
                                    autoComplete="off"
                                    sx={{
                                        color: 'white',
                                        borderColor: 'white',
                                    }}
                                />

                                <TextField
                                    required
                                    id="password"
                                    label="password"
                                    defaultValue="password"
                                    variant="outlined"
                                    name="loginPassword"
                                    value={loginPassword}
                                    onChange={onInputChange}
                                    autoComplete="off"
                                />

                                <Button variant="contained"
                                    onClick={loginPartition}
                                    sx={{
                                        backgroundColor: '#001d3d',
                                        color: 'white',
                                        fontSize: '1rem',
                                        fontWeight: 'bold',
                                        padding: '5px',
                                        '&:hover': {
                                            backgroundColor: 'white',
                                            color: '#001d3d',
                                        },
                                        '&:active': {
                                            backgroundColor: 'white',
                                            color: '#001d3d',
                                        },
                                    }}
                                >
                                    Login
                                </Button>
                            </Stack>
                        </Grid>
                    </Grid>
                </Box>
            </DrawerComponent>
        )
    }


    // SYSTEM PARTITION
    return (
        <DrawerComponent
            name={ `Particiones - ${ currentPartition.name }` }
        >
            <Fab color="primary" aria-label="back"
                onClick={backReports}
                sx={{
                    position: 'absolute',
                    bottom: 32,
                    left: 32,
                    backgroundColor: '#ffc300',
                }}
            >
                <ArrowLeftIcon />
            </Fab>
            <Grid item
                    sx={{
                        bgcolor: '#003566',
                        width: '95%',
                        margin: '0',
                        borderRadius: '10px',
                    }}
                    lg={4}
                >
                    <Stack
                        spacing={2}
                        direction="column"
                        sx={{
                            justifyContent: 'center',
                            alignItems: 'center',
                        }}
                    >
                        <Stack
                            spacing={2}
                            direction="row"
                            sx={{
                                display: 'flex',
                                justifyContent: 'space-around',
                                alignItems: 'center',

                            }}
                        >
                            <IconButton
                                sx={{
                                    color: 'whites',
                                    borderRadius: '10px',
                                }}
                            >
                                <SourceOutlinedIcon 
                                    sx={{
                                        fontSize: '5rem',
                                        color: 'white',
                                    }}
                                />
                            </IconButton>
                            <Typography
                                sx={{
                                    fontSize: '2.5rem',
                                    fontWeight: 'bold',
                                    color: 'white',
                                }}
                            >
                                {currentPartition.name}
                            </Typography>
                            <Button
                                onClick={onLogout}
                                sx={{
                                    backgroundColor: 'red',
                                    color: 'white',
                                    fontSize: '1rem',
                                    fontWeight: 'bold',
                                    padding: '5px',
                                    '&:hover': {
                                        backgroundColor: 'red',
                                        color: '#001d3d',
                                    },
                                    '&:active': {
                                        backgroundColor: 'red',
                                        color: '#001d3d',
                                    },
                                }}
                            >
                                Log Out
                            </Button>
                            
                        </Stack>
                        {
                            (currentPartition.fileSystem && currentPartition.fileSystem.length > 0) 
                            ? (<Folder fld={currentPartition.fileSystem} />)
                            : null
                        }
                    </Stack>
            </Grid>
        </DrawerComponent>
    )
}
import { Box, Grid, Stack, Typography, Avatar, Button, Fab} from '@mui/material'
import { useNavigate } from 'react-router-dom';
import DiscFullIcon from '@mui/icons-material/DiscFull';
import AutoAwesomeMotionIcon from '@mui/icons-material/AutoAwesomeMotion';
import ArrowLeftIcon from '@mui/icons-material/ArrowLeft';
import { useFile } from '../../hooks';
import { DrawerComponent } from '../../Layout';

export const Disk = () => {
    const navigate = useNavigate();
    const { currentDisk, onCurrentPartition } = useFile();

    const handlePartition = (partition) => {
        onCurrentPartition(partition);
        navigate('/partition');
    }

    const backReports = () => {
        navigate('/file-system');
    }

    return (
        <DrawerComponent
            name="Discos"
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
                    bgcolor: '#003566',
                    width: '95%',
                    margin: '0',
                    borderRadius: '10px',
                    padding: '10px',
                }}
            >
                <Stack
                    sx={{
                        backgroundColor: '#001d3d',
                        bgcolor: '001d3d',
                        justifyContent: 'center',
                    }}
                >
                    <Stack
                        spacing={2}
                        direction="row"
                        sx={{
                            justifyContent: 'left',
                            alignItems: 'center',
                        }}
                    >
                        <DiscFullIcon 
                            sx={{
                                color: '#ffd60a',
                                fontSize: '3rem',
                            }}
                        />

                        <Typography
                            sx={{
                                fontSize: '1.5rem',
                                fontWeight: 'bold',
                                color: 'white',
                            }}
                        >
                            {currentDisk.name}
                        </Typography>
                    </Stack>

                    <Grid container
                        spacing={1}
                        sx={{
                            justifyContent: 'center',
                            alignItems: 'center',
                            padding: '10px',
                        }}
                        direction="row"
                    >
                        {
                            currentDisk.partitions.map((part, index) => (
                                <Grid item
                                    lg={6}
                                    key={index}
                                >
                                    <Stack
                                        direction="row"
                                        spacing={2}
                                        sx={{
                                            backgroundColor: '#001d3d',
                                            bgcolor: '001d3d',
                                            justifyContent: 'center',
                                            alignItems: 'center',
                                        }}
                                        
                                    >
                                        <Button
                                            onClick={() => handlePartition(part)}
                                            sx={{
                                                backgroundColor: '#003566',
                                                color: 'white',
                                                fontSize: '1rem',
                                                fontWeight: 'bold',
                                                padding: '10px',
                                                '&:hover': {
                                                    backgroundColor: '#003566',
                                                    color: 'black',
                                                },
                                            }}
                                        >
                                            <Avatar
                                                sx={{
                                                    color: 'black',
                                                    bgcolor: '#ffd60a',
                                                    fontSize: '3rem',
                                                    width: '100px',
                                                    height: '100px',
                                                    '&:hover': {
                                                        color: '#ffd60a',
                                                        bgcolor: 'black'
                                                    },
                                                    '&:active': {
                                                        color: '#ffd60a',
                                                        bgcolor: 'black'
                                                    },
                                                }}
                                            >
                                                <AutoAwesomeMotionIcon 
                                                    sx={{
                                                        fontSize: '3rem',
                                                    }}
                                                />
                                            </Avatar>
                                        </Button>
                                        <Typography
                                            sx={{
                                                fontSize: '1.5rem',
                                                fontWeight: 'bold',
                                                color: 'white',
                                            }}
                                        >
                                            {part.name}
                                        </Typography>
                                    </Stack>
                                </Grid>                            
                            ))
                        }
                    </Grid>
                </Stack>

            </Box>
        </DrawerComponent>
    )
}

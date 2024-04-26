import { Avatar, Box, Button, Grid, Stack, Typography } from '@mui/material';
import { useNavigate } from 'react-router-dom';

import { useFile } from '../hooks';
import { DrawerComponent } from '../Layout';

// import DnsIcon from '@mui/icons-material/Dns';

import StorageIcon from '@mui/icons-material/Storage';

export const FileSystem = () => {

  const navigate = useNavigate();
  const { makefileSystem, file, onCurrentDisk } = useFile();

  const handleMakefileSystem = () => {
    makefileSystem();
  }

  const handleDisk = (disk) => {
    onCurrentDisk(disk);
    navigate('/disk');
  }

  return (
    <DrawerComponent
      name="Sistema de Archivos"
    >
      <Box
        bgcolor='#001d3d'
        sx={{
          backgroundColor: '#001d3d',
          bgcolor: '001d3d',
          justifyContent: 'center',
          alignItems: 'center',
          width: '100%',
          height: '100%',
          padding: '0',
          margin: '0',
          color: 'white',
          borderRadius: '10px',
        }}
      >
        <Box
          sx={{
            fontSize: '1.5rem',
            fontWeight: 'bold',
            padding: '15px',
          }}
        >
          <Stack spacing={2}
          >

          <Stack spacing={2}
            direction="row"
          >
            <Typography
              sx={{
                fontSize: '1.5rem',
                fontWeight: 'bold',
                padding: '20px',
                color: 'white',
              }}
            >
              Sistema de Archivos
            </Typography>
 
            <Button
              size='small'
              variant="contained"
              onClick={handleMakefileSystem}
              sx={{
                backgroundColor: '#ffc300',
                color: 'black',
                fontSize: '1rem',
                fontWeight: 'bold',
                padding: '5px',
                '&:hover': {
                  backgroundColor: '#ffeb3b',
                  color: 'black',
                },
                '&:active': {
                  backgroundColor: '#ffeb3b',
                  color: 'black',
                },

              }}
            >
              Mostrar
            </Button>

          </Stack>
          
          <Grid container
            spacing={1}
          >
            {
              (file) 
              ? (
                file.map((item, index) => (
                  <Grid 
                    item 
                    lg={6} 
                    key={index}
                  >
                    <Stack
                      direction="column"
                      sx={{
                        backgroundColor: '#001d3d',
                        bgcolor: '001d3d',
                        justifyContent: 'center',
                        alignItems: 'center',
                      }}
                    >
                      <Button
                        onClick={() => handleDisk(item)}
                        sx={{
                          color: 'white',
                          fontSize: '1rem',
                          fontWeight: 'bold',
                          padding: '10px',
                          '&:hover': {
                            backgroundColor: '#001d3d',
                            color: 'black',
                          },
                          '&:active': {
                            backgroundColor: '#001d3d',
                            color: 'black',
                          },
                        }}
                      >
                        <Avatar
                          sx={{
                            color: 'white',
                            bgcolor: '#001d3d',
                            fontSize: '3rem',
                            width: '100px',
                            height: '100px',
                            '&:hover': {
                              color: '#001d3d',
                              bgcolor: 'white'
                            },
                          }}
                        >
                          <StorageIcon 
                            sx={{
                              fontSize: '3rem',
                            }}
                          />
                        </Avatar>
                      </Button>
                      <Typography>
                        {item.name}
                      </Typography>
                    </Stack>

                    {/* <Disk disk={item}  /> */}
                  </Grid>
                ))
              )
              : <></>
            }
          </Grid>

          </Stack>
        
        </Box>
      </Box>
    </DrawerComponent>
  )
}
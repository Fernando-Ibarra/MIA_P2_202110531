import { Box, Button, Grid, Stack, Typography } from '@mui/material';
import { DrawerComponent } from '../Layout';
import { useFile } from '../hooks';
import { Disk } from '../components/Disk';

export const FileSystem = () => {
  const { makefileSystem, file } = useFile();

  const handleMakefileSystem = () => {
    makefileSystem();
  }
  console.log(file);
  
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
              variant="contained"
              onClick={handleMakefileSystem}
              sx={{
                backgroundColor: '#ffc300',
                color: 'black',
                fontSize: '1rem',
                fontWeight: 'bold',
                padding: '10px',
              }}
            >
              Mostrar
            </Button>

          </Stack>
          
          <Grid container
            spacing={4}
          >
            {
              (file) 
              ? (
                file.map((item, index) => (
                  <Grid 
                    item 
                    lg={12} 
                    key={index}
                  >
                    <Disk disk={item}  />
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
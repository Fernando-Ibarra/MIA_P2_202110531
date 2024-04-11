import { Box } from '@mui/material';
import { DrawerComponent } from "../Layout";

export const FileSystem = () => {
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
          }}
        >
          <Box
            sx={{
              fontSize: '1.5rem',
              fontWeight: 'bold',
              padding: '20px',
            }}
          >
            Sistema de Archivos
          </Box>
        </Box>
      </DrawerComponent>
    )
}
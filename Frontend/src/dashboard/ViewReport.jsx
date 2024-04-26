import { Fab, Grid, Typography } from '@mui/material';
import { Graphviz } from 'graphviz-react';
import { useReport } from '../hooks';
import { useNavigate } from 'react-router-dom';
import ArrowLeftIcon from '@mui/icons-material/ArrowLeft';


export const ViewReport = () => {

  const { currentReport } = useReport();
  const navigate = useNavigate();


  const backReports = () => {
    navigate('/reports');
  }
    
  return (
    <Grid container
      direction="column"
      justifyContent="center"
      alignItems="center"
      sx={{
        height: '100vh',
        backgroundColor: '#001d3d',
      }}
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
          padding: '20px',
          margin: '20px',
          border: '1px solid #003566',
          borderRadius: '10px',
          backgroundColor: '#003566',
          color: 'white',
          width: '90%',
        }}
      >
        <Typography
          sx={{
            fontSize: '1.5rem',
            fontWeight: 'bold',
            padding: '20px',
            color: 'white',
          }}
          variant="h1"
        >
          { currentReport.name }
        </Typography>
      </Grid>
      <Grid item
        sx={{
          padding: '20px',
          margin: '20px',
          border: '1px solid #003566',
          borderRadius: '10px',
          backgroundColor: '#003566',
          color: 'white',
          width: '90%',
          height: 1000,
        }}
      >
        <Graphviz
          dot={currentReport.content}
          options={{
              width: '100%',
              height: 990,
              zoom: true,
          }}
        />
      </Grid>
    </Grid>
  )
}

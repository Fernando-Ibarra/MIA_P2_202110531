import { Box, Button, Grid, IconButton, Stack, Typography } from '@mui/material';
import { DrawerComponent } from '../Layout';
import { useReport } from '../hooks';
import { useNavigate } from 'react-router-dom';

import SummarizeOutlinedIcon from '@mui/icons-material/SummarizeOutlined';
import ArticleOutlinedIcon from '@mui/icons-material/ArticleOutlined';

export const Reports = () => {

  const navigate = useNavigate();

  const { report, makeReport, onCurrentReport } = useReport();

  const handleReport = () => {
    makeReport();
  }

  const handleReportView = (r) => {
    onCurrentReport(r);
    navigate('/report');
  }

  const handleFileView = (r) => {
    onCurrentReport(r);
    navigate('/file');
  }

  return (
    <DrawerComponent
      name="Reportes"
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
            Reportes
          </Typography>
          <Button
            variant="contained"
            onClick={handleReport}
            sx={{
              backgroundColor: '#ffc300',
              color: 'black',
              fontSize: '1rem',
              fontWeight: 'bold',
              padding: '10px',
            }}
          >
            Listar
          </Button>
        </Stack>
        
        <Grid container
          spacing={3}
          lg={12}
        >
          {
            (report)
            ? (
              report.map((r, i) => {
                if (r.name.includes('.txt')) {
                  return (
                    <Grid item key={i}
                    >
                      <Stack spacing={2}
                        direction="column"
                      >
                        <IconButton aria-label="delete" size="large"
                          onClick={() => handleFileView(r)}
                          sx={{
                            color: '#ffc300'
                          }}
                        >
                          <ArticleOutlinedIcon fontSize="inherit"
                            sx={{
                              fontSize: '3rem',
                            }}
                          />
                        </IconButton>

                        <Typography variant="h5" sx={{ color: 'white' }}>
                          Reporte {r.name}
                        </Typography>
                      </Stack>
                    
                    </Grid>
                  )
                }
                return (
                  <Grid item key={i}>
                    <Stack spacing={2}
                      direction="column"
                    >
                        <IconButton aria-label="delete" size="large"
                          onClick={() => handleReportView(r)}
                          sx={{
                            color: 'white'
                          }}
                        >
                          <SummarizeOutlinedIcon fontSize="inherit"
                            sx={{
                              fontSize: '3rem',
                            }}
                          />
                        </IconButton>

                        <Typography variant="h5" sx={{ color: 'white' }}>
                          Reporte {r.name}
                        </Typography>
                      </Stack>
                  </Grid>
                );
              })
            )
            : null
          }
        </Grid>
        </Stack>
      
      </Box>
      </Box>
    </DrawerComponent>
  )
}

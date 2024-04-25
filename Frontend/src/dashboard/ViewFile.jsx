import { Accordion, AccordionDetails, AccordionSummary, Fab, Grid, Typography } from '@mui/material'
import ArrowDownwardIcon from '@mui/icons-material/ArrowDownward';
import { useReport } from '../hooks';

import { useNavigate } from 'react-router-dom';

import ArrowLeftIcon from '@mui/icons-material/ArrowLeft';

export const ViewFile = () => {

    const navigate = useNavigate();
    const { currentReport } = useReport();

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
                height: 'auto',
              }}
        >
            <Accordion>
                <AccordionSummary
                expandIcon={<ArrowDownwardIcon />}
                aria-controls="panel1-content"
                id="panel1-header"
                >
                <Typography>{ currentReport.name }</Typography>
                </AccordionSummary>
                <AccordionDetails>
                <Typography>
                    { currentReport.content }
                </Typography>
                </AccordionDetails>
            </Accordion>
        </Grid>

        </Grid>
    )
}
